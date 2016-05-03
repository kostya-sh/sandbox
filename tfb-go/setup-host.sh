#!/bin/bash

set -e

if [ -z "$1" -o -z "$2" ] ; then
  echo "Usage: $0 [db|load|app] ip"
  exit 1
fi

type=$1
host=$2

echo "$type $host"

sshcmd="ssh -oStrictHostKeyChecking=no -i key.pem ubuntu@$host"
files=$type-files

if [ -d $files ] ; then
  echo "Copying files"
  rsync -e "ssh -oStrictHostKeyChecking=no -i key.pem" -rcEzit common-files/ $files/ ubuntu@$host:.

  echo "Running setup if necessary"
  $sshcmd "/bin/bash run-setup.sh"
fi
