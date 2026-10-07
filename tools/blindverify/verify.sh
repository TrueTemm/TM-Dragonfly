#!/bin/bash
# verify.sh PROTO GTVER — blind, bidirectional zero-value wire check of our protocol PROTO against the real
# historical gophertunnel GTVER (e.g. `verify.sh 748 v1.42.0`). Run from the repository root with Git Bash.
#
#   OUTGOING  every packet in the vendored server pool -> ConvertFromLatest -> our writer -> decoded by the
#             REAL library's pool: panics, trailing bytes and unknown IDs are reported.
#   INCOMING  every packet in the REAL client pool, zero-encoded -> our client pool + reader -> ConvertToLatest.
#
# BLINDVERIFY_FILL=1 fills every field/slice/Optional with a non-zero value (round #79) — this is what catches
# bugs inside nested structures (AddPlayer's ability layers kicked 1.21.51 live while the zero-value run was
# clean). Expect more noise: filler-only artifacts (enum 0, bitset size, interfaces left nil) appear on EVERY
# protocol and drop out of the baseline diff; Optional sections Dragonfly never populates (biome definitions
# at 844+) can show trailing bytes that real data does not have — confirm those with a realistic packet.
#
# Compare the output against a live-verified protocol's (786 vs v1.45.0 is the reference baseline): anything
# that appears ONLY for the new protocol is either a packet that does not exist there (fine if Dragonfly
# never sends it) or a real wire bug. Zero values do not exercise value-dependent branches — pair this with
# targeted tests (v748/flags_test.go) for those. See HANDOFF.md round #77.
PROTO=$1; GTVER=$2
[ -z "$PROTO" ] || [ -z "$GTVER" ] && { echo "usage: $0 PROTO GTVER"; exit 1; }
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
WORK=${BLINDVERIFY_DIR:-$ROOT/tools/blindverify/out}; mkdir -p "$WORK"
cd "$ROOT/tools/blindverify/reallib" && rm -f go.mod go.sum && go mod init reallib >/dev/null 2>&1 && go mod edit -require=github.com/sandertv/gophertunnel@$GTVER && go mod tidy >/dev/null 2>&1
echo "################ protocol $PROTO vs real $GTVER ################"
echo "=== OUTGOING: our ConvertFromLatest bytes decoded by real $GTVER ==="
cd "$ROOT" && TMDRAGONFLY_BLINDVERIFY=1 DUMPDIR="$WORK" go test ./multiversion/conformance/ -run "TestBlindVerifySend/$PROTO\$" -count=1 2>&1 | grep -E "FAIL|panic|\.go:[0-9]+:" | head -5
cd "$ROOT/tools/blindverify/reallib" && go run . decsend "$WORK/send_$PROTO.txt"
echo "=== INCOMING: real $GTVER client bytes decoded by our pool ==="
go run . encrecv "$WORK/recv_$PROTO.txt"
cd "$ROOT" && TMDRAGONFLY_BLINDVERIFY=1 DUMPDIR="$WORK" go test ./multiversion/conformance/ -run "TestBlindVerifyRecv/$PROTO\$" -count=1 -v 2>&1 | grep -E "RECV|FAIL|panic" | head -60
