#!/bin/bash

S=$1
SSH="ssh -i key.pem ubuntu@$S"

cat setup-db.sh | $SSH

cat create.sql | $SSH mysql --user root
