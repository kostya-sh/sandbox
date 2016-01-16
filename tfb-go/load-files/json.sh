#!/bin/bash

URL="http://${APPHOST:-localhost}:8080/json"

echo $URL
for i in 1 2 4 8 16 32 64 128 256 ; do
    ./wrk -t$i -c$i -d15s "$URL" | awk -f extract-wrk.awk
done
