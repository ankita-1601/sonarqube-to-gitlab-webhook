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
```