package consensus

import (
	"CoralDB/common"
	"CoralDB/util"
	"bytes"
	"encoding/json"
	"errors"
)

// 定义信息类型
const (
	VOTE = iota
	BlockProposal
	BLOCK
)

type VoteMessage struct {
	Signature  []byte `json:"signature"`
	Round      uint64 `json:"Round"`
	Step       int    `json:"step"`
	VRF        []byte `json:"vrf"`
	Proof      []byte `json:"proof"`
	ParentHash []byte `json:"parentHash"`
	Hash       []byte `json:"hash"`
}

func (v *VoteMessage) Serialize() ([]byte, error) {
	return json.Marshal(v)
}

func (v *VoteMessage) Deserialize(data []byte) error {
	return json.Unmarshal(data, v)
}

func (v *VoteMessage) VerifySign() error {
	pubkey := v.RecoverPubkey()
	data := bytes.Join([][]byte{
		common.Uint2Bytes(v.Round),
		common.Uint2Bytes(uint64(v.Step)),
		v.VRF,
		v.Proof,
		v.ParentHash,
		v.Hash,
	}, nil)
	return pubkey.VerifySign(data, v.Signature)
}

func (v *VoteMessage) Sign(priv *util.PrivateKey) ([]byte, error) {
	data := bytes.Join([][]byte{
		common.Uint2Bytes(v.Round),
		common.Uint2Bytes(uint64(v.Step)),
		v.VRF,
		v.Proof,
		v.ParentHash,
		v.Hash,
	}, nil)
	sign, err := priv.Sign(data)
	if err != nil {
		return nil, err
	}
	v.Signature = sign
	return sign, nil
}

func (v *VoteMessage) RecoverPubkey() *util.PublicKey {
	return util.RecoverPubkey(v.Signature)
}

type Proposal struct {
	Round  uint64 `json:"Round"`
	Hash   []byte `json:"hash"`
	Prior  []byte `json:"prior"`
	VRF    []byte `json:"vrf"` // vrf of user's Sortition hash
	Proof  []byte `json:"proof"`
	Pubkey []byte `json:"public_key"`
}

func (b *Proposal) Serialize() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Proposal) Deserialize(data []byte) error {
	return json.Unmarshal(data, b)
}

func (b *Proposal) PublicKey() *util.PublicKey {
	return &util.PublicKey{Pk: b.Pubkey}
}

func (b *Proposal) Address() common.Address {
	return common.BytesToAddress(b.Pubkey)
}

func (b *Proposal) Verify(weight uint64, m []byte) error {
	// verify vrf
	pubkey := b.PublicKey()
	if err := pubkey.VerifyVRF(b.Proof, m); err != nil {
		return err
	}

	// verify priority
	subusers := subUsers(util.ExpectedBlockProposers, weight, b.VRF)
	if bytes.Compare(maxPriority(b.VRF, subusers), b.Prior) != 0 {
		return errors.New("max priority mismatch")
	}

	return nil
}
