package main

import (
	"encoding/base64"
	"fmt"
	"github.com/libp2p/go-libp2p/core/crypto"
	"log"
	"os"
)

func loadOrCreatePrivateKey(filepath string) (crypto.PrivKey, error) {
	privKey, err := loadPrivateKeyFromFile(filepath)
	if err == nil {
		log.Println("private key loaded successfully")
		return privKey, nil
	}

	log.Println("file not found, generating a new private key")
	privKey, _, err = crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	if err := savePrivateKeyToFile(filepath, privKey); err != nil {
		return nil, fmt.Errorf("failed to save private key to file: %v", err)
	}

	return privKey, nil
}

func loadPrivateKeyFromFile(filepath string) (crypto.PrivKey, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	privKeyBytes, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}
	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private key: %v", err)
	}

	return privKey, nil
}

func savePrivateKeyToFile(filepath string, privKey crypto.PrivKey) error {
	privKeyBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %v", err)
	}
	privKeyBase64 := base64.StdEncoding.EncodeToString(privKeyBytes)
	err = os.WriteFile(filepath, []byte(privKeyBase64), 0600)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	log.Println("private key saved to file")
	return nil
}
