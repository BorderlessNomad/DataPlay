#!/bin/bash

# This is setup script for App (Go server)
# 1. Install Ubuntu base image or dataplay-ubuntu-base (recommended)
# 2. Run this script as 'user'
#
# Note: Installing from pre-configured base image is highly recommended
#	e.g.
#		dataplay-go-master for Master/Producer instance

# Go
install_go () {
	wget http://golang.org/dl/go1.3.linux-amd64.tar.gz
	tar xf go1.3.linux-amd64.tar.gz
	echo "export GOROOT=\$HOME/go" >> $HOME/.profile
	echo "PATH=$PATH:\$GOROOT/bin" >> $HOME/.profile
	source $HOME/.profile
	mkdir $HOME/gocode
	echo "export GOPATH=\$HOME/gocode" >> $HOME/.profile
	echo "PATH=\$PATH:\$GOPATH/bin" >> $HOME/.profile
	source $HOME/.profile
	mkdir $HOME/www
}

export_variables () {
	# @todo make POSTGRES & REDIS dynamic
	HOST_POSTGRES="109.231.121.12"
	HOST_REDIS="109.231.121.13:6379"

	echo "export DATABASE=$HOST_POSTGRES" >> $HOME/.profile
	echo "export redishost=$HOST_REDIS" >> $HOME/.profile
	source $HOME/.profile
}

if [ "$(id -u)" == "0" ]; then
	echo "Error: This script must as normal user" 1>&2
	exit 1
fi

# As User
echo "---- Running as User ----"

echo "1. ---- Install GO ----"
install_go
echo "2. ---- Export Variables ----"
export_variables

echo "---- Completed ----"
