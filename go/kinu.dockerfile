FROM golang:1.18 as builder

WORKDIR /
RUN git clone https://github.com/kinu/whitespace
WORKDIR /whitespace
RUN go mod init github.com/kinu/whitespace
RUN go build

FROM scratch as runner

COPY --from=builder /whitespace/whitespace /
ENTRYPOINT ["/whitespace"]
