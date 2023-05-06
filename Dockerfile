# syntax=docker/dockerfile:1
FROM golang:1.19.4-alpine3.16

ARG CROWSNET_CONFIG=/opt/crowsnest/config.yaml

RUN mkdir /opt/src

WORKDIR /opt/src

COPY . .

RUN go mod download

ENV CGO_ENABLED=0

ENV GOOS=linux

RUN go build -o crowsnest .

CMD [ "./crowsnest config $CROWSNET_CONFIG" ]