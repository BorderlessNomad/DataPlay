#!/bin/bash

# This is application bootstrap script for App in Master (producer) node
# 1). Install Ubuntu base image or dataplay-ubuntu-base
# 2). Run go.sh (base server setup)
# 3). Run this script as 'sudo'
#
# [Note: running dataplay-go-master image will do 1) & 2) (recommended)]

APP="dataplay"
SOURCE="https://github.com/playgenhub/DataPlay/archive/"
BRANCH="develop"
DEST="/home/ubuntu/www"
DIR="dataplay"
START="start.sh"
LOG="ouput.log"
EXCHANGE="playgen-prod"
REQUEST_QUEUE="dataplay-request-prod"
REQUEST_KEY="api-request-prod"
REQUEST_TAG="consumer-request-prod"
RESPONSE_QUEUE="dataplay-response-prod"
RESPONSE_KEY="api-response-prod"
RESPONSE_TAG="consumer-response-prod"
MODE="2" # Master mode

# Kill any running process
echo "SHUTDOWN RUNING APP.."
killall -9 $APP

cd $DEST
echo "Fetching latest ZIP"
wget -Nq $SOURCE$BRANCH.zip -O $BRANCH.zip
echo "Extracting from $BRANCH.zip"
unzip -oq $BRANCH.zip
if [ -d $DIR ]; then
	rm -r $DIR
fi
mkdir $DIR
echo "Moving files from DataPlay-$BRANCH/ to $DIR"
mv -f DataPlay-$BRANCH/* $DIR
cd $DIR
chmod u+x $START
echo "Starting $START in Mode=$MODE"
nohup sh $START --mode=$MODE --exchange="$EXCHANGE" --requestqueue="$REQUEST_QUEUE" --requestkey="$REQUEST_KEY" --reqtag="$REQUEST_TAG" --responsequeue="$RESPONSE_QUEUE" --responsekey="$RESPONSE_KEY" --restag="$RESPONSE_TAG" &> $LOG &
echo "Done!"
echo "(Note: tail -f $DEST/$DIR/$LOG for more details)"

exit 1
