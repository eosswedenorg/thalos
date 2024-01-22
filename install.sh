#!/bin/bash

if [ $# -lt 1 ]; then
    echo "Usage: $0 <install path>"
    exit 1
fi

INSTALL_DIR=$1

echo -e "\033[34m-\033[0m Installing application in: $INSTALL_DIR"

PROGRAMS=(make go)

missing_prog=0
for prog in ${PROGRAMS[@]}; do
    CMD=$(which $prog)
    if [ -z "$CMD" ]; then
        echo -e "\033[31m!!\033[0m Failed to locate $prog, please install this program"
        missing_prog=1
    fi
done

if [ $missing_prog -ne 0 ]; then
    exit 1
fi

mkdir -p "$INSTALL_DIR"/{bin,logs}
make -e DESTDIR=$INSTALL_DIR PREFIX= CFGDIR= install install-scripts
if [ $? -ne 0 ]; then
    echo -e "\033[31m!!\033[0m Installation failed"
    exit 1
fi

echo -e "\033[32m*\033[0m Done"
