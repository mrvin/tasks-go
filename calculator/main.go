package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const userAccess = "superuser"

type response struct {
	Sum int `json:"sum"`
}

func main() {
	http.HandleFunc("/calc", calculator)

	//nolint:gosec
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// Обработчик возвращающий сумму чисел.
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
		log.Printf("Read body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numbers := splitByNum(string(bytesBody))
	log.Printf("numbers: %v\n", numbers)

	sum := 0
	for _, strNum := range numbers {
		num, err := strconv.Atoi(strNum)
		if err != nil {
			log.Printf("Can't convert: %v", err)
			break
		}
		sum += num
	}
	log.Printf("sum: %d\n", sum)

	jsonSum, err := json.Marshal(response{Sum: sum})
	if err != nil {
		log.Printf("Marshaling response to json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	wByte, err := w.Write(jsonSum)
	if err != nil {
		log.Printf("Write sum to response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Bytes written: %d", wByte)
}

func splitByNum(str string) []string {
	var result []string

	prevPos := 0
	for pos, ch := range str {
		if ch == '+' || ch == '-' {
			result = append(result, str[prevPos:pos])
			prevPos = pos
		}
		if pos == len([]rune(str))-1 {
			result = append(result, str[prevPos:pos+1])
		}
	}

	return result
}
