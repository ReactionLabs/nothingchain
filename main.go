package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const ChainFile = "nothingchain.json"

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	switch os.Args[1] {

	case "mine":
		fs := flag.NewFlagSet("mine", flag.ExitOnError)
		miner := fs.String("miner", "dip", "who earns the 0 NADA reward")
		_ = fs.Parse(os.Args[2:])
		bc := LoadOrCreate(ChainFile)
		b := bc.MinePending(*miner)
		_ = bc.Save(ChainFile)
		fmt.Printf("Mined block #%d\n  hash: %s\n  txs:  %d\n  reward: 0 NADA (paid to %s)\n", b.Index, b.Hash, len(b.Txs), *miner)

	case "send":
		fs := flag.NewFlagSet("send", flag.ExitOnError)
		from := fs.String("from", "", "sender")
		to := fs.String("to", "", "recipient")
		amount := fs.Int64("amount", 0, "must be 0")
		memo := fs.String("memo", "", "profound nothing, on the record")
		_ = fs.Parse(os.Args[2:])
		bc := LoadOrCreate(ChainFile)
		tx := Transaction{From: *from, To: *to, Amount: *amount, Memo: *memo}
		if err := bc.AddTx(tx); err != nil {
			fmt.Println("Rejected:", err)
			os.Exit(1)
		}
		_ = bc.Save(ChainFile)
		fmt.Printf("Queued: %s -> %s : 0 NADA (%q)\nRun `mine` to seal this nothing forever.\n", *from, *to, *memo)

	case "balance":
		if len(os.Args) < 3 {
			fmt.Println("usage: nothingchain balance <address>")
			return
		}
		bc := LoadOrCreate(ChainFile)
		fmt.Printf("%s: %d NADA\n", os.Args[2], bc.Balance(os.Args[2]))

	case "stats":
		bc := LoadOrCreate(ChainFile)
		s := bc.Stats()
		out, _ := json.MarshalIndent(s, "", "  ")
		fmt.Println(string(out))

	case "chain":
		bc := LoadOrCreate(ChainFile)
		out, _ := json.MarshalIndent(bc.Blocks, "", "  ")
		fmt.Println(string(out))

	case "validate":
		bc := LoadOrCreate(ChainFile)
		if err := bc.Validate(); err != nil {
			fmt.Println("INVALID:", err)
			os.Exit(1)
		}
		fmt.Printf("Valid. %d blocks of authentic, untampered nothing.\n", len(bc.Blocks))

	case "serve":
		fs := flag.NewFlagSet("serve", flag.ExitOnError)
		port := fs.String("port", "8080", "port")
		_ = fs.Parse(os.Args[2:])
		serve(*port)

	default:
		usage()
	}
}

func usage() {
	fmt.Println(`NothingChain — the blockchain that solves absolutely nothing.

Usage:
  nothingchain mine   [--miner dip]              mine pending txs, earn 0
  nothingchain send   --from A --to B [--amount 0] [--memo "..."]
  nothingchain balance <address>                 always 0
  nothingchain stats                             supply: 0, utility: none
  nothingchain chain                             print the ledger of nothing
  nothingchain validate                          prove the nothing is real
  nothingchain serve  [--port 8080]              run the node HTTP API`)
}

func serve(port string) {
	bc := LoadOrCreate(ChainFile)
	mux := http.NewServeMux()

	writeJSON := func(w http.ResponseWriter, v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"status": "ok", "utility": "none"})
	})
	mux.HandleFunc("GET /stats", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, bc.Stats())
	})
	mux.HandleFunc("GET /chain", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, bc.Blocks)
	})
	mux.HandleFunc("GET /balance/", func(w http.ResponseWriter, r *http.Request) {
		addr := strings.TrimPrefix(r.URL.Path, "/balance/")
		writeJSON(w, map[string]any{"address": addr, "balance": bc.Balance(addr), "unit": "NADA"})
	})
	mux.HandleFunc("POST /tx", func(w http.ResponseWriter, r *http.Request) {
		var tx Transaction
		if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		if err := bc.AddTx(tx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = bc.Save(ChainFile)
		writeJSON(w, map[string]any{"queued": tx})
	})
	mux.HandleFunc("POST /mine", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Miner string `json:"miner"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Miner == "" {
			body.Miner = "dip"
		}
		b := bc.MinePending(body.Miner)
		_ = bc.Save(ChainFile)
		writeJSON(w, b)
	})

	log.Printf("NothingChain node listening on :%s — solving nothing, 24/7", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
