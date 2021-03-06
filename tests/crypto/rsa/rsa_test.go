package rsatest

import (
	"crypto"
	"errors"
	"fmt"
	"testing"

	"github.com/markity/goutils/crypto/rsa"
)

func TestAll(t *testing.T) {
	if err := rsa.GenerateRSAKeyToFile(2048, "app_private_key.pem", "app_public_key.pem"); err != nil {
		t.Fatal(err)
	}

	pri, err := rsa.LoadPKCS1PrivateKeyFromFile("app_private_key.pem")
	if err != nil {
		t.Fatal(err)
	}

	pub, err := rsa.LoadPublicKeyFromFile("app_public_key.pem")
	if err != nil {
		t.Fatal(err)
	}

	clearText := []byte("神秘的密文, it's very secert")

	cipher, err := rsa.RSAEncrypt(pub, clearText)
	if err != nil {
		t.Fatal(err)
	}

	if clearText2, err := rsa.RSADecrypt(pri, cipher); err != nil {
		t.Fatal(err)
	} else if string(clearText) != string(clearText2) {
		t.Fatal(errors.New(fmt.Sprintf("text does not match: '%s' and '%s'", clearText, clearText2)))
	}

	signature, err := rsa.RSASign(pri, crypto.SHA256, clearText)
	if err != nil {
		t.Fatal(err)
	}

	if !rsa.RSAVerifySignature(pub, crypto.SHA256, signature, clearText) {
		t.Fatal(errors.New("verify sign failed"))
	}
}
