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

const (
	// notDefined string constante
	notDefined    string = "NotDefined"
	emptyResponse string = "emptyResponse"
)

// GitlabGetProjectID func
func GitlabGetProjectID(token string, url string, projectPathWithNamespace string) (string, string, error) {
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
		number := fmt.Sprintf("%d", len(result))
		go log.Printf("[INFO] Number of Projects Found: %s", number)
		// we use projectPathWithNamespace to match entry
		if projectPathWithNamespace != notDefined {
			for _, n := range result {
				if projectPathWithNamespace == n["path_with_namespace"] {
					res = fmt.Sprint(n["id"])
				}
			}
		} else {
			// Using the first entry in array. Maybe it will not work if return more than one project
			s := result[0]
			res = fmt.Sprint(s["id"])
		}
	} else {
		res = emptyResponse
		err := errors.New("Empty Response")
		return "204", res, err
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
