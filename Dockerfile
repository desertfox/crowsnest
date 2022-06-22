# syntax=docker/dockerfile:1
FROM golang:1.17.5-alpine as builder

RUN mkdir /opt/crowsnest

WORKDIR /opt/crowsnest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o crowsnest .

FROM scratch

RUN mkdir /opt/crowsnest

WORKDIR /opt/crowsnest

COPY --from=builder /opt/crowsnest/crowsnest ./

CMD [ "./opt/crowsnestcrowsnest" ]
