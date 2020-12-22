# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2020-12-22
### Added
- new logger using zap
- add more tests
- Go test badge
- Coveralls badge

### Updated
- golang version to 1.15 and all packages

### Changed 
- Change All funcs to return error

## [0.0.1] - 2020-04-03
### Added
- Endpoint for sonarqube `/sonarqube-to-gitlab-webhook/v1/events`
- Integration with Gitlab using token
- Add tests in usecase, appcontext, controller, utils
- Add build script to control tag version and branch version
