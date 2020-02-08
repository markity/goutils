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

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

type curve struct {
	c elliptic.Curve
}

func CurveP256() *curve {
	return &curve{elliptic.P256()}
}

func CurveS256() *curve {
	return &curve{crypto.S256()}
}

// GenerateECCKeyToFile generates ECC private key and public key, and writes them into file
func GenerateECCKeyToFile(cur *curve, privateKeyPath string, publicKeyPath string) error {
	pri, err := ecdsa.GenerateKey(cur.c, rand.Reader)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to generate private key: %v", err))
	}
	priBytes, err := toolFromECDSAPri(pri)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to marshal private key into bytes: %v", err))
	}
	if err := ioutil.WriteFile(privateKeyPath, []byte(hex.EncodeToString(priBytes)), 0600); err != nil {
		return errors.New(fmt.Sprintf("failed to write private key into file: %v", err))
	}

	pub := &pri.PublicKey
	pubBytes, err := toolFromECDSAPub(pub)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to marshal public key into bytes: %v", err))
	}
	if err := ioutil.WriteFile(publicKeyPath, []byte(hex.EncodeToString(pubBytes)), 0600); err != nil {
		return errors.New(fmt.Sprintf("failed to write public key into file: %v", err))
	}

	return nil
}

// LoadECCPrivateKeyFromFile loads ECC private key from file
func LoadECCPrivateKeyFromFile(cur *curve, filePath string) (*ecdsa.PrivateKey, error) {
	priBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	priBytes, err = hex.DecodeString(string(priBytes))
	if err != nil {
		return nil, err
	}

	return toolToECDSAPri(cur.c, priBytes)
}

// LoadECCPublicKeyFromFile Loads ECC public key from file
func LoadECCPublicKeyFromFile(cur *curve, filePath string) (*ecdsa.PublicKey, error) {
	pubBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	pubBytes, err = hex.DecodeString(string(pubBytes))
	if err != nil {
		return nil, err
	}

	return toolToECDSAPub(cur.c, pubBytes)
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
