FROM golang:latest as builder

ENV REALM localhost
ENV USERS ""
ENV DB_DSN ""
ENV PORT 3478

WORKDIR /src

ADD . .

RUN go build -o app/turn cmd/simple-turn/main.go
RUN go build -buildmode=plugin -o app/plugins/env.so internal/plugins/auth/env/env.go
RUN go build -buildmode=plugin -o app/plugins/redis.so internal/plugins/auth/redis/redis.go

FROM alpine

WORKDIR /app
RUN mkdir /app/plugins
COPY --from=builder /src/entrypoint.sh /app/entrypoint.sh
COPY --from=builder /src/app/turn /app/turn
COPY --from=builder /src/app/plugins/env.so /app/plugins/env.so
COPY --from=builder /src/app/plugins/redis.so /app/plugins/redis.so

RUN chmod +x /app/turn
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT /app/entrypoint.sh
