from optparse import OptionParser
import http.server
import os
import os.path
import http
import http.client
import threading
import queue
import re
import datetime
import urllib.parse
import logging
import simplejson
import json

kv = {}

class CustomHandler(http.server.BaseHTTPRequestHandler):
    def do_PUT(self):
        """Respond to a PUT request for storing a key"""
        print("Received PUT request")
        try:
            url = urllib.parse.urlparse(self.path)
            print("url is ", url)
            content_length = int(self.headers.get("Content-length"))
            print("content-length", content_length)
            data = self.rfile.read(content_length)
            print(data, "data var")
            message = simplejson.loads(data)
            print("message value ", message)
            if(content_length != 0):
                successfully_updated = 1
            print("Data Dir", self.server)
            print("Data Dir", self.server.data_dir)
            print("kv ", self.server.kveachinstance)
            ########List to appending dict###########
            for kv in message:
                print(kv)



            #################################################
            # respond to client
            if successfully_updated == 1:
                self.send_response(201)
                self.send_header("Content-type", "application/json")
                self.end_headers()
            else:
                self.send_response(500, "Only {} nodes updated".format(successfully_updated))
                self.end_headers()

        except Exception as err:
            print(err)
            self.send_response(500)
            self.end_headers()

    __key_pattern = re.compile("^[a-zA-Z0-9]+$")

    def valid_key(self, key):
        return self.__key_pattern.match(key)

    def do_POST(self):
        """Respond to a POST request for storing a key"""
        print("Received POST request")


class CustomHttpServer(http.server.HTTPServer):

    def __init__(self, server_address, RequestHandlerClass, data_dir="data"):
        print("Enters this server func 1")
        super(CustomHttpServer, self).__init__(server_address, RequestHandlerClass)
        self.data_dir = data_dir
        if not os.path.exists(self.data_dir): os.mkdir(self.data_dir)
        self.kveachinstance = {}
        print("exit this server func 1")

if __name__ == '__main__':
    print("hi")
    parser = OptionParser()
    parser.add_option("-p", "--port", dest="port", help="local port to listen on", type="int", default=9000)
    parser.add_option("-d", "--data", dest="data_dir", help="data directory", type="string", default="data")
    (options, node_urls) = parser.parse_args()
    httpd = CustomHttpServer(("", options.port), CustomHandler, str(options.port)+options.data_dir)
    print(("Server Starts - {}:{}".format("localhost", options.port)))
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        pass
    httpd.server_close()
    folder = str(options.port)+options.data_dir
    with open(folder+'/result.json', 'w') as fp:
        json.dump(kv, fp)
    print(("Server Stops - {}:{}".format("localhost", options.port)))