FROM golang:1.24-bookworm AS builder
RUN apt-get update && apt-get install -y --no-install-recommends \
        gcc libwebp-dev \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive \
    DOCKERMODE=true \
    DISPLAY=:99 \
    SERVER_PORT=8000 \
    POOL_SIZE=2 \
    BROWSER_PATH=/usr/bin/google-chrome \
    CF=http://localhost:8000 \
    MANGABALL_PROXY_POOL_FILE=/app/p.txt

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates libwebp7 wget \
        curl \
        python3 python3-pip \
        xvfb \
        fonts-liberation \
        libasound2t64 \
        libatk-bridge2.0-0 \
        libatk1.0-0 \
        libcups2 \
        libdbus-1-3 \
        libgdk-pixbuf2.0-0 \
        libnspr4 \
        libnss3 \
        libx11-xcb1 \
        libxcomposite1 \
        libxdamage1 \
        libxfixes3 \
        libxrandr2 \
        libxss1 \
        libxtst6 \
        xdg-utils \
    && wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
    && apt-get install -y ./google-chrome-stable_current_amd64.deb \
    && rm google-chrome-stable_current_amd64.deb \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY api/py/requirements.txt ./api/py/
RUN pip3 install --no-cache-dir --break-system-packages -r ./api/py/requirements.txt

COPY api/py/ ./api/py/
COPY p.txt /app/p.txt
COPY --from=builder /app/main .
COPY start.sh .
RUN chmod +x start.sh

EXPOSE 3400 3500 8000

CMD ["./start.sh"]