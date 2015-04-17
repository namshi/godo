FROM golang

MAINTAINER Alessandro Nadalin "alessandro.nadalin@gmail.com"

RUN go get github.com/codegangsta/gin
RUN go get github.com/stretchr/testify/assert
RUN go get golang.org/x/tools/cmd/godoc
RUN go get golang.org/x/crypto/ssh
RUN go get github.com/coreos/fleet/ssh
RUN go get github.com/mgutz/ansi
RUN go get github.com/kvz/logstreamer/src/pkg/logstreamer/
RUN go get github.com/codegangsta/cli
RUN go get gopkg.in/yaml.v2

COPY . /src
WORKDIR /src

CMD go build -o godo main.go