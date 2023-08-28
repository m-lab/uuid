FROM golang:1.20 as uuid-creator-build
COPY . /go/src/github.com/m-lab/uuid/
RUN cd /go/src/github.com/m-lab/uuid && CGO_ENABLED=0 go install -v ./cmd/create-uuid-prefix-file

FROM alpine:3.16
COPY --from=uuid-creator-build /go/bin/create-uuid-prefix-file /
WORKDIR /
# Make sure /create-uuid-prefix-file can run (has no missing external dependencies).
RUN /create-uuid-prefix-file -h 2> /dev/null
ENTRYPOINT ["/create-uuid-prefix-file"]
