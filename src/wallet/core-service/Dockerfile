FROM golang:1.16-alpine

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV CGO_ENABLED 1
ENV GOFLAGS -mod=vendor
ENV GOOS=linux
ENV GOARCH=amd64

RUN export GO111MODULE=on

CMD ["core-service-app"]

RUN go build -t core-service-app .