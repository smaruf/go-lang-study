FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o flight-tracker .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/flight-tracker .
COPY routes.json .
COPY static/ static/
EXPOSE 8080
CMD ["./flight-tracker"]
