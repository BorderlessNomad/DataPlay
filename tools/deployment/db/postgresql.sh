#!/bin/bash

# This is setup script for PostreSQL Database server.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

APP_HOST=$(ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}')

PGPOOL_API_HOST=$(ss-get --timeout 360 pgpool.hostname)
PGPOOL_API_PORT="1937"

JCATASCOPIA_REPO="109.231.126.62"
JCATASCOPIA_DASHBOARD="109.231.122.112"

timestamp () {
	date +"%F %T,%3N"
}

setuphost () {
	HOSTNAME=$(hostname)
	HOSTLOCAL="127.0.1.1"
	echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts
}

install_postgres () {
	apt-get update
	apt-get install -y axel postgresql postgresql-contrib postgresql-client libpq-dev
	apt-get autoclean
	service postgresql restart
}

setup_database () {
	cd # /var/lib/postgresql

	DB_USER="playgen"
	DB_PASSWORD="aDam3ntiUm"
	DB_NAME="dataplay"
	DB_VERSION="9.4"

	# Create a PostgreSQL user named 'playgen' with 'aDam3ntiUm' as the password and
	# then create a database 'dataplay' owned by the 'playgen' role.
	psql --command "CREATE USER $DB_USER WITH SUPERUSER PASSWORD '$DB_PASSWORD';" && \
	createdb -O $DB_USER $DB_NAME

	# Adjust PostgreSQL configuration so that remote connections to the database are possible.
	# From Flexiant clusters & PlayGen dev IP
	echo "host    all             all             109.231.121.0/24        md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf
	echo "host    all             all             109.231.122.0/24        md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf
	echo "host    all             all             109.231.123.0/24        md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf
	echo "host    all             all             109.231.124.0/24        md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf
	echo "host    all             all             109.231.125.0/24        md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf
	echo "host    all             all             109.231.126.0/24        md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf
	echo "host    all             all             213.122.181.2/32        md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf
	echo "host    all             all             149.11.102.50/32        md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf
	echo "host    all             all             134.60.64.0/24        md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf
	echo "host    all             all             192.168.0.0/24        md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf

	# And add 'listen_addresses' to '/etc/postgresql/$DB_VERSION/main/postgresql.conf'
	echo "listen_addresses='*'" >> /etc/postgresql/$DB_VERSION/main/postgresql.conf
}

import_data () {
	cd # /var/lib/postgresql

	MAX_RETRIES="200"

	DB_HOST="localhost"
	DB_PORT="5432"
	DB_USER="playgen"
	DB_PASSWORD="aDam3ntiUm"
	DB_NAME="dataplay"

	LASTDATE=$(date +%Y-%m-%d) # Today
	BACKUP_HOST="109.231.121.72" # Flexiant
	BACKUP_PORT="8080"
	BACKUP_DIR="postgresql/$LASTDATE-daily"
	BACKUP_USER="playgen"
	BACKUP_PASS="D@taP1aY"
	BACKUP_FILE="$DB_NAME.sql.gz"

	echo "$DB_HOST:$DB_PORT:$DB_NAME:$DB_USER:$DB_PASSWORD" > .pgpass && chmod 0600 .pgpass

	i="1"
	if [[ $i -ge $MAX_RETRIES ]]; then
		echo >&2 "Error: Unable to fetch '$BACKUP_FILE' from backup server."; exit 1;
	fi
	until axel -a "http://$BACKUP_USER:$BACKUP_PASS@$BACKUP_HOST:$BACKUP_PORT/$BACKUP_DIR/$BACKUP_FILE"; do
		LASTDATE=$(date +%Y-%m-%d --date="$LASTDATE -1 days") # Decrement by 1 Day
		BACKUP_DIR="postgresql/$LASTDATE-daily"
		echo "Latest backup not available, try fetching $LASTDATE"
		i=$[$i+1]
	done

	gunzip -vk $BACKUP_FILE
	#nohup psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f $DB_NAME.sql > $DB_NAME.import.log 2>&1&
	psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f $DB_NAME.sql
}

update_iptables () {
	iptables -A INPUT -p tcp --dport 5432 -j ACCEPT # PostgreSQL listener

	iptables-save
}

setup_pgpool_access() {
	DB_VERSION="9.4"
	DB_HOST="localhost"
	DB_PORT="5432"
	DB_USER="playgen"
	DB_PASSWORD="aDam3ntiUm"
	DB_NAME="dataplay"
	PGPOOL_VERSION="3.4.3"

	cp /var/lib/postgresql/.pgpass ~/.pgpass

	mkdir ~/pgpool-local

	echo "host    all             playgen         $PGPOOL_API_HOST/32       md5" >> /etc/postgresql/$DB_VERSION/main/pg_hba.conf

	service postgresql restart

	apt-get install -y postgresql-9.4-pgpool2

	#apt-get install -y build-essential

	wget http://www.pgpool.net/download.php?f=pgpool-II-$PGPOOL_VERSION.tar.gz -O pgpool-II-$PGPOOL_VERSION.tar.gz

	tar -xvzf pgpool-II-$PGPOOL_VERSION.tar.gz

	cp pgpool-II-$PGPOOL_VERSION/src/sql/pgpool_adm/pgpool_adm.sql.in ~/pgpool-local/pgpool_adm.sql
	cp pgpool-II-$PGPOOL_VERSION/src/sql/pgpool-recovery/pgpool-recovery.sql.in ~/pgpool-local/pgpool-recovery.sql
	cp pgpool-II-$PGPOOL_VERSION/src/sql/pgpool-regclass/pgpool-regclass.sql.in ~/pgpool-local/pgpool-regclass.sql

	###
	# ls -al /usr/lib/postgresql/9.4/lib/ | grep 'pgpool'
	# -rw-r--r-- 1 root root   14432 Nov  6 21:15 pgpool_adm.so
	# -rw-r--r-- 1 root root   14040 Nov  6 21:15 pgpool-recovery.so
	# -rw-r--r-- 1 root root    9944 Nov  6 21:15 pgpool-regclass.so
	###
	sed -i "s/MODULE_PATHNAME/\/usr\/lib\/postgresql\/$DB_VERSION\/lib\/pgpool_adm/g" ~/pgpool-local/pgpool_adm.sql
	# Note: error on line # 45 & 51 should retrun integer
	sed -i "43,51s/record/integer/" ~/pgpool-local/pgpool_adm.sql

	sed -i "s/MODULE_PATHNAME/\/usr\/lib\/postgresql\/$DB_VERSION\/lib\/pgpool-recovery/g" ~/pgpool-local/pgpool-recovery.sql
	sed -i "s/\$libdir\/pgpool-recovery/\/usr\/lib\/postgresql\/$DB_VERSION\/lib\/pgpool-recovery/g" ~/pgpool-local/pgpool-recovery.sql

	sed -i "s/MODULE_PATHNAME/\/usr\/lib\/postgresql\/$DB_VERSION\/lib\/pgpool-regclass/g" ~/pgpool-local/pgpool-regclass.sql

	psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f ~/pgpool-local/pgpool_adm.sql
	psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f ~/pgpool-local/pgpool-recovery.sql
	psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f ~/pgpool-local/pgpool-regclass.sql

	service postgresql restart
}

inform_pgpool () {
	retries=0
	until curl -H "Content-Type: application/json" -X POST -d "{\"ip\":\"$APP_HOST\"}" http://$PGPOOL_API_HOST:$PGPOOL_API_PORT; do
		echo "[$(timestamp)] PGPOOL Server is not up yet, retry... [$(( retries++ ))]"
		sleep 5
	done
}

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install PostgresSQL ----"
install_postgres

echo "[$(timestamp)] ---- 3. Setup Database ----"
su postgres -c "$(typeset -f setup_database); setup_database" # Run function as user 'postgres'

echo "[$(timestamp)] ---- 4. Restart PostgreSQL as root ----"
service postgresql restart

echo "[$(timestamp)] ---- 5. Import Data ----"
su postgres -c "$(typeset -f import_data); import_data" # Run function as user 'postgres'

echo "[$(timestamp)] ---- 6. Setup pgpool access ----"
setup_pgpool_access

echo "[$(timestamp)] ---- 7. Inform pgpool (Add) ----"
inform_pgpool

echo "[$(timestamp)] ---- 8. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
