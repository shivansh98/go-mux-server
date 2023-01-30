package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type quote struct {
	Id     int32  `json:"id"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Text   string `json:"text,omitempty"`
}

type request struct {
	Action string `json:"action"`
	quote
}

var store map[int32]quote

func init() {
	store = make(map[int32]quote, 1000)
	store[34] = quote{
		Id:     34,
		Title:  "Hey there",
		Author: "shivansh Tamrakar",
		Text:   "hheeheheheheheheh",
	}
}

func addQuotes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got a request for /addQuotes from :", r.RemoteAddr)
	req := request{}
	body, err := io.ReadAll(r.Body)
	checkErr(err)
	err = json.Unmarshal(body, &req)
	if err != nil {
		fmt.Print("Got an error in unmarshal ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.quote.Id != 0 {
		store[req.quote.Id] = req.quote
	}
	w.Write([]byte("Success"))
}

func getQuotes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got a request for /getQuotes from :", r.RemoteAddr)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(store)
	if err != nil {
		resp = []byte(fmt.Sprintf("{\"status\":\"got an err %v\"}", err))
		w.Write(resp)
	}
	w.Write(resp)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func deleteQuote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got a request for /deleteQuotes from :", r.RemoteAddr)
	req := request{}
	body, err := io.ReadAll(r.Body)
	checkErr(err)
	err = json.Unmarshal(body, &req)
	if err != nil {
		fmt.Print("Got an error in unmarshal ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.Id != 0 {
		if _, present := store[req.quote.Id]; present {
			delete(store, req.Id)
			w.Write([]byte("Successfully deleted"))
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}
	w.Write([]byte("Id not found"))
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/quotes", getQuotes)
	mux.HandleFunc("/write-quote", addQuotes)
	mux.HandleFunc("/delete-quote", deleteQuote)
	fmt.Println("The server is up and running at 8888 port on localhost")
	err := http.ListenAndServe("localhost:8888", mux)
	checkErr(err)

}
