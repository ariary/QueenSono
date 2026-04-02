# Echo Reply Mode — Design Spec

**Date:** 2026-04-02

## Problem

QueenSono currently embeds data exclusively in ICMP echo **requests** (type 8). The `golang.org/x/net/icmp` package does not prevent the kernel from auto-replying to echo requests — the kernel mirrors the request payload unchanged. This means echo **replies** (type 0) cannot carry arbitrary server-injected data with the existing approach.

This blocks a "reverse" communication channel where the server sends data to the client via echo replies.

## Goal

Add an echo reply mode where:

- The **client** (`qssender`) sends a trigger echo request and collects data from incoming echo replies.
- The **server** (`qsreceiver`) receives the trigger and sends data chunks embedded in custom echo replies.

## Protocol

```
Client (qssender)              Server (qsreceiver)
      |                               |
      |-- echo request "QS_READY" --> |  trigger
      |<-- echo reply "N" -----------|  announce chunk count
      |<-- echo reply "0,chunk0" ----|
      |<-- echo reply "1,chunk1" ----|
      |           ...                 |
      |<-- echo reply "N-1,chunkN" --|
```

- The trigger payload `"QS_READY"` identifies this as a new-mode exchange.
- Data chunks use the existing `QueenSonoMarshall` format (`"index,data"`).
- The server uses `ipv4.ICMPTypeEchoReply` in crafted packets sent via `WriteTo` on a raw `ip4:icmp` connection.

## Kernel Auto-Reply Handling (Option A — Marker Filtering)

On Linux, the kernel automatically sends an echo reply (type 0) with the same payload as the received echo request before the raw socket app can act. This means the client may receive two echo replies for each request:
1. The kernel's auto-reply (payload = `"QS_READY"`, not a valid QueenSono message).
2. The server app's custom reply (payload = `"N"` or `"index,data"`).

The client filters incoming echo replies: any packet whose data does not parse correctly via `QueenSonoUnmarshall` (or is not the size announcement) is silently discarded.

## New Code

### `pkg/icmp/reply.go` (new file)

Two functions:

**`ServeWithEchoReply(listenAddr, remoteAddr string, data string, chunkSize, delay int)`**
- Listens for an echo request with payload `"QS_READY"` on `listenAddr`.
- Extracts the peer address from the received packet.
- Chunks `data` using `Chunks` + `QueenSonoMarshall`.
- Announces the chunk count by sending an echo reply with `strconv.Itoa(len(chunks))`.
- Sends each chunk as an echo reply with the appropriate delay.

**`TriggerAndReceiveReplies(listenAddr, remoteAddr string, chunkSize, delay int) string`**
- Opens a raw `ip4:icmp` listener on `listenAddr`.
- Sends a single echo request with payload `"QS_READY"` to `remoteAddr`.
- Reads the first parseable echo reply to learn the chunk count N.
- Reads N more echo replies, filtering out non-QueenSono packets.
- Reassembles and returns the full data string.

### `cmd/server/main.go`

New subcommand `reply-send [data]` under `qsreceiver`:
- Flags: `--listen / -l`, `--size / -s`, `--delay / -d`
- Calls `icmp.ServeWithEchoReply`.

### `cmd/client/main.go`

New subcommand `receive` under `qssender`:
- Flags: `--remote / -r` (required), `--listen / -l`, `--size / -s`, `--delay / -d`
- Calls `icmp.TriggerAndReceiveReplies` and prints the result.

## Out of Scope

- Encryption support in the new mode (can be added later).
- IPv6 support.
- File output flag for the client receive command (can be added later).
