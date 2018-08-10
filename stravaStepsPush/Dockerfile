FROM golang:1.8 as builder
WORKDIR /go/src/github.com/andjos/strava-api/tree/master/stravaStepsPush
RUN go get -d -v github.com/gin-gonic/gin
RUN go get -d -v github.com/strava/go.strava
RUN go get -d -v github.com/xconstruct/go-pushbullet

COPY main.go  .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/andjos/strava-api/tree/master/stravaStepsPush .
CMD ["./main"]

