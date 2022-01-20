#!/bin/bash

if [ $# -lt 1 ]; then
    echo "Usage: $0 <install path>"
    exit 1
fi

INSTALL_DIR=$1

echo "Installing application in: $INSTALL_DIR"

mkdir -p "$INSTALL_DIR"/{bin,logs}
install -m 750 -t "${INSTALL_DIR}/bin" build/eosio-ship-trace-reader

if [ ! -f "${INSTALL_DIR}/config.json" ]; then
    install -T -m 600 config.example.json "${INSTALL_DIR}/config.json"
fi

install -m 750 -t "${INSTALL_DIR}" scripts/start.sh scripts/stop.sh

echo "Done"
