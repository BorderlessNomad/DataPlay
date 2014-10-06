#!/bin/bash

# This is setup script for Cassandra Single-Node server.

set -ex

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
	apt-get install -y oracle-java7-installer && \
	apt-get autoclean

	echo "export JAVA_HOME=/usr/lib/jvm/java-7-oracle" > /home/ubuntu/.profile
}

install_cassandra () {
	echo "deb http://debian.datastax.com/community stable main" | sudo tee -a /etc/apt/sources.list.d/cassandra.sources.list && \
	curl -L http://debian.datastax.com/debian/repo_key | sudo apt-key add - && \
	apt-get update && \
	apt-get install -y cassandra

	service cassandra restart > cassandra-service.log & # Start Cassandara in background
	echo "Waiting for Cassandra initialize..."
	sleep 1
	while ! grep -m1 '...done.' < cassandra-service.log ; do
		sleep 1
	done
	echo "Cassandra is UP!"

	echo "export CASSANDRA_CONFIG=/etc/cassandra" > /home/ubuntu/.profile
	. /home/ubuntu/.profile

	# nodetool status # Verify that DataStax Community is running
}

test_cassandra () {
	if [[ -n $(cqlsh -e "CREATE KEYSPACE test WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };") ]]; then
		echo >&2 "Error: Cassandra is not running.";
		exit 1
	fi
}

configure_cassandra () {
	IP=`ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}'`

	# sed -i -e "s/num_tokens/\#num_tokens/" /etc/cassandra/cassandra.yaml # Disable virtual nodes
	sed -i -e "s/^listen_address.*/listen_address: $IP/" /etc/cassandra/cassandra.yaml # Listen on IP of the container
	sed -i -e "s/^rpc_address.*/rpc_address: $IP/" /etc/cassandra/cassandra.yaml # Enable Remote connections
	# sed -i -e "s/^broadcast_rpc_address.*/broadcast_rpc_address: $IP/" /etc/cassandra/cassandra.yaml # Enable Remote connections
	sed -i -e "s/- seeds: \"127.0.0.1\"/- seeds: \"$IP\"/" /etc/cassandra/cassandra.yaml # Be your own seed

	# With virtual nodes disabled, we need to manually specify the token
	# sed -i -e "s/# JVM_OPTS=\"$JVM_OPTS -Djava.rmi.server.hostname=<public name>\"/ JVM_OPTS=\"$JVM_OPTS -Djava.rmi.server.hostname=$IP\"/" /etc/cassandra/cassandra-env.sh
	# echo "JVM_OPTS=\"\$JVM_OPTS -Dcassandra.initial_token=0\"" >> /etc/cassandra/cassandra-env.sh

	# netstat -an | grep 9160.*LISTEN

	service cassandra restart
}

install_opscenter () {
	apt-get install -y opscenter && \
	service opscenterd start

	# Connect using http://<IP>:8888
}

update_iptables () {
	# iptables -A INPUT -p tcp --dport 7000 -j ACCEPT # Internode communication (not used if TLS enabled) Used internal by Cassandra
	iptables -A INPUT -p tcp --dport 7199 -j ACCEPT # JMX
	iptables -A INPUT -p tcp --dport 8888 -j ACCEPT # OpsCenter
	iptables -A INPUT -p tcp --dport 9042 -j ACCEPT # CQL
	iptables -A INPUT -p tcp --dport 9160 -j ACCEPT # Thift client API

	iptables-save
}

if [ "$(id -u)" != "0" ]; then
	echo "Error: This script must be run as root" 1>&2
	exit 1
fi

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install Oracle Java 7 ----"
install_java

echo "[$(timestamp)] ---- 3. Install Cassandra ----"
install_cassandra

echo "[$(timestamp)] ---- 4. Configure Cassandra ----"
configure_cassandra

echo "[$(timestamp)] ---- 5. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
