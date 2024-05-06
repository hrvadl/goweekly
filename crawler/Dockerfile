FROM golang:latest AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY ./ ./
RUN cd cmd/crawler && CGO_ENABLED=0 GOOS=linux go build .

FROM alpine:latest AS final
WORKDIR /app
COPY --from=builder /app/cmd/crawler/crawler .
CMD [ "./crawler" ]
