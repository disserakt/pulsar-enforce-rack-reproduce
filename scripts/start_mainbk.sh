#!/bin/bash

bin/pulsar initialize-cluster-metadata \
  --cluster pulsar-cluster \
  --zookeeper zoo1:2181 \
  --web-service-url localhost:8080 \
  --configuration-store zoo1:2181 \
  --broker-service-url pulsar://broker1:6650 \ > logs/init-metadata.log 2>&1

export BOOKIE_CONF=/pulsar/conf/bookkeeper.conf
cp /conf/bookkeeper.conf /pulsar/conf/bookkeeper.conf
echo "bookieId=$(hostname):3181" >> /pulsar/conf/bookkeeper.conf
cat /pulsar/conf/bookkeeper.conf

bin/bookkeeper bookie > logs/bookkeeper.log 2>&1
