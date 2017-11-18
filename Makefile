# This is the sample Makefile provided by professor.
# All rules here can be used as , make <rule-name>
# Example: make clean
# HW6
# Author: Sarita Joshi and Akshaya Khare
# CS 5600
CC=gcc -g -O3
all: 	clean

default: check
clean:
	rm -rf hw5  *.o *.dat
hw5.o:hw5.c
		${CC} -c -Wall -o hw5.o hw5.c -O3 -lm
hw5:hw5.o
		${CC} -g -o hw5 hw5.o -O3 -lm
build:hw5
calculate:
	./hw5
check:clean build calculate