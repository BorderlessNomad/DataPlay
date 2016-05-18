#!/bin/bash

LOADBALANCER='109.231.121.115'

SERVERS=(
	"FrontendNode:109.231.121.95"
	"LoadBalancer:$LOADBALANCER"
	"Redis:109.231.121.157"
	"MasterNode:109.231.121.150"
	"Cassandra:109.231.121.151"
	"PgPool:109.231.121.5"
	"PostgreSql:109.231.121.123"
)

LOCUST_CONFIGS=(
	'dataplay-correlated.py'
	'dataplay-correlated_generate.py'
	'dataplay-related.py'
	'dataplay-search.py'
)

LOG_DIR=./logs

if [ ! -e $LOG_DIR ]
then
	mkdir $LOG_DIR
fi

echo clearing logs from previous runs
rm -f $LOG_DIR/* || true

CONFIG_DURATION=600
CONFIG_IOSTAT_DURATION=630
CONFIG_INTERVAL=120
CONFIG_CLIENTS=500
CONFIG_HATCH_RATE=50

LOCUST=/usr/local/bin/locust

LOCUST_ARGS="--host=http://$LOADBALANCER --no-web --clients $CONFIG_CLIENTS --hatch-rate $CONFIG_HATCH_RATE --print-stats --only-summary -f "

if [ ! -e $LOCUST -a -x $LOCUST ] 
then
	echo "Locust does not exist or is not executable"
	exit 1
fi

for CONFIG in ${LOCUST_CONFIGS[@]}
do
	echo running $CONFIG
	for SERVER in ${SERVERS[@]}
	do
		SERVER_NAME="${SERVER%%:*}"
		SERVER_IP="${SERVER##*:}"
		echo "starting iostat on remote server $SERVER_IP ($SERVER_NAME)"
		nohup timeout $CONFIG_IOSTAT_DURATION ssh -oStrictHostKeyChecking=no ubuntu@$SERVER_IP iostat -c -d 20 >> $LOG_DIR/$CONFIG.$SERVER_NAME.log &
    	done

	echo "executing $CONFIG for $CONFIG_DURATION seconds"
	timeout $CONFIG_DURATION $LOCUST $LOCUST_ARGS ${CONFIG} >> $LOG_DIR/$CONFIG.log 2>&1
	echo "sleeping for $CONFIG_INTERVAL seconds"
	sleep $CONFIG_INTERVAL || true
done
