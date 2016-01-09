#!/bin/bash

S=$1

scp -i key.pem load-files/* ubuntu@$S:.
