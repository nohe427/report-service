# Use base golang image from Docker Hub
FROM golang:1.12 as builder

WORKDIR .
RUN mkdir -p github.com/nohe427/report_service && cd github.com/nohe427/report_service
WORKDIR github.com/nohe427/report_service

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /main

FROM alpine
RUN apk add ca-certificates && update-ca-certificates
COPY --from=builder /main /main
CMD ["/main"]