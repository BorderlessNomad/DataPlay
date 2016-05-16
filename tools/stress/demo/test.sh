#!/bin/bash

LOCUST=/usr/local/bin/locust

IP=`ifconfig eth0 | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p'`

LOCUST_ARGS="--no-web --host=http://$IP --print-stats --clients=1000 --hatch-rate=100 -f "

LOCUST_CONFIGS=(
	'dataplay-home.py'
	'dataplay-search.py'
	'dataplay-news.py'
	'dataplay-related.py'
	'dataplay-correlated.py'
	'dataplay-correlated_generate.py'
	'dataplay-activities.py'
)

CONFIG_DURATION=600

CONFIG_INTERVAL=300


if [ ! -e $LOCUST -a -x $LOCUST ] 
then
	echo "Locust does not exist or is not executable"
	exit 1
fi

while true
do
	for CONFIG in ${LOCUST_CONFIGS[@]}
	do
		echo "Executing $CONFIG for $CONFIG_DURATION seconds"
		exec timeout $CONFIG_DURATION $LOCUST $LOCUST_ARGS ${CONFIG}
		echo "Sleeping for $CONFIG_INTERVAL seconds"
		exec sleep $CONFIG_INTERVAL
	done
done
