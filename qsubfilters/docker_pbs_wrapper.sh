#!/bin/sh

FILE_UUID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
FILE_NAME="/tmp/$FILE_UUID.dockerrun"

while read i
do
echo $i >> $FILE_NAME
done

echo "docker run -v $FILE_NAME:/job.sh ubuntu:14.04 /bin/bash /job.sh"