package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database struct {
	sync.RWMutex
	data map[string]dollars
}

func (db *database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/list":
		db.list(w, req)
	case "/price":
		db.price(w, req)
	case "/create":
		db.create(w, req)
	case "/update":
		db.update(w, req)
	case "/delete":
		db.delete(w, req)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such page: %s\n", req.URL)
	}
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	db.RLock()
	defer db.RUnlock()
	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	db.RLock()
	defer db.RUnlock()
	item := req.URL.Query().Get("item")
	price, ok := db.data[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid price: %s\n", priceStr)
		return
	}
	db.data[item] = dollars(price)
	fmt.Fprintf(w, "item %s created with price %s\n", item, dollars(price))
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid price: %s\n", priceStr)
		return
	}
	if _, ok := db.data[item]; !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	db.data[item] = dollars(price)
	fmt.Fprintf(w, "price of item %s updated to %s\n", item, dollars(price))
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	item := req.URL.Query().Get("item")
	if _, ok := db.data[item]; !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	delete(db.data, item)
	fmt.Fprintf(w, "item %s deleted\n", item)
}

func main() {
	db := &database{data: map[string]dollars{"shoes": 50, "socks": 5}}
	log.Fatal(http.ListenAndServe("localhost:8000", db))
}
