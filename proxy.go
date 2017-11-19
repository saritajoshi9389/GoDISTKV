package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"encoding/json"
	"bytes"
	"time"
	"sync"
)

type Key struct {
	Encoding string `json:"encoding"`
	Data     string `json:"data"`
}
type Value struct {
	Encoding string `json:"encoding"`
	Data     string `json:"data"`
}

type MyData struct {
	Key  `json:"key"`
	Value  `json:"value"`
}

type MyDatas struct {
	dataList MyData
}

type ErrorResponse struct {
	RCode    int
	RMessage string
}

type SetResponse struct {
	KeysFailed []Key `json:"keys_failed"`
	KeysAdded  int        `json:"keys_added"`
}

type MakeQueryRequest struct {
	Encoding string `json:"encoding"`
	Data     string `json:"data"`
}

type QueryResponse struct {
	Value bool `json:"value"`
	Key   struct {
		      Data     string `json:"data"`
		      Encoding string `json:"encoding"`
	      } `json:"key"`
}

type FetchResponse struct {
	Value struct {
		      Data     string `json:"data"`
		      Encoding string `json:"encoding"`
	      } `json:"value"`
	Key   struct {
		      Data     string `json:"data"`
		      Encoding string `json:"encoding"`
	      } `json:"key"`
}

var url string

const SUCCESS int = 200
const PARTIAL_SUCCESS int = 206
const INTERNAL_SERVER_ERROR int = 500
const OTHER_ERROR int = 501

func handler(w http.ResponseWriter, r *http.Request,
total_servers int, server_list []string,
ip_list []string, port_list []string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
	contents, _ := ioutil.ReadAll(r.Body)
	if (r.URL.Path == "/set") {
		set_handler(w, r, total_servers, ip_list, port_list, contents)

	} else if (r.URL.Path == "/fetch" && r.Method == "POST") {
		fetch_handler(w, r, total_servers, ip_list, port_list, contents)

	} else if (r.URL.Path == "/fetch" && r.Method == "GET") {
		fetch_handler_all(w, r, total_servers, ip_list, port_list)

	} else if (r.URL.Path == "/query") {
		query_handler(w, r, total_servers, ip_list, port_list, contents)

	} else {
		error_handler(w, &ErrorResponse{RCode: OTHER_ERROR, RMessage: "invalid_api_key"})

	}
}

func fetch_handler_all(w http.ResponseWriter, r *http.Request,
total_servers int, ip_list []string,
port_list []string) {
	var i = 0
	resps := make([]*http.Response, 0)
	respsChan := make(chan *http.Response)
	var wg sync.WaitGroup
	wg.Add(total_servers)
	for i < total_servers {
		url = strings.Join([]string{"http://", string(ip_list[i]), ":", string(port_list[i]), r.URL.Path}, "")
		response, err := http.NewRequest("GET", url, nil)
		if err != nil {
			os.Exit(2)
		} else {
			go func(response *http.Request) {
				defer wg.Done()
				response.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp_received, err := client.Do(response)
				if err != nil {
					panic(err)
				} else {
					respsChan <- resp_received
				}
				time.Sleep(time.Second * 2)
			}(response)
		}
		i++

	}
	go func() {
		for response := range respsChan {
			resps = append(resps, response)
		}
	}()
	wg.Wait()
	send_message, r_code := format_fetch_response(resps)
	success_handler(w, send_message, r_code)
}
func query_handler_jsonToObj(total_servers int,
contents []uint8) (map[int][]MakeQueryRequest, int) {
	var d []MakeQueryRequest
	err1 := json.Unmarshal(contents, &d)
	if err1 != nil {
		os.Exit(1)
	}
	server_ele := 0
	struct_map := make(map[int][]MakeQueryRequest)
	for _, elem := range d {
		temp_struct := MakeQueryRequest{
			Encoding:  elem.Encoding,
			Data: elem.Data,
		}
		index := hash_function(elem.Data) % total_servers
		struct_map[index] = append(struct_map[index], temp_struct)
		server_ele ++
	}
	return struct_map, server_ele
}
func query_handler(w http.ResponseWriter, r *http.Request,
total_servers int, ip_list []string,
port_list []string, contents []uint8) {
	struct_map, _ := query_handler_jsonToObj(total_servers, contents)
	i := 0
	var wg sync.WaitGroup
	wg.Add(len(struct_map))
	respsChan := make(chan *http.Response)
	resps := make([]*http.Response, 0)
	for i < total_servers {
		if val, ok := struct_map[i]; ok {
			json_obj, _ := json.Marshal(val)
			url = strings.Join([]string{"http://", string(ip_list[i]), ":", string(port_list[i]), r.URL.Path}, "")
			response, err := http.NewRequest("POST", url, bytes.NewBuffer(json_obj))
			if err != nil {
				os.Exit(2)
			} else {
				go func(response *http.Request) {
					defer response.Body.Close()
					defer wg.Done()
					response.Header.Set("Content-Type", "application/json")
					client := &http.Client{}
					resp_received, err := client.Do(response)
					if err != nil {
						panic(err)
					} else {
						respsChan <- resp_received
					}
					time.Sleep(time.Second * 2)
				}(response)
			}
		}
		i++

	}
	go func() {
		for response := range respsChan {
			resps = append(resps, response)
		}
	}()
	wg.Wait()
	send_message, r_code := format_query_response(resps)
	success_handler(w, send_message, r_code)
}

func format_query_response(responses []*http.Response) ([]byte, int) {
	output := make([]QueryResponse, 0)
	code := SUCCESS
	for _, response := range responses {
		if response.StatusCode >= SUCCESS {
			body, error := ioutil.ReadAll(response.Body)
			if error != nil {
				log.Fatal(error)
			}
			var back_response []QueryResponse
			json.Unmarshal(body, &back_response)
			output = append(output, back_response...)
		} else {
			code = PARTIAL_SUCCESS
		}
		response.Body.Close()
	}
	body, err := json.Marshal(output)
	if err != nil {
		return nil, INTERNAL_SERVER_ERROR
	}
	return body, code

}

func fetch_handler_jsonToObj(total_servers int,
contents []uint8) (map[int][]MakeQueryRequest, int) {
	var d []MakeQueryRequest
	err1 := json.Unmarshal(contents, &d)
	if err1 != nil {
		os.Exit(1)
	}
	server_ele := 0
	struct_map := make(map[int][]MakeQueryRequest)
	for _, elem := range d {
		temp_struct := MakeQueryRequest{
			Encoding:  elem.Encoding,
			Data: elem.Data,
		}
		index := hash_function(elem.Data) % total_servers
		struct_map[index] = append(struct_map[index], temp_struct)
		server_ele ++
	}
	return struct_map, server_ele
}

func fetch_handler(w http.ResponseWriter, r *http.Request,
total_servers int, ip_list []string,
port_list []string, contents []uint8) {
	struct_map, _ := fetch_handler_jsonToObj(total_servers, contents)
	i := 0
	var wg sync.WaitGroup
	wg.Add(len(struct_map))
	respsChan := make(chan *http.Response)
	resps := make([]*http.Response, 0)
	for i < total_servers {
		if val, ok := struct_map[i]; ok {
			json_obj, _ := json.Marshal(val)
			url = strings.Join([]string{"http://", string(ip_list[i]), ":", string(port_list[i]), r.URL.Path}, "")
			response, err := http.NewRequest("POST", url, bytes.NewBuffer(json_obj))
			if err != nil {
				os.Exit(2)
			} else {
				go func(response *http.Request) {
					defer response.Body.Close()
					defer wg.Done()
					response.Header.Set("Content-Type", "application/json")
					client := &http.Client{}
					resp_received, err := client.Do(response)
					if err != nil {
						panic(err)
					} else {
						respsChan <- resp_received
					}
					time.Sleep(time.Second * 2)
				}(response)
			}

		}
		i++

	}
	go func() {
		for response := range respsChan {
			resps = append(resps, response)
		}
	}()
	wg.Wait()
	send_message, r_code := format_fetch_response(resps)
	success_handler(w, send_message, r_code)
}

func format_fetch_response(responses []*http.Response) ([]byte, int) {
	output := make([]FetchResponse, 0)
	code := SUCCESS
	for _, response := range responses {
		if response.StatusCode >= SUCCESS {
			body, error := ioutil.ReadAll(response.Body)
			if error != nil {
				log.Fatal(error)
			}
			var back_response []FetchResponse
			json.Unmarshal(body, &back_response)
			output = append(output, back_response...)
		} else {
			code = PARTIAL_SUCCESS
		}
		response.Body.Close()
	}
	body, err := json.Marshal(output)
	if err != nil {
		return nil, INTERNAL_SERVER_ERROR
	}
	return body, code

}
func set_handler_jsonToObj(total_servers int,
contents []uint8) (map[int][]MyData, int) {
	var d []MyData
	err1 := json.Unmarshal(contents, &d)
	if err1 != nil {
		os.Exit(1)
	}
	server_ele := 0
	struct_map := make(map[int][]MyData)
	for _, elem := range d {
		temp_struct := MyData{
			Key: Key{
				Encoding:  elem.Key.Encoding,
				Data: elem.Key.Data,
			},
			Value: Value{
				Encoding: elem.Value.Encoding,
				Data: elem.Value.Data,
			},
		}
		index := hash_function(elem.Key.Data) % total_servers
		struct_map[index] = append(struct_map[index], temp_struct)
		server_ele ++
	}
	return struct_map, server_ele
}
func set_handler(w http.ResponseWriter, r *http.Request,
total_servers int, ip_list []string,
port_list []string, contents []uint8) {
	struct_map, _ := set_handler_jsonToObj(total_servers, contents)
	i := 0
	var wg sync.WaitGroup
	wg.Add(len(struct_map))
	respsChan := make(chan *http.Response)
	resps := make([]*http.Response, 0)
	for i < total_servers {
		if val, ok := struct_map[i]; ok {
			json_obj, _ := json.Marshal(val)
			url = strings.Join([]string{"http://", string(ip_list[i]), ":", string(port_list[i]), r.URL.Path}, "")
			response, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json_obj))
			if err != nil {
				os.Exit(2)
			} else {
				go func(response *http.Request) {
					defer response.Body.Close()
					defer wg.Done()
					response.Header.Set("Content-Type", "application/json")
					client := &http.Client{}
					resp_received, err := client.Do(response)
					if err != nil {
						panic(err)
					} else {
						respsChan <- resp_received
					}
					time.Sleep(time.Second * 2)
				}(response)
			}

		}
		i++
	}
	go func() {
		for response := range respsChan {
			resps = append(resps, response)
		}
	}()
	wg.Wait()
	send_message, r_code := format_set_response(resps)
	success_handler(w, send_message, r_code)
}
func format_set_response(responses []*http.Response) ([]byte, int) {
	failed_map := make([]Key, 0)
	count_of_keys := 0
	code := SUCCESS
	for _, response := range responses {
		if response.StatusCode >= SUCCESS {
			body, error := ioutil.ReadAll(response.Body)
			if error != nil {
				log.Fatal(error)
			}
			var back_response SetResponse
			json.Unmarshal(body, &back_response)
			count_of_keys += back_response.KeysAdded
			failed_map = append(failed_map, back_response.KeysFailed...)
		} else {
			code = PARTIAL_SUCCESS
		}
		response.Body.Close()
	}
	res := SetResponse{KeysAdded: count_of_keys, KeysFailed: failed_map}
	if len(failed_map) > 0 {
		code = PARTIAL_SUCCESS
	}
	body, err := json.Marshal(res)
	if err != nil {
		return nil, INTERNAL_SERVER_ERROR
	}
	return body, code
}

func success_handler(w http.ResponseWriter, reply []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(reply)
}

func error_handler(w http.ResponseWriter, e *ErrorResponse) {
	resp, error := json.Marshal(e)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.RCode)
	w.Write(resp)
}

func hash_function(str string) (int) {
	i := 0
	sum := 0
	for i = 0; i < len(str); i++ {
		sum = sum + int(str[i])
	}
	return sum
}

func distribute_servers(length int, server_list []string) ([]string, []string) {

	var ip_list = make([]string, length)
	var port_list = make([]string, length)
	for i := 0; i < length; i++ {
		ip_port := strings.Split(server_list[i], ":")
		ip_list[i] = ip_port[0]
		port_list[i] = ip_port[1]
	}
	return ip_list, port_list
}
func main() {
	arg := os.Args[1:]
	server_list := arg[1:]
	total_servers := len(server_list)
	ip_list, port_list := distribute_servers(total_servers, server_list)
	if arg[0] != "-p" {
		fmt.Println("Incorrect flag variable, exiting....")
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, total_servers, server_list, ip_list, port_list)
	})
	fmt.Println("Proxy up and running!!!")
	err := http.ListenAndServe("localhost:8100", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
