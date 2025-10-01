FROM golang:1.25-alpine AS build

ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -ldflags="-s -w" -o ./alertmanager-webhook-signal


FROM scratch AS runtime

COPY --from=build /app/alertmanager-webhook-signal /alertmanager-webhook-signal

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER appuser:appuser

ENTRYPOINT ["/alertmanager-webhook-signal"]
