package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// curl -i -X GET 'http://127.0.0.1:8000/calc' -H "User-Access: SuperUser" -d '2+5+7+8'

const userAccess = "superuser"

type response struct {
	sum int `json:"sum"`
}

func main() {
	http.HandleFunc("/calc", calculator)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// Обработчик возвращающий сумму
func calculator(w http.ResponseWriter, r *http.Request) {
	usr := r.Header.Get("User-Access")
	if !strings.EqualFold(usr, userAccess) {
		errMsg := "Incorrect User-Access: " + usr
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}

	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Read body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numbers := strings.Split(string(bytesBody), "+")

	sum := 0
	for _, strNum := range numbers {
		num, err := strconv.Atoi(strNum)
		if err != nil {
			log.Println("Can't convert: %v", err)
			break
		}
		sum += num
	}
	log.Printf("sum: %d\n", sum)

	jsonSum, err := json.Marshal(response{sum: sum})
	if err != nil {
		log.Println("Marshaling response to json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(string(jsonSum))

	w.Header().Set("Content-Type", "application/json")
	wByte, err := w.Write(jsonSum)
	if err != nil {
		log.Println("Write sum to response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Bytes written: %d", wByte)
}
