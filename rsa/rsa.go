package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	rrsa "crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func GenerateRsaKeyToFile(bits int, privateKeyPath string, publicKeyPath string) (err error) {
	pri, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to generate private key: %v", err))
		return
	}
	priDerStream := x509.MarshalPKCS1PrivateKey(pri)
	priBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: priDerStream,
	}
	priFile, err := os.Create(privateKeyPath)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to create file for private key: %v", err))
		return
	}
	if err = pem.Encode(priFile, priBlock); err != nil {
		err = errors.New(fmt.Sprintf("failed to encode private key into file: %v", err))
		return
	}

	pub := &pri.PublicKey
	pubDerStream, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to marshal public key: %v", err))
		return
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubDerStream,
	}
	pubFile, err := os.Create(publicKeyPath)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to create file for public key: %v", err))
		return
	}
	if err = pem.Encode(pubFile, pubBlock); err != nil {
		err = errors.New(fmt.Sprintf("failed to encode public key into file: %v", err))
		return
	}

	return nil
}

// LoadPKCS1PrivateKeyFromFile loads a rsa private key from a file
func LoadPKCS1PrivateKeyFromFile(filePath string) (*rrsa.PrivateKey, error) {
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

// LoadPublicKeyFromFile loads a rsa public key from a file
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

	pubKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("the public key is not rsa.PublicKey")
	}

	return pubKey, nil
}

// RsaEncrypt encrypts clear text by using rsa public key, and returns the cipher text
func RsaEncrypt(publicKey *rrsa.PublicKey, clearText []byte) ([]byte, error) {
	return rrsa.EncryptPKCS1v15(rand.Reader, publicKey, clearText)
}

// RsaDecrypt decrypts cipher text by using rsa private key, and returns the clear text
func RsaDecrypt(privateKey *rrsa.PrivateKey, cipherText []byte) ([]byte, error) {
	return rrsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
}

// RsaSign signs clear text by using rsa private key, and returns the signature
func RsaSign(privateKey *rrsa.PrivateKey, hash crypto.Hash, clearText []byte) ([]byte, error) {
	h := hash.New()
	h.Write(clearText)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, hash, h.Sum(nil))
}

// RsaVerifySign verifies signature by using public key and clear text
func RsaVerifySign(publicKey *rrsa.PublicKey, hash crypto.Hash, signature []byte, clearText []byte) bool {
	h := hash.New()
	h.Write(clearText)
	if err := rsa.VerifyPKCS1v15(publicKey, hash, h.Sum(nil), signature); err != nil {
		return false
	}
	return true
}
