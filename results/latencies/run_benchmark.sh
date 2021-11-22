#!/bin/sh

NB_WORKERS=10
NB_OVENS=10

for WORKER in $(seq 1 1 $NB_WORKERS)
do
for OVEN in $(seq 1 1 $NB_OVENS)
do
    echo "calling with $WORKER workers and $OVEN ovens"
    go run ../../cmd/pizzeria/main.go 1 2 5 1 $WORKER $OVEN 200
done
done