############################
# STEP 1 build executable binary
############################
FROM golang:alpine as builder
RUN apk update &&                               \
    apk add --no-cache ca-certificates git 
WORKDIR /src
COPY ./go.mod ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o app-cron ./cmd/crontab/main.go && \
    CGO_ENABLED=0 go build -o app-api ./cmd/api/main.go && \
    CGO_ENABLED=0 go build -o app-db ./cmd/db/main.go && \
    CGO_ENABLED=0 go build -o app-dbinit ./cmd/dbinit/main.go && \
    CGO_ENABLED=0 go build -o app-getaeskey ./cmd/getaeskey/main.go

############################
# STEP 2 build a small image
############################
FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/init/* ./init/
COPY --from=builder /src/app-* ./
CMD ["./app-cron"]