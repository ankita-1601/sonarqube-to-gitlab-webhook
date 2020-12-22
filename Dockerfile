FROM golang:1.15.0-alpine3.12 AS golang

RUN apk add --no-cache git
RUN mkdir -p /builds/go/src/github.com/betorvs/sonarqube-to-gitlab-webhook/
ENV GOPATH /go
COPY . /builds/go/src/github.com/betorvs/sonarqube-to-gitlab-webhook/
ENV CGO_ENABLED 0
RUN cd /builds/go/src/github.com/betorvs/sonarqube-to-gitlab-webhook/ && TESTRUN=true go test ./... && go build

FROM alpine:3.12
WORKDIR /
VOLUME /tmp
RUN apk add --no-cache ca-certificates
RUN update-ca-certificates
RUN mkdir -p /app
RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -G app -S -D -h /app app && \
    chmod 755 /app
COPY --from=golang /builds/go/src/github.com/betorvs/sonarqube-to-gitlab-webhook/sonarqube-to-gitlab-webhook /app

EXPOSE 9090
RUN chmod +x /app/sonarqube-to-gitlab-webhook
WORKDIR /app    
USER app
CMD ["/app/sonarqube-to-gitlab-webhook"]
