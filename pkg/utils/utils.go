package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatal(err)
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to base64 encoding
func PublicKeyToBase64(pub *rsa.PublicKey) string {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Fatal(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	pubEnc := b64.RawStdEncoding.EncodeToString(pubBytes)
	return pubEnc
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
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

// Base64ToPublicKey base64 encoded to public key
func Base64ToPublicKey(pub string) *rsa.PublicKey {
	pubDec, err := b64.RawStdEncoding.DecodeString(pub)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	block, _ := pem.Decode(pubDec)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
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

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}

// Split message in multiple chunk before encrypt it to avoid the key size limitation
func ChunkAndEncrypt(msg []byte, pub *rsa.PublicKey) []byte {
	hashSize := sha512.Size //Change it if you change hash function in EncryptWithPublicKey
	msgLen := len(msg)
	step := pub.Size() - 2*hashSize - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes := EncryptWithPublicKey(msg[start:finish], pub)

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes
}

// EncryptWithPublicKey encrypts data with public key and encode it with base64
func Base64EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) string {
	mEncrypted := ChunkAndEncrypt(msg, pub)
	mEncoded := b64.RawStdEncoding.EncodeToString(mEncrypted)

	return mEncoded
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	return plaintext
}

// Decrypt a message which has been encrypted using ChunkAndEncrypt
func DecryptChunked(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	msgLen := len(ciphertext)
	step := priv.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes := DecryptWithPrivateKey(ciphertext[start:finish], priv)

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes
}

// DecryptWithPrivateKey decrypts base64 encoded data with private key
func Base64DecryptWithPrivateKey(ciphertextEnc string, priv *rsa.PrivateKey) []byte {
	//decode
	ciphertext, err := b64.RawStdEncoding.DecodeString(ciphertextEnc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//decrypt
	plaintext := DecryptChunked(ciphertext, priv)
	return plaintext
}
