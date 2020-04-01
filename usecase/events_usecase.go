package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/domain"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/gateway/gitlabclient"
)

const (
	// notDefined string constante
	notDefined    string = "NotDefined"
	emptyResponse string = "emptyResponse"
)

// GitlabCommit func
func GitlabCommit(event *domain.Events) {
	// func GitlabCommit(projectName string, revision string, url string, status string) {
	//event.Revision, event.Project.URL, event.Status
	// test projectName
	project, projectPathWithNamespace := testProjectName(event.Project.Name)
	// project := strings.Split(projectName, "/")
	urlget := fmt.Sprintf("%s/api/v4/projects?search=%s", config.GitlabURL, project)
	// fmt.Printf("INFO: %s, %s", project, urlget)
	statusCode, projectID, err := gitlabclient.GitlabGetProjectID(config.GitlabToken, urlget, projectPathWithNamespace)
	if err != nil {
		log.Printf("[ERROR]: %s, %s", statusCode, err)
	}
	posturl := fmt.Sprintf("%s/api/v4/projects/%s/repository/commits/%s/comments", config.GitlabURL, projectID, event.Revision)

	form, err := parseForm(event)
	if err != nil {
		form = fmt.Sprintf("SONARQUBE Status: %s , URL: %s", event.Status, event.Project.URL)
		log.Printf("[ERROR] Cannot parse event fields %s", event.Project.Name)
	}
	extraParams := map[string]string{
		"note": form,
	}

	if projectID != emptyResponse && !shouldPost(event) {
		go gitlabclient.GitlabPostComment(posturl, extraParams)
	} else {
		log.Printf("[INFO] ProjectID not found for Project Name: %s ", event.Project.Name)
	}

}

func parseForm(event *domain.Events) (string, error) {
	var form string
	var finalForm string
	form = fmt.Sprintf("SONARQUBE Report: URL: %s  ", event.Project.URL)
	if event.Branch.URL != "" {
		form = fmt.Sprintf("SONARQUBE Report: URL: %s  ", event.Branch.URL)
	}
	var qualityGateway string
	qualityGateway += fmt.Sprintf("Quality Gateway Name: %s  \n", event.QualityGate.Name)
	qualityGateway += fmt.Sprintf("Quality Gateway Status: %s  \n", event.QualityGate.Status)
	if !shouldSendQualityReport(event) {
		for _, v := range event.QualityGate.Conditions {
			qualityGateway += fmt.Sprintf("Metric Name: %s  \n", v.Metric)
			qualityGateway += fmt.Sprintf("Operator: %s  \n", v.Operator)
			if v.Value != "" {
				qualityGateway += fmt.Sprintf("Value: %s  \n", v.Value)
			}
			qualityGateway += fmt.Sprintf("Error Threshold: %s  \n", v.ErrorThreshold)
			qualityGateway += fmt.Sprintf("Status: %s  \n", v.Status)
		}
	}
	if qualityGateway != "" {
		finalForm = fmt.Sprintf("%s,  \n%s", form, qualityGateway)
	}

	// fmt.Println(finalForm)
	return finalForm, nil
}

func shouldPost(event *domain.Events) bool {
	var disabledGitlabPost bool
	disabledGitlabPost = false
	if len(event.Properties) != 0 {
		value := event.Properties["sonar.analysis.disabledGitlabPost"]
		if value == "true" {
			disabledGitlabPost = true
			// fmt.Println("should not post")
		}
	}
	return disabledGitlabPost
}

func shouldSendQualityReport(event *domain.Events) bool {
	var disabledQualityReport bool
	disabledQualityReport = false

	if len(event.Properties) != 0 {
		value := event.Properties["sonar.analysis.disabledQualityReport"]
		if value == "true" {
			disabledQualityReport = true
			// fmt.Println("should not post")
		}
	}
	return disabledQualityReport
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
