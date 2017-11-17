FROM golang:1.9.2-alpine

ARG APP_VERSION

RUN mkdir -p /usr/share/streaming

WORKDIR /usr/share/streaming

COPY config/docker.toml ./config/config.toml

COPY ui ./ui

COPY data ./data

ADD build/streaming-$APP_VERSION /usr/share/streaming/streaming

EXPOSE 8001

ENTRYPOINT /usr/share/streaming/streaming -config=/usr/share/streaming/config/config.toml
