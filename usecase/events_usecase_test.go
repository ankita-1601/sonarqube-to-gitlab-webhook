package usecase

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/appcontext"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/domain"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/utiltests"
	"github.com/stretchr/testify/assert"
)

func TestParseForm(t *testing.T) {
	eventByte := []byte(`{
		"serverUrl": "https://sonar.example.com",
		"taskId": "AXE0r5sCrP-5C-x-gqlN",
		"status": "SUCCESS",
		"analysedAt": "2020-04-01T07:40:53+0000",
		"revision": "2ba0cf4b9572f8d46f6b8f9a77f545faec64cf82",
		"changedAt": "2020-04-01T07:40:53+0000",
		"project": {
		  "key": "greatuser/test",
		  "name": "greatuser/test",
		  "url": "https://sonar.example.com/dashboard?id=greatuser%2Ftest"
		},
		"branch": {
		  "name": "test_branch",
		  "type": "BRANCH",
		  "isMain": false,
		  "url": "https://sonar.example.com/dashboard?id=greatuser%2Ftest&branch=test_branch"
		},
		"qualityGate": {
		  "name": "Example",
		  "status": "OK",
		  "conditions": [
			{
			  "metric": "new_reliability_rating",
			  "operator": "GREATER_THAN",
			  "status": "NO_VALUE",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_security_rating",
			  "operator": "GREATER_THAN",
			  "status": "NO_VALUE",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_maintainability_rating",
			  "operator": "GREATER_THAN",
			  "status": "NO_VALUE",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_coverage",
			  "operator": "LESS_THAN",
			  "status": "NO_VALUE",
			  "errorThreshold": "70"
			},
			{
			  "metric": "new_duplicated_lines_density",
			  "operator": "GREATER_THAN",
			  "status": "NO_VALUE",
			  "errorThreshold": "3"
			}
		  ]
		},
		"properties": {}
	  }`)
	var event domain.Events
	_ = json.Unmarshal(eventByte, &event)
	form := parseForm(&event)
	assert.Contains(t, form, "OK")
}

func TestParseFormComplete(t *testing.T) {
	completeByte := []byte(`{
		"serverUrl": "https://sonar.example.com",
		"taskId": "AXE0vIlqrP-5C-x-gqlQ",
		"status": "SUCCESS",
		"analysedAt": "2020-04-01T07:54:11+0000",
		"revision": "1d1cf5c299274e72420a39e4f0ad2cbf424d2719",
		"changedAt": "2020-04-01T07:54:11+0000",
		"project": {
		  "key": "great-project",
		  "name": "great-project",
		  "url": "https://sonar.example.com/dashboard?id=great-project"
		},
		"branch": {
		  "name": "869",
		  "type": "PULL_REQUEST",
		  "isMain": false,
		  "url": "https://sonar.example.com/dashboard?id=great-project&pullRequest=869"
		},
		"qualityGate": {
		  "name": "Example",
		  "status": "ERROR",
		  "conditions": [
			{
			  "metric": "new_reliability_rating",
			  "operator": "GREATER_THAN",
			  "value": "1",
			  "status": "OK",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_security_rating",
			  "operator": "GREATER_THAN",
			  "value": "1",
			  "status": "OK",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_maintainability_rating",
			  "operator": "GREATER_THAN",
			  "value": "1",
			  "status": "OK",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_coverage",
			  "operator": "LESS_THAN",
			  "value": "90.96385542168674",
			  "status": "OK",
			  "errorThreshold": "70"
			},
			{
			  "metric": "new_duplicated_lines_density",
			  "operator": "GREATER_THAN",
			  "value": "12.13768115942029",
			  "status": "ERROR",
			  "errorThreshold": "3"
			}
		  ]
		},
		"properties": {}
	  }`)
	var event domain.Events
	_ = json.Unmarshal(completeByte, &event)
	form := parseForm(&event)
	assert.Contains(t, form, "ERROR")
}

func TestShouldPostAndQuality(t *testing.T) {
	byteWithProperties := []byte(`{
		"serverUrl": "https://sonar.example.com",
		"taskId": "AXE1WPblrP-5C-x-gql1",
		"status": "SUCCESS",
		"analysedAt": "2020-04-01T10:45:53+0000",
		"revision": "c2f696f5a5edfc21058ce7984308447afc24c800",
		"changedAt": "2020-04-01T10:45:53+0000",
		"project": {
		  "key": "greatuser/test",
		  "name": "greatuser/test",
		  "url": "https://sonar.example.com/dashboard?id=greatuser%2Ftest"
		},
		"branch": {
		  "name": "test_anotherBranch",
		  "type": "BRANCH",
		  "isMain": false,
		  "url": "https://sonar.example.com/dashboard?id=greatuser%2Ftest&branch=test_anotherBranch"
		},
		"qualityGate": {
		  "name": "Example",
		  "status": "OK",
		  "conditions": [
			{
			  "metric": "new_reliability_rating",
			  "operator": "GREATER_THAN",
			  "value": "1",
			  "status": "OK",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_security_rating",
			  "operator": "GREATER_THAN",
			  "value": "1",
			  "status": "OK",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_maintainability_rating",
			  "operator": "GREATER_THAN",
			  "value": "1",
			  "status": "OK",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_coverage",
			  "operator": "LESS_THAN",
			  "status": "NO_VALUE",
			  "errorThreshold": "70"
			},
			{
			  "metric": "new_duplicated_lines_density",
			  "operator": "GREATER_THAN",
			  "status": "NO_VALUE",
			  "errorThreshold": "3"
			}
		  ]
		},
		"properties": {
		  "sonar.analysis.disabledGitlabPost": "true",
		  "sonar.analysis.disabledQualityReport": "true",
		  "sonar.analysis.projectID": "123"
		}
	  }`)
	var event domain.Events
	_ = json.Unmarshal(byteWithProperties, &event)
	properties := shouldPost(&event)
	assert.True(t, properties)
	quality := shouldSendQualityReport(&event)
	assert.True(t, quality)
	shouldUseProjectID, projectID, err := shouldUseProjectID(&event)
	assert.NoError(t, err)
	expected := 123
	assert.True(t, shouldUseProjectID)
	assert.Equal(t, expected, projectID)
}

func TestShouldPostAndQualityEmpty(t *testing.T) {
	byteWithProperties := []byte(`{
		"serverUrl": "https://sonar.example.com",
		"taskId": "AXE1WPblrP-5C-x-gql1",
		"status": "SUCCESS",
		"analysedAt": "2020-04-01T10:45:53+0000",
		"revision": "c2f696f5a5edfc21058ce7984308447afc24c800",
		"changedAt": "2020-04-01T10:45:53+0000",
		"project": {
		  "key": "greatuser/test",
		  "name": "greatuser/test",
		  "url": "https://sonar.example.com/dashboard?id=greatuser%2Ftest"
		},
		"branch": {
		  "name": "test_anotherBranch",
		  "type": "BRANCH",
		  "isMain": false,
		  "url": "https://sonar.example.com/dashboard?id=greatuser%2Ftest&branch=test_anotherBranch"
		},
		"qualityGate": {
		  "name": "Example",
		  "status": "OK",
		  "conditions": [
			{
			  "metric": "new_reliability_rating",
			  "operator": "GREATER_THAN",
			  "value": "1",
			  "status": "OK",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_security_rating",
			  "operator": "GREATER_THAN",
			  "value": "1",
			  "status": "OK",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_maintainability_rating",
			  "operator": "GREATER_THAN",
			  "value": "1",
			  "status": "OK",
			  "errorThreshold": "1"
			},
			{
			  "metric": "new_coverage",
			  "operator": "LESS_THAN",
			  "status": "NO_VALUE",
			  "errorThreshold": "70"
			},
			{
			  "metric": "new_duplicated_lines_density",
			  "operator": "GREATER_THAN",
			  "status": "NO_VALUE",
			  "errorThreshold": "3"
			}
		  ]
		},
		"properties": {}
	  }`)
	var event domain.Events
	_ = json.Unmarshal(byteWithProperties, &event)
	properties := shouldPost(&event)
	assert.False(t, properties)
	quality := shouldSendQualityReport(&event)
	assert.False(t, quality)
	shouldUseProjectID, projectID, err := shouldUseProjectID(&event)
	assert.NoError(t, err)
	expected := 0
	assert.False(t, shouldUseProjectID)
	assert.Equal(t, expected, projectID)
}

func TestCheckProjectName(t *testing.T) {
	// Expected: projectGroup/projectName
	testProjectName1 := "projectGroupFoo/projectNameBar"
	project1, projectPathWithNamespace1 := checkProjectName(testProjectName1)
	assert.Contains(t, project1, "projectNameBar")
	assert.Contains(t, projectPathWithNamespace1, "projectGroupFoo")

	// Expected: com.example.local:namewithoutspaces:name-with-dashes-or-not
	testProjectName2 := "com.example.local:projectfoobar:name-with-dashes-or-not"
	project2, projectPathWithNamespace2 := checkProjectName(testProjectName2)
	assert.Contains(t, project2, "projectfoobar")
	assert.Contains(t, projectPathWithNamespace2, notDefined)

	// only a single name
	testProjectName3 := "projectFooBar"
	project3, projectPathWithNamespace3 := checkProjectName(testProjectName3)
	assert.Contains(t, project3, "projectFooBar")
	assert.Contains(t, projectPathWithNamespace3, notDefined)
}

func TestValidateWebhook(t *testing.T) {
	longString := string("2aeccc9c03b36fea59ebec69")
	wrongLongString := string("1656ca38de0749bb8620814f")
	bodyString := string("body")
	timestamp := "1580475458"
	baseString := fmt.Sprintf("v0:%s:%s", timestamp, bodyString)
	signature := "38918cc13bf5e8b9ad7d2b7d85595d5d19af21783b86dbcfe4716422277ae404"
	if ValidateWebhook(signature, baseString, longString) != true {
		t.Fatalf("Invalid 1.1 testValidateWebhook %s, %s", signature, longString)
	}
	if ValidateWebhook(signature, baseString, wrongLongString) == true {
		t.Fatalf("Invalid 1.2 testValidateWebhook %s, %s", signature, wrongLongString)
	}
}

func TestSearchProjectID(t *testing.T) {
	config.GitlabURL = "http://gitlab.example.local"
	config.GitlabToken = "aaaa-bbb-ssss"
	repo := utiltests.RepositoryMock{}
	appcontext.Current.Add(appcontext.Repository, repo)
	projectID, err := searchProjectID("foobar", "foobar/foobarsystem")
	assert.NoError(t, err)
	expected := 1
	if utiltests.GetGitlabCalls != expected {
		t.Fatalf("Invalid 1.1 TestSearchProjectID %d", utiltests.GetGitlabCalls)
	}
	expectedID := 2
	assert.Equal(t, expectedID, projectID)
}

func TestSearchProjectIDWithoutNamespace(t *testing.T) {
	config.GitlabURL = "http://gitlab.example.local"
	config.GitlabToken = "aaaa-bbb-ssss"
	repo := utiltests.RepositoryMock{}
	appcontext.Current.Add(appcontext.Repository, repo)
	projectID, err := searchProjectID("fakesystem", notDefined)
	assert.NoError(t, err)
	expectedID := 1
	assert.Equal(t, expectedID, projectID)
}

func TestGitlabPostCommit(t *testing.T) {
	config.GitlabToken = "aaaa"
	config.GitlabURL = "https://gitlab.example.local"
	repo := utiltests.RepositoryMock{}
	appcontext.Current.Add(appcontext.Repository, repo)
	extraParams := map[string]string{
		"note": "nice comment",
	}
	gitlabPostCommit(utiltests.FakeURL, extraParams)
	expected := 1
	if utiltests.GitlabPostCommentCalls != expected {
		t.Fatalf("Invalid 1.1 TestGitlabPostCommit %d", utiltests.GitlabPostCommentCalls)
	}
}
