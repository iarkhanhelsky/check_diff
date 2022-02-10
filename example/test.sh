#!/bin/bash

check_diff="./check_diff"
if [[ ! -z $1 ]]; then
  check_diff=$1
fi

for f in diffs/*.diff; do
  diff -U3 <(jq < "$f.gitlab.json")      <($check_diff -i $f -f gitlab      | jq)
  diff -U3 <(jq < "$f.phabricator.json") <($check_diff -i $f -f phabricator | jq)
  diff -U3 "$f.stdout.txt"       <($check_diff -i $f -f stdout --no-color)
done