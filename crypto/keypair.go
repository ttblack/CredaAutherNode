package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ttblack/CredaAutherNode/crypto/secp256k1"
)

const PrivateKeyLength = 32

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

type Keypair struct {
	public  *ecdsa.PublicKey
	private *ecdsa.PrivateKey
}

func NewKeypairFromString(priv string) (*Keypair, error) {
	pk, err := HexToECDSA(priv)
	if err != nil {
		return nil, err
	}

	return &Keypair{
		public:  pk.Public().(*ecdsa.PublicKey),
		private: pk,
	}, nil
}

func NewKeypair(prvkey *ecdsa.PrivateKey) (*Keypair, error) {
	return &Keypair{
		public:  prvkey.Public().(*ecdsa.PublicKey),
		private: prvkey,
	}, nil
}

// CommonAddress returns the Ethereum address in the common.Address Format
func (kp *Keypair) CommonAddress() common.Address {
	return crypto.PubkeyToAddress(*kp.public)
}

// PrivateKey returns the keypair's private key
func (kp *Keypair) PrivateKey() *ecdsa.PrivateKey {
	return kp.private
}

// HexToECDSA parses a secp256k1 private key.
func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return ToECDSA(b)
}

// ToECDSA creates a private key with the given D value.
func ToECDSA(d []byte) (*ecdsa.PrivateKey, error) {
	return toECDSA(d, true)
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = S256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
}

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return secp256k1.S256()
}
