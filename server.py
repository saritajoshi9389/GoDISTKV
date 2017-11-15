from optparse import OptionParser
import http.server
import os
import os.path
import http
import http.client
import re
import urllib.parse
import simplejson
import json


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
            if (content_length != 0):
                successfully_updated = 1
            print("Data Dir", self.server)
            print("Data Dir", self.server.data_dir)
            print("kv ", self.server.kveachinstance)
            code = 200
            failed_inputs = []
            count = 0
            ########List to appending dict###########
            try:
                for kv in message:
                    print(kv["key"], kv["value"], self, self.server, self.server.kveachinstance)
                    print("haha1.....", frozenset(kv["key"].items()))
                    result = self.server.kveachinstance.set_value(frozenset(kv["key"].items()),
                                                                  frozenset(kv["value"].items()))
                    print("haha2")
                    print(result)
                    if not result:
                        print("haha3")
                        failed_inputs.append(kv["key"])
                    else:
                        count += 1
                    print("result baby", self.server.kveachinstance.get_value(frozenset(kv["key"].items())))
            except TypeError:
                return {"error": True, "message": "Bad request"}, 400
            except KeyError:
                return {"error": True, "message": "Bad request"}, 400
            if count < len(message):
                code = 206
            return {"Number of keys added": count, "Number of keys failed": failed_inputs}, code

        #################################################
        # respond to client
        # if successfully_updated == 1:
        #     self.send_response(201)
        #     self.send_header("Content-type", "application/json")
        #     self.end_headers()
        # else:
        #     self.send_response(500, "Only {} nodes updated".format(successfully_updated))
        #     self.end_headers()

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
    def __init__(self, server_address, RequestHandlerClass, data_dir="data", data_map=''):
        print("Enters this server func 1")
        super(CustomHttpServer, self).__init__(server_address, RequestHandlerClass)
        self.data_dir = data_dir
        if not os.path.exists(self.data_dir): os.mkdir(self.data_dir)
        self.kveachinstance = data_map
        print("exit this server func 1")


class DataInstance:
    def __init__(self):
        self.data = {}

    def get_value(self, key):
        if key in self.data:
            return self.data[key]
        return None

    def set_value(self, key, value):
        print("enters set")
        print(key)
        print(value)
        print(self)
        self.data[key] = value
        print(self.data.items())
        return True

    def search(self, key):
        return key in self.data


if __name__ == '__main__':
    print("hi")
    parser = OptionParser()
    parser.add_option("-p", "--port", dest="port", help="local port to listen on", type="int", default=9000)
    parser.add_option("-d", "--data", dest="data_dir", help="data directory", type="string", default="data")
    (options, node_urls) = parser.parse_args()
    initial_data = DataInstance()
    httpd = CustomHttpServer(("", options.port), CustomHandler, str(options.port) + options.data_dir, initial_data)
    print(("Server Starts - {}:{}".format("localhost", options.port)))
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        pass
    httpd.server_close()
    folder = str(options.port) + options.data_dir
    with open(folder + '/result.json', 'a') as fp:
        print(initial_data.data)
        json.dump(str(initial_data.data), fp)
    print(("Server Stops - {}:{}".format("localhost", options.port)))
