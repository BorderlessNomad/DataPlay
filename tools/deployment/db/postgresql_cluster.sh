#!/bin/bash

# This is setup script for PostreSQL Cluster server.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
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

install_pgpool () {
	DB_USER="playgen"
	DB_PASSWORD="aDam3ntiUm"

	yum install -y http://yum.postgresql.org/9.4/redhat/rhel-7-x86_64/pgdg-centos94-9.4-1.noarch.rpm
	yum install -y http://www.pgpool.net/yum/rpms/3.4/redhat/rhel-7-x86_64/pgpool-II-release-3.4-1.noarch.rpm

	yum update -y

	yum install -y pgpool-II-94 pgpool-II-94-extensions postgresql94

	/usr/pgsql-9.4/bin/postgresql94-setup initdb

	systemctl start postgresql-9.4
	systemctl enable postgresql-9.4

	cp /etc/pgpool-II-94/pcp.conf.sample /etc/pgpool-II-94/pcp.conf
	echo "$DB_USER:`pg_md5 $DB_PASSWORD`" >> /etc/pgpool-II-94/pcp.conf

	cp /etc/pgpool-II-94/pgpool.conf.sample /etc/pgpool-II-94/pgpool.conf
	# Connections
	sed -i "s/^listen_addresses = .localhost./listen_addresses = '*'/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^port = .*/port = 5432/" /etc/pgpool-II-94/pgpool.conf
	# Logs
	sed -i "s/^log_destination = .stderr./log_destination = 'syslog'/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^log_connections = off/log_connections = on/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^log_hostname =.*/log_hostname = on/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^log_statement = off/log_statement = on/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^log_per_node_statement = off/log_per_node_statement = on/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^log_standby_delay = 'none'/log_standby_delay = 'always'/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^syslog_facility =.*/syslog_facility = 'daemon.info'/" /etc/pgpool-II-94/pgpool.conf
	# Health check
	sed -i "s/^health_check_period =.*/health_check_period = 10/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^health_check_user =.*/health_check_user = 'admin'/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^health_check_password =.*/health_check_password = 'password123'/" /etc/pgpool-II-94/pgpool.conf
	# Pools
	sed -i "s/^enable_pool_hba = off/enable_pool_hba = on/" /etc/pgpool-II-94/pgpool.conf
	# Master/Slave Mode + Streaming Replication
	sed -i "s/^master_slave_mode = off/master_slave_mode = on/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^master_slave_sub_mode =.*/master_slave_sub_mode = 'stream'/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^sr_check_period = 0/sr_check_period = 10/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^sr_check_user =.*/sr_check_user = '$DB_USER'/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^sr_check_password =.*/sr_check_password = '$DB_PASSWORD'/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^delay_threshold = 0/delay_threshold = 10000000/" /etc/pgpool-II-94/pgpool.conf
	# Watchdog
	sed -i "s/^use_watchdog =.*/use_watchdog = on/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^delegate_IP =.*/delegate_IP = '10.32.243.250'/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^netmask 255.255.255.0/netmask 255.255.255.128/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^heartbeat_device0 =.*/heartbeat_device0 = 'eth0'/" /etc/pgpool-II-94/pgpool.conf
	# LOAD BALANCING MODE
	sed -i "s/^load_balance_mode = off/load_balance_mode = on/" /etc/pgpool-II-94/pgpool.conf
	# FAILOVER AND FAILBACK
	sed -i "s@^failover_command = ''@failover_command = '/etc/pgpool-II-94/failover_stream.sh %d %H'@"
	# ONLINE RECOVERY
	sed -i "s/^recovery_user = 'nobody'/recovery_user = 'admin'/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^recovery_password = ''/recovery_password = 'password123'/" /etc/pgpool-II-94/pgpool.conf
	sed -i "s/^recovery_1st_stage_command = ''/recovery_1st_stage_command = 'basebackup.sh'/" /etc/pgpool-II-94/pgpool.conf

	cp /etc/pgpool-II-94/pool_hba.conf.sample /etc/pgpool-II-94/pool_hba.conf
	echo "host    all         all         0.0.0.0/0             md5" >> /etc/pgpool-II-94/pool_hba.conf
}

update_iptables () {
	yum install -y firewalld
	systemctl start firewalld.service
	systemctl enable firewalld.service

	firewall-cmd --permanent --add-port=5432/tcp
	firewall-cmd --permanent --add-port=80/tcp

	firewall-cmd --reload
}

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install pgpool-II ----"
install_pgpool

echo "[$(timestamp)] ---- 3. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
