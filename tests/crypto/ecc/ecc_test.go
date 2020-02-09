package ecctest

import (
	"errors"
	"fmt"
	"testing"

	"github.com/markity/goutils/crypto/ecc"
)

func TestP256(t *testing.T) {
	curve := ecc.CurveP256()
	if err := ecc.GenerateECDSAKeyToFile(curve, "private", "public"); err != nil {
		t.Fatal(err)
	}

	ecdsaPri, err := ecc.LoadECDSAPrivateKeyFromFile(curve, "private")
	if err != nil {
		t.Fatal(err)
	}

	ecdsaPub, err := ecc.LoadECDSAPublicKeyFromFile(curve, "public")

	clearText := []byte("神秘的密文, I am secret")

	cipherText, err := ecc.ECCEncrypt(ecdsaPub, clearText)
	if err != nil {
		t.Fatal(err)
	}

	clearText2, err := ecc.ECCDecrypt(ecdsaPri, cipherText)
	if string(clearText) != string(clearText2) {
		t.Fatal(errors.New(fmt.Sprintf("text does not match: '%s' and '%s'", clearText, clearText2)))
	}

	signature, err := ecc.ECCSign(ecdsaPri, clearText)
	if err != nil {
		t.Fatal(err)
	}

	if ecc.ECCVerifySignature(ecdsaPub, signature, clearText) != true {
		t.Fatal("verify signature failed")
	}
}

func TestS256(t *testing.T) {
	curve := ecc.CurveS256()
	if err := ecc.GenerateECDSAKeyToFile(curve, "private", "public"); err != nil {
		t.Fatal(err)
	}

	ecdsaPri, err := ecc.LoadECDSAPrivateKeyFromFile(curve, "private")
	if err != nil {
		t.Fatal(err)
	}

	ecdsaPub, err := ecc.LoadECDSAPublicKeyFromFile(curve, "public")

	clearText := []byte("神秘的密文, I am secret")

	cipherText, err := ecc.ECCEncrypt(ecdsaPub, clearText)
	if err != nil {
		t.Fatal(err)
	}

	clearText2, err := ecc.ECCDecrypt(ecdsaPri, cipherText)
	if string(clearText) != string(clearText2) {
		t.Fatal(errors.New(fmt.Sprintf("text does not match: '%s' and '%s'", clearText, clearText2)))
	}

	signature, err := ecc.ECCSign(ecdsaPri, clearText)
	if err != nil {
		t.Fatal(err)
	}

	if ecc.ECCVerifySignature(ecdsaPub, signature, clearText) != true {
		t.Fatal("verify signature failed")
	}
}
