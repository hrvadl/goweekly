FROM golang:latest AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY ./ ./
RUN cd cmd/core && CGO_ENABLED=0 GOOS=linux go build .

FROM alpine:latest AS final
WORKDIR /app
COPY --from=builder /app/cmd/core/core .
CMD [ "./core" ]
