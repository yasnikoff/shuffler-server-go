#!/usr/bin/env bash
echo "building..."
Port=7474
go build

if [ $? -eq 0 ]
then
    echo "ok"
    echo "generating test data..."
    cd ./test_data
    ./create.sh
    cd ..
else
    echo "failed. aborting"
    exit
fi


if [ $? -eq 0 ]
then
    echo "ok"
    echo "listening on the port " $Port
    ./shuffler-server -p $Port  "repo 1" "./test_data/root/repo 1/" "repo 2" "./test_data/root/repo 2/"

else
    echo "failed. aborting..."
fi
echo "done"


