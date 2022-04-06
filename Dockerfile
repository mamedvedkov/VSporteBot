FROM golang:1.18-alpine AS builder
LABEL stage=builder

ENV CGO_ENABLED 0

ENV TZ=Europe/Moscow

RUN apk --no-cache add ca-certificates tzdata && \
    cp -r -f /usr/share/zoneinfo/$TZ /etc/localtime

WORKDIR /app

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./vendor ./vendor
COPY ./go.mod ./
COPY ./go.sum ./
COPY ./credentials.json /
COPY ./inn.json /
COPY ./compensation.pdf /
COPY ./szInstruction.pdf /

RUN go build -mod=vendor -o /vsportebot ./cmd/vsportebot

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/localtime /etc/localtime
COPY --from=builder /vsportebot /vsportebot
COPY --from=builder /credentials.json /credentials.json
COPY --from=builder /inn.json /inn.json
COPY --from=builder /compensation.pdf /compensation.pdf
COPY --from=builder /szInstruction.pdf /szInstruction.pdf

ENTRYPOINT ["/vsportebot"]
