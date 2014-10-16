#!/bin/bash

# Cassandra source backup script (should run as CRON)
###
# # Daily backup of Cassandra data on 22:30
# 30      22      *       *       *       /root/cassandra_backup.sh >> /root/cassandra_backup.log
#
# # Sync it with Backup server
# *       *       *       *       *       /usr/bin/rsync -a --no-group --no-owner --no-perms -v -z -r --delete /var/lib/cassandra/backups/ root@108.61.197.87:~/backups/cassandara/ >> /var/lib/cassandra/rsync.log
###


set -ex

timestamp () {
	date +"%F %T,%3N"
}

echo "[$(timestamp)] ---- Started ----"

KEYSPACE="dp"
NODETOOL=$(nodetool snapshot $KEYSPACE)
TIMESTAMP=${NODETOOL#*: }
SOURCE="/var/lib/cassandra/data/$KEYSPACE"
TABLES=`ls -l $SOURCE | egrep '^d' | awk '{print $9}'`
BACKUP="/var/lib/cassandra/backups"
DATE=$(date +%Y-%m-%d)

for table in $TABLES; do
	mkdir -p $BACKUP/$KEYSPACE/$table
	cp -R $SOURCE/$table/snapshots/$TIMESTAMP/. $BACKUP/$KEYSPACE/$table
done

mkdir -p $BACKUP/$DATE && cd $BACKUP/$DATE
tar -czvf $KEYSPACE-data.tar.gz -C $BACKUP/$KEYSPACE .
# tar -czvf $KEYSPACE-schema.tar.gz -C $BACKUP/$KEYSPACE .

rm -rf $BACKUP/$KEYSPACE

echo "[$(timestamp)] ---- Completed ----"

exit 0
