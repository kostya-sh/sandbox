#!/bin/bash

if [ ! -f setup.done ] ; then
    wget -c https://storage.googleapis.com/golang/go1.5.2.linux-amd64.tar.gz
    tar xf go1.5.2.linux-amd64.tar.gz
    mv go go1.5

    mkdir bin

    echo 'export GOPATH=$HOME' >> ~/.profile
    echo 'export GOROOT=$HOME/go1.5' >> ~/.profile
    echo 'export PATH=$GOROOT/bin:$PATH' >> ~/.profile

    sudo apt-get update
    sudo apt-get -y install git

    touch setup.done
fi
