#This multistage docker build is designed to build a small image
#with compilation happening on an intermediate image.
FROM golang:1.12.4 AS tmp

WORKDIR /opt/jfrog-stats/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o jfrog-stats -v main.go

FROM alpine:3.10

COPY --from=tmp /opt/jfrog-stats/jfrog-stats /opt/jfrog-stats/jfrog-stats
WORKDIR /opt/jfrog-stats/

CMD ["./jfrog-stats"]