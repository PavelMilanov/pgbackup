FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY . /build

RUN go build .

FROM alpine

ARG USER_DOCKER
ARG UID_DOCKER

ENV USER_DOCKER="$USER_DOCKER"
ENV UID_DOCKER="$UID_DOCKER"
ENV TZ=Europe/Moscow

WORKDIR /app

COPY --from=builder /build/pgbackup /app

COPY templates/ /app/templates/
COPY static/ /app/static/

RUN apk --update --no-cache add postgresql-client tzdata && \
    addgroup -g ${UID_DOCKER} ${USER_DOCKER} && \
    adduser -u ${UID_DOCKER} -G ${USER_DOCKER} -s /bin/sh -D -H ${USER_DOCKER} && \
    chown -R ${USER_DOCKER}:${USER_DOCKER} /app

VOLUME [ "/app/data" ]
VOLUME [ "/app/dumps" ]

EXPOSE 8080

USER ${USER_DOCKER}:${USER_DOCKER}

ENTRYPOINT [ "./pgbackup" ]
