FROM golang:1.23-alpine AS builder

RUN apk --update --no-cache add gcc musl-dev

WORKDIR /build

COPY . /build

ARG VERSION

ENV VERSION="${VERSION}"
ENV CGO_ENABLED=1

RUN go mod tidy && go build -ldflags="-s -w -X 'config.VERSION=${VERSION}'"

FROM alpine:3.20

ARG USER_DOCKER=pgbackup
ARG UID_DOCKER=10000

ENV USER_DOCKER="$USER_DOCKER"
ENV UID_DOCKER="$UID_DOCKER"

ENV TZ=Europe/Moscow
ENV GIN_MODE=release

WORKDIR /app

COPY --from=builder /build/pgbackup /app

COPY templates/ /app/templates/
COPY static/ /app/static/

VOLUME [ "/app/dumps" ]
VOLUME [ "/app/data" ]

RUN apk --update --no-cache add postgresql-client tzdata sqlite-libs && \
    rm -rf /var/cache/apk/ && \
    addgroup -g ${UID_DOCKER} ${USER_DOCKER} && \
    adduser -u ${UID_DOCKER} -G ${USER_DOCKER} -s /bin/sh -D -H ${USER_DOCKER} && \
    chown -R ${USER_DOCKER}:${USER_DOCKER} /app


EXPOSE 8080/tcp

ENTRYPOINT ["./pgbackup" ]

USER ${USER_DOCKER}:${USER_DOCKER}
