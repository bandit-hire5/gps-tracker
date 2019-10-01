#!/usr/bin/env bash

start=$(date +%s.%N)

docker build -t gps-tracker .

dur=$(echo "$(date +%s.%N) - $start" | bc)

LC_NUMERIC="en_US.UTF-8" printf "\n-----Execution time: %.6f seconds-----\n\n" $dur
