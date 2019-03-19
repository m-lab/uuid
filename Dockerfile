FROM golang:1.12-alpine as uuid-creator-build
RUN apk add --no-cache git
COPY . /go/src/github.com/m-lab/uuid/
RUN go get -v github.com/m-lab/uuid/cmd/create-uuid-prefix-file

FROM alpine:3.9
COPY --from=uuid-creator-build /go/bin/create-uuid-prefix-file /
WORKDIR /
ENTRYPOINT ["/create-uuid-prefix-file"]
