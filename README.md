# Sonarqube to Gitlab Webhook

TravisCi [![TravisCI Build Status](https://travis-ci.org/betorvs/sonarqube-to-gitlab-webhook.svg?branch=master)](https://travis-ci.org/betorvs/sonarqube-to-gitlab-webhook)


## Go Installation

Install go


### Run

```sh
go build
./sonarqube-to-gitlab-webhook
```

## Environment Variables

You need to configure these for local tests or real deployment.

Configure these environment variables:
* **GITLAB_URL** : Gitlab URL. Example: https://gitlab.domain
* **SONARQUBE_SECRET** : Secret created in Webhook in Sonarqube configuration. Example: LONGHASH
* **GITLAB_TOKEN** : Gitlab Personal Token with api access.


## Deploy Kubernetes

```sh
kubectl create ns sonarqube-webhook
kubectl create secret generic sonarqube-webhook --from-literal=sonarqubeSecret=LONGHASH --from-literal=gitlabToken=xxx-9X-zxczxczxczxc -n sonarqube-webhook --dry-run -o yaml > sonarqube-secret.yaml
kubectl apply -f sonarqube-secret.yaml
kubeclt apply -f deployment.yaml
```

## Example Job and sonar config

.gitlab-ci.yml

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

sonar-project.properties

```
sonar.projectKey=projectGroup/projectName
sonar.analysis.disabledGitlabPost=false
sonar.analysis.disabledQualityReport=true
```

Use `sonar.analysis.disabledGitlabPost` equal `true` to disable post in GitLab.   
Use `sonar.analysis.disabledQualityReport` equal `true` to remove full quality report (it will print only Quality Gateway Name and Quality Gateway Status)


## Example of commit


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


