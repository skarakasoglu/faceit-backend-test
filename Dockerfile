FROM golang:1.19-alpine as builder

RUN apk update
RUN apk add make

ADD . /app/faceit-backend-test

WORKDIR /app/faceit-backend-test

#ENV GO111MODULE=on
#ENV GOFLAGS="-mod=vendor"

#RUN go mod vendor

#RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o faceit-backend-test ./cmd/main.go
RUN make build

FROM golang:1.19-alpine as runner
RUN apk add --no-cache ca-certificates bash
COPY --from=builder /app/faceit-backend-test/faceit-backend-test /usr/bin/faceit-backend-test

EXPOSE 8800

ENTRYPOINT ["/usr/bin/faceit-backend-test"]