package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateShort(t *testing.T) {
	tt := []struct {
		name   string
		data   string
		status int
		body   string
	}{
		{name: "first", data: "http://quququ.com", status: 201, body: "http://localhost:8080/52fdfc07"},
		{name: "second", data: "httpslkdfaklsdfj.com", status: 400, body: "Не корректный URL"},
		{name: "third", data: "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/29450/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/81e6f378-83b9-405d-962c-00a9de1f15c0/", status: 201, body: "http://localhost:8080/2182654f"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "localhost:8080/", bytes.NewBufferString(tc.data))
			if err != nil {
				t.Fatalf("Error %v", err)
			}
			rec := httptest.NewRecorder()
			CreateShort(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			byteBody, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
			}
			if string(byteBody) != tc.body {
				t.Errorf("Ожидалось %v, получено %v", tc.body, string(byteBody))
			}
			if res.StatusCode != tc.status {
				t.Errorf("Ожидалось %v, получено %v", tc.status, res.StatusCode)
			}
		})
	}
}

func TestReturnLong(t *testing.T) {
	tt := []struct {
		name     string
		id       string
		location string
		status   int
	}{
		{name: "first", id: "52fdfc07", location: "http://quququ.com", status: 307},
		{name: "second", id: "52fdfc0", location: "Не верный идентификатор", status: 400},
		{name: "third", id: "2182654f", location: "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/29450/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/81e6f378-83b9-405d-962c-00a9de1f15c0/", status: 307},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://localhost:8080/"+tc.id, nil)
			if err != nil {
				t.Fatalf("Error %v", err)
			}
			rec := httptest.NewRecorder()

			ReturnLong(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			if res.Header["Location"][0] != tc.location {
				t.Errorf("Ожидалось %v. Получено %v", tc.location, res.Header["Location"][0])
			}
			if res.StatusCode != tc.status {
				t.Errorf("Ожидалось %v. Получено %v", tc.status, res.StatusCode)
			}
		})
	}
}

func TestAPIShorten(t *testing.T) {
	tt := []struct {
		name   string
		data   string
		status int
		body   string
	}{
		{name: "first", data: `{"url":"http://quququ.com"}`, status: 201, body: `{"result":"http://localhost:8080/163f5f0f"}`},
		{name: "second", data: `{"url":"httpslkdfaklsdfj.com"}`, status: 400, body: "Не корректный URL"},
		{name: "third", data: `{"url":"https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/29450/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/81e6f378-83b9-405d-962c-00a9de1f15c0/"}`, status: 201, body: `{"result":"http://localhost:8080/9a621d72"}`},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "localhost:8080/api/shorten", bytes.NewBufferString(tc.data))
			if err != nil {
				t.Fatalf("Error %v", err)
			}
			rec := httptest.NewRecorder()
			APIShorten(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			byteBody, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
			}
			if string(byteBody) != tc.body {
				t.Errorf("Ожидалось %v, получено %v", tc.body, string(byteBody))
			}
			if res.StatusCode != tc.status {
				t.Errorf("Ожидалось %v, получено %v", tc.status, res.StatusCode)
			}
		})
	}
}
