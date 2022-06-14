FROM wspace-corpus/rust AS builder

RUN git clone https://github.com/wspace/zrneely-rust whitespace
WORKDIR /whitespace
RUN cargo test
RUN RUSTFLAGS='-C target-feature=+crt-static' cargo build --release --target x86_64-unknown-linux-gnu

FROM scratch

COPY --from=builder /whitespace/target/x86_64-unknown-linux-gnu/release/whitespace /
ENTRYPOINT ["/whitespace"]
