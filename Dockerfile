FROM golang:1.21-alpine AS go-builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY backend/ ./backend/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./backend/cmd/server

FROM node:18-alpine AS node-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM alpine:latest
RUN apk add --no-cache sqlite
WORKDIR /app
COPY --from=go-builder /app/main ./main
COPY --from=node-builder /app/build/* ./static/
RUN mkdir -p /app/data

EXPOSE 8080
ENV DB_PATH=/app/data/anime.db
ENV STATIC_PATH=/app/static
ENV PORT=8080

CMD ["./main"] 