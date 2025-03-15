package main

import (
	"encoding/base64"
	"testing"

	"github.com/jamesruan/sodium"
)

func TestEncryptSodiumSecret_ValidKey(t *testing.T) {
	secretValue := "mysecret"
	keyPair := sodium.MakeBoxKP()
	publicKey := base64.StdEncoding.EncodeToString(keyPair.PublicKey.Bytes)
	encryptedValue, err := encryptSodiumSecret(secretValue, publicKey)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	newValue, err := sodium.Bytes(encryptedBytes).SealedBoxOpen(keyPair)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(newValue) != secretValue {
		t.Fatalf("Expected %s, got %s", secretValue, newValue)
	}
}

func TestEncryptSodiumSecret_InvalidKey(t *testing.T) {
	secretValue := "mysecret"
	keyPair := sodium.MakeBoxKP()
	publicKey := base64.StdEncoding.EncodeToString(keyPair.PublicKey.Bytes)
	encryptedValue, err := encryptSodiumSecret(secretValue, publicKey)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	invalid := sodium.MakeBoxKP()
	_, err = sodium.Bytes(encryptedBytes).SealedBoxOpen(invalid)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}
