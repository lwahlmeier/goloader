FROM alpine:3.12

COPY run.sh /run.sh
RUN touch env.sh && chmod 755 /run.sh
COPY build/goloader /goloader

EXPOSE 8080/tcp

ENTRYPOINT ["/run.sh"]
CMD ["/goloader"]

