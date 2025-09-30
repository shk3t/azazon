#!/usr/bin/env bash

topics=(order_created order_confirmed order_cancelling order_canceled)

for topic in ${topics[@]}; do
    kafka-topics.sh \
        --bootstrap-server localhost:9092 \
        --delete \
        --topic $topic

    kafka-topics.sh \
        --bootstrap-server localhost:9092 \
        --create \
        --topic $topic \
        --partitions 3 \
        --replication-factor 1
done
