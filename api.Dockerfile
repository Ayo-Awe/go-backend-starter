FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT_SHA
ARG BUILD_TIME

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-X 'github.com/ayo-awe/go-backend-starter/internal/vcs.version=${VERSION}' -X 'github.com/ayo-awe/go-backend-starter/internal/vcs.commitSHA=${COMMIT_SHA}' -X 'github.com/ayo-awe/go-backend-starter/internal/vcs.buildTime=${BUILD_TIME}'" -o api ./cmd/server

FROM --platform=$TARGETPLATFORM gcr.io/distroless/static-debian11

WORKDIR /root/

COPY --from=builder /app/api .

EXPOSE 4000

CMD ["./api"]