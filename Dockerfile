FROM golang:1.23-alpine AS builder

RUN apk --update --no-cache add gcc musl-dev

WORKDIR /build

COPY . /build

ENV CGO_ENABLED=1

RUN go mod tidy && go build .

FROM alpine:3.20

ARG USER_DOCKER=pgbackup
ARG UID_DOCKER=10000

ENV USER_DOCKER="$USER_DOCKER"
ENV UID_DOCKER="$UID_DOCKER"

ENV TZ=Europe/Moscow

WORKDIR /app

COPY --from=builder /build/pgbackup /app

COPY templates/ /app/templates/
COPY static/ /app/static/

RUN apk --update --no-cache add postgresql-client tzdata sqlite-libs && \
    addgroup -g ${UID_DOCKER} ${USER_DOCKER} && \
    adduser -u ${UID_DOCKER} -G ${USER_DOCKER} -s /bin/sh -D -H ${USER_DOCKER} && \
    chown -R ${USER_DOCKER}:${USER_DOCKER} /app

VOLUME [ "/app/dumps" ]
VOLUME [ "/app/data" ]

EXPOSE 8080

USER ${USER_DOCKER}:${USER_DOCKER}

ENTRYPOINT ["./pgbackup" ]
