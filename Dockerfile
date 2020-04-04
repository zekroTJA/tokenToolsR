FROM golang:1.14-alpine AS build
LABEL maintainer="zekro <contact@zekro.de>"

WORKDIR /build

RUN apk add --update \
        nodejs npm git

ADD . .

RUN cd web &&\
    npm ci &&\
    npm run build

RUN go build \
        -v \
        -o tt \
        -ldflags "\
            -X github.com/zekroTJA/tokenToolsR/internal/static.AppDate=$(date -u '+%Y-%m-%d_%I:%M:%S%p') \
            -X github.com/zekroTJA/tokenToolsR/internal/static.AppVersion=$(git describe --tags --abbrev=0) \
            -X github.com/zekroTJA/tokenToolsR/internal/static.AppCommit=$(git rev-parse --short HEAD) \
            -X github.com/zekroTJA/tokenToolsR/internal/static.AppProd=TRUE" \
        ./cmd/tokentools/*.go


FROM alpine:latest AS final

COPY --from=build /build/tt /app/tt
COPY --from=build /build/web/build /app/web

RUN mkdir -p /etc/certs
RUN chmod +x /app/tt

ENV TLS="false"

EXPOSE 8080

ENTRYPOINT ["/app/tt", "-web", "/app/web"]

CMD ["-addr", "localhost:8080"]
