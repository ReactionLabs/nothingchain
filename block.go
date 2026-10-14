package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Transaction moves nothing, from someone, to someone, beautifully.
type Transaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"` // MUST be 0. Enforced by consensus.
	Memo   string `json:"memo"`
}

func (tx Transaction) Validate() error {
	if tx.From == "" || tx.To == "" {
		return fmt.Errorf("tx needs a from and a to, even if nothing moves")
	}
	if tx.Amount != 0 {
		return fmt.Errorf("invalid amount %d: NothingChain only moves nothing (0 NADA)", tx.Amount)
	}
	return nil
}

// Block is a batch of nothing, sealed with wasted computation.
type Block struct {
	Index     int           `json:"index"`
	Timestamp int64         `json:"timestamp"`
	Txs       []Transaction `json:"txs"`
	PrevHash  string        `json:"prev_hash"`
	Nonce     int64         `json:"nonce"`
	Reward    int64         `json:"reward"` // always 0
	Hash      string        `json:"hash"`
}

func (b *Block) computeHash() string {
	payload, _ := json.Marshal(struct {
		Index     int           `json:"index"`
		Timestamp int64         `json:"timestamp"`
		Txs       []Transaction `json:"txs"`
		PrevHash  string        `json:"prev_hash"`
		Nonce     int64         `json:"nonce"`
		Reward    int64         `json:"reward"`
	}{b.Index, b.Timestamp, b.Txs, b.PrevHash, b.Nonce, b.Reward})
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

// Mine performs Proof of Nothing: burn CPU until the hash looks important.
func (b *Block) Mine(difficulty int) {
	target := strings.Repeat("0", difficulty)
	for {
		b.Hash = b.computeHash()
		if strings.HasPrefix(b.Hash, target) {
			return
		}
		b.Nonce++
	}
}

func NewBlock(index int, prevHash string, txs []Transaction) *Block {
	if txs == nil {
		txs = []Transaction{}
	}
	return &Block{
		Index:     index,
		Timestamp: time.Now().Unix(),
		Txs:       txs,
		PrevHash:  prevHash,
		Reward:    0, // always
	}
}

func GenesisBlock() *Block {
	b := NewBlock(0, "0", []Transaction{
		{From: "nobody", To: "everybody", Amount: 0, Memo: "In nothing we trust."},
	})
	b.Timestamp = 1728000000 // fixed, so genesis is deterministic-ish and noble
	b.Mine(4)
	return b
}
