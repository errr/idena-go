package types

import (
	"idena-go/common"
	"idena-go/crypto/sha3"
	"idena-go/rlp"
	"math/big"
	"sync/atomic"
)

type Network = int32

type Seed [32]byte

func (h Seed) Bytes() []byte { return h[:] }

type EmptyBlockHeader struct {
	ParentHash common.Hash
	Height     uint64
	RootHash   common.Hash
}

type ProposedHeader struct {
	ParentHash     common.Hash
	Height         uint64
	Time           *big.Int `json:"timestamp"        gencodec:"required"`
	TxHash         common.Hash // hash of tx hashes
	ProposerPubKey []byte
	Root           common.Hash
}

type Header struct {
	EmptyBlockHeader *EmptyBlockHeader `rlp:"nil"`
	ProposedHeader   *ProposedHeader   `rlp:"nil"`
}

type VoteHeader struct {
	Round      uint64
	Step       uint16
	ParentHash common.Hash
	VotedHash  common.Hash
}

type Block struct {
	Header *Header

	BlockSeed Seed

	SeedProof []byte

	Body *Body

	// caches
	hash atomic.Value
}

type Body struct {
	Transactions []*Transaction
}

type Transaction struct {
	// pubkey of node that requests approving
	PubKey    []byte
	Signature []byte

	// caches
	hash atomic.Value
}

// Transactions is a Transaction slice type for basic sorting.
type Transactions []*Transaction

type Vote struct {
	Header          *VoteHeader
	Signature       []byte

	// caches
	hash atomic.Value
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.Header.Hash()
	b.hash.Store(v)
	return v
}

func (b *Block) IsEmpty() bool {
	return b.Header.EmptyBlockHeader != nil
}

func (b *Block) Seed() Seed {
	return b.BlockSeed
}
func (b *Block) Height() uint64 {
	if b.IsEmpty() {
		return b.Header.EmptyBlockHeader.Height
	}
	return b.Header.ProposedHeader.Height
}

func (h *Header) Hash() common.Hash {
	if h.ProposedHeader != nil {
		return h.ProposedHeader.Hash()
	}
	return h.EmptyBlockHeader.Hash()
}

func (h *Header) Height() uint64 {
	if h.ProposedHeader != nil {
		return h.ProposedHeader.Height
	}
	return h.EmptyBlockHeader.Height
}

func (h *Header) ParentHash() common.Hash {
	if h.ProposedHeader != nil {
		return h.ProposedHeader.ParentHash
	}
	return h.EmptyBlockHeader.ParentHash
}

func (h *ProposedHeader) Hash() common.Hash {
	return rlpHash(h)
}
func (h *EmptyBlockHeader) Hash() common.Hash {
	return rlpHash(h)
}

func (h *VoteHeader) Hash() common.Hash {
	return rlpHash(h)
}

func (v *Vote) Hash() common.Hash {

	if hash := v.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	h := v.Header.Hash()
	v.hash.Store(h)
	return h
}

func (tx *Transaction) Hash() common.Hash {

	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	h := rlpHash(tx.PubKey)
	tx.hash.Store(h)
	return h
}

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}