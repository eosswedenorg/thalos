#!/bin/bash
BIN=bin/eosio-ship-trace-reader

DIR=$(dirname $(realpath $0))
cd "$DIR"

date
./stop.sh
timestamp=`date +%s`
$BIN -p ./eosio-ship-trace-reader.pid 2> logs/$timestamp.log &
rm -f out.log
ln -s logs/$timestamp.log out.log
