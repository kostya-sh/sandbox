#!/bin/bash

S=$1
SSH="ssh -i key.pem ubuntu@$S"

cat setup-app.sh | $SSH

scp -i key.pem -r src ubuntu@$S:.
