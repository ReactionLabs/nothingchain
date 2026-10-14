# NothingChain — Spec

## (1) End goal
A fully working blockchain in Go that solves absolutely nothing, on purpose, perfectly.

It mines real blocks, validates a real chain, runs a real node API, and moves
exactly zero value forever. The joke is the architecture is legit. The utility
is zero.

Ticker: **NADA**
Consensus: **Proof of Nothing (PoN)** — SHA-256 proof-of-work, wasted correctly.
Total supply: 0. Max supply: 0. Circulating supply: 0. Forever.

## (2) What done looks like
- Genesis block: "In nothing we trust."
- Blocks mined via PoW, hash-linked, difficulty enforced.
- Transactions validate — but amount MUST be 0. Try to send 1 NADA and the
  chain rejects you: there is nothing to send.
- Block reward: 0 NADA, paid instantly to nobody in particular.
- CLI (first-class): `mine`, `send`, `balance`, `stats`, `chain`, `validate`
- HTTP API (first-class): `GET /health`, `GET /stats`, `GET /chain`,
  `GET /balance/{addr}`, `POST /tx`, `POST /mine`
- Balance of every address, no matter what happens: **0**
- Chain validation proves the nothing is authentic and untampered.
- Builds with only the Go standard library. No dependencies. Nothing to install.

## (3) What done does NOT look like
- No token sale, no utility, no roadmap, no whitepaper promises.
- No smart contracts that do something. (NoVM is planned: it will do nothing.)
- No price. It cannot pump. It cannot dump. FDV is $0 and that's a feature.
- Not a fork of anything. Written from scratch, like Moneta, but useless.
- It does not solve payments, identity, supply chain, AI, gaming, or storage.
  It solves nothing, and it solves it completely.

## Economics (audited by nobody)
| Metric | Value |
|---|---|
| Total supply | 0 NADA |
| Inflation | 0% |
| Block reward | 0 NADA |
| Tx fee | 0 NADA |
| Value transferred, all time | 0 |
| Downside | None. You can't lose nothing. |
