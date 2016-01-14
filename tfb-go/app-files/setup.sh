#!/bin/bash

wget -c https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz
tar xf go1.5.3.linux-amd64.tar.gz
mv go go1.5

wget -c https://storage.googleapis.com/golang/go1.6beta2.linux-amd64.tar.gz
tar xf go1.6beta2.linux-amd64.tar.gz
mv go go1.6

mkdir bin
