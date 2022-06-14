FROM alpine AS builder

RUN apk add git g++
RUN git clone https://github.com/noia1/Whitespace-Interpreter
WORKDIR /Whitespace-Interpreter
RUN g++ -O3 -Wall -Werror -static -o ws Lexer.cc Parser.cc

FROM scratch

COPY --from=builder /Whitespace-Interpreter/ws /
ENTRYPOINT ["/ws"]
