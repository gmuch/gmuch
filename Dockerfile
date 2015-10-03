FROM golang:1.5

ENV NOTMUCH_VERSION=0.20 GOAPP=github.com/gmuch/gmuch GO15VENDOREXPERIMENT=1

RUN apt-get update \
  && apt-get install -y libxapian-dev libgmime-2.6-dev libtalloc-dev \
    zlib1g-dev python-sphinx --no-install-recommends \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

RUN curl -L http://notmuchmail.org/releases/notmuch-${NOTMUCH_VERSION}.tar.gz | tar -C /tmp -xzf- \
  && cd /tmp/notmuch-${NOTMUCH_VERSION} \
  && ./configure --prefix=/usr \
  && make install \
  && cd \
  && rm -rf /tmp/notmuch-${NOTMUCH_VERSION}

EXPOSE 8000
EXPOSE 8001
EXPOSE 8002

ADD . /go/src/${GOAPP}
WORKDIR /go/src/${GOAPP}

RUN go install ./cmd/...
