#!/bin/bash

wget -c https://storage.googleapis.com/golang/go1.5.4.linux-amd64.tar.gz
tar xf go1.5.4.linux-amd64.tar.gz
mv go go1.5

wget -c https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
tar xf go1.6.2.linux-amd64.tar.gz
mv go go1.6

git clone https://github.com/golang/go.git go-tip
cd go-tip/src
GOROOT_BOOTSTRAP=~/go1.6 ./make.bash
cd $HOME

mkdir bin
