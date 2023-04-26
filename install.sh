#!/bin/bash

if [ $# -lt 1 ]; then
    echo "Usage: $0 <install path>"
    exit 1
fi

INSTALL_DIR=$1

echo "Installing application in: $INSTALL_DIR"

mkdir -p "$INSTALL_DIR"/{bin,logs}
install -m 750 -t "${INSTALL_DIR}/bin" build/thalos-server

if [ ! -f "${INSTALL_DIR}/config.yml" ]; then
    install -T -m 600 config.example.yml "${INSTALL_DIR}/config.yml"
fi

install -m 750 -t "${INSTALL_DIR}" scripts/start.sh scripts/stop.sh

echo "Done"
