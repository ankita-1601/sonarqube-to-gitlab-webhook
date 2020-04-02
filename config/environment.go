package config

import (
	"os"
)

var (
	// Port to be listened by application
	Port string
	// SonarqubeSecret secret to generate the header
	SonarqubeSecret string
	// GitlabToken string which contains token to talk with gitlab
	GitlabToken string
	// GitlabURL string which contains gitlab url
	GitlabURL string
)

// GetEnv func
func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func init() {
	Port = GetEnv("PORT", "9090")
	SonarqubeSecret = GetEnv("SONARQUBE_SECRET", "Absent")
	GitlabToken = GetEnv("GITLAB_TOKEN", "Absent")
	GitlabURL = GetEnv("GITLAB_URL", "Absent")
	// SonarqubeSecret = os.Getenv("SONARQUBE_SECRET")
	// if SonarqubeSecret == "" {
	// 	log.Fatal("variable SONARQUBE_SECRET not defined")
	// }
	// GitlabToken = os.Getenv("GITLAB_TOKEN")
	// if GitlabToken == "" {
	// 	log.Fatal("variable GITLAB_TOKEN not defined")
	// }
	// GitlabURL = os.Getenv("GITLAB_URL")
	// if GitlabURL == "" {
	// 	log.Fatal("variable GITLAB_URL not defined")
	// }
}
