FROM golang:1.10.4 as builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64
RUN chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/robinjoseph08/hello

COPY Gopkg.toml Gopkg.toml
COPY Gopkg.lock Gopkg.lock
RUN dep ensure -vendor-only

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -ldflags '-w -s' -o ./bin/hello ./cmd/hello && strip ./bin/hello

FROM alpine:3.8

RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/robinjoseph08/hello/bin /bin

CMD ["hello"]
