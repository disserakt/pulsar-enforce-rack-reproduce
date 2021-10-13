#!/bin/bash

set -uxe

until bin/pulsar-admin --admin-url http://broker1:8080 bookies set-bookie-rack -b bookie1:3181 -r rack1; do echo 'retry'; done
bin/pulsar-admin --admin-url http://broker1:8080 bookies set-bookie-rack -b bookie2:3181 -r rack1
bin/pulsar-admin --admin-url http://broker1:8080 bookies set-bookie-rack -b bookie3:3181 -r rack1
bin/pulsar-admin --admin-url http://broker1:8080 bookies set-bookie-rack -b bookie4:3181 -r rack2
bin/pulsar-admin --admin-url http://broker1:8080 bookies set-bookie-rack -b bookie5:3181 -r rack2
bin/pulsar-admin --admin-url http://broker1:8080 bookies set-bookie-rack -b bookie6:3181 -r rack2
bin/pulsar-admin --admin-url http://broker1:8080 bookies set-bookie-rack -b bookie7:3181 -r rack3
bin/pulsar-admin --admin-url http://broker1:8080 bookies set-bookie-rack -b bookie8:3181 -r rack3
bin/pulsar-admin --admin-url http://broker1:8080 bookies set-bookie-rack -b bookie9:3181 -r rack3

