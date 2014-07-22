#!/bin/bash

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
	apt-get install -y build-essential sudo openssh-server screen gcc curl git mercurial bzr make binutils bison wget python-software-properties xclip htop zip
}

# Node.js
install_nodejs () {
	apt-add-repository -y ppa:chris-lea/node.js
	apt-get update
	apt-get install -y python g++ make nodejs
	npm install -g grunt grunt-cli
}

update_iptables () {
	iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 3000
	iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 443 -j REDIRECT --to-port 3443
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

echo "---- Completed ----"
