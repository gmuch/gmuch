FROM golang:1.5

ENV GO15VENDOREXPERIMENT=1 NOTMUCH_VERSION=0.20.2

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    libxapian-dev \
    libgmime-2.6-dev \
    libtalloc-dev \
    zlib1g-dev \
    python-sphinx \
    xz-utils \
	&& rm -rf /var/lib/apt/lists/*

RUN mkdir -p /usr/src \
  && curl -L notmuch.tar.gz https://github.com/notmuch/notmuch/archive/${NOTMUCH_VERSION}.tar.gz | tar -C /usr/src -xzf - \
  && cd /usr/src/notmuch-${NOTMUCH_VERSION} \
  && ./configure --prefix=/usr \
  && make install

COPY . /go/src/github.com/gmuch/gmuch
WORKDIR /go/src/github.com/gmuch/gmuch
RUN make test
RUN make install

VOLUME /mail
EXPOSE 8000
EXPOSE 8001
EXPOSE 8002
CMD ["gmuch", "-db-path", "/mail"]
