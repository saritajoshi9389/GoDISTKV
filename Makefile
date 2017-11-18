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
run:
	./hw5
check:clean build run
