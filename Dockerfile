FROM golang:1.17.5-alpine as builder

WORKDIR /opt/crowsnest

COPY ./ ./

RUN go build -o crowsnest .

FROM golang:1.17.5-alpine

WORKDIR /root/

COPY --from=builder /opt/crowsnest/crowsnest ./

CMD [ "./crowsnest" ]