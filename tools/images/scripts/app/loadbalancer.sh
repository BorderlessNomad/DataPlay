#!/bin/bash

set -ex

# This is setup script for App Master Load Balancer
# 1. Install Ubuntu base image or dataplay-ubuntu-base (recommended)
# 2. Run this script as 'sudo'
#
# Note: Installing from pre-configured base image is highly recommended
#	e.g.
#		dataplay-load-master

HOSTNAME=$(hostname)
HOSTLOCAL="127.0.1.1"
echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts

apt-get update
apt-get -y upgrade
apt-get install -y build-essential sudo openssh-server screen gcc curl git make binutils bison wget python-software-properties htop zip

# HAProxy 1.5
add-apt-repository -y ppa:vbernat/haproxy-1.5
apt-get update
apt-get install -y haproxy

# *:1936 playgen:D@taP1aY
exit 0
