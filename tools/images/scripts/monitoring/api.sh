#!/bin/bash

# This is application bootstrap script for API Monitoring probe
# 1). Install Ubuntu base image or dataplay-ubuntu-base
# 2). Run go.sh (base server setup)
# 3). Run this script as 'sudo'
#
# [Note: running dataplay-go-master image will do 1) & 2) (recommended)]

APP="dataplay-monitoring"
REPO="DataPlay-Monitoring"
SOURCE="https://github.com/playgenhub/$REPO/archive/"
BRANCH="master"
DEST="/home/ubuntu/www"
START="run.sh"
LOG="ouput.log"

# Kill any running process
echo "SHUTDOWN RUNING APP.."
killall -9 $APP

cd $DEST
echo "Fetching latest ZIP"
wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE$BRANCH.zip -O $BRANCH.zip
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
echo "Starting $START"
nohup sh $START &> $LOG &
echo "Done!"
echo "(Note: tail -f $DEST/$APP/$LOG for more details)"

exit 1
