#!/bin/bash

# Usage: ./qsproxy.sh <ip_listening>

LISTEN=$1

while [[ 1 ]]; do
    ##Get curl command
    RECEIVE_OUTPUT=$(../../qsreceiver receive truncated 1 -l $LISTEN)
    #echo $RECEIVE_OUTPUT
    REMOTE=$(echo $RECEIVE_OUTPUT | cut -d ':' -f 2 | cut -d "," -f 1 | cut -d ' ' -f 2 )
    #echo "REMOTE: $REMOTE"

    CURL_CMD=$(echo $RECEIVE_OUTPUT | cut -d ':' -f 4-)
    #CURL_CMD="${CURL_CMD} -s"
    CURL_CMD="${CURL_CMD} -s -v --stderr -"
    #echo $CURL_CMD

    ##Execute it
    CURL_OUTPUT=$($CURL_CMD)
    #echo $CURL_OUTPUT

    ##Return output to remote
    ../../qssender send "$CURL_OUTPUT" -d 1 -l $LISTEN -r $REMOTE -s 1000 -N
done