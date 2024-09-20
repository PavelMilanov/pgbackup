FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY . /build

RUN go build .

FROM alpine

RUN apk --update --no-cache add postgresql-client

ARG USER_DOCKER
ARG UID_DOCKER
ARG GID_DOCKER

ENV USER_DOCKER="$USER_DOCKER"
ENV UID_DOCKER="$UID_DOCKER"
ENV GID_DOCKER="$GID_DOCKER"

WORKDIR /app

COPY --from=builder /build/pgbackup /app

COPY templates/ /app/templates/
COPY static/ /app/static/

RUN addgroup -g ${GID_DOCKER} ${USER_DOCKER} && \
    adduser -u ${UID_DOCKER} -G ${USER_DOCKER} -s /bin/sh -D -H ${USER_DOCKER} && \
    chown -R ${USER_DOCKER}:${USER_DOCKER} /app

VOLUME [ "/app/data" ]
VOLUME [ "/app/dumpls" ]

EXPOSE 8080

USER ${USER_DOCKER}:${USER_DOCKER}

ENTRYPOINT [ "./pgbackup" ]