FROM golang:1.10.8-stretch

WORKDIR /go/src/github.com/catherinetcai/gsuite-aws-sso
RUN mkdir /release
COPY . .
RUN dep ensure -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /release/gsuite-aws-sso

FROM scratch
COPY --from=builder release/gsuite-aws-sso /bin/gsuite-aws-sso
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/ca-bundle.pem

ENV PATH=/bin
ENV TMPDIR=/
ENTRYPOINT ["/bin/gsuite-aws-sso"]
