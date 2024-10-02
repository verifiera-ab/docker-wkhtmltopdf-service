FROM golang:1.23.2-bullseye as builder

WORKDIR /app

COPY /app/go.mod ./
COPY /app/go.sum ./

RUN go mod download

COPY /app/*.go ./

RUN go build -v -o app

FROM ubuntu:22.04

MAINTAINER Oleg Lemesenko <oleg@verifiera.se>

RUN apt-get update
RUN apt-get install -y --no-install-recommends wkhtmltopdf

RUN mkdir /app && \
    mkdir /app/ssl



COPY --from=builder /app/app /app/app

RUN useradd -rm appuser
RUN chown -R appuser:appuser /home/appuser

VOLUME /app/ssl

USER appuser
WORKDIR /app
EXPOSE 3000

# whether to run under SSL or not
ENV SECURE false

CMD ["/app/app"]

