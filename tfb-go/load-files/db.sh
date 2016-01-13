#!/bin/bash

URL="http://${LOADHOST:-localhost}:8080/db"

echo $URL
for i in 1 2 4 8 16 32 64 128 ; do
    ./wrk -t$i -c$i -d15s "$URL" | awk -f extract-wrk.awk
done