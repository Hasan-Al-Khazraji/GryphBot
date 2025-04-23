# docker build -t gryphbot:latest .
# docker run -e DISCORD_TOKEN=... -e GEMINI_API_KEY=... gryphbot:latest

FROM --platform=$BUILDPLATFORM golang:1.24-bullseye AS builder

ENV CGO_ENABLED=0 \
    GOFLAGS="-buildvcs=false"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -tags netgo -ldflags="-s -w" -o /gryphbot ./main.go

FROM scratch AS final

ENV DISCORD_TOKEN= \
    GEMINI_API_KEY=

COPY --from=builder /gryphbot /gryphbot

EXPOSE 443

ENTRYPOINT ["/gryphbot"]
