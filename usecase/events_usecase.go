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

// GitlabCommit func
func GitlabCommit(projectName string, revision string, url string, status string) (string, error) {
	project := strings.Split(projectName, "/")
	urlget := fmt.Sprintf("%s/api/v4/projects?search=%s", config.GitlabURL, project[1])
	statusCode, projectID, err := gitlabclient.GitlabGetProjectID(config.GitlabToken, urlget)
	if err != nil {
		log.Printf("[ERROR]: %s, %s", statusCode, err)
	}
	posturl := fmt.Sprintf("%s/api/v4/projects/%s/repository/commits/%s/comments", config.GitlabURL, projectID, revision)
	form := fmt.Sprintf("SONARQUBE Status: %s , URL: %s", status, url)
	extraParams := map[string]string{
		"note": form,
	}
	go gitlabclient.GitlabPostComment(posturl, extraParams)

	return "OK", nil
}

// ValidateWebhook func to validate auth from Sonarqube
func ValidateWebhook(header string, message string) bool {
	mac := hmac.New(sha256.New, []byte(config.SonarqubeSecret))
	if _, err := mac.Write([]byte(message)); err != nil {
		log.Printf("mac.Write(%v) failed\n", message)
		return false
	}
	calculatedMAC := hex.EncodeToString(mac.Sum(nil))
	fmt.Printf("Headers: %s ; %s", calculatedMAC, header)
	return hmac.Equal([]byte(calculatedMAC), []byte(header))
}
