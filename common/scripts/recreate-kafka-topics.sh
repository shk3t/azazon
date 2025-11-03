#!/usr/bin/env bash

source $(dirname $0)/envs.sh

BOOTSTRAP_SERVER="localhost:9092"

if [[ " $@ " =~ ( --minikube ) ]]; then
    kafka_pod="${kafka_cluster_name}-${kafka_nodepool_name}-0"

    if ! kubectl exec -n $namespace $kafka_pod -- bin/kafka-broker-api-versions.sh --bootstrap-server $BOOTSTRAP_SERVER > /dev/null 2>&1; then
        echo "ERROR: Could not connect to Kafka on ${BOOTSTRAP_SERVER}."
        exit 1
    fi

    for topic in ${topics[@]}; do
        kubectl exec -n $namespace $kafka_pod -- \
            bin/kafka-topics.sh \
                --bootstrap-server $BOOTSTRAP_SERVER \
                --delete \
                --topic $topic \
                --if-exists

        kubectl exec -n $namespace $kafka_pod -- \
            bin/kafka-topics.sh \
                --bootstrap-server $BOOTSTRAP_SERVER \
                --create \
                --topic $topic \
                --partitions 3 \
                --replication-factor 1
    done
else
    for topic in ${topics[@]}; do
        if ! kafka-broker-api-versions.sh --bootstrap-server $BOOTSTRAP_SERVER > /dev/null 2>&1; then
            echo "ERROR: Could not connect to Kafka on ${BOOTSTRAP_SERVER}."
            exit 1
        fi

        kafka-topics.sh \
            --bootstrap-server $BOOTSTRAP_SERVER \
            --delete \
            --topic $topic \
            --if-exists

        kafka-topics.sh \
            --bootstrap-server $BOOTSTRAP_SERVER \
            --create \
            --topic $topic \
            --partitions 3 \
            --replication-factor 1
    done
fi
