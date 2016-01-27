#!/bin/bash

URL="http://${APPHOST:-localhost}:8080"
NAME=$1

echo $URL
for i in 1 2 4 8 16 32 64 128 256 ; do
    curl "$URL/profile/start?f=${NAME}_$i"
    if [ "$2" == "-trace" ] ; then
        curl "$URL/trace/start?f=${NAME}_$i"
    fi
    ./wrk -t$i -c$i -d15s "$URL/json" | awk -f extract-wrk.awk
    curl "$URL/profile/stop"
    if [ "$2" == "-trace" ] ; then
        curl "$URL/trace/stop"
    fi
    sleep 3s
    echo ""
done
