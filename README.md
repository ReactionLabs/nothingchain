# NothingChain

A blockchain in Go that solves absolutely nothing.

Real proof-of-work. Real hash-linked blocks. Real validation.
Zero value, zero utility, zero supply — by design, forever.

## Quickstart

```bash
go build -o nothingchain .
./nothingchain mine --miner dip
./nothingchain send --from dip --to cole --amount 0 --memo "for nothing"
./nothingchain mine --miner dip
./nothingchain balance cole     # 0 NADA. Always.
./nothingchain balance anyone   # Also 0. It's fair.
./nothingchain validate
./nothingchain stats
./nothingchain serve --port 8080
```

Try to cheat it:

```bash
./nothingchain send --from dip --to cole --amount 1
# Rejected: NothingChain only moves nothing (0 NADA)
```

## API

```bash
curl localhost:8080/stats
curl localhost:8080/balance/anybody
curl -X POST localhost:8080/tx  -d '{"from":"a","to":"b","amount":0,"memo":"hi"}'
curl -X POST localhost:8080/mine -d '{"miner":"dip"}'
```

## Why

Every other chain promises to solve something. This one keeps its promise.

- Cannot pump. Cannot dump. FDV $0 at genesis, $0 forever.
- Inflation-proof: 0% forever.
- The downside is nothing, because you hold nothing.

See SPEC.md for the full end goal / done / not-done.
