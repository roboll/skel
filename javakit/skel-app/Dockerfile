###############################################################################
# {{ .name}}-app
###############################################################################
FROM jeanblanchard/java:serverjre-8

RUN apk add --update bash && rm -rf /var/cache/apk/*

ADD ./{{ .name }}-app-*.tar /dist

RUN ln -s /dist/{{ .name }}-app-* /opt/app

EXPOSE 8080
ENTRYPOINT ["/opt/app/bin/{{ .name }}-app"]
