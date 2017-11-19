# pid=$(lsof -i:9000 -t); kill -TERM $pid || kill -KILL $pid
# pid1=$(lsof -i:9001 -t); kill -TERM $pid1 || kill -KILL $pid1
# pid2=$(lsof -i:9002 -t); kill -TERM $pid2 || kill -KILL $pid2

#!/bin/bash
value=9000

for number in {1..3}
do
temp=$(expr $value + $number)
# echo $temp
pid=$(lsof -i:$temp);kill -TERM $pid || kill -KILL $pid
# python3 server.py -p (value)
done
exit 0
