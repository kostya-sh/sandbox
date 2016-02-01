#!/bin/bash

URL="http://${APPHOST:-localhost}:8080"
NAME=$1
THREADS=$2
TRACE=$3

echo $URL
for i in $THREADS ; do
    curl "$URL/profile/start?f=${NAME}_$i$TRACE"
    if [ "$TRACE" == "-trace" ] ; then
        curl "$URL/trace/start?f=${NAME}_$i"
    fi
    sleep 1s
    ./wrk -t$i -c$i -d15s "$URL/json" | awk -f extract-wrk.awk
    sleep 1s
    curl "$URL/profile/stop"
    if [ "$TRACE" == "-trace" ] ; then
        curl "$URL/trace/stop"
    fi
    sleep 3s
    echo ""
done
