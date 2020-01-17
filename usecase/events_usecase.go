package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/gateway/gitlabclient"
)

const (
	// notDefined string constante
	notDefined    string = "NotDefined"
	emptyResponse string = "emptyResponse"
)

// GitlabCommit func
func GitlabCommit(projectName string, revision string, url string, status string) {
	// test projectName
	project, projectPathWithNamespace := testProjectName(projectName)
	// project := strings.Split(projectName, "/")
	urlget := fmt.Sprintf("%s/api/v4/projects?search=%s", config.GitlabURL, project)
	// fmt.Printf("INFO: %s, %s", project, urlget)
	statusCode, projectID, err := gitlabclient.GitlabGetProjectID(config.GitlabToken, urlget, projectPathWithNamespace)
	if err != nil {
		log.Printf("[ERROR]: %s, %s", statusCode, err)
	}
	posturl := fmt.Sprintf("%s/api/v4/projects/%s/repository/commits/%s/comments", config.GitlabURL, projectID, revision)
	form := fmt.Sprintf("SONARQUBE Status: %s , URL: %s", status, url)
	extraParams := map[string]string{
		"note": form,
	}
	if projectID != emptyResponse {
		go gitlabclient.GitlabPostComment(posturl, extraParams)
	} else {
		log.Printf("[INFO] ProjectID not found for ProjectName: %s ", projectName)
	}

}

// ValidateWebhook func to validate auth from Sonarqube
func ValidateWebhook(header string, message string) bool {
	mac := hmac.New(sha256.New, []byte(config.SonarqubeSecret))
	if _, err := mac.Write([]byte(message)); err != nil {
		log.Printf("mac.Write(%v) failed\n", message)
		return false
	}
	calculatedMAC := hex.EncodeToString(mac.Sum(nil))
	// fmt.Printf("Headers: %s ; %s", calculatedMAC, header)
	return hmac.Equal([]byte(calculatedMAC), []byte(header))
}

func testProjectName(projectName string) (string, string) {
	var project string
	var projectPathWithNamespace string
	if strings.Contains(projectName, "/") {
		// fmt.Printf("string %s", name)
		// Expected: projectGroup/projectName
		// we use projectName as name to search in gitlab api
		projectSlice := strings.Split(projectName, "/")
		project = projectSlice[1]
		projectPathWithNamespace = projectName

	} else if strings.Contains(projectName, ":") {
		// Expected: com.example.local:namewithoutspaces:name-with-dashes-or-not
		// We use namewithoutspaces as project name to search in gitlab api
		projectSlice := strings.Split(projectName, ":")
		project = projectSlice[1]
		projectPathWithNamespace = notDefined
	} else {
		project = projectName
		projectPathWithNamespace = notDefined
	}
	return project, projectPathWithNamespace
}
