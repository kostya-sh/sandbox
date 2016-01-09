#!/bin/bash

sudo -i

apt-get update
apt-get -y install mysql-server
apt-get -y install sysstat

sed -i 's|\[mysqld\]|\[mysqld\]\
lower_case_table_names = 1\
character-set-server=utf8\
collation-server=utf8_general_ci|g' /etc/mysql/my.cnf

sed -i 's|bind-address.*=.*127.0.0.1|bind-address = 0.0.0.0|g' /etc/mysql/my.cnf

service mysql restart
