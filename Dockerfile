FROM node:20 AS node-builder
WORKDIR /app
COPY . .
ENV CI=true
ENV PNPM_CONFIG_CONFIRM=true
RUN corepack enable && pnpm install --frozen-lockfile && pnpm build

FROM golang:1.23 as go-builder
WORKDIR /app
# Copy go.mod and go.sum first for better layer caching
COPY --from=node-builder /app/go.mod /app/go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY --from=node-builder /app .
RUN go run github.com/a-h/templ/cmd/templ generate && go build -o tmp/main ./cmd/

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=go-builder /app/tmp/main .
COPY --from=go-builder /app/assets ./assets
COPY --from=go-builder /app/db ./db
EXPOSE 8080
CMD ["./main"]
