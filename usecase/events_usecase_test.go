package usecase

import (
	"encoding/json"
	"testing"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/domain"
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
	form, err := parseForm(&event)
	assert.NoError(t, err)
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
	form, err := parseForm(&event)
	assert.NoError(t, err)
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
		  "sonar.analysis.disabledQualityReport": "true"
		}
	  }`)
	var event domain.Events
	_ = json.Unmarshal(byteWithProperties, &event)
	properties := shouldPost(&event)
	assert.True(t, properties)
	quality := shouldSendQualityReport(&event)
	assert.True(t, quality)
	// assert.Contains(t, properties, "True")

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

}
