FROM golang:1.4

MAINTAINER Alessandro Nadalin "alessandro.nadalin@gmail.com"

RUN go get github.com/codegangsta/gin
RUN go get github.com/stretchr/testify/assert
RUN go get golang.org/x/tools/cmd/godoc
RUN go get golang.org/x/crypto/ssh
RUN go get github.com/coreos/fleet/ssh
RUN go get github.com/mgutz/ansi
RUN go get github.com/kvz/logstreamer
RUN go get github.com/codegangsta/cli
RUN go get gopkg.in/yaml.v2
RUN go get github.com/mitchellh/gox
RUN gox -build-toolchain

COPY . /go/src/github.com/namshi/godo/
WORKDIR /go/src/github.com/namshi/godo/

CMD godo gox --output=build/{{.OS}}_{{.Arch}}/{{.Dir}}