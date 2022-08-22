FROM golang:1.18-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY pkg ./pkg
ENV GOARCH=amd64
RUN go build -o /go/bin/api cmd/api/main.go


FROM alpine:latest  
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
EXPOSE 8080
ENV TZ=Asia/Bangkok
ENV APP_MODE=production

COPY --from=build /go/bin/api /app/api

CMD ["/app/api"]  