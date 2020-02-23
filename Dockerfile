FROM golang:latest

WORKDIR /go
ADD . /go

# app environment
ENV SLACK_API_TOKEN=YOUR_SLACK_API_LEGACYTOKEN
ENV GOPATH=

CMD ["make" , "init"]

ENTRYPOINT go run ./cmd/slacktheft
