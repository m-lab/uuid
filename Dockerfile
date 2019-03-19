FROM golang:1.12-alpine as uuid-creator-build
ADD . /go/src/github.com/m-lab/uuid/
RUN apk add git
RUN go get -v github.com/m-lab/uuid/cmd/create-uuid-prefix-file

FROM alpine
COPY --from=uuid-creator-build /go/bin/create-uuid-prefix-file /
WORKDIR /
ENTRYPOINT ["/create-uuid-prefix-file"]
