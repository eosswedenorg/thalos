#!/bin/bash

if [ $# -lt 1 ]; then
    echo "Usage: $0 <install path>"
    exit 1
fi

INSTALL_DIR=$1

echo "Installing application in: $INSTALL_DIR"

if [ -f "config.json" ];then
    CONFIG_FILE=./config.json
else :
    CONFIG_FILE=config.example.json
fi

mkdir -p "$INSTALL_DIR"/{bin,logs}
install -m 750 -t "${INSTALL_DIR}/bin" build/eosio-ship-trace-reader

install -T -m 600 ${CONFIG_FILE} "${INSTALL_DIR}/config.json"
install -m 750 -t "${INSTALL_DIR}" scripts/start.sh scripts/stop.sh

echo "Done"
