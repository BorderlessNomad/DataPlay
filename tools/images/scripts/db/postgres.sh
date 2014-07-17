#!/bin/bash

# This is setup script for DB (PostreSQL server)
# 1. Install Ubuntu base image or dataplay-ubuntu-base (recommended)
# 2. Run this script as 'sudo'
#
# Note: Installing from pre-configured base image is highly recommended
#	e.g.
#		dataplay-posgresql for DB instance

HOSTNAME=$(hostname)
HOSTLOCAL="127.0.1.1"
echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts

apt-get update
apt-get -y dist-upgrade
apt-get install -y build-essential sudo openssh-server screen gcc curl git make binutils bison wget python-software-properties htop zip

# Install latest PostgresSQL with GIS support
apt-get install -y postgresql postgresql-contrib postgresql-9.3-postgis-2.1 postgresql-client libpq-dev
apt-get clean && sudo rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# sudo -u postgres

################## Extras ####################
# pg_dump -v -cC -f dataplay.sql -h 10.0.0.2 -U playgen dataplay
# pg_dump -v -cC -f dataplay.sql -h localhost -U playgen dataplay
# gzip -vk dataplay.sql
# scp dataplay.sql.gz ubuntu@109.231.121.12:/home/ubuntu
# gunzip -vk dataplay.sql.gz
# psql -h localhost -U playgen -d dataplay -f dataplay.sql
