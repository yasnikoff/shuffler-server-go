#!/usr/bin/env bash
echo "testing..."
go test ../../...
#echo "skipping tests"

if [ $? -eq 0 ]
then
    echo "installing..."
    go install

else
    echo $test_result
    echo "tests failed"
fi
echo "done"
