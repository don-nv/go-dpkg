#!/bin/bash


cmd=$1
if [ -z "$cmd" ]; then
  echo missing \$1 \(cmd\) argument
  exit 1
fi

rps=$2
if [ -z "$rps" ]; then
  echo missing \$2 \(rps\) argument
  exit 1
fi
if [ "$rps" -lt 1 ]; then
  echo invalid \$2 \(rps\) argument \< 1
  exit 1
fi

printf "command: %s\nrps: %d\n" "$cmd" "$rps"
env sleep 1s

delay=$(bc <<< "1/$rps")

for ((i = 0; ; i++)); do
    echo "run: $i"

    $cmd &

    env sleep "$delay"
done