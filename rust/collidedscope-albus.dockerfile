FROM wspace-corpus/rust AS builder

RUN git clone https://github.com/collidedscope/albus
WORKDIR /albus
RUN RUSTFLAGS='-C target-feature=+crt-static' cargo build --release --target x86_64-unknown-linux-gnu

FROM scratch

COPY --from=builder /albus/target/x86_64-unknown-linux-gnu/release/albus /
ENTRYPOINT ["/albus"]
