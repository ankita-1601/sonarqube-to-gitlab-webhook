# Sonarqube to Gitlab Webhook

![Go Test](https://github.com/github.com/betorvs/sonarqube-to-gitlab-webhook/sonarqube-to-gitlab-webhook/workflows/Go%20Test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/betorvs/sonarqube-to-gitlab-webhook/badge.svg?branch=master)](https://coveralls.io/github/betorvs/sonarqube-to-gitlab-webhook?branch=master)

## Environment variables

```sh
export PORT=9090

export APP_NAME=sonarqube-to-gitlab-webhook

export LOG_LEVEL=INFO
```

### Configuration

You need to configure these for local tests or real deployment.

Configure these environment variables:
* **GITLAB_URL** : Gitlab URL. Example: https://gitlab.domain
* **SONARQUBE_SECRET** : Secret created in Webhook in Sonarqube configuration. Example: LONGHASH
* **GITLAB_TOKEN** : Gitlab Personal Token with api access.


### Dependency Management
The project is using [Go Modules](https://blog.golang.org/using-go-modules) for dependency management
Module: github.com/betorvs/sonarqube-to-gitlab-webhook

### Test and coverage

Run the tests

```sh
TESTRUN=true go test ./... -coverprofile=cover.out

go tool cover -html=cover.out
```

Install [golangci-lint](https://github.com/golangci/golangci-lint#install) and run lint:

```sh
golangci-lint run
```

### Docker Build

```sh
docker build .
```

### Deploy Kubernetes

```sh
kubectl create ns sonarqube-webhook
kubectl create secret generic sonarqube-webhook --from-literal=sonarqubeSecret=LONGHASH --from-literal=gitlabToken=xxx-9X-zxczxczxczxc -n sonarqube-webhook --dry-run=client -o yaml > sonarqube-secret.yaml
kubectl apply -f sonarqube-secret.yaml
kubectl apply -f deployment.yaml
```

## Example Job and sonar config

Creates a step inside your `.gitlab-ci.yml`:

```
services:
  - docker:18.09.7-dind

stages:
- test

image:
  name: sonarsource/sonar-scanner-cli:latest
  entrypoint: [""]
variables:
  SONAR_HOST_URL: "https://sonar.example.local"
  GIT_DEPTH: 0
sonarqube-check:
  stage: test
  script:
    - sonar-scanner -Dsonar.qualitygate.wait=true
  allow_failure: true
```

### Extra options for sonar-project.properties

```
sonar.projectKey=projectGroup/projectName
sonar.analysis.disabledGitlabPost=false
sonar.analysis.disabledQualityReport=true
```

Use `sonar.analysis.disabledGitlabPost` equal `true` to disable post in GitLab.   
Use `sonar.analysis.disabledQualityReport` equal `true` to remove full quality report (it will print only Quality Gateway Name and Quality Gateway Status)


### Configure Project ID

Add this configuration in sonar-propject.properties:

```
sonar.analysis.projectID="10"
```


## Example of commit

Commit example in [gitlab.com](https://gitlab.com/betorvs/sonarqube-webhook-test/-/commit/278fccb1c6c68b9725e1a40315b39f0ebb3e93a3#note_472095580)


# SONARQUBE REPORT  
URL: [Report Link](https://sonar.example.com/dashboard?id=greatuser%2Ftest&branch=test)  
  
## Quality Gateway  
 - Name: TEST  
 - Status: OK  
### Quality Gateway Conditions  
#### Metric Name: reliability_rating  
 - Operator: GREATER_THAN  
 - Value: 1  
 - Error Threshold: 1  
 - Status: **OK** :+1:  
#### Metric Name: security_rating  
 - Operator: GREATER_THAN  
 - Value: 1  
 - Error Threshold: 1  
 - Status: **OK** :+1:  
#### Metric Name: maintainability_rating  
 - Operator: GREATER_THAN  
 - Value: 1  
 - Error Threshold: 1  
 - Status: **OK** :+1:  
#### Metric Name: coverage  
 - Operator: LESS_THAN  
 - Value: 0.0  
 - Error Threshold: 70  
 - Status: **OK** :+1:  
#### Metric Name: duplicated_lines_density  
 - Operator: GREATER_THAN  
 - Value: 0.0  
 - Error Threshold: 3  
 - Status: **OK** :+1:


## References

### Golang Spell
The project was initialized using [Golang Spell](https://github.com/danilovalente/golangspell).

### Architectural Model
The Architectural Model adopted to structure the application is based on The Clean Architecture.
Further details can be found here: [The Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html) and in the Clean Architecture Book.
