package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
)

var db = make(map[string]string)

func WriteInDB(p string) {
	file, err := os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Попытка чтения: %v\n", err)
	}
	err1 := json.Unmarshal(data, &db)
	if err != nil {
		log.Fatal(err1)
	}
}

func longToShort(long string) string {
	fmt.Println(db)
	b := make([]byte, 4)
	rand.Read(b)
	s := fmt.Sprintf("%x", b)
	short := s
	lg := string(long)
	db[short] = lg
	filePath := os.Getenv("FILE_STORAGE_PATH")
	if filePath != "" {
		filename := filePath
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		data, err := json.Marshal(&db)
		if err != nil {
			fmt.Println(err)
		}
		// добавляем перенос строки
		data = append(data, '\n')
		_, err = file.Write(data)
		if err != nil {
			fmt.Println("1")
		}
	}
	return short
}

func CreateShort(w http.ResponseWriter, r *http.Request) {
	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = "localhost:8080"
	}
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
		w.Write([]byte("http://" + addr + "/" + short))
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
func APIShorten(w http.ResponseWriter, r *http.Request) {
	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = "localhost:8080"
	}
	defer r.Body.Close()
	long, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Internal Server Error %v", err)
	}
	//if r.Header.Get("Content-type") != "application/json" {
	//	fmt.Println("Now only json!")
	//	return
	//}
	in := map[string]string{}
	if err := json.Unmarshal(long, &in); err != nil {
		panic(err)
	}
	for _, ur := range in {
		u, err := url.ParseRequestURI(ur)
		if err != nil || u == nil {
			w.WriteHeader(400)
			w.Write([]byte("Не корректный URL"))
		} else {
			short := longToShort(ur)
			out := map[string]string{"result": "http://" + addr + "/" + short}
			s, err := json.Marshal(out)
			if err != nil {
				panic(err)
			}
			w.Header().Set("content-type", "application/json")
			w.Header().Add("Accept", "application/json")
			w.WriteHeader(201)
			w.Write([]byte(string(s)))
		}
	}

}
