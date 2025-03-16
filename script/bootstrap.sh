#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=core
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}