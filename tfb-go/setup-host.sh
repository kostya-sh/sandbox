#!/bin/bash

set -e

if [ -z "$1" -o -z "$2" ] ; then
  echo "Usage: $0 [db|load|app] ip"
  exit 1
fi

type=$1
host=$2

echo "$type $host"

sshcmd="ssh -i key.pem ubuntu@$host"
files=$type-files

if [ -d $files ] ; then
  echo "Copying files"
  rsync -e "ssh -i key.pem" -rcEzi $files/ ubuntu@$host:.

  if [ -f $files/sudo-setup.sh ] ; then
    echo "Running sudo-setup if necessary"
    $sshcmd "if [ ! -f sudo-setup.done ] ; then sudo /bin/bash sudo-setup.sh && touch sudo-setup.done ; fi"
  fi

  if [ -f $files/setup.sh ] ; then
    echo "Running setup if necessary"
    $sshcmd "if [ ! -f setup.done ] ; then /bin/bash setup.sh && touch setup.done ; fi"
  fi

  if [ -d $files/sysconfig ] ; then
    echo "Copying system config files"
    $sshcmd sudo cp -ruv sysconfig/* /
  fi
fi
