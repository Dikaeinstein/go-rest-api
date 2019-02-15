# Build binary from official Go image
FROM golang:1.11.4-stretch AS build
ARG app=/go/src/github.com/dikaeinstein/go-rest-api
COPY . /${app}
WORKDIR /${app}
RUN go get -u github.com/kardianos/govendor && govendor sync
RUN go build -o /go-rest-api .

# Put the binary onto Heroku image
FROM heroku/heroku:16
COPY --from=build /go-rest-api /go-rest-api
RUN useradd -m myuser
USER myuser
CMD ["/go-rest-api"]
