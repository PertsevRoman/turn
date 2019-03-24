FROM golang:latest as builder

ENV REALM localhost
ENV USERS username=password
ENV UDP_PORT 3478

RUN mkdir /build
WORKDIR /src

ADD . .

RUN go build -o /build/turn cmd/simple-turn/main.go
RUN go build -buildmode=plugin -o internal/plugins/auth/env/env.go /build/plugins/env.so
RUN go build -buildmode=plugin -o internal/plugins/auth/redis/redis.go /build/plugins/redis.so

FROM alpine

WORKDIR /app
COPY --from=builder /build/turn /app/turn
COPY --from=builder /build/plugins/env.so /app/plugins/env.so
COPY --from=builder /build/plugins/redis.so /app/plugins/redis.so
RUN chmod +x /app/turn

ENTRYPOINT [ "/app/turn" ]
