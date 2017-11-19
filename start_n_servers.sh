#!/bin/bash
value=9000

#for number in {1..99}
 for number in {1..3}
do
temp=$(expr $value + $number)
# echo $temp
python3 server.py -p $temp &
done
#python3 server.py -p 9100
exit 0
pid=$(lsof -i:9000 -t); kill -TERM $pid || kill -KILL $pid
pid1=$(lsof -i:9001 -t); kill -TERM $pid1 || kill -KILL $pid1
pid2=$(lsof -i:9002 -t); kill -TERM $pid2 || kill -KILL $pid2
