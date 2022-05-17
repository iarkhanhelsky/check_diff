#!/bin/bash
awk 'BEGIN {OFS=":"} {print FILENAME, NR, $0 }' $1