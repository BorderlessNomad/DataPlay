#!/bin/bash

# This is application bootstrap script for App in Master (producer) node
# 1). Install Ubuntu base image or dataplay-ubuntu-base
# 2). Run go.sh (base server setup)
# 3). Run this script as 'sudo'
#
# [Note: running dataplay-go-master image will do 1) & 2) (recommended)]

APP="dataplay"
REPO="DataPlay"
SOURCE="https://github.com/playgenhub/$REPO/archive/"
BRANCH="develop"
DEST="/home/ubuntu/www"
START="start.sh"
LOG="ouput.log"

QUEUE_USERNAME="playgen"
QUEUE_PASSWORD="aDam3ntiUm"
QUEUE_SERVER="109.231.121.13"
QUEUE_PORT="5672"
QUEUE_ADDRESS="amqp://$QUEUE_USERNAME:$QUEUE_PASSWORD@$QUEUE_SERVER:$QUEUE_PORT/"
QUEUE_EXCHANGE="playgen-prod"

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
if [ -d $APP ]; then
	rm -r $APP
fi
mkdir $APP
echo "Moving files from $REPO-$BRANCH/ to $APP"
mv -f $REPO-$BRANCH/* $APP
cd $APP
chmod u+x $START
echo "Starting $START in Mode=$MODE"
nohup sh $START --mode=$MODE --uri="$QUEUE_ADDRESS" --exchange="$QUEUE_EXCHANGE" --requestqueue="$REQUEST_QUEUE" --requestkey="$REQUEST_KEY" --reqtag="$REQUEST_TAG" --responsequeue="$RESPONSE_QUEUE" --responsekey="$RESPONSE_KEY" --restag="$RESPONSE_TAG" &> $LOG &
echo "Done!"
echo "(Note: tail -f $DEST/$APP/$LOG for more details)"

exit 1
