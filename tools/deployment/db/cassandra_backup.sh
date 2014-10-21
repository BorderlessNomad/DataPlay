#!/bin/bash

# Cassandra source backup script (should run as CRON)
#
###
# # Daily backup of Cassandra data on 22:30
# 30      22      *       *       *       /root/cassandra_backup.sh >> /root/cassandra_backup.log
#
# # Sync it with Backup server
# # VULTR
# 45      22      *       *       *       /usr/bin/rsync --stats --progress --whole-file -a --no-group --no-owner --no-perms -v -z -r /var/lib/cassandra/backups/ root@108.61.197.87:~/backups/cassandra/ >> /var/lib/cassandra/rsync.vultr.log
# # Flexiant
# 45      22      *       *       *       /usr/bin/rsync --stats --progress --whole-file -a --no-group --no-owner --no-perms -v -z -r /var/lib/cassandra/backups/ ubuntu@109.231.121.85:~/backups/cassandra/ >> /var/lib/cassandra/rsync.flexiant.log
###
#	1. run $ nodetool snapshot dp (using CRON)
#		output is,
#			Requested creating snapshot(s) for [dp] with snapshot name [1413378613098]
#			Snapshot directory: 1413378613098
#	2. extract timestamp info from output of 1. e.g. Snapshot directory: 1413378613098
#	3. for each dir in /var/lib/cassandra/data/dp copy content of snapshots/<timestamp>
#		e.g. /var/lib/cassandra/data/dp/response-3b35cc404d6311e497ddbd0e0515b177/snapshots/1413378613098
#	4. compress the dir and place it in response-3b35cc404d6311e497ddbd0e0515b177
###

set -ex

timestamp () {
	date +"%F %T,%3N"
}

echo "[$(timestamp)] ---- Started ----"

HOST="172.17.0.78" # Local
KEYSPACE="dp"
NODETOOL=$(nodetool -h $HOST snapshot $KEYSPACE)
TIMESTAMP=${NODETOOL#*: }
SOURCE="/var/lib/cassandra/data/$KEYSPACE"
TABLES=`ls -l $SOURCE | egrep '^d' | awk '{print $9}'`
BACKUP="/var/lib/cassandra/backups"
DATE=$(date +%Y-%m-%d)

mkdir -p $BACKUP/$DATE && cd $BACKUP/$DATE

# Schema
cqlsh $HOST -e "DESCRIBE KEYSPACE $KEYSPACE;" > $BACKUP/$DATE/$KEYSPACE-schema.cql 2>&1

# Tables
for table in $TABLES; do
	mkdir -p $BACKUP/$KEYSPACE/$table
	cp -R $SOURCE/$table/snapshots/$TIMESTAMP/. $BACKUP/$KEYSPACE/$table
done

tar -czvf $KEYSPACE-data.tar.gz -C $BACKUP/$KEYSPACE .

rm -rf $BACKUP/$KEYSPACE

echo "[$(timestamp)] ---- Completed ----"

exit 0
