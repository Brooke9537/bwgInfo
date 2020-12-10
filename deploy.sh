#!/bin/bash
set -x
GOOS=linux GOARCH=amd64 go build main.go
mv main bwgInfo
ssh root@95.169.11.215 'supervisorctl stop bwgInfo'
scp bwgInfo root@95.169.11.215:/opt/
ssh root@95.169.11.215 'supervisorctl start bwgInfo'
