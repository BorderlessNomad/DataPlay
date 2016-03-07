#!/bin/bash

# This is setup script for PostreSQL Cluster server.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

URL="https://raw.githubusercontent.com"
USER="playgenhub"
REPO="DataPlay"
BRANCH="master"
SOURCE="$URL/$USER/$REPO/$BRANCH"

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

	echo "deb http://apt.postgresql.org/pub/repos/apt/ wily-pgdg main" > /etc/apt/sources.list.d/pgdg.list
	apt-get install -y wget ca-certificates rsyslog
	wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
	apt-get update
	apt-get install -y axel postgresql-client-9.4 pgpool2

	# Provides UDP syslog reception
	sed -i 's/^#module(load="imudp")/module(load="imudp")/g' /etc/rsyslog.conf
	sed -i 's/^#input(type="imudp" port="514")/input(type="imudp" port="514")/g' /etc/rsyslog.conf

	# Provides TCP syslog reception
	sed -i 's/^#module(load="imtcp")/module(load="imtcp")/g' /etc/rsyslog.conf
	sed -i 's/^#input(type="imtcp" port="514")/input(type="imtcp" port="514")/g' /etc/rsyslog.conf

	echo "
	# Save PgPool-II log to pgpool.log
	local0.*                                                /var/log/pgpool.log" >> /etc/rsyslog.conf

	service rsyslog restart

	echo "$DB_USER:`pg_md5 $DB_PASSWORD`" >> /etc/pgpool2/pcp.conf

	echo "host    all         all         0.0.0.0/0             md5" >> /etc/pgpool2/pool_hba.conf

	pg_md5 -m -u $DB_USER $DB_PASSWORD # Generate pool_passwd

	chown postgres.postgres /etc/pgpool2/pool_passwd

	service pgpool2 restart
}

setup_pgpool_api () {
	command -v pgpool >/dev/null 2>&1 || { echo >&2 'Error: Command "pgpool" not found!'; exit 1; }

	command -v npm >/dev/null 2>&1 || { echo >&2 'Error: Command "npm" not found!'; exit 1; }

	command -v pm2 >/dev/null 2>&1 || { echo >&2 "Error: 'pm2' is not installed!"; exit 1; }

	command -v coffee >/dev/null 2>&1 || { echo >&2 'Error: "coffee-script" is not installed!'; exit 1; }

	cd /root && mkdir -p pgpool-api && cd pgpool-api

	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/db/api/app.coffee && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/db/api/package.json && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/db/api/cluster.json && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/db/api/pgpool.conf.template

	npm install

	coffee -cb app.coffee > app.js

	pm2 startup

	pm2 start app.js --name="pgpool-api" -o output.log -e errors.log

	pm2 save

	###
	# curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X POST -d '{"ip":"109.231.124.33"}' http://109.231.124.33:1937
	# curl -i -H "Accept: application/json" -X DELETE http://109.231.124.33:1937/109.231.124.33
	###
}

setup_pgpoolAdmin () {
	yum install -y httpd

	systemctl start httpd.service
	systemctl enable httpd.service

	yum install -y php php-fpm

	systemctl restart httpd.service

	yum install -y pgpoolAdmin

	chown -R root:root /var/www/html/pgpoolAdmin/

	chcon -R -t httpd_sys_content_rw_t /var/www/html/pgpoolAdmin/templates_c
	chcon -R -t httpd_sys_content_rw_t /var/www/html/pgpoolAdmin/conf/pgmgt.conf.php
	chcon -R -t httpd_sys_content_rw_t /etc/pgpool-II-94/pgpool.conf
	chcon -R -t httpd_sys_content_rw_t /etc/pgpool-II-94/pcp.conf
	chcon -R -t httpd_sys_content_rw_t /var/log/pgpool.log
}

update_iptables () {
	iptables -A INPUT -p tcp --dport 1937 -j ACCEPT # PgPool-II API
	iptables -A INPUT -p tcp --dport 9999 -j ACCEPT # PgPool-II listener

	iptables-save
}

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install pgpool-II ----"
install_pgpool

echo "[$(timestamp)] ---- 3. Setup pgpool API ----"
setup_pgpool_api

echo "[$(timestamp)] ---- 4. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
