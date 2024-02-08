#!/bin/sh

##TARGET=/data/app-release.apk
##DESTINATION=/client/build/client.apk
##
##while [ ! -f $TARGET ]; do
##  echo "Waiting for $TARGET to be created..."
##  sleep 1
##done
##
##cp $TARGET $DESTINATION
##echo "Copied $TARGET to $DESTINATION"
serve -s build -l 8081