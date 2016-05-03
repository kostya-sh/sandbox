#!/bin/bash

set -e

if [ -z "$1" -o -z "$2" -o -z "$3" ] ; then
  echo "Usage: $0 dbIP loadIP appIP"
  exit 1
fi

./setup-host.sh db $1 &
./setup-host.sh load $2 &
./setup-host.sh app $3 &

wait

echo ""
echo "To login:"
echo ""

echo "db:"
echo "ssh -oStrictHostKeyChecking=no -i $PWD/key.pem ubuntu@$1"
echo "load:"
echo "ssh -oStrictHostKeyChecking=no -i $PWD/key.pem ubuntu@$2"
echo "app:"
echo "ssh -oStrictHostKeyChecking=no -i $PWD/key.pem ubuntu@$3"
