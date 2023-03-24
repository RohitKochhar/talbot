# Stage 1: Builder
FROM golang:latest AS builder

WORKDIR /src

COPY ./go.mod /src/
COPY ./go.sum /src/

RUN go mod download

COPY ./cmd/api /src

# Run unit tests before building to stop build if tests fail
RUN go test -v

RUN CGO_ENABLED=0 GOOS=linux go build -o entrypoint


# Stage 2: Certs
FROM docker.io/library/alpine@sha256:686d8c9dfa6f3ccfc8230bc3178d23f84eeaf7e457f36f271ab1acc53015037c AS tools

RUN apk add --no-cache \
    ca-certificates

# Stage 3: Runner
FROM scratch

COPY --from=tools /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /src/entrypoint /

EXPOSE 4000

ENTRYPOINT [ "./entrypoint" ]
