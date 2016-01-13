#/bin/bash

cat mysql-create.sql | mysql -u root

sudo -u postgres psql template1 < postgres-create-db.sql
sudo -u benchmarkdbuser psql hello_world < postgres-create.sql
