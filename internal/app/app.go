package app

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

var db = make(map[string]string)

func longToShort(long string) string {
	b := make([]byte, 4)
	rand.Read(b)
	s := fmt.Sprintf("%x", b)
	short := s
	lg := string(long)
	db[short] = lg
	return short
}

func CreateShort(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	long, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Internal Server Error %v", err)
	}
	u, err := url.ParseRequestURI(string(long))
	if err != nil || u == nil {
		w.WriteHeader(400)
		w.Write([]byte("Не корректный URL"))
	} else {
		short := longToShort(string(long))
		w.WriteHeader(201)
		w.Write([]byte("http://localhost:8080/" + short))
	}

}

func ReturnLong(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	long := db[id]
	if long == "" {
		w.Header().Set("Location", "Не верный идентификатор")
		w.WriteHeader(400)
	} else {
		w.Header().Set("Location", long)
		w.WriteHeader(307)
	}
}
