#!/bin/bash

PIDFILE="$(pwd)/thalos.pid"

if [ -f "$PIDFILE" ]; then
    pid=$(cat "$PIDFILE")
    echo $pid
    kill -s INT $pid
    rm -r "$PIDFILE"
    echo -ne "Stopping process"
    while true; do
        [ ! -d "/proc/$pid/fd" ] && break
        echo -ne "."
        sleep 1
    done
    echo -ne "\rProcesss stopped. \n"
fi
