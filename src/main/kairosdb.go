package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Datapoint struct {
	Name      string            `json:"name"`
	Timestamp int64             `json:"timestamp"`
	Value     int64             `json:"value"`
	Tags      map[string]string `json:"tags"`
}

type Kairosdb struct {
	client *http.Client
	host   string
}

func (kdb *Kairosdb) MsTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// AddDatapints add datapoints to configured kairosdb instance
func (kdb *Kairosdb) AddDatapoints(datapoints []Datapoint) error {

	json, err := json.Marshal(datapoints)
	if err != nil {
		return err
	}

	_, tpErr := kdb.TimedPost("/api/v1/datapoints", string(json))
	return tpErr
}

func (kdb *Kairosdb) TimedGet(path string) (int64, error) {
	req, err := http.NewRequest("GET", kdb.host+path, bytes.NewBuffer([]byte("")))
	if err != nil {
		return 0, err
	}
	return kdb.timedRequest(req)
}

func (kdb *Kairosdb) TimedPost(path string, body string) (int64, error) {

	req, err := http.NewRequest("POST", kdb.host+path, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return 0, err
	}

	//requests are always JSON
	req.Header.Set("Content-Type", "application/json")

	return kdb.timedRequest(req)
}

func (kdb *Kairosdb) timedRequest(req *http.Request) (int64, error) {

	start := kdb.MsTime()

	//create HTTP client

	resp, err := kdb.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" && resp.Status != "204 No Content" {
		return 0, errors.New("Response was non-200: " + resp.Status)
	}

	end := kdb.MsTime()

	return (end - start), nil
}
