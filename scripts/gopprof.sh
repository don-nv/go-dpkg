#!/bin/bash


pprof_url=$1
if [ -z "$pprof_url" ]; then
  echo missing \$1 argument, pprof_url
  exit 1
fi

pprof_seconds=$2
if [ -z "$pprof_seconds" ]; then
  pprof_seconds=15
fi
if [ "$pprof_seconds" -lt 1 ]; then
    pprof_seconds=15
fi

printf "pprof_url: %s\npprof_seconds: %d (\$2)\n" "$pprof_url" $pprof_seconds

exec go tool pprof \
  -call_tree \
  -http :9110 \
  -seconds $pprof_seconds \
  -source_path=. \
  -trim_path=/go/src/ "$pprof_url"