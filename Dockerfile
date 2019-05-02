FROM alpine:latest

WORKDIR /build

ADD build/nigiri-chopsticks-linux-amd64 /build/chopsticks

EXPOSE 3000

ENTRYPOINT ["/build/chopsticks"]