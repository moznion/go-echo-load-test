package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type MetadataResp struct {
	TaskARN string `json:"TaskARN"`
}

func main() {
	metadataURI := os.Getenv("ECS_CONTAINER_METADATA_URI")

	hc := http.Client{Timeout: 3 * time.Second}
	resp, err := hc.Get(metadataURI + "/task")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var r MetadataResp
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Fatal(err)
	}
	taskID := strings.Split(r.TaskARN, "/")[2]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(taskID))
	})
	log.Fatal(http.ListenAndServe(":9000", nil))
}
