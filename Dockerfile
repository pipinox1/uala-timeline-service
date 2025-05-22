FROM golang:1.23-alpine AS build
RUN apk update && apk add git
RUN mkdir /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app

RUN go mod tidy
RUN go build -o /go/src/app/main cmd/app/main.go

FROM alpine

COPY --from=build /go/src/app/main /go/src/app/
COPY --from=build /go/src/app/config /go/src/app/config
RUN ls -la /go/src/app/config

EXPOSE 8080
ENTRYPOINT ["/go/src/app/main"]