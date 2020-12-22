package gitlabclient

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/appcontext"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/utils"
)

// Repository struct
type Repository struct {
	Client *http.Client
}

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
	logger := config.GetLogger
	defer logger().Sync()
	if resp.StatusCode <= 204 {
		var links string
		logger().Debugf("%s", resp.Header.Get("Links"))
		if resp.Header.Get("Links") != "" {
			links = utils.CleanNextLinksHeader(resp.Header.Get("Links"))
		}
		logger().Debug("Response" + resp.Status + " and Next Links " + links)
		return bodyText, links, nil
	}
	logger().Debugf("Response " + resp.Status)
	if resp.StatusCode == 404 {
		return []byte(`[]`), "", nil
	}
	defer resp.Body.Close()
	return []byte{}, "", nil
}

// GitlabPostComment func
func (repo Repository) GitlabPostComment(url string, params map[string]string) (err error) {
	logger := config.GetLogger
	defer logger().Sync()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		logger().Error("Cannot close writer")
		return err
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		logger().Errorf("Cannot create a post " + url)
		return err
	}
	req.Header.Add("Private-Token", config.Values.GitlabToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := repo.Client.Do(req)
	if err != nil {
		logger().Errorf("Cannot make a post " + url)
		return err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger().Errorf("Cannot read response %s", url)
		return err
	}
	s := string(bodyText)
	logger().Debugf("Gitlab Commit " + resp.Status + " , " + s)
	defer resp.Body.Close()
	return nil
}

func initclient() appcontext.Component {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	return Repository{Client: &client}
}

func init() {
	if config.Values.TestRun {
		return
	}
	appcontext.Current.Add(appcontext.Repository, initclient)
	if appcontext.Current.Count() != 0 {
		logger := config.GetLogger
		defer logger().Sync()
		logger().Debug("Repository initiated")
	}
}
