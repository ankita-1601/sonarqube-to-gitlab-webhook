package gitlabclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
)

// GitlabGetProjectID func
func GitlabGetProjectID(token string, url string) (string, string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Private-Token", config.GitlabToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
	var result []map[string]interface{}
	json.Unmarshal(bodyText, &result)
	s := result[0]
	res := fmt.Sprint(s["id"])
	defer resp.Body.Close()
	return resp.Status, res, nil
}

// GitlabPostComment func
func GitlabPostComment(url string, params map[string]string) (string, string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		log.Printf("[ERROR]: %s", err)
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Private-Token", config.GitlabToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
	s := string(bodyText)
	go log.Printf("[INFO]: %s, %s", resp.Status, s)
	defer resp.Body.Close()
	return resp.Status, s, nil

}
