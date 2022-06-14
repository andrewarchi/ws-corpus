FROM alpine AS builder

RUN apk add git make g++
RUN git clone https://github.com/wspace/marcellippmann-whitepp Whitepp
WORKDIR /Whitepp
RUN make

FROM scratch

COPY --from=builder /Whitepp/bin/White++ /
ENTRYPOINT ["/White++"]
