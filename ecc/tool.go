package ecc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
)

var errInvalidPubkey = errors.New("invalid secp256k1 public key")
var secp256k1N, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

func toolFromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(pub.Curve, pub.X, pub.Y)
}

func toolToECDSAPub(curve elliptic.Curve, pubBytes []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(curve, pubBytes)
	if x == nil {
		return nil, errInvalidPubkey
	}
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}

func toolFromECDSAPri(pri *ecdsa.PrivateKey) []byte {
	if pri == nil {
		return nil
	}
	return math.PaddedBigBytes(pri.D, pri.Params().BitSize/8)
}

func toolToECDSAPri(curve elliptic.Curve, priBytes []byte) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	if 8*len(priBytes) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(priBytes)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(priBytes)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
}
