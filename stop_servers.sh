pid=$(lsof -i:9000 -t); kill -TERM $pid || kill -KILL $pid
pid1=$(lsof -i:9001 -t); kill -TERM $pid1 || kill -KILL $pid1
pid2=$(lsof -i:9002 -t); kill -TERM $pid2 || kill -KILL $pid2

