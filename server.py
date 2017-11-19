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
    # def check_dict(dict)
    def do_PUT(self):
        """Respond to a PUT request for storing a key"""
        print("Received PUT request")
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
                return_message, return_code = self.process_valid_put_request(
                    message)
            else:
                return_code = 501
                return_message = {"errors": [{"error": "invalid_api_key"}]}
            self._set_headers(return_code)
            self.wfile.write(simplejson.dumps(return_message).encode())
        except Exception as err:
            print(err)
            self.send_response(500)
            self.end_headers()

    def check_binary(self, str):
        for c in str:
            if c not in ('1', '0'):
                return False
        return True

    def verify_dict(self, kv):
        correct_json = True
        if not (kv["key"]["encoding"] == kv["value"]["encoding"]):
            correct_json = False
        if (kv["key"]["encoding"] == "binary"):
            if self.check_binary(kv["key"]["data"]):
                if not self.check_binary(kv["value"]["data"]):
                    correct_json = False
            else:
                correct_json = False

        if not(kv["key"]["encoding"] in ("binary", "string") and
               kv["value"]["encoding"] in ("binary", "string")):
            correct_json = False
        return correct_json

    def process_valid_put_request(self, message):
        code = 200
        failed_inputs = []
        count = 0
        try:
            for kv in message:
                temp_store = self.server.kveachinstance.get_value(
                    (frozenset(kv["key"].items())))
                flag = True
                if temp_store:
                    print("temp_store is -> ", (dict(temp_store)
                                                ["data"]), kv["value"]["data"])
                    if (dict(temp_store)["data"]) == kv["value"]["data"]:
                        flag = False

                correct_json = self.verify_dict(kv)
                if correct_json:
                    result = self.server.kveachinstance.set_value(frozenset(kv["key"].items()),
                                                                  frozenset(kv["value"].items()))
                else:
                    result = False
                if not result:
                    failed_inputs.append(kv["key"])
                else:
                    if flag:
                        count += 1
        except TypeError:
            return {"error": True, "message": "Bad request"}, 400
        except KeyError:
            return {"error": True, "message": "Bad request"}, 400
        if count < len(message):
            code = 206
        return {"keys_failed": failed_inputs, "keys_added": count}, code

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
            if content_type != 'application/json':
                self.send_response(400)
                self.end_headers()
                return
            if self.path == "/query":
                return_message, return_code = self.query_results(message)
            elif self.path == "/fetch":
                return_message, return_code = self.fetch_results(message)
            else:
                return_code = 501
                return_message = {"errors": [{"error": "invalid_api_key"}]}
            self._set_headers(return_code)
            self.wfile.write(simplejson.dumps(return_message).encode())
        except Exception as err:
            print(err)
            self.send_response(500)
            self.end_headers()

    def verify_dict_get(self, k):
        correct_json = True
        if (k["encoding"] == "binary"):
            if not self.check_binary(k["data"]):
                correct_json = False

        if not(k["encoding"] in ("binary", "string")):
            correct_json = False
        return correct_json

    def fetch_results(self, message):
        code = 200
        try:
            result = []
            for k in message:
                if k["encoding"] in ("binary", "string"):
                    if self.server.kveachinstance.get_value(frozenset(k.items())) is None:
                        result = [{
                            "key": k,
                            "value": {}
                        }]
                    else:
                        result.append(
                            {"key": k,
                             "value": dict(self.server.kveachinstance.get_value(frozenset(k.items())))
                             })
                else:
                    result = [{None}]
        except TypeError:
            return {"error": True, "message": "Bad request"}, 400
        except KeyError:
            return {"error": True, "message": "Bad request"}, 400
        return result, code

    def query_results(self, message):
        code = 200
        try:
            result = []
            for k in message:
                correct_json = self.verify_dict_get(k)
                if correct_json:
                    if self.server.kveachinstance.get_value(frozenset(k.items())) is None:
                        result = [{
                            "key": k,
                            "value": False
                        }]
                    else:
                        result.append({"key": k, "value": True})
                else:
                    result = [{None}]
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
        print("Received GET request")
        return_code = 404
        return_message = None
        if self.path == '/fetch':
            return_message, return_code = self.fetch_all()
        else:
            return_code = 501
            return_message = {"invalid_api_key"}
        # send response
        self._set_headers(return_code)
        self.wfile.write(simplejson.dumps(return_message).encode())

    def fetch_all(self):
        code = 200
        output = self.server.kveachinstance.get_all()
        return output, code


class CustomHttpServer(http.server.HTTPServer):
    def __init__(self, server_address, RequestHandlerClass, data_dir="data", data_map=''):
        super(CustomHttpServer, self).__init__(
            server_address, RequestHandlerClass)
        self.data_dir = data_dir
        # if not os.path.exists(self.data_dir):
        #     os.mkdir(self.data_dir)
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
    parser = OptionParser()
    parser.add_option("-p", "--port", dest="port",
                      help="local port to listen on", type="int", default=9000)
    parser.add_option("-d", "--data", dest="data_dir",
                      help="data directory", type="string", default="data")
    (options, node_urls) = parser.parse_args()
    initial_data = DataInstance()
    httpd = CustomHttpServer(("", options.port), CustomHandler, str(
        options.port) + options.data_dir, initial_data)
    print(("Server Starts - {}:{}".format("localhost", options.port)))
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        pass
    httpd.server_close()
    folder = str(options.port) + options.data_dir
    with open(folder + '/result.json', 'a') as fp:
        json.dump(str(initial_data.data), fp)
    print(("Server Stops - {}:{}".format("localhost", options.port)))
