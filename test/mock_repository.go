package test

import (
	"fmt"
	"strings"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/appcontext"
)

const (
	// FakeURL const to be used in tests
	FakeURL = "http://gilab.local"
	// FakeToken const
	FakeToken = "aaa-bbb-ccc-x0"
	// FakeSonarqubeSecret const
	FakeSonarqubeSecret = "aaaaa-wwwww-ssss"
)

var (
	// GetGitlabCalls int
	GetGitlabCalls int
	// GitlabPostCommentCalls int
	GitlabPostCommentCalls int
	// TestBadJSON string
	TestBadJSON = `{"serverUrl":"http://sonar.local", error`
	// TestValidJSON string
	TestValidJSON = `{
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
	  }`
	// TestValidJSONComplete string
	TestValidJSONComplete = `{
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
	  }`
	// TestValidJSONFake string
	TestValidJSONFake = `{
		"serverUrl": "https://sonar.example.com",
		"taskId": "AXE0vIlqrP-5C-x-gqlQ",
		"status": "SUCCESS",
		"analysedAt": "2020-04-01T07:54:11+0000",
		"revision": "1d1cf5c299274e72420a39e4f0ad2cbf424d2719",
		"changedAt": "2020-04-01T07:54:11+0000",
		"project": {
		  "key": "fake/fakesystem",
		  "name": "fake/fakesystem",
		  "url": "https://sonar.example.com/dashboard?id=fake"
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
	  }`
	// TestValidJSONFakeInvalid string
	TestValidJSONFakeInvalid = `{
		"serverUrl": "https://sonar.example.com",
		"taskId": "AXE0vIlqrP-5C-x-gqlQ",
		"status": "SUCCESS",
		"analysedAt": "2020-04-01T07:54:11+0000",
		"revision": "1d1cf5c299274e72420a39e4f0ad2cbf424d2719",
		"changedAt": "2020-04-01T07:54:11+0000",
		"project": {
		  "key": "fake/fakesystem",
		  "name": "fake/fakesystem",
		  "url": "https://testinvalid.sonar.example.com/dashboard?id=fake"
		},
		"branch": {
		  "name": "869",
		  "type": "PULL_REQUEST",
		  "isMain": false,
		  "url": "https://testinvalid.sonar.example.com/dashboard?id=great-project&pullRequest=869"
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
	  }`

	// TestValidJSONWithProperties string
	TestValidJSONWithProperties = `{
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
	  }`
	// TestValidJSONWithoutProperties string
	TestValidJSONWithoutProperties = `{
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
	  }`
	// TestValidJSONWithInvalidProjectID string
	TestValidJSONWithInvalidProjectID = `{
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
		  "sonar.analysis.projectID": "12aa3"
		}
	  }`
	// TestValidJSONWithUnknownProjectID string
	TestValidJSONWithUnknownProjectID = `{
		"serverUrl": "https://sonar.example.com",
		"taskId": "AXE1WPblrP-5C-x-gql1",
		"status": "SUCCESS",
		"analysedAt": "2020-04-01T10:45:53+0000",
		"revision": "c2f696f5a5edfc21058ce7984308447afc24c800",
		"changedAt": "2020-04-01T10:45:53+0000",
		"project": {
		  "key": "wronguser/failedtest",
		  "name": "wronguser/test",
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
	  }`
	// TestValidJSONWithDuplicateProject string
	TestValidJSONWithDuplicateProject = `{
		"serverUrl": "https://sonar.example.com",
		"taskId": "AXE1WPblrP-5C-x-gql1",
		"status": "SUCCESS",
		"analysedAt": "2020-04-01T10:45:53+0000",
		"revision": "c2f696f5a5edfc21058ce7984308447afc24c800",
		"changedAt": "2020-04-01T10:45:53+0000",
		"project": {
		  "key": "teste",
		  "name": "teste",
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
	  }`
)

// RepositoryMock struct is used for Mock Repository requests
type RepositoryMock struct {
}

// GetGitlab func Mock
func (repo RepositoryMock) GetGitlab(gitlabURL, gitlabToken string) ([]byte, string, error) {
	if strings.Contains(gitlabURL, "http://gitlaberror.example.local") {
		return []byte{}, "", fmt.Errorf("invalid project id %s", gitlabURL)
	}
	if strings.Contains(gitlabURL, "http://gitlabinvalid.example.local") {
		bodyBytes := []byte(`[{aaaa: 123, asd: "test"}]`)
		return bodyBytes, "", nil
	}
	if strings.Contains(gitlabURL, "http://gitlabempty.example.local") {
		bodyBytes := []byte(`[]`)
		return bodyBytes, "", nil
	}
	//
	if strings.Contains(gitlabURL, "http://gitlabduplicate.example.local") {
		bodyBytes := []byte(`[{"id":3,"description":"Test system","name":"test","name_with_namespace":"test / testsystem","path":"test","path_with_namespace":"test/testsystem","created_at":"2018-08-20T18:54:13.601Z","default_branch":"master","tag_list":["test"],"ssh_url_to_repo":"git@gitlab.local:test/testsystem.git","http_url_to_repo":"https://gitlab.local/test/testsystem.git","web_url":"https://gitlab.local/test/testsystem","readme_url":"https://gitlab.local/test/testsystem/blob/master/README.md","avatar_url":"https://gitlab.local/uploads/-/system/project/avatar/2/calculator-icon.png","star_count":10,"forks_count":1,"last_activity_at":"2020-03-24T10:16:13.983Z","namespace":{"id":3,"name":"test","path":"test","kind":"group","full_path":"test","parent_id":null,"avatar_url":"/uploads/-/system/group/avatar/3/beans.png","web_url":"https://gitlab.local/groups/test"},"_links":{"self":"https://gitlab.local/api/v4/projects/3","merge_requests":"https://gitlab.local/api/v4/projects/3/merge_requests","repo_branches":"https://gitlab.local/api/v4/projects/3/repository/branches","labels":"https://gitlab.local/api/v4/projects/3/labels","events":"https://gitlab.local/api/v4/projects/3/events","members":"https://gitlab.local/api/v4/projects/3/members"},"empty_repo":false,"archived":false,"visibility":"internal","resolve_outdated_diff_discussions":false,"container_registry_enabled":true,"issues_enabled":false,"merge_requests_enabled":true,"wiki_enabled":false,"jobs_enabled":true,"snippets_enabled":true,"issues_access_level":"disabled","repository_access_level":"enabled","merge_requests_access_level":"enabled","wiki_access_level":"disabled","builds_access_level":"enabled","snippets_access_level":"enabled","shared_runners_enabled":true,"lfs_enabled":true,"creator_id":2,"import_status":"none","ci_default_git_depth":null,"public_jobs":true,"build_timeout":14400,"auto_cancel_pending_pipelines":"enabled","build_coverage_regex":"","ci_config_path":"","shared_with_groups":[{"group_id":59,"group_name":"test external","group_full_path":"testexternal","group_access_level":11,"expires_at":null}],"only_allow_merge_if_pipeline_succeeds":false,"request_access_enabled":false,"only_allow_merge_if_all_discussions_are_resolved":true,"remove_source_branch_after_merge":false,"printing_merge_request_link_enabled":true,"merge_method":"rebase_merge","suggestion_commit_message":"","auto_devops_enabled":false,"auto_devops_deploy_strategy":"continuous","autoclose_referenced_issues":true,"permissions":{"project_access":null,"group_access":{"access_level":21,"notification_level":3}},"approvals_before_merge":1,"mirror":false,"external_authorization_classification_label":null,"packages_enabled":false,"service_desk_enabled":false,"service_desk_address":null,"marked_for_deletion_at":null},{"id":4,"description":"Test component","name":"test","name_with_namespace":"test / testcomponent","path":"test","path_with_namespace":"test/testcomponent","created_at":"2018-08-20T18:54:13.601Z","default_branch":"master","tag_list":["test"],"ssh_url_to_repo":"git@gitlab.local:test/testcomponent.git","http_url_to_repo":"https://gitlab.local/test/testcomponent.git","web_url":"https://gitlab.local/test/testcomponent","readme_url":"https://gitlab.local/test/testcomponent/blob/master/README.md","avatar_url":"https://gitlab.local/uploads/-/component/project/avatar/2/calculator-icon.png","star_count":10,"forks_count":1,"last_activity_at":"2020-03-24T10:16:13.983Z","namespace":{"id":3,"name":"test","path":"test","kind":"group","full_path":"test","parent_id":null,"avatar_url":"/uploads/-/system/group/avatar/3/beans.png","web_url":"https://gitlab.local/groups/test"},"_links":{"self":"https://gitlab.local/api/v4/projects/4","merge_requests":"https://gitlab.local/api/v4/projects/4/merge_requests","repo_branches":"https://gitlab.local/api/v4/projects/4/repository/branches","labels":"https://gitlab.local/api/v4/projects/4/labels","events":"https://gitlab.local/api/v4/projects/4/events","members":"https://gitlab.local/api/v4/projects/4/members"},"empty_repo":false,"archived":false,"visibility":"internal","resolve_outdated_diff_discussions":false,"container_registry_enabled":true,"issues_enabled":false,"merge_requests_enabled":true,"wiki_enabled":false,"jobs_enabled":true,"snippets_enabled":true,"issues_access_level":"disabled","repository_access_level":"enabled","merge_requests_access_level":"enabled","wiki_access_level":"disabled","builds_access_level":"enabled","snippets_access_level":"enabled","shared_runners_enabled":true,"lfs_enabled":true,"creator_id":2,"import_status":"none","ci_default_git_depth":null,"public_jobs":true,"build_timeout":14400,"auto_cancel_pending_pipelines":"enabled","build_coverage_regex":"","ci_config_path":"","shared_with_groups":[{"group_id":59,"group_name":"test external","group_full_path":"testexternal","group_access_level":11,"expires_at":null}],"only_allow_merge_if_pipeline_succeeds":false,"request_access_enabled":false,"only_allow_merge_if_all_discussions_are_resolved":true,"remove_source_branch_after_merge":false,"printing_merge_request_link_enabled":true,"merge_method":"rebase_merge","suggestion_commit_message":"","auto_devops_enabled":false,"auto_devops_deploy_strategy":"continuous","autoclose_referenced_issues":true,"permissions":{"project_access":null,"group_access":{"access_level":21,"notification_level":3}},"approvals_before_merge":1,"mirror":false,"external_authorization_classification_label":null,"packages_enabled":false,"service_desk_enabled":false,"service_desk_address":null,"marked_for_deletion_at":null}]
	`)
		return bodyBytes, "", nil
	}

	GetGitlabCalls++
	// returnByteProjects := []byte(`[{"id":1,"description":"Fake system","name":"fake","name_with_namespace":"fake / fakesystem","path":"fake","path_with_namespace":"fake/fakesystem","created_at":"2018-08-20T18:54:13.601Z","default_branch":"master","tag_list":["Fake"],"ssh_url_to_repo":"git@gitlab.local:fake/fakesystem.git","http_url_to_repo":"https://gitlab.local/fake/fakesystem.git","web_url":"https://gitlab.local/fake/fakesystem","readme_url":"https://gitlab.local/fake/fakesystem/blob/master/README.md","avatar_url":"https://gitlab.local/uploads/-/system/project/avatar/1/calculator-icon.png","star_count":10,"forks_count":1,"last_activity_at":"2020-03-24T10:16:13.983Z","namespace":{"id":3,"name":"fake","path":"fake","kind":"group","full_path":"fake","parent_id":null,"avatar_url":"/uploads/-/system/group/avatar/3/cacao_beans.png","web_url":"https://gitlab.local/groups/fake"},"_links":{"self":"https://gitlab.local/api/v4/projects/1","merge_requests":"https://gitlab.local/api/v4/projects/1/merge_requests","repo_branches":"https://gitlab.local/api/v4/projects/1/repository/branches","labels":"https://gitlab.local/api/v4/projects/1/labels","events":"https://gitlab.local/api/v4/projects/1/events","members":"https://gitlab.local/api/v4/projects/1/members"},"empty_repo":false,"archived":false,"visibility":"internal","resolve_outdated_diff_discussions":false,"container_registry_enabled":true,"issues_enabled":false,"merge_requests_enabled":true,"wiki_enabled":false,"jobs_enabled":true,"snippets_enabled":true,"issues_access_level":"disabled","repository_access_level":"enabled","merge_requests_access_level":"enabled","wiki_access_level":"disabled","builds_access_level":"enabled","snippets_access_level":"enabled","shared_runners_enabled":true,"lfs_enabled":true,"creator_id":2,"import_status":"none","ci_default_git_depth":null,"public_jobs":true,"build_timeout":14400,"auto_cancel_pending_pipelines":"enabled","build_coverage_regex":"","ci_config_path":"","shared_with_groups":[{"group_id":58,"group_name":"Fake external","group_full_path":"fakeexternal","group_access_level":10,"expires_at":null}],"only_allow_merge_if_pipeline_succeeds":false,"request_access_enabled":false,"only_allow_merge_if_all_discussions_are_resolved":true,"remove_source_branch_after_merge":false,"printing_merge_request_link_enabled":true,"merge_method":"rebase_merge","suggestion_commit_message":"","auto_devops_enabled":false,"auto_devops_deploy_strategy":"continuous","autoclose_referenced_issues":true,"permissions":{"project_access":null,"group_access":{"access_level":20,"notification_level":3}},"approvals_before_merge":1,"mirror":false,"external_authorization_classification_label":null,"packages_enabled":false,"service_desk_enabled":false,"service_desk_address":null,"marked_for_deletion_at":null}]`)
	// returnByteProjects := []byte(`[{"id":1,"description":"Fake system","name":"fake","name_with_namespace":"fake / fakesystem","path":"fake","path_with_namespace":"fake/fakesystem","created_at":"2018-08-20T18:54:13.601Z","default_branch":"master","tag_list":["Fake"],"ssh_url_to_repo":"git@gitlab.local:fake/fakesystem.git","http_url_to_repo":"https://gitlab.local/fake/fakesystem.git","web_url":"https://gitlab.local/fake/fakesystem","readme_url":"https://gitlab.local/fake/fakesystem/blob/master/README.md","avatar_url":"https://gitlab.local/uploads/-/system/project/avatar/1/calculator-icon.png","star_count":10,"forks_count":1,"last_activity_at":"2020-03-24T10:16:13.983Z","namespace":{"id":3,"name":"fake","path":"fake","kind":"group","full_path":"fake","parent_id":null,"avatar_url":"/uploads/-/system/group/avatar/3/cacao_beans.png","web_url":"https://gitlab.local/groups/fake"},"_links":{"self":"https://gitlab.local/api/v4/projects/1","merge_requests":"https://gitlab.local/api/v4/projects/1/merge_requests","repo_branches":"https://gitlab.local/api/v4/projects/1/repository/branches","labels":"https://gitlab.local/api/v4/projects/1/labels","events":"https://gitlab.local/api/v4/projects/1/events","members":"https://gitlab.local/api/v4/projects/1/members"},"empty_repo":false,"archived":false,"visibility":"internal","resolve_outdated_diff_discussions":false,"container_registry_enabled":true,"issues_enabled":false,"merge_requests_enabled":true,"wiki_enabled":false,"jobs_enabled":true,"snippets_enabled":true,"issues_access_level":"disabled","repository_access_level":"enabled","merge_requests_access_level":"enabled","wiki_access_level":"disabled","builds_access_level":"enabled","snippets_access_level":"enabled","shared_runners_enabled":true,"lfs_enabled":true,"creator_id":2,"import_status":"none","ci_default_git_depth":null,"public_jobs":true,"build_timeout":14400,"auto_cancel_pending_pipelines":"enabled","build_coverage_regex":"","ci_config_path":"","shared_with_groups":[{"group_id":58,"group_name":"Fake external","group_full_path":"fakeexternal","group_access_level":10,"expires_at":null}],"only_allow_merge_if_pipeline_succeeds":false,"request_access_enabled":false,"only_allow_merge_if_all_discussions_are_resolved":true,"remove_source_branch_after_merge":false,"printing_merge_request_link_enabled":true,"merge_method":"rebase_merge","suggestion_commit_message":"","auto_devops_enabled":false,"auto_devops_deploy_strategy":"continuous","autoclose_referenced_issues":true,"permissions":{"project_access":null,"group_access":{"access_level":20,"notification_level":3}},"approvals_before_merge":1,"mirror":false,"external_authorization_classification_label":null,"packages_enabled":false,"service_desk_enabled":false,"service_desk_address":null,"marked_for_deletion_at":null},{"id":2,"description":"Foobar system","name":"foobar","name_with_namespace":"foobar / foobarsystem","path":"foobar","path_with_namespace":"foobar/foobarsystem","created_at":"2018-08-20T18:54:13.601Z","default_branch":"master","tag_list":["foobar"],"ssh_url_to_repo":"git@gitlab.local:foobar/foobarsystem.git","http_url_to_repo":"https://gitlab.local/foobar/foobarsystem.git","web_url":"https://gitlab.local/foobar/foobarsystem","readme_url":"https://gitlab.local/foobar/foobarsystem/blob/master/README.md","avatar_url":"https://gitlab.local/uploads/-/system/project/avatar/2/calculator-icon.png","star_count":10,"forks_count":1,"last_activity_at":"2020-03-24T10:16:13.983Z","namespace":{"id":3,"name":"foobar","path":"foobar","kind":"group","full_path":"foobar","parent_id":null,"avatar_url":"/uploads/-/system/group/avatar/3/beans.png","web_url":"https://gitlab.local/groups/foobar"},"_links":{"self":"https://gitlab.local/api/v4/projects/2","merge_requests":"https://gitlab.local/api/v4/projects/2/merge_requests","repo_branches":"https://gitlab.local/api/v4/projects/2/repository/branches","labels":"https://gitlab.local/api/v4/projects/2/labels","events":"https://gitlab.local/api/v4/projects/2/events","members":"https://gitlab.local/api/v4/projects/2/members"},"empty_repo":false,"archived":false,"visibility":"internal","resolve_outdated_diff_discussions":false,"container_registry_enabled":true,"issues_enabled":false,"merge_requests_enabled":true,"wiki_enabled":false,"jobs_enabled":true,"snippets_enabled":true,"issues_access_level":"disabled","repository_access_level":"enabled","merge_requests_access_level":"enabled","wiki_access_level":"disabled","builds_access_level":"enabled","snippets_access_level":"enabled","shared_runners_enabled":true,"lfs_enabled":true,"creator_id":2,"import_status":"none","ci_default_git_depth":null,"public_jobs":true,"build_timeout":14400,"auto_cancel_pending_pipelines":"enabled","build_coverage_regex":"","ci_config_path":"","shared_with_groups":[{"group_id":59,"group_name":"foobar external","group_full_path":"foobarexternal","group_access_level":11,"expires_at":null}],"only_allow_merge_if_pipeline_succeeds":false,"request_access_enabled":false,"only_allow_merge_if_all_discussions_are_resolved":true,"remove_source_branch_after_merge":false,"printing_merge_request_link_enabled":true,"merge_method":"rebase_merge","suggestion_commit_message":"","auto_devops_enabled":false,"auto_devops_deploy_strategy":"continuous","autoclose_referenced_issues":true,"permissions":{"project_access":null,"group_access":{"access_level":21,"notification_level":3}},"approvals_before_merge":1,"mirror":false,"external_authorization_classification_label":null,"packages_enabled":false,"service_desk_enabled":false,"service_desk_address":null,"marked_for_deletion_at":null}]`)
	returnByteProjects := []byte(`[{"id":1,"description":"Fake system","name":"fake","name_with_namespace":"fake / fakesystem","path":"fake","path_with_namespace":"fake/fakesystem","created_at":"2018-08-20T18:54:13.601Z","default_branch":"master","tag_list":["Fake"],"ssh_url_to_repo":"git@gitlab.local:fake/fakesystem.git","http_url_to_repo":"https://gitlab.local/fake/fakesystem.git","web_url":"https://gitlab.local/fake/fakesystem","readme_url":"https://gitlab.local/fake/fakesystem/blob/master/README.md","avatar_url":"https://gitlab.local/uploads/-/system/project/avatar/1/calculator-icon.png","star_count":10,"forks_count":1,"last_activity_at":"2020-03-24T10:16:13.983Z","namespace":{"id":3,"name":"fake","path":"fake","kind":"group","full_path":"fake","parent_id":null,"avatar_url":"/uploads/-/system/group/avatar/3/cacao_beans.png","web_url":"https://gitlab.local/groups/fake"},"_links":{"self":"https://gitlab.local/api/v4/projects/1","merge_requests":"https://gitlab.local/api/v4/projects/1/merge_requests","repo_branches":"https://gitlab.local/api/v4/projects/1/repository/branches","labels":"https://gitlab.local/api/v4/projects/1/labels","events":"https://gitlab.local/api/v4/projects/1/events","members":"https://gitlab.local/api/v4/projects/1/members"},"empty_repo":false,"archived":false,"visibility":"internal","resolve_outdated_diff_discussions":false,"container_registry_enabled":true,"issues_enabled":false,"merge_requests_enabled":true,"wiki_enabled":false,"jobs_enabled":true,"snippets_enabled":true,"issues_access_level":"disabled","repository_access_level":"enabled","merge_requests_access_level":"enabled","wiki_access_level":"disabled","builds_access_level":"enabled","snippets_access_level":"enabled","shared_runners_enabled":true,"lfs_enabled":true,"creator_id":2,"import_status":"none","ci_default_git_depth":null,"public_jobs":true,"build_timeout":14400,"auto_cancel_pending_pipelines":"enabled","build_coverage_regex":"","ci_config_path":"","shared_with_groups":[{"group_id":58,"group_name":"Fake external","group_full_path":"fakeexternal","group_access_level":10,"expires_at":null}],"only_allow_merge_if_pipeline_succeeds":false,"request_access_enabled":false,"only_allow_merge_if_all_discussions_are_resolved":true,"remove_source_branch_after_merge":false,"printing_merge_request_link_enabled":true,"merge_method":"rebase_merge","suggestion_commit_message":"","auto_devops_enabled":false,"auto_devops_deploy_strategy":"continuous","autoclose_referenced_issues":true,"permissions":{"project_access":null,"group_access":{"access_level":20,"notification_level":3}},"approvals_before_merge":1,"mirror":false,"external_authorization_classification_label":null,"packages_enabled":false,"service_desk_enabled":false,"service_desk_address":null,"marked_for_deletion_at":null},{"id":2,"description":"Foobar system","name":"foobar","name_with_namespace":"foobar / foobarsystem","path":"foobar","path_with_namespace":"foobar/foobarsystem","created_at":"2018-08-20T18:54:13.601Z","default_branch":"master","tag_list":["foobar"],"ssh_url_to_repo":"git@gitlab.local:foobar/foobarsystem.git","http_url_to_repo":"https://gitlab.local/foobar/foobarsystem.git","web_url":"https://gitlab.local/foobar/foobarsystem","readme_url":"https://gitlab.local/foobar/foobarsystem/blob/master/README.md","avatar_url":"https://gitlab.local/uploads/-/system/project/avatar/2/calculator-icon.png","star_count":10,"forks_count":1,"last_activity_at":"2020-03-24T10:16:13.983Z","namespace":{"id":3,"name":"foobar","path":"foobar","kind":"group","full_path":"foobar","parent_id":null,"avatar_url":"/uploads/-/system/group/avatar/3/beans.png","web_url":"https://gitlab.local/groups/foobar"},"_links":{"self":"https://gitlab.local/api/v4/projects/2","merge_requests":"https://gitlab.local/api/v4/projects/2/merge_requests","repo_branches":"https://gitlab.local/api/v4/projects/2/repository/branches","labels":"https://gitlab.local/api/v4/projects/2/labels","events":"https://gitlab.local/api/v4/projects/2/events","members":"https://gitlab.local/api/v4/projects/2/members"},"empty_repo":false,"archived":false,"visibility":"internal","resolve_outdated_diff_discussions":false,"container_registry_enabled":true,"issues_enabled":false,"merge_requests_enabled":true,"wiki_enabled":false,"jobs_enabled":true,"snippets_enabled":true,"issues_access_level":"disabled","repository_access_level":"enabled","merge_requests_access_level":"enabled","wiki_access_level":"disabled","builds_access_level":"enabled","snippets_access_level":"enabled","shared_runners_enabled":true,"lfs_enabled":true,"creator_id":2,"import_status":"none","ci_default_git_depth":null,"public_jobs":true,"build_timeout":14400,"auto_cancel_pending_pipelines":"enabled","build_coverage_regex":"","ci_config_path":"","shared_with_groups":[{"group_id":59,"group_name":"foobar external","group_full_path":"foobarexternal","group_access_level":11,"expires_at":null}],"only_allow_merge_if_pipeline_succeeds":false,"request_access_enabled":false,"only_allow_merge_if_all_discussions_are_resolved":true,"remove_source_branch_after_merge":false,"printing_merge_request_link_enabled":true,"merge_method":"rebase_merge","suggestion_commit_message":"","auto_devops_enabled":false,"auto_devops_deploy_strategy":"continuous","autoclose_referenced_issues":true,"permissions":{"project_access":null,"group_access":{"access_level":21,"notification_level":3}},"approvals_before_merge":1,"mirror":false,"external_authorization_classification_label":null,"packages_enabled":false,"service_desk_enabled":false,"service_desk_address":null,"marked_for_deletion_at":null},{"id":3,"description":"Test system","name":"test","name_with_namespace":"test / testsystem","path":"test","path_with_namespace":"test/testsystem","created_at":"2018-08-20T18:54:13.601Z","default_branch":"master","tag_list":["test"],"ssh_url_to_repo":"git@gitlab.local:test/testsystem.git","http_url_to_repo":"https://gitlab.local/test/testsystem.git","web_url":"https://gitlab.local/test/testsystem","readme_url":"https://gitlab.local/test/testsystem/blob/master/README.md","avatar_url":"https://gitlab.local/uploads/-/system/project/avatar/2/calculator-icon.png","star_count":10,"forks_count":1,"last_activity_at":"2020-03-24T10:16:13.983Z","namespace":{"id":3,"name":"test","path":"test","kind":"group","full_path":"test","parent_id":null,"avatar_url":"/uploads/-/system/group/avatar/3/beans.png","web_url":"https://gitlab.local/groups/test"},"_links":{"self":"https://gitlab.local/api/v4/projects/3","merge_requests":"https://gitlab.local/api/v4/projects/3/merge_requests","repo_branches":"https://gitlab.local/api/v4/projects/3/repository/branches","labels":"https://gitlab.local/api/v4/projects/3/labels","events":"https://gitlab.local/api/v4/projects/3/events","members":"https://gitlab.local/api/v4/projects/3/members"},"empty_repo":false,"archived":false,"visibility":"internal","resolve_outdated_diff_discussions":false,"container_registry_enabled":true,"issues_enabled":false,"merge_requests_enabled":true,"wiki_enabled":false,"jobs_enabled":true,"snippets_enabled":true,"issues_access_level":"disabled","repository_access_level":"enabled","merge_requests_access_level":"enabled","wiki_access_level":"disabled","builds_access_level":"enabled","snippets_access_level":"enabled","shared_runners_enabled":true,"lfs_enabled":true,"creator_id":2,"import_status":"none","ci_default_git_depth":null,"public_jobs":true,"build_timeout":14400,"auto_cancel_pending_pipelines":"enabled","build_coverage_regex":"","ci_config_path":"","shared_with_groups":[{"group_id":59,"group_name":"test external","group_full_path":"testexternal","group_access_level":11,"expires_at":null}],"only_allow_merge_if_pipeline_succeeds":false,"request_access_enabled":false,"only_allow_merge_if_all_discussions_are_resolved":true,"remove_source_branch_after_merge":false,"printing_merge_request_link_enabled":true,"merge_method":"rebase_merge","suggestion_commit_message":"","auto_devops_enabled":false,"auto_devops_deploy_strategy":"continuous","autoclose_referenced_issues":true,"permissions":{"project_access":null,"group_access":{"access_level":21,"notification_level":3}},"approvals_before_merge":1,"mirror":false,"external_authorization_classification_label":null,"packages_enabled":false,"service_desk_enabled":false,"service_desk_address":null,"marked_for_deletion_at":null},{"id":4,"description":"Test component","name":"test","name_with_namespace":"test / testcomponent","path":"test","path_with_namespace":"test/testcomponent","created_at":"2018-08-20T18:54:13.601Z","default_branch":"master","tag_list":["test"],"ssh_url_to_repo":"git@gitlab.local:test/testcomponent.git","http_url_to_repo":"https://gitlab.local/test/testcomponent.git","web_url":"https://gitlab.local/test/testcomponent","readme_url":"https://gitlab.local/test/testcomponent/blob/master/README.md","avatar_url":"https://gitlab.local/uploads/-/component/project/avatar/2/calculator-icon.png","star_count":10,"forks_count":1,"last_activity_at":"2020-03-24T10:16:13.983Z","namespace":{"id":3,"name":"test","path":"test","kind":"group","full_path":"test","parent_id":null,"avatar_url":"/uploads/-/system/group/avatar/3/beans.png","web_url":"https://gitlab.local/groups/test"},"_links":{"self":"https://gitlab.local/api/v4/projects/4","merge_requests":"https://gitlab.local/api/v4/projects/4/merge_requests","repo_branches":"https://gitlab.local/api/v4/projects/4/repository/branches","labels":"https://gitlab.local/api/v4/projects/4/labels","events":"https://gitlab.local/api/v4/projects/4/events","members":"https://gitlab.local/api/v4/projects/4/members"},"empty_repo":false,"archived":false,"visibility":"internal","resolve_outdated_diff_discussions":false,"container_registry_enabled":true,"issues_enabled":false,"merge_requests_enabled":true,"wiki_enabled":false,"jobs_enabled":true,"snippets_enabled":true,"issues_access_level":"disabled","repository_access_level":"enabled","merge_requests_access_level":"enabled","wiki_access_level":"disabled","builds_access_level":"enabled","snippets_access_level":"enabled","shared_runners_enabled":true,"lfs_enabled":true,"creator_id":2,"import_status":"none","ci_default_git_depth":null,"public_jobs":true,"build_timeout":14400,"auto_cancel_pending_pipelines":"enabled","build_coverage_regex":"","ci_config_path":"","shared_with_groups":[{"group_id":59,"group_name":"test external","group_full_path":"testexternal","group_access_level":11,"expires_at":null}],"only_allow_merge_if_pipeline_succeeds":false,"request_access_enabled":false,"only_allow_merge_if_all_discussions_are_resolved":true,"remove_source_branch_after_merge":false,"printing_merge_request_link_enabled":true,"merge_method":"rebase_merge","suggestion_commit_message":"","auto_devops_enabled":false,"auto_devops_deploy_strategy":"continuous","autoclose_referenced_issues":true,"permissions":{"project_access":null,"group_access":{"access_level":21,"notification_level":3}},"approvals_before_merge":1,"mirror":false,"external_authorization_classification_label":null,"packages_enabled":false,"service_desk_enabled":false,"service_desk_address":null,"marked_for_deletion_at":null}]
`)
	return returnByteProjects, "", nil
}

// GitlabPostComment func mock
func (repo RepositoryMock) GitlabPostComment(url string, params map[string]string) error {
	if strings.Contains(params["note"], "testinvalid") {
		return fmt.Errorf("test invalid %s", url)
	}
	GitlabPostCommentCalls++
	return nil
}

// Initclient func returns a RepositoryMock interface
func Initclient() appcontext.Component {

	return RepositoryMock{}
}
