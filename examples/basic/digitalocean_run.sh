#!/bin/bash

set -e

apt-get install -y git

# create a temp directory
dir=`mktemp -d` && cd $dir

# download and unpack go
wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
tar -zxvf go1.4.2.linux-amd64.tar.gz

mkdir gocode

# setup the environment
export GOROOT=$PWD/go
export PATH=$PATH:$GOROOT/bin
export GOPATH=$PWD/gocode

go get github.com/mattn/goreman
go get github.com/dronemill/eventsocket

cd `go list -f '{{.Dir}}' github.com/dronemill/eventsocket`
cd examples/basic

$GOPATH/bin/goreman start