#!/bin/sh

case "$1" in
'cpu')
  shift
  go test -bench="$1" "$2" -test.cpuprofile "$3".cpu.pprof -test.outputdir ./pprof -test.count "$4"
  ;;

'mem')
  shift
  go test -bench="$1" "$2" -cpu 4 -test.benchmem -test.memprofile "$3".mem.pprof -test.outputdir ./pprof -test.count "$4"
  ;;

esac
