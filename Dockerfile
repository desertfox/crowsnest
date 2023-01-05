# syntax=docker/dockerfile:1
FROM golang:1.19.4-alpine3.16 as builder

RUN mkdir /opt/crowsnest

WORKDIR /opt/crowsnest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o crowsnest .

FROM scratch

WORKDIR /root

COPY --from=builder /opt/crowsnest/crowsnest ./

CMD [ "./crowsnest" ]