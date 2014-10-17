#!/bin/bash

# This is setup script for Cassandra Single-Node server.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo "Error: This script must be run as root" 1>&2
	exit 1
fi

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

	echo "export JAVA_HOME=/usr/lib/jvm/java-7-oracle" >> /etc/profile.d/dataplay.sh

	. /etc/profile
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

	echo "export CASSANDRA_CONFIG=/etc/cassandra" >> /etc/profile.d/dataplay.sh

	. /etc/profile

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

backup_cassandra () {
	echo "TODO: manually"
	# On Source:
	# cd /var/lib/cassandra/data/dp/
	#
	# For each column family
	#
	# tar -zcf author.tar.gz 1412773254804/
	# scp author.tar.gz ubuntu@109.231.121.88:/home/ubuntu/cassandra-backup
	#
	# tar -zcf entity.tar.gz 1412773254804/
	# scp entity.tar.gz ubuntu@109.231.121.88:/home/ubuntu/cassandra-backup
	#
	# tar -zcf image.tar.gz 1412773254804/
	# scp image.tar.gz ubuntu@109.231.121.88:/home/ubuntu/cassandra-backup
	#
	# tar -zcf keyword.tar.gz 1412773254804/
	# scp keyword.tar.gz ubuntu@109.231.121.88:/home/ubuntu/cassandra-backup
	#
	# tar -zcf related.tar.gz 1412773254804/
	# scp related.tar.gz ubuntu@109.231.121.88:/home/ubuntu/cassandra-backup
	#
	# tar -zcf response.tar.gz 1412773254804/
	# scp response.tar.gz ubuntu@109.231.121.88:/home/ubuntu/cassandra-backup
	#
	#
	# On Destination:
	#	1. Make sure that schema defination exists
	#	2. nodetool snapshot dp
	#	3. sevice cassandra stop
	#	4. clean commitlog, cd /var/lib/cassandra/data/commitlog/ && rm -r *.log
	#	5. extract files to individual column dir
	# 		e.g. /var/lib/cassandra/data/dp/keyword-<UUID>
	# 		tar -zxf /home/ubuntu/cassandra-backup/keyword.tar.gz -C .
	#		mv 1412773254804/* .
	#		rm -r 1412773254804
}

export_variables () {
	. /etc/profile

	su - ubuntu -c ". /etc/profile"
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

echo "[$(timestamp)] ---- 6. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
