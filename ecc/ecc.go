package ecc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// GenerateECCKeyToFile generates ECC private key and public key, and writes them into file
func GenerateECCKeyToFile(curve elliptic.Curve, privateKeyPath string, publicKeyPath string) error {
	pri, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to generate private key: %v", err))
	}
	priFile, err := os.Create(privateKeyPath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create file for private key: %v", err))
	}
	if _, err := priFile.Write([]byte(hex.EncodeToString(toolFromECDSAPri(pri)))); err != nil {
		return errors.New(fmt.Sprintf("failed to write private key into file: %v", err))
	}
	if err := priFile.Close(); err != nil {
		return errors.New(fmt.Sprintf("failed to close private key file: %v", err))
	}

	pub := &pri.PublicKey
	pubFile, err := os.Create(publicKeyPath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create file for public key: %v", err))
	}
	if _, err := pubFile.Write([]byte(hex.EncodeToString(toolFromECDSAPub(pub)))); err != nil {
		return errors.New(fmt.Sprintf("failed to write public key into file: %v", err))
	}
	if err := pubFile.Close(); err != nil {
		return errors.New(fmt.Sprintf("failed to close public key file: %v", err))
	}

	return nil
}

// LoadECCPrivateKeyFromFile loads ECC private key from file
func LoadECCPrivateKeyFromFile(curve elliptic.Curve, filePath string) (*ecdsa.PrivateKey, error) {
	priBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	priBytes, err = hex.DecodeString(string(priBytes))
	if err != nil {
		return nil, err
	}

	return toolToECDSAPri(curve, priBytes)
}

// LoadECCPublicKeyFromFile Loads ECC public key from file
func LoadECCPublicKeyFromFile(curve elliptic.Curve, filePath string) (*ecdsa.PublicKey, error) {
	pubBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	pubBytes, err = hex.DecodeString(string(pubBytes))
	if err != nil {
		return nil, err
	}

	return toolToECDSAPub(curve, pubBytes)
}

// ECCEncrypt encrypts clear text by using public key, and returns the cipher text
func ECCEncrypt(ecdsaPubilcKey *ecdsa.PublicKey, clearText []byte) ([]byte, error) {
	eciesPublicKey := ecies.ImportECDSAPublic(ecdsaPubilcKey)
	return ecies.Encrypt(rand.Reader, eciesPublicKey, clearText, nil, nil)
}

// ECCDecrypt decrypts cipher text by using private key, and returns the clear text
func ECCDecrypt(ecdsaPrivateKey *ecdsa.PrivateKey, cipherText []byte) ([]byte, error) {
	eciesPrivateKey := ecies.ImportECDSA(ecdsaPrivateKey)
	return eciesPrivateKey.Decrypt(cipherText, nil, nil)
}

// ECCSign signs clear text by using private key, and returns the signature
func ECCSign(ecdsaPrivateKey *ecdsa.PrivateKey, clearText []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, ecdsaPrivateKey, clearText)
	if err != nil {
		return nil, err
	}

	rText, _ := r.MarshalText()
	sText, _ := s.MarshalText()

	var buf bytes.Buffer
	buf.Write(rText)
	buf.Write([]byte("+"))
	buf.Write(sText)
	return buf.Bytes(), nil
}

// ECCVerifySignature verifies signature by using public key and clear text
func ECCVerifySignature(ecdsaPublicKey *ecdsa.PublicKey, signature []byte, clearText []byte) bool {
	var r, s big.Int
	rs := bytes.Split(signature, []byte("+"))
	if err := r.UnmarshalText(rs[0]); err != nil {
		return false
	}
	if err := s.UnmarshalText(rs[1]); err != nil {
		return false
	}

	return ecdsa.Verify(ecdsaPublicKey, clearText, &r, &s)
}
