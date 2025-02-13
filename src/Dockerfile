FROM golang:1.24-alpine AS builder

RUN apk --update --no-cache add gcc musl-dev

WORKDIR /

COPY ./go.mod ./ 
COPY ./go.sum ./

RUN go mod download

COPY . ./

ARG VERSION

ENV VERSION="${VERSION}"
ENV CGO_ENABLED=1

RUN go install -ldflags="-s -w -X 'github.com/PavelMilanov/pgbackup/config.VERSION=${VERSION}'"

FROM alpine:3.20

ARG USER_DOCKER=pgbackup
ARG UID_DOCKER=10000

ENV USER_DOCKER="$USER_DOCKER"
ENV UID_DOCKER="$UID_DOCKER"

ENV TZ=Europe/Moscow
ENV GIN_MODE=release

WORKDIR /app

COPY --from=builder /go/bin/pgbackup /app/

COPY templates/ /app/templates/
COPY static/ /app/static/

RUN apk --update --no-cache add postgresql-client tzdata sqlite-libs curl gzip && \
    addgroup -g ${UID_DOCKER} ${USER_DOCKER} && \
    adduser -u ${UID_DOCKER} -G ${USER_DOCKER} -s /bin/sh -D -H ${USER_DOCKER} && \
    chown -R ${USER_DOCKER}:${USER_DOCKER} /app

EXPOSE 8080/tcp

HEALTHCHECK --interval=10m --timeout=3s --start-period=5s --retries=3 CMD curl -f http://localhost:8080/api/check || exit 1

ENTRYPOINT ["./pgbackup" ]

USER ${USER_DOCKER}:${USER_DOCKER}
