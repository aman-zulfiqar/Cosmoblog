## CosmoBlog IBC Demo

**CosmoBlog** is a minimal two-chain IBC demo built with Ignite CLI and Cosmos SDK.  
It runs a single blockchain binary (`planetd`) as two networks (**earth** and **mars**) and uses an IBC-enabled `blog` module to send posts between them.

---

## Project Layout

- **`planet/`** – Cosmos SDK chain source (Ignite scaffolded).
  - `app/` – app wiring, IBC routing, module registration.
  - `x/blog/` – blog module with IBC `IbcPostPacketData`/acks, sent/timeout tracking.
  - `readme.md` – Ignite’s default chain README.
- **`earth.yml`** – Ignite config to run `planetd` as chain **earth**.
- **`mars.yml`** – Ignite config to run `planetd` as chain **mars**.

You always run node/IBC commands from inside `planet/`, using `earth.yml` / `mars.yml` as configs.

---

## Prerequisites

- **Go**: 1.21+  
  - Check: `go version`
- **Ignite CLI**: latest
  - Check: `ignite version`
  - Install (macOS, typical): `brew install ignite/cli/ignite`
- **Hermes relayer app via Ignite** (for IBC):

```bash
ignite app install -g github.com/ignite/apps/hermes
```

If you get a "target directory ... is not empty" error, clear the stale clone:

```bash
rm -rf "$HOME/.ignite/apps/github.com/ignite/apps"
ignite app install -g github.com/ignite/apps/hermes
```

---

## Build & Test

```bash
cd "/Users/muhammadamman/CodeLab/completed projects/CosmoBlog/planet"

go test ./...
```

This runs unit tests for `planet/app` and `x/blog` (including IBC packet logic).

---

## Running the Two Chains

Run these in **two separate terminals**.

- **Terminal 1 – earth**

```bash
cd "/Users/muhammadamman/CodeLab/completed projects/CosmoBlog/planet"

ignite chain serve -c ../earth.yml
```

This starts:

- **Chain ID**: `earth`
- **RPC**: `http://0.0.0.0:26657`
- **API**: `http://0.0.0.0:1317`
- **gRPC**: `0.0.0.0:9090`
- **Faucet**: `http://0.0.0.0:4500`

- **Terminal 2 – mars**

```bash
cd "/Users/muhammadamman/CodeLab/completed projects/CosmoBlog/planet"

ignite chain serve -c ../mars.yml
```

This starts:

- **Chain ID**: `mars`
- **RPC**: `http://0.0.0.0:26659`
- **API**: `http://0.0.0.0:1318`
- **gRPC**: `0.0.0.0:9092`
- **Faucet**: `http://0.0.0.0:4501`

**Ports in use**:  
If you see `bind: address already in use` (e.g. on `:4500`), either:

```bash
lsof -nP -iTCP:4500 -sTCP:LISTEN
kill <PID>
```

or just go to the terminal where `ignite chain serve` is running and press `q` to stop it.

---

## IBC Blog Module Overview

- **Packet type**: `IbcPostPacketData` (see `proto/planet/blog/v1/packet.proto` and `x/blog/types/packet.pb.go`)
  - Fields: `title`, `content`, `creator`.
- **Sender-side tx**: `MsgSendIbcPost` implemented in `x/blog/keeper/msg_server_ibc_post.go`.
  - Builds `IbcPostPacketData` and calls `TransmitIbcPostPacket(...)`.
- **IBC hooks** in `x/blog/keeper/ibc_post.go`:
  - **`TransmitIbcPostPacket`** – wraps `ChannelKeeper.SendPacket`.
  - **`OnRecvIbcPostPacket`** – creates a new `Post` on the receiving chain and sets `IbcPostPacketAck.post_id`.
  - **`OnAcknowledgementIbcPostPacket`** – records successful sends into `SentPost`.
  - **`OnTimeoutIbcPostPacket`** – records timeouts into `TimeoutPost`.

The app-level IBC wiring (ports, channels, callbacks) is handled by `x/blog/module_ibc.go` and `app/ibc.go`.

---

## Setting Up the Hermes Relayer

With both chains running, configure Hermes via Ignite in a **third** terminal:

```bash
cd "/Users/muhammadamman/CodeLab/completed projects/CosmoBlog/planet"

ignite relayer hermes configure \
  earth "http://localhost:26657" "http://localhost:9090" \
  mars  "http://localhost:26659" "http://localhost:9092" \
  --chain-a-faucet "http://0.0.0.0:4500" \
  --chain-b-faucet "http://0.0.0.0:4501" \
  --chain-a-port-id "blog" \
  --chain-b-port-id "blog" \
  --channel-version "blog-1"
```

This:

- Funds relayer accounts on **earth** and **mars** using the faucets.
- Creates:
  - IBC clients on each side,
  - a connection,
  - a channel `channel-0` bound to port `blog` on both chains.

Start the relayer:

```bash
ignite relayer hermes start earth mars
```

This long-running process watches for IBC packet/ack/timeout events and ferries proofs between chains.

**Resetting relayer state** (if configs are broken or you want a clean slate):

```bash
rm -rf "$HOME/.ignite/relayer"
```

---

## Sending a Blog Post over IBC

The `planetd` binary is installed to `$(go env GOPATH)/bin/planetd` by Ignite.  
You can use it directly, or via `ignite tx` wrappers. Below is the direct CLI.

- **On earth: send an IBC post to mars**

```bash
planetd tx blog send-ibc-post blog channel-0 \
  "Hello from earth" \
  "This post was sent over IBC to mars" \
  --from alice \
  --chain-id earth \
  --node http://localhost:26657 \
  --home "$HOME/.earth" \
  --yes
```

The tx builds an `IbcPostPacketData` and calls `TransmitIbcPostPacket`. The packet is committed on `earth` and then the relayer delivers it to `mars`.

- **On mars: inspect received state**

List the available blog queries:

```bash
planetd q blog --help
```

Typical examples (exact subcommands may vary depending on your scaffold):

```bash
planetd q blog list-post \
  --node http://localhost:26659 \
  --chain-id mars \
  --home "$HOME/.mars"

planetd q blog list-sent-post \
  --node http://localhost:26657 \
  --chain-id earth \
  --home "$HOME/.earth"
```

You should see the post created on `mars` with the same `title`/`content`, and `SentPost`/`TimeoutPost` tracking on `earth` depending on ack/timeout outcomes.

---

## Relayer / IBC Notes

- **Relayer role**:
  - Off-chain process that:
    - watches for `SendPacket` events on both chains,
    - builds `MsgRecvPacket`, `MsgAcknowledgement`, and `MsgTimeout` txs with proofs,
    - submits them to the counterparty chain.
  - Required for **liveness** (no relayer → packets just sit in commitments).
- **Security model**:
  - The relayer is **not trusted** for correctness; chains verify all proofs using light clients.
  - A malicious relayer can delay or skip packets, but cannot forge them.

---

## Where to Look in the Code

- **IBC packet types & encoding**
  - `proto/planet/blog/v1/packet.proto`
  - `x/blog/types/packet.pb.go`
  - `x/blog/types/packet_ibc_post.go`
- **IBC blog keeper logic**
  - `x/blog/keeper/msg_server_ibc_post.go`
  - `x/blog/keeper/ibc_post.go`
- **Module wiring & IBC hooks**
  - `x/blog/module_ibc.go`
  - `app/ibc.go`

These files together show the full IBC lifecycle: packet construction, send, recv, ack, and timeout handling for cross-chain blog posts.


