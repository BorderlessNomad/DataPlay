#!/bin/bash

# This is setup script for Cassandra Single-Node server.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

IP=`ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}'`

timestamp () {
	date +"%F %T,%3N"
}

setuphost () {
	HOSTNAME=$(hostname)
	HOSTLOCAL="127.0.1.1"
	echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts
}

install_java () {
	echo oracle-java7-installer shared/accepted-oracle-license-v1-1 select true | debconf-set-selections && \
	apt-add-repository -y ppa:webupd8team/java && \
	apt-get update && \
	apt-get install -y axel oracle-java7-installer && \
	apt-get autoclean

	echo "export JAVA_HOME=/usr/lib/jvm/java-7-oracle" >> /etc/profile.d/dataplay.sh

	. /etc/profile
}

restart_cassandra() {
	service cassandra restart >> cassandra-service.log & # Start Cassandara in background
	echo "Waiting for Cassandra restart..."
	while ! grep -m1 '...done.' < cassandra-service.log ; do
		sleep 1
	done
	echo "Cassandra is UP!"
}

install_cassandra () {
	echo "deb http://debian.datastax.com/community stable main" | sudo tee -a /etc/apt/sources.list.d/cassandra.sources.list && \
	curl -L http://debian.datastax.com/debian/repo_key | sudo apt-key add - && \
	apt-get update && \
	apt-get install -y cassandra

	restart_cassandra

	echo "export CASSANDRA_CONFIG=/etc/cassandra" >> /etc/profile.d/dataplay.sh

	. /etc/profile

	# nodetool status # Verify that DataStax Community is running
}

configure_cassandra () {
	# sed -i -e "s/num_tokens/\#num_tokens/" /etc/cassandra/cassandra.yaml # Disable virtual nodes
	sed -i -e "s/^listen_address.*/listen_address: $IP/" /etc/cassandra/cassandra.yaml # Listen on IP of the container
	sed -i -e "s/^rpc_address.*/rpc_address: $IP/" /etc/cassandra/cassandra.yaml # Enable Remote connections
	# sed -i -e "s/^broadcast_rpc_address.*/broadcast_rpc_address: $IP/" /etc/cassandra/cassandra.yaml # Enable Remote connections
	sed -i -e "s/- seeds: \"127.0.0.1\"/- seeds: \"$IP\"/" /etc/cassandra/cassandra.yaml # Be your own seed

	# With virtual nodes disabled, we need to manually specify the token
	# sed -i -e "s/# JVM_OPTS=\"$JVM_OPTS -Djava.rmi.server.hostname=<public name>\"/ JVM_OPTS=\"$JVM_OPTS -Djava.rmi.server.hostname=$IP\"/" /etc/cassandra/cassandra-env.sh
	# echo "JVM_OPTS=\"\$JVM_OPTS -Dcassandra.initial_token=0\"" >> /etc/cassandra/cassandra-env.sh

	# netstat -an | grep 9160.*LISTEN

	restart_cassandra
}

install_opscenter () {
	apt-get install -y opscenter && \
	service opscenterd start

	# Connect using http://<IP>:8888
}

export_variables () {
	. /etc/profile

	su - ubuntu -c ". /etc/profile"
}

import_data () {
	LASTDATE=$(date +%Y-%m-%d) # Today
	BACKUP_HOST="109.231.121.85"
	BACKUP_PORT="8080"
	BACKUP_DIR="cassandra/$LASTDATE"
	BACKUP_USER="playgen"
	BACKUP_PASS="D@taP1aY"
	BACKUP_SCHEMA_FILE="dp-schema.cql"
	BACKUP_DATA_FILE="dp-data.tar.gz"
	KEYSPACE="dp"
	CASSANDRA_DIR="/var/lib/cassandra"
	DATA_DIR="$CASSANDRA_DIR/data"
	LOG_DIR="$CASSANDRA_DIR/commitlog"
	SOURCE_DIR="/tmp/cassandra-data"
	MAX_RETRIES="60"

	i="1"
	until [[ $i -lt $MAX_RETRIES ]] && axel -a "http://$BACKUP_USER:$BACKUP_PASS@$BACKUP_HOST:$BACKUP_PORT/$BACKUP_DIR/$BACKUP_SCHEMA_FILE"; do
		LASTDATE=$(date +%Y-%m-%d --date="$LASTDATE -1 days") # Decrement by 1 Day
		BACKUP_DIR="cassandra/$LASTDATE"
		echo "Latest $BACKUP_SCHEMA_FILE backup not available, trying $LASTDATE"
		i=$[$i+1]
	done
	if [[ $i -gt $MAX_RETRIES ]]; then
		echo >&2 "Error: Unable to fetch '$BACKUP_SCHEMA_FILE' from backup server."; exit 1;
	fi

	j="1"
	until [[ $j -lt $MAX_RETRIES ]] && axel -a "http://$BACKUP_USER:$BACKUP_PASS@$BACKUP_HOST:$BACKUP_PORT/$BACKUP_DIR/$BACKUP_DATA_FILE"; do
		LASTDATE=$(date +%Y-%m-%d --date="$LASTDATE -1 days") # Decrement by 1 Day
		BACKUP_DIR="cassandra/$LASTDATE"
		echo "Latest $BACKUP_DATA_FILE backup not available, trying $LASTDATE"
		j=$[$j+1]
	done
	if [[ $j -gt $MAX_RETRIES ]]; then
		echo >&2 "Error: Unable to fetch '$BACKUP_DATA_FILE' from backup server."; exit 1;
	fi

	cqlsh $IP -f $BACKUP_SCHEMA_FILE

	service cassandra stop

	mkdir -p $SOURCE_DIR
	tar -xzvf $BACKUP_DATA_FILE -C $SOURCE_DIR
	SOURCE_TABLES=`ls -l $SOURCE_DIR | egrep '^d' | awk '{print $9}'`
	for table in $SOURCE_TABLES; do
		table_name=$(echo $table | awk -F'-' '{print $1}')
		mv $SOURCE_DIR/$table/* $DATA_DIR/$KEYSPACE/$table_name-*
	done

	rm -rf $LOG_DIR/*.log

	restart_cassandra

	# sleep 5
	# nodetool -h $IP repair $KEYSPACE

	rm -rf $SOURCE_DIR
}

update_iptables () {
	# iptables -A INPUT -p tcp --dport 7000 -j ACCEPT # Internode communication (not used if TLS enabled) Used internal by Cassandra
	iptables -A INPUT -p tcp --dport 7199 -j ACCEPT # JMX
	iptables -A INPUT -p tcp --dport 8888 -j ACCEPT # OpsCenter
	iptables -A INPUT -p tcp --dport 9042 -j ACCEPT # CQL
	iptables -A INPUT -p tcp --dport 9160 -j ACCEPT # Thift client API

	iptables-save
}

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install Oracle Java 7 ----"
install_java

echo "[$(timestamp)] ---- 3. Install Cassandra ----"
install_cassandra

echo "[$(timestamp)] ---- 4. Configure Cassandra ----"
configure_cassandra

echo "[$(timestamp)] ---- 5. Export Variables ----"
export_variables

echo "[$(timestamp)] ---- 6. Import Data ----"
import_data

echo "[$(timestamp)] ---- 7. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
