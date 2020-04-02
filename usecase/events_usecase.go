package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/domain"
)

const (
	// notDefined string constante
	notDefined    string = "NotDefined"
	emptyResponse string = "emptyResponse"
)

// GitlabCommit func
func GitlabCommit(event *domain.Events) {
	useProjectID, projectWithID, err := shouldUseProjectID(event)
	if err != nil {
		log.Printf("[ERROR]: Cannot parse project ID %v", err)
	}
	// test if it already have projectID from sonarProperties
	if !useProjectID {
		project, projectPathWithNamespace := checkProjectName(event.Project.Name)
		parsedID, err := searchProjectID(project, projectPathWithNamespace)
		if err != nil {
			log.Printf("[ERROR] Cannot find projectID %s %s", project, projectPathWithNamespace)
		}
		projectWithID = parsedID
	}
	posturl := fmt.Sprintf("%s/api/v4/projects/%v/repository/commits/%s/comments", config.GitlabURL, projectWithID, event.Revision)
	form := parseForm(event)
	if err != nil {
		form = fmt.Sprintf("SONARQUBE Status: %s , URL: %s", event.Status, event.Project.URL)
		log.Printf("[ERROR] Cannot parse event fields %s", event.Project.Name)
	}
	extraParams := map[string]string{
		"note": form,
	}
	if projectWithID != 0 && !shouldPost(event) {
		go gitlabPostCommit(posturl, extraParams)
	} else {
		log.Printf("[INFO] ProjectID not found for Project Name: %s ", event.Project.Name)
	}
}

func gitlabPostCommit(posturl string, extraParams map[string]string) {
	if config.GitlabURL != "Absent" && config.GitlabToken != "Absent" {
		repo := domain.GetRepository()
		err := repo.GitlabPostComment(posturl, extraParams)
		if err != nil {
			log.Printf("[ERROR] Cannot post a gitlab commit")
		}
	}
}

func searchProjectID(project, projectPathWithNamespace string) (projectWithID int, err error) {
	if config.GitlabURL == "Absent" || config.GitlabToken == "Absent" {
		err = errors.New("Gitlab API not configured")
		return projectWithID, err
	}
	// urlget := fmt.Sprintf("%s/api/v4/projects?search=%s", config.GitlabURL, project)
	var result []domain.GitlabProject
	repo := domain.GetRepository()
	// statusCode, projectID, err := repo.GitlabGetProjectID(config.GitlabToken, urlget, projectPathWithNamespace)
	urlget := fmt.Sprintf("%s/api/v4/projects?pagination=keyset&per_page=100&order_by=id&sort=asc&search=%s", config.GitlabURL, project)
	bodyBytes, _, err := repo.GetGitlab(urlget, config.GitlabToken)
	if err != nil {
		// log.Printf("[ERROR]: %s, %s", statusCode, err)
		return projectWithID, err
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return projectWithID, err
	}
	projectsLength := len(result)
	if projectsLength == 0 {
		err = errors.New(emptyResponse)
		return projectWithID, err
	}
	if projectPathWithNamespace != notDefined {
		for _, n := range result {
			if projectPathWithNamespace == n.PathWithNamespace {
				projectWithID = n.ID
			}
		}
	} else {
		// Using the first entry in array. Maybe it will not work if return more than one project
		s := result[0]
		projectWithID = s.ID
	}
	return projectWithID, nil
}

func parseForm(event *domain.Events) string {
	var form string
	var finalForm string
	finalForm = "# SONARQUBE REPORT  \n"
	form = fmt.Sprintf("URL: [Project Link](%s)  \n", event.Project.URL)
	if event.Branch.URL != "" {
		form = fmt.Sprintf("URL: [Report Link](%s)  \n", event.Branch.URL)
	}
	var qualityGateway string
	qualityGateway += "## Quality Gateway  \n"
	qualityGateway += fmt.Sprintf(" - Name: %s  \n", event.QualityGate.Name)
	qualityGateway += fmt.Sprintf(" - Status: **%s**  \n", event.QualityGate.Status)
	if !shouldSendQualityReport(event) {
		qualityGateway += "### Quality Gateway Conditions  \n"
		for _, v := range event.QualityGate.Conditions {
			qualityGateway += fmt.Sprintf("#### Metric Name: %s  \n", v.Metric)
			qualityGateway += fmt.Sprintf(" - Operator: %s  \n", v.Operator)
			if v.Value != "" {
				qualityGateway += fmt.Sprintf(" - Value: %s  \n", v.Value)
			}
			qualityGateway += fmt.Sprintf(" - Error Threshold: %s  \n", v.ErrorThreshold)
			thumbsUP := ":+1:"
			if v.Status != "OK" {
				thumbsUP = ":-1:"
			}
			qualityGateway += fmt.Sprintf(" - Status: **%s** %s  \n", v.Status, thumbsUP)
		}
	}
	if qualityGateway != "" {
		finalForm += fmt.Sprintf("%s  \n%s", form, qualityGateway)
	}

	// fmt.Println(finalForm)
	return finalForm
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

func shouldUseProjectID(event *domain.Events) (shouldUseProjectID bool, projectID int, err error) {
	shouldUseProjectID = false
	if len(event.Properties) != 0 {
		if event.Properties["sonar.analysis.projectID"] != "" {
			value := event.Properties["sonar.analysis.projectID"]
			shouldUseProjectID = true
			projectID, err = strconv.Atoi(value)
			if err == nil {
				return true, projectID, err
			}
			fmt.Printf("[ERROR] Cannot convert sonar.analysis.projectID %s  \n", event.Project.Name)
			return false, 0, err
		}
	}
	fmt.Printf("[INFO] Cannot found sonar.analysis.projectID %s  \n", event.Project.Name)
	return shouldUseProjectID, projectID, err
}

// ValidateWebhook func to validate auth from Sonarqube
func ValidateWebhook(header string, message string, mysigning string) bool {
	mac := hmac.New(sha256.New, []byte(mysigning))
	if _, err := mac.Write([]byte(message)); err != nil {
		log.Printf("mac.Write(%v) failed\n", message)
		return false
	}
	calculatedMAC := hex.EncodeToString(mac.Sum(nil))
	// fmt.Printf("Headers: %s ; %s  \n", calculatedMAC, header)
	return hmac.Equal([]byte(calculatedMAC), []byte(header))
}

func checkProjectName(projectName string) (string, string) {
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
