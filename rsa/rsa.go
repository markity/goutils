package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	rrsa "crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

func LoadPrivateKeyFromFile(filePath string) (*rrsa.PrivateKey, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(fileBytes)
	if block == nil {
		return nil, errors.New("incorrect private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func LoadPublicKeyFromFile(filePath string) (*rrsa.PublicKey, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(fileBytes)
	if block == nil {
		return nil, errors.New("incorrect public key")
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pubKeyInterface.(*rsa.PublicKey), nil
}

func RsaEncrypt(publicKey *rrsa.PublicKey, data []byte) ([]byte, error) {
	return rrsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
}

func RsaDecrypt(privateKey *rrsa.PrivateKey, data []byte) ([]byte, error) {
	return rrsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
}
