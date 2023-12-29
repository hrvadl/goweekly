FROM golang:latest as builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY ./ ./
RUN cd cmd/goweekly && CGO_ENABLED=0 GOOS=linux go build .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/cmd/goweekly/goweekly .
CMD [ "./goweekly" ]
