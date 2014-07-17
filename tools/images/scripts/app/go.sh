#!/bin/bash

# This is setup script for App (Go server)
# 1. Install Ubuntu base image or dataplay-ubuntu-base (recommended)
# 2. Run this script as 'sudo'
#
# Note: Installing from pre-configured base image is highly recommended
#	e.g.
#		dataplay-go-master for Master/Producer instance

HOSTNAME=$(hostname)
HOSTLOCAL="127.0.1.1"
echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts

apt-get update
apt-get -y dist-upgrade
apt-get install -y build-essential sudo openssh-server screen gcc curl git mercurial bzr make binutils bison wget python-software-properties xclip htop zip

# Node.js
apt-add-repository -y ppa:chris-lea/node.js
apt-get update
apt-get install -y python g++ make nodejs
npm install -g grunt-cli

# Go
wget http://golang.org/dl/go1.3.linux-amd64.tar.gz
tar xf go1.3.linux-amd64.tar.gz
echo "export GOROOT=\$HOME/go" >> ~/.profile
echo "PATH=$PATH:\$GOROOT/bin" >> ~/.profile
source ~/.profile
mkdir ~/gocode
echo "export GOPATH=\$HOME/gocode" >> ~/.profile
echo "PATH=\$PATH:\$GOPATH/bin" >> ~/.profile
source ~/.profile
mkdir www
cd www

echo "export DATABASE=109.231.121.12" >> ~/.profile
echo "export redishost=109.231.121.13:6379" >> ~/.profile
source ~/.profile

# TODO: git clone & ./start.sh --DBHost="109.231.121.12"
# sudo iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 3000
# sudo iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 443 -j REDIRECT --to-port 3443
