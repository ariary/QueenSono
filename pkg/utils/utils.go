package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

// GenerateKeyPair generates a new RSA key pair.
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatal(err)
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes encodes a private key as PEM bytes.
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
}

// PublicKeyToBase64 encodes a public key as base64-encoded PEM.
func PublicKeyToBase64(pub *rsa.PublicKey) string {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Fatal(err)
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	return base64.RawStdEncoding.EncodeToString(pubBytes)
}

// BytesToPrivateKey decodes PEM bytes into a private key.
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block) //nolint:staticcheck
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil) //nolint:staticcheck
		if err != nil {
			log.Fatal(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

// Base64ToPublicKey decodes a base64-encoded PEM public key.
func Base64ToPublicKey(pub string) *rsa.PublicKey {
	pubDec, err := base64.RawStdEncoding.DecodeString(pub)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	block, _ := pem.Decode(pubDec)
	enc := x509.IsEncryptedPEMBlock(block) //nolint:staticcheck
	b := block.Bytes
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil) //nolint:staticcheck
		if err != nil {
			log.Fatal(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Fatal(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Fatal("not ok")
	}
	return key
}

// EncryptWithPublicKey encrypts data with RSA-OAEP using SHA-512.
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}

// ChunkAndEncrypt encrypts msg in RSA-sized blocks to work around key-size limits.
func ChunkAndEncrypt(msg []byte, pub *rsa.PublicKey) []byte {
	hashSize := sha512.Size
	msgLen := len(msg)
	step := pub.Size() - 2*hashSize - 2
	var encryptedBytes []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		encryptedBytes = append(encryptedBytes, EncryptWithPublicKey(msg[start:finish], pub)...)
	}
	return encryptedBytes
}

// Base64EncryptWithPublicKey encrypts and base64-encodes the message.
func Base64EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) string {
	return base64.RawStdEncoding.EncodeToString(ChunkAndEncrypt(msg, pub))
}

// DecryptWithPrivateKey decrypts a single RSA-OAEP block.
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	return plaintext
}

// DecryptChunked decrypts a message encrypted with ChunkAndEncrypt.
func DecryptChunked(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	msgLen := len(ciphertext)
	step := priv.PublicKey.Size()
	var decryptedBytes []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		decryptedBytes = append(decryptedBytes, DecryptWithPrivateKey(ciphertext[start:finish], priv)...)
	}
	return decryptedBytes
}

// Base64DecryptWithPrivateKey base64-decodes then decrypts the message.
func Base64DecryptWithPrivateKey(ciphertextEnc string, priv *rsa.PrivateKey) []byte {
	ciphertext, err := base64.RawStdEncoding.DecodeString(ciphertextEnc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return DecryptChunked(ciphertext, priv)
}
