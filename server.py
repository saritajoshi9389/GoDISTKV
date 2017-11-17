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
        # print("Received PUT request")
        return_code = 200
        return_message = None
        try:
            url = urllib.parse.urlparse(self.path)
            content_length = int(self.headers.get("Content-length"))
            content_type = self.headers.get("Content-type")
            data = self.rfile.read(content_length)
            message = simplejson.loads(data)
            if content_type != 'application/json':
                self.send_response(400)
                self.end_headers()
                return
            if self.path == "/set":
                return_message, return_code = self.process_valid_put_request(message)
            # print(return_message, return_code, "debug", simplejson.dumps(return_message).encode())
            self._set_headers(return_code)
            self.wfile.write(simplejson.dumps(return_message).encode())
        except Exception as err:
            print(err)
            self.send_response(500)
            self.end_headers()

    def process_valid_put_request(self, message):
        code = 200
        failed_inputs = []
        count = 0
        try:
            for kv in message:
                # print(kv["key"], kv["value"], self, self.server, self.server.kveachinstance)
                # print("haha1.....", frozenset(kv["key"].items()))
                temp_store = self.server.kveachinstance.get_value((frozenset(kv["key"].items())))
                flag = True
                if (temp_store):
                    flag = False
                result = self.server.kveachinstance.set_value(frozenset(kv["key"].items()),
                                                              frozenset(kv["value"].items()))
                # print("haha2")
                # print(result)
                if not result:
                    # print("haha3")
                    failed_inputs.append(kv["key"])
                else:
                    if flag:
                        count += 1
                        # print("result baby", self.server.kveachinstance.get_value(frozenset(kv["key"].items())))
        except TypeError:
            return {"error": True, "message": "Bad request"}, 400
        except KeyError:
            return {"error": True, "message": "Bad request"}, 400
        if count < len(message):
            code = 206
        return {"keys_failed": failed_inputs, "keys_added": count}, code

    __key_pattern = re.compile("^[a-zA-Z0-9]+$")

    def valid_key(self, key):
        return self.__key_pattern.match(key)

    def do_POST(self):
        """Respond to a POST request for storing a key"""
        print("Received POST request")
        return_code = 200
        return_message = None
        try:
            url = urllib.parse.urlparse(self.path)
            content_length = int(self.headers.get("Content-length"))
            content_type = self.headers.get("Content-type")
            data = self.rfile.read(content_length)
            message = simplejson.loads(data)
            # if content_type != 'application/json':
            #     self.send_response(400)
            #     self.end_headers()
            #     return
            if self.path == "/query":
                return_message, return_code = self.query_results(message)
            elif self.path == "/fetch":
                return_message, return_code = self.fetch_results(message)
            else:
                return_code = 501
                return_message = {"errors": [{"error": "invalid_api_key"}]}
            # print(return_message, return_code, "debug", simplejson.dumps(return_message).encode())
            self._set_headers(return_code)
            self.wfile.write(simplejson.dumps(return_message).encode())
        except Exception as err:
            print(err)
            self.send_response(500)
            self.end_headers()

    def fetch_results(self, message):
        code = 200
        try:
            for k in message:
                print(frozenset(k["key"].items()), "frozen")
                print(self.server.kveachinstance.get_value(frozenset(k["key"].items())))
                if self.server.kveachinstance.get_value(frozenset(k["key"].items())) is None:
                    result = [{
                        "key": k["key"],
                        "value": {}
                    }]
                else:
                    result = [
                        {
                            "key": k["key"],
                            "value": dict(self.server.kveachinstance.get_value(frozenset(k["key"].items())))
                        }
                    ]
                    # print("result baby", self.server.kveachinstance.get_value(frozenset(kv["key"].items())))
        except TypeError:
            return {"error": True, "message": "Bad request"}, 400
        except KeyError:
            return {"error": True, "message": "Bad request"}, 400
        return result, code

    def query_results(self, message):
        code = 200
        try:
            for k in message:
                print(frozenset(k["key"].items()), "frozen")
                print(self.server.kveachinstance.get_value(frozenset(k["key"].items())))
                if self.server.kveachinstance.get_value(frozenset(k["key"].items())) is None:
                    result = [{
                        "key": k["key"],
                        "value": False
                    }]
                else:
                    result = [
                        {
                            "key": k["key"],
                            "value": True
                        }
                    ]
                    # print("result baby", self.server.kveachinstance.get_value(frozenset(kv["key"].items())))
        except TypeError:
            return {"error": True, "message": "Bad request"}, 400
        except KeyError:
            return {"error": True, "message": "Bad request"}, 400
        return result, code

    def _set_headers(self, code=200):
        self.send_response(code)
        self.send_header('Content-type', 'application/json')
        self.end_headers()

    def do_GET(self):
        code = 404
        message = None
        if self.path == '/fetch':
            message, code = self.fetch_all()
        # send response
        self._set_headers(code)
        self.wfile.write(simplejson.dumps(message).encode())

    def fetch_all(self):
        code = 200
        output = self.server.kveachinstance.get_all()
        return output, code


class CustomHttpServer(http.server.HTTPServer):
    def __init__(self, server_address, RequestHandlerClass, data_dir="data", data_map=''):
        super(CustomHttpServer, self).__init__(server_address, RequestHandlerClass)
        self.data_dir = data_dir
        if not os.path.exists(self.data_dir): os.mkdir(self.data_dir)
        self.kveachinstance = data_map


class DataInstance:
    def __init__(self):
        self.data = {}

    def get_value(self, key):
        if key in self.data:
            return self.data[key]
        return None

    def get_all(self):
        var = [
            {
                "key": dict(key),
                "value": dict(self.data[key])
            } for key in self.data
            ]
        return var

    def set_value(self, key, value):
        self.data[key] = value
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
