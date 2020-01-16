package gitlabclient

import (
	"bytes"
	"encoding/json"
	"errors"
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
	err = json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
	// fmt.Printf("INFO: %s", result)
	var res string
	if len(result) > 0 {
		s := result[0]
		res = fmt.Sprint(s["id"])
		go log.Printf("[INFO] ProjectID: %s, %s", resp.Status, res)
	} else {
		res = "emptyResponse"
		err := errors.New("Empty Response")
		return "", res, err
	}

	defer resp.Body.Close()
	return resp.Status, res, nil
}

// GitlabPostComment func
func GitlabPostComment(url string, params map[string]string) {
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
	go log.Printf("[INFO] Gitlab Commit: %s, %s", resp.Status, s)
	defer resp.Body.Close()

}
