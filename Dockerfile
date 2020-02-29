# Start by building the application.
FROM golang:buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN bash scripts/install

## This beacuse of net package use dynamic linking libc
## https://stackoverflow.com/a/36308464/4567832 
RUN go build -tags netgo -a -o /go/bin/app

# Now copy it into our base image.
FROM alpine:latest
COPY --from=build /go/bin/app /
CMD ["/app"]