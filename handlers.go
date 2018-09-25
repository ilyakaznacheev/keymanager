package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	keyIndex        = "keyman:index"
	keyKey          = "keyman:key"
	statusValud     = "valid"
	statusInvalid   = "invalid"
	statusNotExists = "not exsts"
)

type keyNumber struct {
	Key string `json:"key"`
}

type keyStatus struct {
	Key    string `json:"key"`
	Status string `json:"status"`
}

type keyInfo struct {
	Remains string `json:"remains"`
}

// NewKey generates new key
func NewKey(w http.ResponseWriter, r *http.Request) {
	keyChan := make(chan string, 1)

	log.Println("new key requested")

	var index string
	var intVal int

	action := func(client *redisClient) {
		// get last created index
		val, err := client.client.Get(keyIndex).Result()
		if err != nil {
			// not found means this is a initial request - create initial index
			index = "0"
			intVal = 0
		} else {
			// increment index and set it back
			intVal, err = strconv.Atoi(val)
			if err != nil {
				log.Fatal("error during data conversion:", val)
			}

			intVal++
			index = strconv.Itoa(intVal)
		}

		// set key status
		err = client.client.Set(keyIndex, index, 0).Err()
		if err != nil {
			log.Println("error during data sending:", val)
		}

		nextKey := getNextKey(intVal)

		err = client.client.Set(keyKey+nextKey, strconv.FormatBool(true), 0).Err()
		if err != nil {
			log.Println("error during new key sending")
		}

		keyChan <- nextKey
	}
	redClient.input <- redisAction{action}

	nextKey := <-keyChan
	close(keyChan)

	keyResult := keyNumber{nextKey}
	answer, _ := json.Marshal(keyResult)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(answer))
}

// CancelKey makes key invalid
func CancelKey(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	doneChan := make(chan error, 1)

	var status keyStatus

	index := requestParams["key"]

	log.Println("status change requested for key:", index)

	action := func(client *redisClient) {
		value, err := client.client.Get(keyKey + index).Result()
		if err == nil {
			validityStatus, _ := strconv.ParseBool(value)
			if !validityStatus {
				log.Println("key is already invalid")
				doneChan <- errors.New("key is already invalid")
				return
			}
		}

		err = client.client.Set(keyKey+index, strconv.FormatBool(false), 0).Err()
		if err != nil {
			log.Println("error during validity status change")
			doneChan <- err
			return
		}
		doneChan <- nil

	}
	redClient.input <- redisAction{action}
	err := <-doneChan
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	status = keyStatus{index, statusInvalid}

	answer, _ := json.Marshal(status)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(answer))
}

// CheckKey returns key validity status
func CheckKey(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	index := requestParams["key"]

	var status keyStatus

	log.Println("status requested for key:", index)

	value, err := redClient.client.Get(keyKey + index).Result()
	if err != nil {
		status = keyStatus{index, statusNotExists}
	} else {
		validityStatus, err := strconv.ParseBool(value)
		if err != nil {
			log.Println(err)
			status = keyStatus{index, statusNotExists}
		}

		if validityStatus {
			status = keyStatus{index, statusValud}
		} else {
			status = keyStatus{index, statusInvalid}
		}
	}

	answer, _ := json.Marshal(status)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(answer))
}

// Info returns remaining key number
func Info(w http.ResponseWriter, r *http.Request) {
	val, err := redClient.client.Get(keyIndex).Result()
	if err != nil {
		val = "0"
	}

	log.Println("info requested")

	valInt, _ := strconv.Atoi(val)

	remainKeyNumber := keyInfo{strconv.Itoa(getMaxKeyNumber() - valInt)}
	answer, _ := json.Marshal(remainKeyNumber)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(answer))
}
