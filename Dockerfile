FROM golang:1.25.4-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN go test ./... -count=1

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app 

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
COPY --from=build /bin/app /bin/app

ENTRYPOINT ["/bin/app"]
