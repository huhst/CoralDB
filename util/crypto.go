package util

import (
	"CoralDB/common"
	"crypto"
	"crypto/rand"
	"fmt"
	"github.com/r2ishiguro/vrf/go/vrf_ed25519"
	"golang.org/x/crypto/ed25519"
)

type PublicKey struct {
	Pk ed25519.PublicKey
}

func (pub *PublicKey) Bytes() []byte {
	return pub.Pk
}

func (pub *PublicKey) Address() common.Address {
	return common.BytesToAddress(pub.Pk)
}

func (pub *PublicKey) VerifySign(m, sign []byte) error {
	signature := sign[ed25519.PublicKeySize:]
	if ok := ed25519.Verify(pub.Pk, m, signature); !ok {
		return fmt.Errorf("signature invalid")
	}
	return nil
}

func (pub *PublicKey) VerifyVRF(proof, m []byte) error {
	vrf_ed25519.ECVRF_verify(pub.Pk, proof, m)
	return nil
}

type PrivateKey struct {
	Sk ed25519.PrivateKey
}

func (priv *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{priv.Sk.Public().(ed25519.PublicKey)}
}

func (priv *PrivateKey) Sign(m []byte) ([]byte, error) {
	sign, err := priv.Sk.Sign(rand.Reader, m, crypto.Hash(0))
	if err != nil {
		return nil, err
	}
	pubkey := priv.Sk.Public().(ed25519.PublicKey)
	return append(pubkey, sign...), nil
}

func (priv *PrivateKey) Evaluate(m []byte) (value, proof []byte, err error) {
	proof, err = vrf_ed25519.ECVRF_prove(priv.PublicKey().Pk, priv.Sk, m)
	if err != nil {
		return
	}
	value = vrf_ed25519.ECVRF_proof2hash(proof)
	return
}

func NewKeyPair() (*PublicKey, *PrivateKey, error) {
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	return &PublicKey{pk}, &PrivateKey{sk}, nil
}

func RecoverPubkey(sign []byte) *PublicKey {
	pubkey := sign[:ed25519.PublicKeySize]
	return &PublicKey{pubkey}
}
