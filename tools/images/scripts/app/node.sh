#!/bin/bash

# This is setup script for Compute/Node instance.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo "Error: This script must be run as root" 1>&2
	exit 1
fi

GO_VERSION="go1.3.3"
DEST="/home/ubuntu/www"
APP="dataplay"
WWW="www-src"
REPO="DataPlay"
BRANCH="develop"

# LOADBALANCER="109.231.121.27"
LOADBALANCER=$(ss-get --timeout 360 loadbalancer.hostname)

# DATABASE_HOST="109.231.121.13"
DATABASE_HOST=$(ss-get --timeout 360 postgres.hostname)
DATABASE_PORT="5432"

# REDIS_HOST="109.231.121.13"
REDIS_HOST=$(ss-get --timeout 360 redis_rabbitmq.hostname)
REDIS_PORT="6379"

# QUEUE_HOST="109.231.121.13"
QUEUE_HOST=$(ss-get --timeout 360 redis_rabbitmq.hostname)
QUEUE_PORT="5672"

# CASSANDRA_HOST="109.231.121.13"
CASSANDRA_HOST=$(ss-get --timeout 360 redis_rabbitmq.hostname)
CASSANDRA_PORT="9042"

timestamp () {
	date +"%F %T,%3N"
}

setuphost () {
	HOSTNAME=$(hostname)
	HOSTLOCAL="127.0.1.1"
	echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts
}

install_go () {
	apt-get install -y mercurial bzr

	mkdir -p /home/ubuntu && cd /home/ubuntu
	mkdir -p gocode && mkdir -p www

	wget -Nq https://storage.googleapis.com/golang/$GO_VERSION.linux-amd64.tar.gz
	tar xf $GO_VERSION.linux-amd64.tar.gz

	echo "export GOROOT=/home/ubuntu/go" >> /home/ubuntu/.profile
	echo "PATH=\$PATH:\$GOROOT/bin" >> /home/ubuntu/.profile

	echo "export GOPATH=/home/ubuntu/gocode" >> /home/ubuntu/.profile
	echo "PATH=\$PATH:\$GOPATH/bin" >> /home/ubuntu/.profile

	. /home/ubuntu/.profile
}

export_variables () {
	echo "export DP_LOADBALANCER=$LOADBALANCER" >> /home/ubuntu/.profile
	echo "export DP_DATABASE_HOST=$DATABASE_HOST" >> /home/ubuntu/.profile
	echo "export DP_DATABASE_PORT=$DATABASE_PORT" >> /home/ubuntu/.profile
	echo "export DP_REDIS_HOST=$REDIS_HOST" >> /home/ubuntu/.profile
	echo "export DP_REDIS_PORT=$REDIS_PORT" >> /home/ubuntu/.profile
	echo "export DP_QUEUE_HOST=$QUEUE_HOST" >> /home/ubuntu/.profile
	echo "export DP_QUEUE_PORT=$QUEUE_PORT" >> /home/ubuntu/.profile
	echo "export DP_CASSANDRA_HOST=$CASSANDRA_HOST" >> /home/ubuntu/.profile
	echo "export DP_CASSANDRA_PORT=$CASSANDRA_PORT" >> /home/ubuntu/.profile

	. /home/ubuntu/.profile
}

run_node () {
	SOURCE="https://github.com/playgenhub/$REPO/archive/"
	START="start.sh"
	LOG="ouput.log"

	QUEUE_USERNAME="playgen"
	QUEUE_PASSWORD="aDam3ntiUm"
	QUEUE_ADDRESS="amqp://$QUEUE_USERNAME:$QUEUE_PASSWORD@$QUEUE_HOST:$QUEUE_PORT/"
	QUEUE_EXCHANGE="playgen-prod"

	REQUEST_QUEUE="dataplay-request-prod"
	REQUEST_KEY="api-request-prod"
	REQUEST_TAG="consumer-request-prod"
	RESPONSE_QUEUE="dataplay-response-prod"
	RESPONSE_KEY="api-response-prod"
	RESPONSE_TAG="consumer-response-prod"

	MODE="1" # Node mode

	# Kill any running process
	if ps ax | grep -v grep | grep $APP > /dev/null; then
		echo "SHUTDOWN RUNING APP..."
		killall -9 $APP
	fi

	cd $DEST
	echo "Fetching latest ZIP"
	wget -Nq $SOURCE$BRANCH.zip -O $BRANCH.zip
	echo "Extracting from $BRANCH.zip"
	unzip -oq $BRANCH.zip
	if [ -d $APP ]; then
		rm -r $APP
	fi
	mkdir -p $APP
	echo "Moving files from $REPO-$BRANCH/ to $APP"
	mv -f $REPO-$BRANCH/* $APP
	cd $APP
	chmod u+x $START
	echo "Starting $START in Mode=$MODE"
	nohup sh $START --mode=$MODE --uri="$QUEUE_ADDRESS" --exchange="$QUEUE_EXCHANGE" --requestqueue="$REQUEST_QUEUE" --requestkey="$REQUEST_KEY" --reqtag="$REQUEST_TAG" --responsequeue="$RESPONSE_QUEUE" --responsekey="$RESPONSE_KEY" --restag="$RESPONSE_TAG" > $LOG 2>&1&
	echo "Done! $ sudo tail -f $DEST/$APP/$LOG for more details"
}

update_iptables () {
	# Accept direct connections to Compute Instance
	iptables -A INPUT -p tcp --dport 3000 -j ACCEPT # HTTP
	iptables -A INPUT -p tcp --dport 3443 -j ACCEPT # HTTPS

	# NAT Redirect
	iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 3000
	iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 443 -j REDIRECT --to-port 3443

	iptables-save
}

if [ "$(id -u)" != "0" ]; then
	echo "Error: This script must be run as root" 1>&2
	exit 1
fi

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install GO ----"
install_go

echo "[$(timestamp)] ---- 3. Export Variables ----"
export_variables

echo "[$(timestamp)] ---- 4. Run Compute (Node) Server ----"
run_node

echo "[$(timestamp)] ---- 5. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
