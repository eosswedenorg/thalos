#!/bin/bash

if [ $# -lt 1 ]; then
    echo "Usage: $0 <install path>"
    exit 1
fi

INSTALL_DIR=$1

echo "Installing application in: $INSTALL_DIR"

mkdir -p "$INSTALL_DIR"/{bin,logs}

make -e DESTDIR=$INSTALL_DIR PREFIX= CFGDIR= install install-scripts

echo "Done"
