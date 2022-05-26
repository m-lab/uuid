FROM golang:1.17-alpine as uuid-creator-build
RUN apk add --no-cache git
COPY . /go/src/github.com/m-lab/uuid/
RUN go install -v github.com/m-lab/uuid/cmd/create-uuid-prefix-file

FROM alpine:3.16
COPY --from=uuid-creator-build /go/bin/create-uuid-prefix-file /
WORKDIR /
ENTRYPOINT ["/create-uuid-prefix-file"]
