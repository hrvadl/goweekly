FROM golang:latest AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy && go get ./...
COPY ./ ./
RUN cd cmd/translator && CGO_ENABLED=0 GOOS=linux go build .

FROM alpine:latest AS final
WORKDIR /app
COPY --from=builder /app/cmd/translator/translator .
CMD [ "./translator" ]
