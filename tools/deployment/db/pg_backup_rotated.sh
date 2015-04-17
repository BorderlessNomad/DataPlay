#!/bin/bash

# PostgreSQL source backup script (should run as CRON under user 'postgres')
#
# Copy ssh keys before RSYNC
# ssh-keygen
# ssh-copy-id -i ~/.ssh/id_rsa.pub ubuntu@109.231.122.208
# ssh-copy-id -i ~/.ssh/id_rsa.pub root@108.61.197.87
#
###
# # Daily backup of PostgreSQL data on 22:30
# 30      22      *       *       *       PGPASSWORD=aDam3ntiUm /var/lib/postgresql/pg_backup_rotated.sh >> /var/lib/postgresql/output.log &
#
# # Sync it with Backup servers
# # VULTR
# 45      22      *       *       *       /usr/bin/rsync --stats --progress --whole-file -a --no-group --no-owner --no-perms -v -z -r /var/lib/postgresql/backups/ root@108.61.197.87:~/backups/postgresql/ >> /var/lib/postgresql/rsync.vultr.log
# # Flexiant
# 45      22      *       *       *       /usr/bin/rsync --stats --progress --whole-file -a --no-group --no-owner --no-perms -v -z -r /var/lib/postgresql/backups/ ubuntu@109.231.122.208:~/backups/postgresql/ >> /var/lib/postgresql/rsync.flexiant.log
###

###########################
####### LOAD CONFIG #######
###########################
while [ $# -gt 0 ]; do
	case $1 in
		-c)
			CONFIG_FILE_PATH="$2"
			shift 2
			;;
		*)
			${ECHO} "Unknown Option \"$1\"" 1>&2
			exit 2
			;;
	esac
done

if [ -z $CONFIG_FILE_PATH ] ; then
	SCRIPTPATH=$(cd ${0%/*} && pwd -P)
	CONFIG_FILE_PATH="${SCRIPTPATH}/pg_backup.config"
fi

if [ ! -r ${CONFIG_FILE_PATH} ] ; then
	echo "Could not load config file from ${CONFIG_FILE_PATH}" 1>&2
	exit 1
fi

source "${CONFIG_FILE_PATH}"

###########################
#### PRE-BACKUP CHECKS ####
###########################
# Make sure we're running as the required backup user
if [ $BACKUP_USER != "" -a "$(id -un)" != "$BACKUP_USER" ] ; then
	echo "This script must be run as $BACKUP_USER. Exiting." 1>&2
	exit 1
fi


###########################
### INITIALISE DEFAULTS ###
###########################
if [ ! $HOSTNAME ]; then
	HOSTNAME="localhost"
fi;

if [ ! $USERNAME ]; then
	USERNAME="postgres"
fi;

###########################
#### START THE BACKUPS ####
###########################
function perform_backups() {
	SUFFIX=$1
	FINAL_BACKUP_DIR=$BACKUP_DIR"`date +\%Y-\%m-\%d`$SUFFIX/"

	echo "Making backup directory in $FINAL_BACKUP_DIR"

	if ! mkdir -p $FINAL_BACKUP_DIR; then
		echo "Cannot create backup directory in $FINAL_BACKUP_DIR. Go and fix it!" 1>&2
		exit 1;
	fi;

	###########################
	### SCHEMA-ONLY BACKUPS ###
	###########################
	for SCHEMA_ONLY_DB in ${SCHEMA_ONLY_LIST//,/ }; do
		SCHEMA_ONLY_CLAUSE="$SCHEMA_ONLY_CLAUSE or datname ~ '$SCHEMA_ONLY_DB'"
	done

	SCHEMA_ONLY_QUERY="select datname from pg_database where false $SCHEMA_ONLY_CLAUSE order by datname;"

	echo -e "\n\nPerforming schema-only backups"
	echo -e "--------------------------------------------\n"

	SCHEMA_ONLY_DB_LIST=`psql -h "$HOSTNAME" -U "$USERNAME" -At -c "$SCHEMA_ONLY_QUERY" postgres`

	echo -e "The following databases were matched for schema-only backup:\n${SCHEMA_ONLY_DB_LIST}\n"

	for DATABASE in $SCHEMA_ONLY_DB_LIST; do
		echo "Schema-only backup of $DATABASE"

		if ! pg_dump -Fp --no-acl --no-owner -h -s -h "$HOSTNAME" -U "$USERNAME" "$DATABASE" | gzip > $FINAL_BACKUP_DIR"$DATABASE"_SCHEMA.sql.gz.in_progress; then
			echo "[!!ERROR!!] Failed to backup database schema of $DATABASE" 1>&2
		else
			mv $FINAL_BACKUP_DIR"$DATABASE"_SCHEMA.sql.gz.in_progress $FINAL_BACKUP_DIR"$DATABASE"_SCHEMA.sql.gz
		fi
	done

	###########################
	###### FULL BACKUPS #######
	###########################
	for SCHEMA_ONLY_DB in ${SCHEMA_ONLY_LIST//,/ }; do
		EXCLUDE_SCHEMA_ONLY_CLAUSE="$EXCLUDE_SCHEMA_ONLY_CLAUSE and datname !~ '$SCHEMA_ONLY_DB'"
	done

	FULL_BACKUP_QUERY="select datname from pg_database where not datistemplate and datallowconn $EXCLUDE_SCHEMA_ONLY_CLAUSE order by datname;"

	echo -e "\n\nPerforming full backups"
	echo -e "--------------------------------------------\n"

	for DATABASE in `psql -h "$HOSTNAME" -U "$USERNAME" -At -c "$FULL_BACKUP_QUERY" postgres`; do
		if [ $ENABLE_PLAIN_BACKUPS = "yes" ]
		then
			echo "Plain backup of $DATABASE"

			if ! pg_dump -Fp --no-acl --no-owner -h "$HOSTNAME" -U "$USERNAME" "$DATABASE" | gzip > $FINAL_BACKUP_DIR"$DATABASE".sql.gz.in_progress; then
				echo "[!!ERROR!!] Failed to produce plain backup database $DATABASE" 1>&2
			else
				mv $FINAL_BACKUP_DIR"$DATABASE".sql.gz.in_progress $FINAL_BACKUP_DIR"$DATABASE".sql.gz
			fi
		fi

		if [ $ENABLE_CUSTOM_BACKUPS = "yes" ]
		then
			echo "Custom backup of $DATABASE"

			if ! pg_dump -Fc --no-acl --no-owner -h "$HOSTNAME" -U "$USERNAME" "$DATABASE" -f $FINAL_BACKUP_DIR"$DATABASE".custom.in_progress; then
				echo "[!!ERROR!!] Failed to produce custom backup database $DATABASE"
			else
				mv $FINAL_BACKUP_DIR"$DATABASE".custom.in_progress $FINAL_BACKUP_DIR"$DATABASE".custom
			fi
		fi
	done

	echo -e "\nAll database backups complete!"
}

# MONTHLY BACKUPS
DAY_OF_MONTH=`date +%d`

if [ $DAY_OF_MONTH -eq 1 ]; then
	# Delete all expired monthly directories
	find $BACKUP_DIR -maxdepth 1 -name "*-monthly" -exec rm -rf '{}' ';'

	perform_backups "-monthly"
fi

# WEEKLY BACKUPS
DAY_OF_WEEK=`date +%u` #1-7 (Monday-Sunday)
EXPIRED_DAYS=`expr $((($WEEKS_TO_KEEP * 7) + 1))`

if [ $DAY_OF_WEEK = $DAY_OF_WEEK_TO_KEEP ]; then
	# Delete all expired weekly directories
	find $BACKUP_DIR -maxdepth 1 -mtime +$EXPIRED_DAYS -name "*-weekly" -exec rm -rf '{}' ';'

	perform_backups "-weekly"
fi

# DAILY BACKUPS
# Delete daily backups 7 days old or more
find $BACKUP_DIR -maxdepth 1 -mtime +$DAYS_TO_KEEP -name "*-daily" -exec rm -rf '{}' ';'

perform_backups "-daily"
