FROM golang:1.13.5

WORKDIR /data/soter/order

ADD . /data/soter/order

EXPOSE 8301

ENTRYPOINT  ["./soter-order-service"]

