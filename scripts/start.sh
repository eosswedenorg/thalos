#!/bin/bash
BIN=bin/thalos-server

DIR=$(dirname $(realpath $0))
cd "$DIR"

date
./stop.sh
timestamp=`date +%s`
$BIN -p ./thalos.pid 2> logs/$timestamp.log &
rm -f out.log
ln -s logs/$timestamp.log out.log
