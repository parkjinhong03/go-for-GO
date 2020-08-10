FROM alpine
MAINTAINER Park, Jinhong <jinhong0719@naver.com>

ARG CONSUL_ADDRESS
COPY ./api-gateway ./api-gateway
ENTRYPOINT ["./api-gateway"]
ENV CONSUL_ADDRESS=$CONSUL_ADDRESS
