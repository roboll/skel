FROM jeanblanchard/java:serverjre-8

RUN apk add --update bash && rm -rf /var/cache/apk/*

ADD ./{{ .Repo }}-ws-*.tar /dist

RUN ln -s /dist/{{ .Repo }}-ws-* /opt/app

ENTRYPOINT ["/opt/app/bin/{{ .Repo }}-ws"]
