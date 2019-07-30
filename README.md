# Sonarqube to Gitlab Webhook

TravisCi [![TravisCI Build Status](https://travis-ci.org/betorvs/sonarqube-to-gitlab-webhook.svg?branch=master)](https://travis-ci.org/betorvs/sonarqube-to-gitlab-webhook)


## Go Installation

Install go

Install dependencies [dep](https://golang.github.io/dep/docs/installation.html)

### Configure

```sh
dep ensure -update -v
```

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

