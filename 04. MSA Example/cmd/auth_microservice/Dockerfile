FROM ubuntu:latest
MAINTAINER Park Jinhong (jinhong0719@naver.com)

RUN apt-get update -y && \
    apt-get upgrade -y && \
    apt-get install mysql-server git-core sudo -y && \
    git clone https://github.com/parkjinhong03/go-for-GO.git

ENV CMDPATH "/go-for-GO/04. MSA Example/cmd/"
ENV NATS localhost:4222

ENTRYPOINT  service mysql start && \
            cd "$CMDPATH" && \
            git pull && \
            ./auth_microservice/main
