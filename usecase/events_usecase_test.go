package usecase

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/appcontext"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/domain"
	localtest "github.com/betorvs/sonarqube-to-gitlab-webhook/test"
	"github.com/stretchr/testify/assert"
)

func TestParseForm(t *testing.T) {
	eventByte := []byte(localtest.TestValidJSON)
	var event domain.Events
	_ = json.Unmarshal(eventByte, &event)
	form := parseForm(&event)
	assert.Contains(t, form, "OK")
}

func TestParseFormComplete(t *testing.T) {
	completeByte := []byte(localtest.TestValidJSONComplete)
	var event domain.Events
	_ = json.Unmarshal(completeByte, &event)
	form := parseForm(&event)
	assert.Contains(t, form, "ERROR")
}

func TestShouldPostAndQuality(t *testing.T) {
	byteWithProperties := []byte(localtest.TestValidJSONWithProperties)
	var event domain.Events
	_ = json.Unmarshal(byteWithProperties, &event)
	properties := shouldPost(&event)
	assert.True(t, properties)
	quality := shouldSendQualityReport(&event)
	assert.True(t, quality)
	shouldUseProjectID1, projectID, err := shouldUseProjectID(&event)
	assert.NoError(t, err)
	expected := 123
	assert.True(t, shouldUseProjectID1)
	assert.Equal(t, expected, projectID)
	byteWithProperties1 := []byte(localtest.TestValidJSONWithInvalidProjectID)
	var event1 domain.Events
	_ = json.Unmarshal(byteWithProperties1, &event1)
	_, _, err1 := shouldUseProjectID(&event1)
	assert.Error(t, err1)
}

func TestShouldPostAndQualityEmpty(t *testing.T) {
	byteWithProperties := []byte(localtest.TestValidJSONWithoutProperties)
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
	res1 := ValidateWebhook(signature, baseString, longString)
	assert.True(t, res1)
	res2 := ValidateWebhook(signature, baseString, wrongLongString)
	assert.False(t, res2)
}

func TestSearchProjectID(t *testing.T) {
	config.Values.GitlabURL = "Absent"
	appcontext.Current.Add(appcontext.Repository, localtest.Initclient)
	_, err := searchProjectID("foobar", "foobar/foobarsystem")
	assert.Error(t, err)

	config.Values.GitlabURL = "http://gitlab.example.local"
	config.Values.GitlabToken = "aaaa-bbb-ssss"
	appcontext.Current.Add(appcontext.Repository, localtest.Initclient)
	projectID, err := searchProjectID("foobar", "foobar/foobarsystem")
	assert.NoError(t, err)
	expected := 1
	assert.Equal(t, localtest.GetGitlabCalls, expected)
	// if localtest.GetGitlabCalls != expected {
	// 	t.Fatalf("Invalid 1.1 TestSearchProjectID %d", localtest.GetGitlabCalls)
	// }
	expectedID := 2
	// expected2 := 2
	assert.Equal(t, expectedID, projectID)
	config.Values.GitlabURL = "http://gitlaberror.example.local"
	_, err2 := searchProjectID("foo2bar", "foobar/foobarsystem")
	assert.Error(t, err2)
	// assert.Equal(t, localtest.GetGitlabCalls, expected2) gitlabinvalid
	// invalid json
	config.Values.GitlabURL = "http://gitlabinvalid.example.local"
	appcontext.Current.Add(appcontext.Repository, localtest.Initclient)
	_, err3 := searchProjectID("foobar", "foobar/foobarsystem")
	assert.Error(t, err3)
	// empty response
	config.Values.GitlabURL = "http://gitlabempty.example.local"
	_, err4 := searchProjectID("foobar", "foobar/foobarsystem")
	assert.Error(t, err4)
}

func TestSearchProjectIDWithoutNamespace(t *testing.T) {
	config.Values.GitlabURL = "http://gitlab.example.local"
	config.Values.GitlabToken = "aaaa-bbb-ssss"
	appcontext.Current.Add(appcontext.Repository, localtest.Initclient)
	projectID, err := searchProjectID("fakesystem", notDefined)
	assert.NoError(t, err)
	expectedID := 1
	assert.Equal(t, expectedID, projectID)
}

func TestGitlabPostCommit(t *testing.T) {
	appcontext.Current.Add(appcontext.Logger, localtest.InitMockLogger)
	config.Values.GitlabToken = "aaaa"
	config.Values.GitlabURL = "https://gitlab.example.local"
	appcontext.Current.Add(appcontext.Repository, localtest.Initclient)
	extraParams := map[string]string{
		"note": "nice comment",
	}
	err1 := gitlabPostCommit(localtest.FakeURL, extraParams)
	assert.NoError(t, err1)
	expected := 1
	assert.Equal(t, localtest.GitlabPostCommentCalls, expected)
	// if localtest.GitlabPostCommentCalls != expected {
	// 	t.Fatalf("Invalid 1.1 TestGitlabPostCommit %d", localtest.GitlabPostCommentCalls)
	// }

	extraParams1 := map[string]string{
		"note": "testinvalid",
	}

	err2 := gitlabPostCommit(localtest.FakeURL, extraParams1)
	assert.Error(t, err2)
}

func TestGitlabCommit(t *testing.T) {
	appcontext.Current.Add(appcontext.Logger, localtest.InitMockLogger)
	//
	byteWithProperties1 := []byte(localtest.TestValidJSONWithInvalidProjectID)
	var event1 domain.Events
	_ = json.Unmarshal(byteWithProperties1, &event1)
	err1 := GitlabCommit(&event1)
	assert.Error(t, err1)
	//
	config.Values.GitlabToken = "aaaa"
	config.Values.GitlabURL = "http://gitlabinvalid.example.local"
	appcontext.Current.Add(appcontext.Repository, localtest.Initclient)
	byteWithProperties2 := []byte(localtest.TestValidJSONWithUnknownProjectID)
	var event2 domain.Events
	_ = json.Unmarshal(byteWithProperties2, &event2)
	err2 := GitlabCommit(&event2)
	assert.Error(t, err2)

	//
	config.Values.GitlabToken = "aaaa"
	config.Values.GitlabURL = "http://gitlab.example.local"
	appcontext.Current.Add(appcontext.Repository, localtest.Initclient)
	byteWithProperties3 := []byte(localtest.TestValidJSONFake)
	var event3 domain.Events
	_ = json.Unmarshal(byteWithProperties3, &event3)
	err3 := GitlabCommit(&event3)
	assert.NoError(t, err3)

	// TestValidJSONFakeInvalid
	config.Values.GitlabToken = "aaaa"
	config.Values.GitlabURL = "http://gitlab.example.local"
	appcontext.Current.Add(appcontext.Repository, localtest.Initclient)
	byteWithProperties4 := []byte(localtest.TestValidJSONFakeInvalid)
	var event4 domain.Events
	_ = json.Unmarshal(byteWithProperties4, &event4)
	err4 := GitlabCommit(&event4)
	assert.Error(t, err4)

	//
	config.Values.GitlabToken = "aaaa"
	config.Values.GitlabURL = "http://gitlabempty.example.local"
	appcontext.Current.Add(appcontext.Repository, localtest.Initclient)
	byteWithProperties5 := []byte(localtest.TestValidJSONWithDuplicateProject)
	var event5 domain.Events
	_ = json.Unmarshal(byteWithProperties5, &event5)
	err5 := GitlabCommit(&event5)
	assert.Error(t, err5)
}
