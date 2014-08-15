#!/bin/bash

set -ex

# This is setup script for App (Go server)
# 1. Install Ubuntu base image or dataplay-ubuntu-base (recommended)
# 2. Run this script as 'sudo'
#
# Note: Installing from pre-configured base image is highly recommended
#	e.g.
#		dataplay-go-master for Master/Producer instance

setuphost () {
	HOSTNAME=$(hostname)
	HOSTLOCAL="127.0.1.1"
	echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts
}

update () {
	apt-get update
	apt-get -y upgrade
}

install_essentials () {
	apt-get install -y build-essential sudo openssh-server screen gcc curl git mercurial bzr make binutils bison wget python-software-properties htop zip
}

# Node.js
install_nodejs () {
	apt-add-repository -y ppa:chris-lea/node.js
	apt-get update
	apt-get install -y python g++ make nodejs
	npm install -g grunt grunt-cli
	if [ ! -d /home/ubuntu/.npm ]; then
		mkdir -p /home/ubuntu/.npm
	fi
	chown -R ubuntu:ubuntu /home/ubuntu/.npm # Fix permissions
}

update_iptables () {
	# Monitoring ports 80, 8080, 4242, 4243, 4245
	iptables -A INPUT -p tcp --dport 80 -j ACCEPT
	iptables -A INPUT -p tcp --dport 443 -j ACCEPT
	iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
	iptables -A INPUT -p tcp --dport 4242 -j ACCEPT
	iptables -A INPUT -p tcp --dport 4243 -j ACCEPT
	iptables -A INPUT -p tcp --dport 4245 -j ACCEPT

	# NAT Redirect
	iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 3000
	iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 443 -j REDIRECT --to-port 3443

	iptables-save
}

install_go () {
	wget http://golang.org/dl/go1.3.linux-amd64.tar.gz
	tar xf go1.3.linux-amd64.tar.gz
	echo "export GOROOT=/home/ubuntu/go" >> /home/ubuntu/.profile
	echo "PATH=$PATH:\$GOROOT/bin" >> /home/ubuntu/.profile
	. /home/ubuntu/.profile
	mkdir /home/ubuntu/gocode
	echo "export GOPATH=/home/ubuntu/gocode" >> /home/ubuntu/.profile
	echo "PATH=\$PATH:\$GOPATH/bin" >> /home/ubuntu/.profile
	. /home/ubuntu/.profile
	if [ ! -d /home/ubuntu/www ]; then
		mkdir /home/ubuntu/www
	fi
}

export_variables () {
	# @todo make POSTGRES & REDIS dynamic
	HOST_POSTGRES=$(ss-get --timeout 360 postgres.hostname)
	HOST_REDIS="109.231.121.13:6379"

	echo "export DATABASE=$HOST_POSTGRES" >> /home/ubuntu/.profile
	echo "export redishost=$HOST_REDIS" >> /home/ubuntu/.profile
	. /home/ubuntu/.profile
}

run_master () {
	APP="dataplay"
	REPO="DataPlay"
	SOURCE="https://github.com/playgenhub/$REPO/archive/"
	BRANCH="develop"
	DEST="/home/ubuntu/www"
	START="start.sh"
	LOG="ouput.log"

	QUEUE_USERNAME="playgen"
	QUEUE_PASSWORD="aDam3ntiUm"
	QUEUE_SERVER="109.231.121.13"
	QUEUE_PORT="5672"
	QUEUE_ADDRESS="amqp://$QUEUE_USERNAME:$QUEUE_PASSWORD@$QUEUE_SERVER:$QUEUE_PORT/"
	QUEUE_EXCHANGE="playgen-prod"

	REQUEST_QUEUE="dataplay-request-prod"
	REQUEST_KEY="api-request-prod"
	REQUEST_TAG="consumer-request-prod"
	RESPONSE_QUEUE="dataplay-response-prod"
	RESPONSE_KEY="api-response-prod"
	RESPONSE_TAG="consumer-response-prod"

	MODE="2" # Master mode

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
	mkdir $APP
	echo "Moving files from $REPO-$BRANCH/ to $APP"
	mv -f $REPO-$BRANCH/* $APP
	cd $APP
	chmod u+x $START
	echo "Starting $START in Mode=$MODE"
	nohup sh $START --mode=$MODE --uri="$QUEUE_ADDRESS" --exchange="$QUEUE_EXCHANGE" --requestqueue="$REQUEST_QUEUE" --requestkey="$REQUEST_KEY" --reqtag="$REQUEST_TAG" --responsequeue="$RESPONSE_QUEUE" --responsekey="$RESPONSE_KEY" --restag="$RESPONSE_TAG" > $LOG 2>&1&
	echo "Done! $ sudo tail -f $DEST/$APP/$LOG for more details"
}

if [ "$(id -u)" != "0" ]; then
	echo "Error: This script must be run as root" 1>&2
	exit 1
fi

# As root
echo "---- Running as Root ----"

echo "1. ---- Setup Host ----"
setuphost

echo "2. ---- Update system ----"
update

echo "3. ---- Install essential packages ----"
install_essentials

echo "4. ---- Install Node.js ----"
install_nodejs

echo "5. ---- Update IPTables rules ----"
update_iptables

echo "6. ---- Install GO ----"
install_go

echo "7. ---- Export Variables ----"
export_variables

echo "8. ---- Run Master ----"
run_master

echo "---- Completed ----"

exit 0
