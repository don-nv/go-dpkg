#!/bin/bash

src=$1
filter=$2
with_leaks=$3

echo source: "$src"

if [  -n  "$filter" ]; then
  echo filter: "$filter"
fi

# check compilation errors.
command go build -o /dev/null "$src" || exit 1

# do escape analysis.
result_m2=$(go build -o /dev/null -gcflags "-l -m -m" "$src" \
  |& grep "$filter" \
  | sort --field-separator ":" --stable -k1,1 -k2,2n -k3,3n --ignore-case
)


result=""
ln_wanted=0

while read -r ln
do
  if [[ "$ln" == *"escapes"* || "$ln" == *"moved"* ]]; then
    ln_wanted=1
  fi
  if [[ -z "$with_leaks" && "$ln" == *"leaks"* || "$ln" == *"not escape"* ]]; then
    ln_wanted=0
  fi

  if [ $ln_wanted -eq 1 ]; then
      result=$(printf "%s\n%s" "$result" "$ln")
  fi

done <<< "$result_m2"

echo "$result"