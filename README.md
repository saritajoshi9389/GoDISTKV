# Implement a Distributed Key-Value Store
## CS5600: HW6, Computer Systems, Fall 2017
### Author: Sarita Joshi and Akshaya Khare

# Goal
For this assignment, we have implemented an in-memory distributed key-value (KV) store. 
The KV store should be able to handle data larger than any one node's memory capacity.
That is, at any given time, a single node might not have all the data

### Implemented in Python3 and Golang

# Important deliverables in this submission

1) README    

        A README with all the descriptions about the design and implementation of this system along 
        with the steps to execute
        Makefile (as per the Sample Makefile provided by Professor)
2) Program files

        a) proxy.go  
            A proxy/coordinator process keeps track of available servers and data stored in those servers. 
            A client connects to the proxy/coordinator process to learn the address of a server that 
            it should connect for performing any operations.
            The proxy server also acts as a load-balancer and ensures a uniform workload distribution 
            among various servers.
            For this phase of the assignment, we have a hash-function that takes the given key
            (string/binary), calculates the server number by a simple mod operation
            In future, this can be replaced by consistency hashing technique
                
            Command to run: go run proxy.go <ip:port> <ip:port> <ip:port>
            Here, ip and port (input parameters) belong to the server ip and server port that is already up 
            and running
                
        b) server.py
            A server program that accepts get/set requests from the clients and returns a valid response. 
            (Future work: The server will communicate with it's peer processes (spread across the network) 
            to maintain a consistent view of the key-value database.)
                
            Command to run: python3 server.py -p <port>
            Here, port (input parameters) belong to the server port
            We have also provided a -d option that can be used to create a specific directory in which 
            all the data can be stored
            Such a technique, in future, will be useful to implement persistent storage and will allow 
            roll-back scenario
                
                
3) Bash Scripts

            To allow easy testing and automation, we have few scripts that can start n number of servers:
            start_n_servers.sh
            stop_servers.sh
            client.sh (For this phase of the assignment, client is simple curl commands 
            that can be run via terminal)
              
           
4) Makefile

            To install all the dependencies
               -    Proxy by default runs on 8080 (hardcoded in the proxy application)
               -    make dependencies :: This will go a `python pip3 install -r requirements.txt` (Not required for ccis machine)
               -    make run_proxy :: creates the proxy binary and runs the same
               -    make stop :: stops the server and the proxy
               -    make run ::  builds executes and runs everything (Testing for 3 servers on port 9001 9002 9003)
               -    make check :: cleans and runs the entire submission.
               -    We have individual scripts to start and stop servers and stop/start proxy
5) requirements.txt
            
            Stores all the python dependencies

NOTE: Client for this assignment is curl command, sample commands can be found in client.sh

# Run and execution steps

As per the problem statement, we support two formats for encoding i.e. "binary" or "string". 
A key or value therefore looks as below:
```json
{
    "encoding": "string",
    "data": "key1"
}
```

As per the piazza discussion (For details, refer the same), we support below REST endpoints:

    /set (Method: PUT) : For storing the key, value pair either insert or update
Input:
```json
[
    {
        "key": {
            "encoding": "string",
            "data": "key1"
        },
        "value": {
            "encoding": "string",
            "data": "data1"
        }
    }
]
```
Output (If all succeeded):
```json
[
    {
        "keys_added": 10,
        "keys_failed": [] 
    }
]
```

    /fetch (Method: POST):  To retrieve the key-value for a given key
Input:
```json
[
    {
      "encoding": "string",
      "data": "key1"
    }
]
```
Output:

```json
[
    {
        "key": {
            "encoding": "string",
            "data": "key1"
        },
        "value": {
            "encoding": "string",
            "data": "data1"
        }
    }
]
```

    /fetch (Method: GET):  To retrieve all the key-value i.e. a simple gellAll feature
     (Same liek above except that it returns all key-value pair)
    /query (Method: POST): To check if the given key is present, if yes, return a boolean true else false
 Input:
```json
[
    {
      "encoding": "string",
      "data": "key1"
    }
]
```   
  
Output:

```json
[
    {
        "key": {
            "encoding": "string",
            "data": "key1"
        },
        "value": true
    }
]
```



# Technique and Assumption
    For this system to work, Golang and Python3 must be installed. Follow, the requirement.txt 
    and pre-requisite mentioned in the instructions.
    This is a simple distributed key-value storage system, implemented using BaseHTTPRequestHandler
    in Python (server side)
    This design, we do not have a communication channel or API to communicate among all the server instances.
    The proxy does the task of sending information to onr of the servers based on the simple hashing technique
    To build the start for the second phase that will allow persistent storage, we have added a simple 
    technique to dump all the  server specific data to a local file as soon as the server dies on keyboard 
    interrupt.
    This can be extended further for signal handling and load the data on frequent intervals when the server
    is up and running again
    The proxy implements all the above mentioned APIs using net/http import in Golang, adds the 
    suitable headers and pass on the request to the server
    The APIs to communicate between client-proxy and proxy-server follow the same contract


# Test
    We have tested our system for multiple values of n (number of servers, max n = 100) and this 
    system works fine for MacOS, ubuntu, CCIS linux box
    You can modify the script to change the for loop until 100 and show work

# Future Scope
    
    Implementing Consistent Hashing technique
    Fault tolerance at server level
    Single point of failure
    Load balancing support
    

