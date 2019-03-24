FROM golang:latest as builder

ENV REALM localhost
ENV USERS username=password
ENV UDP_PORT 3478

WORKDIR /src

ADD . .

RUN go build -o app/turn cmd/simple-turn/main.go
RUN go build -buildmode=plugin -o internal/plugins/auth/env/env.go app/plugins/env.so
RUN go build -buildmode=plugin -o internal/plugins/auth/redis/redis.go app/plugins/redis.so

FROM alpine

WORKDIR /app
RUN mkdir plugins
COPY --from=builder /src/app/turn /app/turn
COPY --from=builder /build/app/plugins/env.so /app/plugins/env.so
COPY --from=builder /build/app/plugins/redis.so /app/plugins/redis.so
RUN chmod +x /app/turn

ENTRYPOINT [ "/app/turn" ]
