#!/bin/bash

set -e

if [ -f sudo-setup.sh ] ; then
  if [ ! -f sudo-setup.done ] ; then
    echo "Running sudo-setup"
    sudo /bin/bash sudo-setup.sh && touch sudo-setup.done
  fi
fi

if [ -f setup.sh ] ; then
  if [ ! -f setup.done ] ; then
    echo "Running setup"
    /bin/bash setup.sh && touch setup.done
  fi
fi

if [ -d sysconfig ] ; then
  echo "Copying system config files"
  #sudo rsync -rcEzi sysconfig/ /
  sudo cp -ruv sysconfig/* /
fi
