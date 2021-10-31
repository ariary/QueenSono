#!/bin/bash
# Usage: Like curl command (replace 'curl' by './qscurl.sh')"
# And use \ behind ' or "
# replace LISTEN address by your listening ip address (must be reachable w/ icmp by remote)
# replace REMOTE address by qsproxy ip address

LISTEN="TOFILL"
REMOTE="TOFILL"

ARGS="$@"
CURL_CMD="curl ${ARGS}"
echo "$CURL_CMD"
##Send curl command
../../qssender send "${CURL_CMD}" -d 1 -l $LISTEN -r $REMOTE -s 1000 -N

##Wait remote for the command output
RECEIVE_CMD_OUTPUT=$(../../qsreceiver receive truncated 1 -l $LISTEN)

CMD_OUTPUT=$(echo $RECEIVE_CMD_OUTPUT | rev | cut -d ':' -f 1 | rev )

echo $CMD_OUTPUT 