FROM golang:1.19-buster AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY pkg ./pkg
ENV GOARCH=amd64
RUN go build -o /go/bin/app cmd/main.go

# Deploy
FROM gcr.io/distroless/base-debian11
COPY --from=build /go/bin/app /app
EXPOSE 3000
USER nonroot:nonroot
CMD ["/app"]
