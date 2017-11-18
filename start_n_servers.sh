#!/bin/bash
value=9000

for number in {1..100}
do
temp=$(expr $value + $number)
echo $temp
# python3 server.py -p (value)
done
exit 0
pid=$(lsof -i:9000 -t); kill -TERM $pid || kill -KILL $pid
pid1=$(lsof -i:9001 -t); kill -TERM $pid1 || kill -KILL $pid1
pid2=$(lsof -i:9002 -t); kill -TERM $pid2 || kill -KILL $pid2
python3 server.py  &
python3 server.py -p 9001 &
python3 server.py -p 9002 &
# python3 server.py -p 9003 &
# python3 server.py -p 9004 &
# python3 server.py -p 9005 &
# python3 server.py -p 9006 &
# python3 server.py -p 9007 &
# python3 server.py -p 9008 &
# python3 server.py -p 9009 &
# python3 server.py -p 9010 &
# python3 server.py -p 9011 &
# python3 server.py -p 9012 &
# python3 server.py -p 9013 &
# python3 server.py -p 9014 &
# python3 server.py -p 9015 &
# python3 server.py -p 9016 &
# python3 server.py -p 9017 &
# python3 server.py -p 9018 &
# python3 server.py -p 9019 &
# python3 server.py -p 9020 &
# python3 server.py -p 9021 &
# python3 server.py -p 9022 &
# python3 server.py -p 9023 &
# python3 server.py -p 9024 &
# python3 server.py -p 9025 &
# python3 server.py -p 9026 &
# python3 server.py -p 9027 &
# python3 server.py -p 9028 &
# python3 server.py -p 9029 &
# python3 server.py -p 9030 &
# python3 server.py -p 9031 &
# python3 server.py -p 9032 &
# python3 server.py -p 9033 &
# python3 server.py -p 9034 &
# python3 server.py -p 9035 &
# python3 server.py -p 9036 &
# python3 server.py -p 9037 &
# python3 server.py -p 9038 &
# python3 server.py -p 9039 &
# python3 server.py -p 9040 &
# python3 server.py -p 9041 &
# python3 server.py -p 9042 &
# python3 server.py -p 9043 &
# python3 server.py -p 9044 &
# python3 server.py -p 9045 &
# python3 server.py -p 9046 &
# python3 server.py -p 9047 &
# python3 server.py -p 9048 &
# python3 server.py -p 9049 &
# python3 server.py -p 9050 &
# python3 server.py -p 9051 &
# python3 server.py -p 9052 &
# python3 server.py -p 9053 &
# python3 server.py -p 9054 &
# python3 server.py -p 9055 &
# python3 server.py -p 9056 &
# python3 server.py -p 9057 &
# python3 server.py -p 9058 &
# python3 server.py -p 9059 &
# python3 server.py -p 9060 &
# python3 server.py -p 9061 &
# python3 server.py -p 9062 &
# python3 server.py -p 9063 &
# python3 server.py -p 9064 &
# python3 server.py -p 9065