FROM golang:latest as builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY ./ ./
RUN cd cmd/translator && CGO_ENABLED=0 GOOS=linux go build .

FROM scratch
WORKDIR /app
COPY --from=builder /app/cmd/translator/translator .
CMD [ "./translator" ]
