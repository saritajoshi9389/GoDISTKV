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
)
import (
	b64 "encoding/base64"
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
	FailedKeys       []string
	CountOfAddedKeys int
}

var url string
const SUCCESS int = 200
const PARTIAL_SUCCESS int = 206
const INTERNAL_SERVER_ERROR int = 500
const OTHER_ERROR int = 405

func handler(w http.ResponseWriter, r *http.Request, total_servers int, server_list []string) {
	fmt.Println("enter handler.............")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
	fmt.Println(r.URL.Path, "hi path")
	if (r.URL.Path == "/set") {
		set_handler(w, r, total_servers, server_list)

	} else if (r.URL.Path == "/fetch") {
		fetch_handler(w, r, total_servers, server_list)

	} else if (r.URL.Path == "/query") {
		query_handler(w, r, total_servers, server_list)

	}
}
func query_handler(w http.ResponseWriter, r *http.Request, total_servers int, server_list []string) {
	//else if r.Method == "PUT"{
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	//fmt.Fprintf(w, string(contents))
	//log.Println(url)
	fmt.Println("abc", contents)
	var d []MyData
	err1 := json.Unmarshal(contents, &d)
	if err1 != nil {
		fmt.Printf("hiiii%s", err1)
		os.Exit(1)
	}

	fmt.Println(d[1].Key, d[1].Value)
}
func fetch_handler(w http.ResponseWriter, r *http.Request, total_servers int, server_list []string) {
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	//fmt.Fprintf(w, string(contents))
	//log.Println(url)
	fmt.Println("abc", contents)
	var d []MyData
	err1 := json.Unmarshal(contents, &d)
	if err1 != nil {
		fmt.Printf("hiiii%s", err1)
		os.Exit(1)
	}

	fmt.Println(d[1].Key, d[1].Value)
}
func set_handler(w http.ResponseWriter, r *http.Request, total_servers int, server_list []string) {
	/////////////////////////
	if (r.URL.Path == "/set") {
		contents, _ := ioutil.ReadAll(r.Body)
		var d []MyData
		err1 := json.Unmarshal(contents, &d)
		if err1 != nil {
			fmt.Printf("hiiii%s", err1)
			os.Exit(1)
		}
		//fmt.Println("URL::",url)
		server_ele := 0
		struct_map := make(map[int][]MyData)
		for _, elem := range d {
			//fmt.Println(elem.Key,elem.Value.Data,server_list[server_ele])
			sEnc := b64.StdEncoding.EncodeToString([]byte(elem.Key.Data))
			val := elem.Key.Data[0]
			fmt.Println("ahhhh", sEnc, val)
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
			index := int(int(val) % total_servers) // changing from 3 to total_servers
			struct_map[index] = append(struct_map[index], temp_struct)
			server_ele ++
		}
		fmt.Println("No of requests ", server_ele)
		i := 0
		var wg sync.WaitGroup
		wg.Add(server_ele)
		respsChan := make(chan *http.Response)
		resps := make([]*http.Response, 0)
		for i < total_servers {
			if val, ok := struct_map[i]; ok {
				json_obj, _ := json.Marshal(val)
				fmt.Println(string(json_obj))
				fmt.Println("temp_struct", val, i)
				url = strings.Join([]string{"http://", "localhost:", string(server_list[i]), r.URL.Path}, "")
				response, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json_obj))
				if err != nil {
					os.Exit(2)
				} else {
					defer response.Body.Close()
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
					}(response)
				}

			}
			i++

		}
		go func() {
			for response := range respsChan {
				fmt.Println("new resp", response)
				resps = append(resps, response)
			}
		}()
		wg.Wait()
		send_message, r_code := format_response(resps)
		fmt.Println("hi result", string(send_message))
		success_handler(w, send_message, r_code)

	}
}
func format_response(responses []*http.Response) ([]byte, int) {
	failed_map := make([]string, 0)
	count_of_keys := 0
	code := SUCCESS
	for _, response := range responses {
		if response.StatusCode >= SUCCESS {
			body, error := ioutil.ReadAll(response.Body)
			fmt.Println("body", string(body))
			if error != nil {
				log.Fatal(error)
			}
			var back_response SetResponse
			json.Unmarshal(body, &back_response)
			count_of_keys += back_response.CountOfAddedKeys
			failed_map = append(failed_map, back_response.FailedKeys...)
		} else {
			code = PARTIAL_SUCCESS
		}
		response.Body.Close()
	}
	res := SetResponse{CountOfAddedKeys: count_of_keys, FailedKeys: failed_map}
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

func main() {
	arg := os.Args[1:]
	server_list := arg[1:]
	total_servers := len(server_list)
	fmt.Println("oyeeeeeeeeeeeeeeeeeeeeee", total_servers, server_list)
	if arg[0] != "-p" {
		fmt.Println("Incorrect flag variable, exiting....")
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, total_servers, server_list)
	})
	fmt.Println("Proxy up and running!!!")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}