# rm -f journey ; docker build -t tmp -f build.Dockerfile . && docker run -it --rm -v $(pwd):/mnt tmp cp journey /mnt/
FROM ubuntu:18.04
RUN apt update &&\
    apt install -y software-properties-common &&\
    add-apt-repository ppa:longsleep/golang-backports &&\
    apt update &&\
    apt install -y golang-go

WORKDIR /opt/journey
COPY . .
RUN go build -a -tags netgo -ldflags '-w' -o journey
