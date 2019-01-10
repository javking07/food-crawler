#!/bin/sh
# wait-for-cassandra.sh

set -e


cmd="$1"

sleep 90 && echo "sleeping for Cassandra"

echo "Cassandra is up - executing command"
exec $cmd