FROM golang:1.18 as builder

WORKDIR /
RUN git clone https://github.com/simomu-github/whitespace_go
WORKDIR /whitespace_go
RUN go test ./...
RUN go build -o releases/ws cmd/ws.go

FROM scratch as runner

COPY --from=builder /whitespace_go/releases/ws /
ENTRYPOINT ["/ws"]
