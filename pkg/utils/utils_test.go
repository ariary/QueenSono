package utils

import (
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	priv, pub := GenerateKeyPair(2048)
	if priv == nil {
		t.Fatal("expected non-nil private key")
	}
	if pub == nil {
		t.Fatal("expected non-nil public key")
	}
	if priv.PublicKey.N.Cmp(pub.N) != 0 {
		t.Fatal("public key does not match private key's public component")
	}
}

func TestEncryptDecryptRoundTrip_Short(t *testing.T) {
	priv, pub := GenerateKeyPair(2048)
	msg := []byte("hello world")
	encrypted := Base64EncryptWithPublicKey(msg, pub)
	decrypted := Base64DecryptWithPrivateKey(encrypted, priv)
	if string(decrypted) != string(msg) {
		t.Fatalf("round-trip mismatch: expected %q, got %q", msg, decrypted)
	}
}

func TestEncryptDecryptRoundTrip_Long(t *testing.T) {
	// 2048-bit key → ~190 bytes usable per block; 500 bytes forces multi-block path.
	priv, pub := GenerateKeyPair(2048)
	msg := make([]byte, 500)
	for i := range msg {
		msg[i] = byte(i % 256)
	}
	encrypted := Base64EncryptWithPublicKey(msg, pub)
	decrypted := Base64DecryptWithPrivateKey(encrypted, priv)
	if string(decrypted) != string(msg) {
		t.Fatal("long message round-trip mismatch")
	}
}

func TestPublicKeyBase64RoundTrip(t *testing.T) {
	_, pub := GenerateKeyPair(2048)
	encoded := PublicKeyToBase64(pub)
	decoded := Base64ToPublicKey(encoded)
	if pub.N.Cmp(decoded.N) != 0 {
		t.Fatal("public key base64 round-trip: modulus mismatch")
	}
	if pub.E != decoded.E {
		t.Fatal("public key base64 round-trip: exponent mismatch")
	}
}
