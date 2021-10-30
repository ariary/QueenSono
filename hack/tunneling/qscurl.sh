#!/bin/bash
# Usage: Like curl command (replace 'curl' by './qscurl.sh')"
# And use \ behind ' or "
# replace LISTEN address by your listening ip address (must be reachable w/ icmp by remote)
# replace REMOTE address by qsproxy ip address

LISTEN="10.10.10.10"
REMOTE="10.10.10.11"

ARGS="$@"
CURL_CMD="curl ${ARGS}"

##Send curl command
../../qssender send "$CURL_CMD" -d 1 -l $LISTEN -r $REMOTE -s 100 -N

##Wait remote for the command output
RECEIVE_CMD_OUTPUT=$(../../qsreceiver receive truncated 1 -l $LISTEN)

CMD_OUTPUT=$(echo $RECEIVE_CMD_OUTPUT | rev | cut -d ':' -f 1 | rev )

echo $CMD_OUTPUT 