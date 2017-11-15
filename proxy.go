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
import b64 "encoding/base64"

type Key struct{
    Encoding string `json:"encoding"`
    Data     string `json:"data"`
}
type Value struct{
    Encoding string `json:"encoding"`
    Data     string `json:"data"`
}


type MyData struct {
       Key  `json:"key"`
       Value  `json:"value"`
}

type MyDatas struct{
  dataList MyData
  }

var url string

//func handler(w http.ResponseWriter, r *http.Request) {
func handler(w http.ResponseWriter, r *http.Request, total_servers int, server_list []string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
	if(r.URL.Path == "/SET") {
		client := &http.Client{}
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
			  fmt.Println("ahhhh",sEnc)
			  temp_struct:= MyData{
						    Key: Key{
							    Encoding:  elem.Key.Encoding,
							    Data: elem.Key.Data,
						    },
						    Value: Value{
							    Encoding: elem.Value.Encoding,
							    Data: elem.Value.Data,
						    },
						}
			index := int(val%3)
			struct_map[index] = append(struct_map[index],temp_struct)
			  server_ele ++
              }
		i :=0
		for i < total_servers {
			if val , ok := struct_map[i]; ok {
				json_obj, _ := json.Marshal(val)
				fmt.Println(string(json_obj))
                  		fmt.Println("temp_struct",val, i)
				url = strings.Join([]string{"http://", "localhost:", string(server_list[i]), r.URL.Path}, "")
				response, err := http.NewRequest(http.MethodPut,url ,bytes.NewBuffer(json_obj))
				if err != nil {
					os.Exit(2)
				} else {
					defer response.Body.Close()
					_, err = client.Do(response)
					cts, err := ioutil.ReadAll(response.Body)
					if err != nil {
						fmt.Printf("%s", err)
						os.Exit(1)
					}
					fmt.Fprintf(w, string(cts))
			}

			}
			i++

		}


	}else if (r.URL.Path == "/FETCH"){
	//else if r.Method == "PUT"{
		contents, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			//fmt.Fprintf(w, string(contents))
			//log.Println(url)
		fmt.Println("abc",contents)
		var d []MyData
		err1 := json.Unmarshal(contents, &d)
		if err1 != nil {
				fmt.Printf("hiiii%s", err1)
				os.Exit(1)
			}

		fmt.Println(d[1].Key, d[1].Value)

	}else if (r.URL.Path == "/QUERY"){
	//else if r.Method == "PUT"{
		contents, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			//fmt.Fprintf(w, string(contents))
			//log.Println(url)
		fmt.Println("abc",contents)
		var d []MyData
		err1 := json.Unmarshal(contents, &d)
		if err1 != nil {
				fmt.Printf("hiiii%s", err1)
				os.Exit(1)
			}

		fmt.Println(d[1].Key, d[1].Value)

	}
}

func main() {
	arg := os.Args[1:]
        server_list := arg[1:]

        total_servers := len(server_list)
        fmt.Println("oyeeeeeeeeeeeeeeeeeeeeee",total_servers,server_list)
        if arg[0] != "-p"{
          fmt.Println("Incorrect flag variable, exiting....")
          return
        }
	//http.HandleFunc("/", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
              handler(w, r, total_servers, server_list)
       })
	fmt.Println("Proxy up and running!!!")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}