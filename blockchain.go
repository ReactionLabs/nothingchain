package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

const Difficulty = 4 // leading hex zeros. Pointless, but consistently pointless.

type Blockchain struct {
	mu      sync.Mutex
	Blocks  []*Block      `json:"blocks"`
	Pending []Transaction `json:"pending"`
}

func NewBlockchain() *Blockchain {
	return &Blockchain{Blocks: []*Block{GenesisBlock()}}
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) AddTx(tx Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.Pending = append(bc.Pending, tx)
	return nil
}

// MinePending seals all pending nothing into a new block and rewards 0.
func (bc *Blockchain) MinePending(miner string) *Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	txs := bc.Pending
	bc.Pending = []Transaction{}
	// The coinbase of nothing: miner earns 0 for their service.
	txs = append([]Transaction{{From: "network", To: miner, Amount: 0, Memo: "block reward: nothing, with thanks"}}, txs...)
	b := NewBlock(len(bc.Blocks), bc.LastBlock().Hash, txs)
	b.Mine(Difficulty)
	bc.Blocks = append(bc.Blocks, b)
	return b
}

// Balance is the hardest function in finance: it is always 0.
func (bc *Blockchain) Balance(addr string) int64 {
	return 0
}

func (bc *Blockchain) TotalTransferred() int64 {
	var total int64
	for _, b := range bc.Blocks {
		for _, tx := range b.Txs {
			total += tx.Amount
		}
	}
	return total // spoiler: 0
}

func (bc *Blockchain) Validate() error {
	target := strings.Repeat("0", Difficulty)
	for i, b := range bc.Blocks {
		for _, tx := range b.Txs {
			if err := tx.Validate(); err != nil {
				return fmt.Errorf("block %d: %w", i, tx.Validate())
			}
		}
		if b.Reward != 0 {
			return fmt.Errorf("block %d: reward is %d, must be nothing", i, b.Reward)
		}
		if b.Hash != b.computeHash() {
			return fmt.Errorf("block %d: hash mismatch, the nothing has been tampered with", i)
		}
		if i > 0 {
			if !strings.HasPrefix(b.Hash, target) {
				return fmt.Errorf("block %d: insufficient proof of nothing", i)
			}
			if b.PrevHash != bc.Blocks[i-1].Hash {
				return fmt.Errorf("block %d: broken link to previous nothing", i)
			}
		}
	}
	return nil
}

func (bc *Blockchain) Save(path string) error {
	data, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadBlockchain(path string) (*Blockchain, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var bc Blockchain
	if err := json.Unmarshal(data, &bc); err != nil {
		return nil, err
	}
	return &bc, nil
}

// LoadOrCreate returns the saved chain, or a fresh genesis.
func LoadOrCreate(path string) *Blockchain {
	if bc, err := LoadBlockchain(path); err == nil && len(bc.Blocks) > 0 {
		if err := bc.Validate(); err == nil {
			return bc
		}
	}
	return NewBlockchain()
}

type Stats struct {
	Height           int    `json:"height"`
	Blocks           int    `json:"blocks"`
	PendingTxs       int    `json:"pending_txs"`
	TotalTransferred int64  `json:"total_transferred"`
	TotalSupply      int64  `json:"total_supply"`
	MaxSupply        int64  `json:"max_supply"`
	BlockReward      int64  `json:"block_reward"`
	Difficulty       int    `json:"difficulty"`
	Utility          string `json:"utility"`
	Valid            bool   `json:"valid"`
}

func (bc *Blockchain) Stats() Stats {
	return Stats{
		Height:           len(bc.Blocks) - 1,
		Blocks:           len(bc.Blocks),
		PendingTxs:       len(bc.Pending),
		TotalTransferred: bc.TotalTransferred(),
		TotalSupply:      0,
		MaxSupply:        0,
		BlockReward:      0,
		Difficulty:       Difficulty,
		Utility:          "none",
		Valid:            bc.Validate() == nil,
	}
}
