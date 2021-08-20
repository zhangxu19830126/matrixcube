#!/bin/bash

for i in {1..100}
do
    make test > 1.log
    v=`tail -n 1 ./1.log | awk {'print $1'}`
    if [ "$v" != "ok" ]
    then
        echo "$i: other error"
        exit
    else
        echo "$i: ok"
    fi
done
