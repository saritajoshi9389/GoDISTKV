# This is the sample Makefile provided by professor.
# All rules here can be used as , make <rule-name>
# Example: make clean
# HW6
# Author: Sarita Joshi and Akshaya Khare
# CS 5600
CC=gcc -g -O3
GO=go build
all: 	clean

default: check
clean:
	rm -rf proxy
proxy:	proxy.go
	$(GO) proxy.go
dependencies: requirement.txt
		pip3 install -r requirement.txt
server:
	./start_n_servers.sh
run_proxy: proxy
	./run_proxy.sh &
run: stop run_proxy
	./start_n_servers.sh && sleep 2
	./client.sh && sleep 2 && printf "\nDone"
stop:
	./stop_proxy.sh &
	 ./stop_servers.sh
check:clean run
