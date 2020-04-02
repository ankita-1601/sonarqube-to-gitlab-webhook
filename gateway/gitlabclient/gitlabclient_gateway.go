package gitlabclient

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/appcontext"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/utils"
)

// const (
// 	// notDefined string constante
// 	notDefined    string = "NotDefined"
// 	emptyResponse string = "emptyResponse"
// )

// Repository struct
type Repository struct {
	Client *http.Client
}

// GitlabGetProjectID func
// func (repo Repository) GitlabGetProjectID(token string, url string, projectPathWithNamespace string) (string, string, error) {
// 	req, err := http.NewRequest("GET", url, nil)
// 	printError(true, err)
// 	req.Header.Add("Private-Token", config.GitlabToken)
// 	resp, err := repo.Client.Do(req)
// 	printError(true, err)
// 	bodyText, err := ioutil.ReadAll(resp.Body)
// 	printError(false, err)
// 	var result []map[string]interface{}
// 	err = json.Unmarshal(bodyText, &result)
// 	printError(false, err)
// 	var res string
// 	if len(result) > 0 {
// 		number := fmt.Sprintf("%d", len(result))
// 		go log.Printf("[INFO] Number of Projects Found: %s", number)
// 		// we use projectPathWithNamespace to match entry
// 		if projectPathWithNamespace != notDefined {
// 			for _, n := range result {
// 				if projectPathWithNamespace == n["path_with_namespace"] {
// 					res = fmt.Sprint(n["id"])
// 				}
// 			}
// 		} else {
// 			// Using the first entry in array. Maybe it will not work if return more than one project
// 			s := result[0]
// 			res = fmt.Sprint(s["id"])
// 		}
// 	} else {
// 		res = emptyResponse
// 		err := errors.New("Empty Response")
// 		return "204", res, err
// 	}

// 	defer resp.Body.Close()
// 	return resp.Status, res, nil
// }

// GetGitlab func
func (repo Repository) GetGitlab(gitlabURL, gitlabToken string) ([]byte, string, error) {
	req, err := http.NewRequest(http.MethodGet, gitlabURL, nil)
	if err != nil {
		return []byte{}, "", err
	}
	req.Header.Add("PRIVATE-TOKEN", gitlabToken)
	resp, err := repo.Client.Do(req)
	if err != nil {
		return []byte{}, "", err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, "", err
	}
	if resp.StatusCode <= 204 {
		var links string
		// fmt.Println(resp.Header.Get("Links"))
		if resp.Header.Get("Links") != "" {
			links = utils.CleanNextLinksHeader(resp.Header.Get("Links"))
		}
		// if config.Debug {
		// 	fmt.Printf("Response: %s \n", resp.Status)
		// 	fmt.Printf("Next Links: %s \n", links)
		// }

		return bodyText, links, nil
	}
	// if config.Debug {
	// 	fmt.Printf("[ERROR] Response: %s \n", resp.Status)
	// }
	if resp.StatusCode == 404 {
		return []byte(`[]`), "", nil
	}
	defer resp.Body.Close()
	return []byte{}, "", nil
}

// GitlabPostComment func
func (repo Repository) GitlabPostComment(url string, params map[string]string) (err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		log.Printf("[ERROR] Cannot close writer: %s", url)
		return err
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Printf("[ERROR] Cannot create a post: %s", url)
		return err
	}
	req.Header.Add("Private-Token", config.GitlabToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := repo.Client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Cannot make a post: %s", url)
		return err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Cannot read response: %s", url)
		return err
	}
	s := string(bodyText)
	log.Printf("[INFO] Gitlab Commit: %s, %s", resp.Status, s)
	defer resp.Body.Close()
	return nil
}

func init() {
	if config.GetEnv("TESTRUN", "false") == "true" {
		return
	}
	client := http.Client{
		Timeout: time.Second * 10,
	}
	appcontext.Current.Add(appcontext.Repository, Repository{Client: &client})
	if appcontext.Current.Count() != 0 {
		log.Println("[INFO] Repository initiated")
	}
}
