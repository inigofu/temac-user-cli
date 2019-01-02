FROM golang:1.11.4 as builder

WORKDIR /go/src/github.com/inigofu/temac-user-cli

COPY . .

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/inigofu/temac-user-cli .

ENTRYPOINT ["./temac-user-cli"]

CMD ["./temac-user-cli"]
