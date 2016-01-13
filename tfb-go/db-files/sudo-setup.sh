#!/bin/bash

export DEBIAN_FRONTEND=noninteractive

apt-get update
apt-get -y install sysstat
apt-get -y install mysql-server
apt-get -y install postgresql

sudo useradd benchmarkdbuser -p benchmarkdbpass
