# syntax=docker/dockerfile:1
FROM golang:1.19.4-alpine3.16 as builder

RUN mkdir /opt/src

WORKDIR /opt/src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o crowsnest .

FROM scratch

WORKDIR /opt/src

COPY --from=builder /opt/src/crowsnest ./

CMD [ "./crowsnest" ]