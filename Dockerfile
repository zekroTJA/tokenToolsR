FROM golang:1.12.6-stretch

LABEL maintainer="zekro <contact@zekro.de>"

RUN apt-get update -y &&\
    apt-get install -y \
        git

ENV PATH="${GOPATH}/bin:${PATH}"

ENV TLS="false"

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR ${GOPATH}/src/github.com/zekroTJA/tokenToolsR

ADD . .

RUN mkdir -p /etc/certs

RUN dep ensure -v

RUN go build \
        -v \
        -o /usr/bin/tt \
        -ldflags "\
            -X github.com/zekroTJA/tokenToolsR/internal/static.AppDate=$(date -u '+%Y-%m-%d_%I:%M:%S%p') \
            -X github.com/zekroTJA/tokenToolsR/internal/static.AppVersion=$(git describe --tags) \
            -X github.com/zekroTJA/tokenToolsR/internal/static.AppCommit=$(git rev-parse HEAD)" \
        ./cmd/tokentools/*.go

EXPOSE 8080

CMD tt \
        -port 8080 \
        -tls-cert /etc/cert/*.cer \
        -tls-key  /etc/cert/*.key \
        -tls="$TLS"