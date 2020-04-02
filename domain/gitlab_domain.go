package domain

import (
	"time"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/appcontext"
)

// GitlabProject struct
type GitlabProject struct {
	ID            int         `json:"id"`
	Description   interface{} `json:"description"`
	DefaultBranch string      `json:"default_branch"`
	Visibility    string      `json:"visibility"`
	SSHURLToRepo  string      `json:"ssh_url_to_repo"`
	HTTPURLToRepo string      `json:"http_url_to_repo"`
	WebURL        string      `json:"web_url"`
	ReadmeURL     string      `json:"readme_url"`
	TagList       []string    `json:"tag_list"`
	Owner         struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"owner"`
	Name                           string    `json:"name"`
	NameWithNamespace              string    `json:"name_with_namespace"`
	Path                           string    `json:"path"`
	PathWithNamespace              string    `json:"path_with_namespace"`
	IssuesEnabled                  bool      `json:"issues_enabled"`
	OpenIssuesCount                int       `json:"open_issues_count"`
	MergeRequestsEnabled           bool      `json:"merge_requests_enabled"`
	JobsEnabled                    bool      `json:"jobs_enabled"`
	WikiEnabled                    bool      `json:"wiki_enabled"`
	SnippetsEnabled                bool      `json:"snippets_enabled"`
	ResolveOutdatedDiffDiscussions bool      `json:"resolve_outdated_diff_discussions"`
	ContainerRegistryEnabled       bool      `json:"container_registry_enabled"`
	CreatedAt                      time.Time `json:"created_at"`
	LastActivityAt                 time.Time `json:"last_activity_at"`
	CreatorID                      int       `json:"creator_id"`
	Namespace                      struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		Kind     string `json:"kind"`
		FullPath string `json:"full_path"`
	} `json:"namespace"`
	ImportStatus                              string        `json:"import_status"`
	Archived                                  bool          `json:"archived"`
	AvatarURL                                 string        `json:"avatar_url"`
	SharedRunnersEnabled                      bool          `json:"shared_runners_enabled"`
	ForksCount                                int           `json:"forks_count"`
	StarCount                                 int           `json:"star_count"`
	RunnersToken                              string        `json:"runners_token"`
	CiDefaultGitDepth                         int           `json:"ci_default_git_depth"`
	PublicJobs                                bool          `json:"public_jobs"`
	SharedWithGroups                          []interface{} `json:"shared_with_groups"`
	OnlyAllowMergeIfPipelineSucceeds          bool          `json:"only_allow_merge_if_pipeline_succeeds"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool          `json:"only_allow_merge_if_all_discussions_are_resolved"`
	RemoveSourceBranchAfterMerge              bool          `json:"remove_source_branch_after_merge"`
	RequestAccessEnabled                      bool          `json:"request_access_enabled"`
	MergeMethod                               string        `json:"merge_method"`
	AutocloseReferencedIssues                 bool          `json:"autoclose_referenced_issues"`
	SuggestionCommitMessage                   interface{}   `json:"suggestion_commit_message"`
	Statistics                                struct {
		CommitCount      int `json:"commit_count"`
		StorageSize      int `json:"storage_size"`
		RepositorySize   int `json:"repository_size"`
		WikiSize         int `json:"wiki_size"`
		LfsObjectsSize   int `json:"lfs_objects_size"`
		JobArtifactsSize int `json:"job_artifacts_size"`
		PackagesSize     int `json:"packages_size"`
	} `json:"statistics"`
	Links struct {
		Self          string `json:"self"`
		Issues        string `json:"issues"`
		MergeRequests string `json:"merge_requests"`
		RepoBranches  string `json:"repo_branches"`
		Labels        string `json:"labels"`
		Events        string `json:"events"`
		Members       string `json:"members"`
	} `json:"_links"`
}

// Repository interface
type Repository interface {
	appcontext.Component
	// GetGitlab return []byte, link, err
	GetGitlab(gitlabURL, gitlabToken string) ([]byte, string, error)
	// GitlabPostComment post a comment in gitlab
	GitlabPostComment(url string, params map[string]string) error
}

// GetRepository func return SensuRepository interface
func GetRepository() Repository {
	return appcontext.Current.Get(appcontext.Repository).(Repository)
}
