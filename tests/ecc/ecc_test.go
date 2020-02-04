package ecctest

import (
	"crypto/elliptic"
	"errors"
	"fmt"
	"testing"

	"github.com/markity/goutils/ecc"
)

func TestAll(t *testing.T) {
	if err := ecc.GenerateECCKeyToFile(elliptic.P256(), "private", "public"); err != nil {
		t.Fatal(err)
	}

	ecdsaPri, err := ecc.LoadECCPrivateKeyFromFile(elliptic.P256(), "private")
	if err != nil {
		t.Fatal(err)
	}

	ecdsaPub, err := ecc.LoadECCPublicKeyFromFile(elliptic.P256(), "public")

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
