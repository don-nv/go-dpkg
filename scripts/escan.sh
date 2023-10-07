#!/bin/sh

go build -gcflags="-m=$1" "$2"
