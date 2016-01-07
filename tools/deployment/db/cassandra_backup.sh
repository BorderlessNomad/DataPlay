#!/bin/bash

# Cassandra source backup script (should run as 'root')
#
# Copy ssh keys before RSYNC
# ssh-keygen
# ssh-copy-id -i ~/.ssh/id_rsa.pub ubuntu@109.231.121.72
# ssh-copy-id -i ~/.ssh/id_rsa.pub root@108.61.197.87
#
###
# # Daily backup of Cassandra data on 22:30
# 30      22      *       *       *       /root/cassandra_backup.sh >> /root/cassandra_backup.log
#
# # Sync it with Backup servers
# # VULTR
# 45      22      *       *       *       /usr/bin/rsync --stats --progress --whole-file -a --no-group --no-owner --no-perms -v -z -r /var/lib/cassandra/backups/ root@108.61.197.87:~/backups/cassandra/ >> /var/lib/cassandra/rsync.vultr.log
# # Flexiant
# 45      22      *       *       *       /usr/bin/rsync --stats --progress --whole-file -a --no-group --no-owner --no-perms -v -z -r /var/lib/cassandra/backups/ ubuntu@109.231.121.72:~/backups/cassandra/ >> /var/lib/cassandra/rsync.flexiant.log
###
#	1. run $ nodetool snapshot dp (using ROOT)
#		output is,
#			Requested creating snapshot(s) for [dp] with snapshot name [1413378613098]
#			Snapshot directory: 1413378613098
#	2. extract timestamp info from output of 1. e.g. Snapshot directory: 1413378613098
#	3. for each dir in /var/lib/cassandra/data/dp copy content of snapshots/<timestamp>
#		e.g. /var/lib/cassandra/data/dp/response-3b35cc404d6311e497ddbd0e0515b177/snapshots/1413378613098
#	4. compress the dir and place it in response-3b35cc404d6311e497ddbd0e0515b177
###

set -x

timestamp () {
	date +"%F %T,%3N"
}

echo "[$(timestamp)] ---- 1. Start ----"

HOST=$(ip route get 8.8.8.8 | awk 'NR==1 {print $NF}')
LOCAL="localhost"
KEYSPACE="dataplay"
NODETOOL=$(nodetool -h $LOCAL snapshot $KEYSPACE)
TIMESTAMP=${NODETOOL#*: }
SOURCE="/var/lib/cassandra/data/$KEYSPACE"
TABLES=`ls -l $SOURCE | egrep '^d' | awk '{print $9}'`
BACKUP="/var/lib/cassandra/backups"
DATE=$(date +%Y-%m-%d)

echo "[$(timestamp)] ---- 2. Creating backup for $BACKUP/$DATE ----"

mkdir -p $BACKUP/$DATE && cd $BACKUP/$DATE

# Schema
cqlsh "$HOST" -e "DESCRIBE KEYSPACE $KEYSPACE;" > $BACKUP/$DATE/$KEYSPACE-schema.cql 2>&1

echo "[$(timestamp)] ---- 3. DESCRIBE KEYSPACE $KEYSPACE; successful ----"

# Tables
for table in $TABLES; do
	mkdir -p $BACKUP/$KEYSPACE/$table
	cp -R $SOURCE/$table/snapshots/$TIMESTAMP/. $BACKUP/$KEYSPACE/$table
done

echo "[$(timestamp)] ---- 4. Moving tables successful ----"

tar -czvf $KEYSPACE-data.tar.gz -C $BACKUP/$KEYSPACE .

echo "[$(timestamp)] ---- 5. $KEYSPACE-data.tar.gz successful ----"

rm -rf $BACKUP/$KEYSPACE

echo "[$(timestamp)] ---- 6. Cleaning BACKUP/$KEYSPACE successful ----"

echo "[$(timestamp)] ---- Completed ----"

exit 0
