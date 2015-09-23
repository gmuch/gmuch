FROM golang:1.5

ENV NOTMUCH_VERSION=0.20 GOAPP=github.com/gmuch/gmuch GO15VENDOREXPERIMENT=1

RUN go get github.com/codegangsta/gin \
  && go get github.com/Masterminds/glide

RUN apt-get update \
  && apt-get install -y libxapian-dev libgmime-2.6-dev libtalloc-dev \
    zlib1g-dev python-sphinx --no-install-recommends

RUN curl -L http://notmuchmail.org/releases/notmuch-${NOTMUCH_VERSION}.tar.gz | tar -C /tmp -xzf- \
  && cd /tmp/notmuch-${NOTMUCH_VERSION} \
  && ./configure --prefix=/usr \
  && make install \
  && cd \
  && rm -rf /tmp/notmuch-${NOTMUCH_VERSION}

ADD . /go/src/${GOAPP}
WORKDIR /go/src/${GOAPP}

RUN go install $(glide nv)
