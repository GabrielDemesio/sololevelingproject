# Build stage
FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build     CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

# Runtime stage
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /out/api /api
ENV PORT=8080
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/api"]