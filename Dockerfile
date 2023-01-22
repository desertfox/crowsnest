# syntax=docker/dockerfile:1
FROM golang:1.19.4-alpine3.16

RUN mkdir /opt/src

WORKDIR /opt/src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o crowsnest .

ENV CROWSNET_CONFIG="./config.yaml"

CMD [ "./crowsnest", "--config", "${CROWSNET_CONFIG}" ]