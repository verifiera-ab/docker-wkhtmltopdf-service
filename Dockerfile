FROM golang:1.14 as builder

WORKDIR /app

COPY /app/go.mod ./
COPY /app/go.sum ./

RUN go mod download

COPY /app/*.go ./

RUN go build -v -o app

FROM debian:jessie

MAINTAINER Potiguar Faga <potz@potz.me>

ENV WKHTML_MAJOR 0.12
ENV WKHTML_MINOR 6.1-2
ENV SECURE=false

# Builds the wkhtmltopdf download URL based on version numbers above
ENV DOWNLOAD_URL "https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/${WKHTML_MAJOR}.${WKHTML_MINOR}/wkhtmltox-${WKHTML_MAJOR}.${WKHTML_MINOR}_linux-jessie-amd64.deb"

# Create system user first so the User ID gets assigned
# consistently, regardless of dependencies added later
RUN useradd -rm appuser && \
    apt-get update && \
    apt-get install -y --no-install-recommends curl ca-certificates \
       fontconfig libfontconfig1 libfreetype6 \
       libpng12-0 libjpeg62-turbo \
       libssl1.0.0 libx11-6 libxext6 libxrender1 \
       xfonts-base xfonts-75dpi && \
    curl -L -o /tmp/wkhtmltox.deb $DOWNLOAD_URL && \
    dpkg -i /tmp/wkhtmltox.deb && \
    rm /tmp/wkhtmltox.deb && \
    apt-get purge -y curl && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir /app && \
    mkdir /app/ssl

COPY --from=builder /app/app /app/app

RUN chown -R appuser:appuser /home/appuser

VOLUME /app/ssl

USER appuser
WORKDIR /app
EXPOSE 3000

# whether to run under SSL or not
ENV SECURE true

CMD ["/app/app"]

