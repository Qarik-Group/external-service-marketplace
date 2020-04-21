FROM scratch
COPY ./esmd /go/bin/esmd
COPY ./esmd.yml /go/bin/esmd.yml
EXPOSE 8080
CMD ["/go/bin/esmd", " -c esmd.yml"]
