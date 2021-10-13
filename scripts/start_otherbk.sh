#!/bin/bash
export BOOKIE_CONF=/pulsar/conf/bookkeeper.conf
cp /conf/bookkeeper.conf /pulsar/conf/bookkeeper.conf
echo "bookieId=$(hostname):3181" >> /pulsar/conf/bookkeeper.conf
cat /pulsar/conf/bookkeeper.conf
bin/bookkeeper bookie > logs/bookkeeper.log 2>&1
