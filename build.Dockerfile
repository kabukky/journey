# rm -f journey ; docker build -t tmp -f build.Dockerfile . && docker run -it --rm -v $(pwd):/mnt tmp cp journey /mnt/
# rm -f journey ; docker build --platform linux/amd64 -t tmp -f build.Dockerfile . && docker run --platform linux/amd64 -it --rm -v $(pwd):/mnt tmp cp journey /mnt/
FROM golang:1.15

WORKDIR /opt/journey
COPY . .
RUN go test ./...
RUN go build -a -tags "noplugins nossl netgo" -ldflags '-w' -o journey
