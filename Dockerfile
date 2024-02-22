#syntax=docker/dockerfile:1

ARG GO_VERSION=1.22
ARG RABBITMQ_VERSION=3.12 

FROM golang:${GO_VERSION}-alpine as base
WORKDIR /app
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

FROM base as cryptoquotes-builder
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    go build -o /bin/cryptoquotes cmd/cryptoquotes/main.go

FROM scratch as cryptoquotes
COPY --from=cryptoquotes-builder /bin/cryptoquotes /bin/cryptoquotes
ARG API_PORT=3000
ENV API_PORT=${API_PORT}
EXPOSE ${API_PORT}
CMD [ "/bin/cryptoquotes" ]

FROM base as attractor-builder
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    go build -o /bin/attractor cmd/attractor/main.go

FROM scratch as attractor
COPY --from=attractor-builder /bin/attractor /bin/attractor
CMD [ "/bin/attractor" ]
