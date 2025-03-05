# syntax=docker/dockerfile:1

FROM alpine:3.21.3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c AS base

FROM golang:1.24.0-alpine3.21@sha256:2d40d4fc278dad38be0777d5e2a88a2c6dee51b0b29c97a764fc6c6a11ca893c AS build
LABEL org.opencontainers.image.source=https://gitea.dreger.lan/sdreger/lib-manager-go
RUN apk add git # required to get build information to be injected with '-buildvcs=true'

ADD go.mod go.sum /lib-manager-api/
WORKDIR /lib-manager-api
RUN go mod download
ADD . /lib-manager-api
RUN CGO_ENABLED=0 go build -buildvcs=true -o app github.com/sdreger/lib-manager-go/cmd/api

FROM base AS final

ARG UID=1001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser
COPY --from=build /lib-manager-api/app /bin/
ENTRYPOINT [ "/bin/app" ]
