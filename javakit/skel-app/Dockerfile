FROM jeanblanchard/java:serverjre-8

RUN apk add --update bash && rm -rf /var/cache/apk/*

ADD ./{{ .Repo }}-app-*.tar /dist

RUN ln -s /dist/{{ .Repo }}-app-* /opt/app

ENTRYPOINT ["/opt/app/bin/{{ .Repo }}-app"]
