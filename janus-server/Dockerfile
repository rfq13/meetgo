# Multi-stage build untuk optimasi ukuran image dan penggunaan RAM
# Stage 1: Build stage
FROM ubuntu:20.04 AS builder

# Set environment variables to avoid interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install build dependencies
RUN apt-get update && apt-get install -y \
    wget \
    gnupg \
    curl \
    build-essential \
    pkg-config \
    libmicrohttpd-dev \
    libjansson-dev \
    libnice-dev \
    libssl-dev \
    libsrtp2-dev \
    libsofia-sip-ua-dev \
    libopus-dev \
    libspeexdsp-dev \
    libogg-dev \
    libvpx-dev \
    libcurl4-openssl-dev \
    liblua5.3-dev \
    libconfig-dev \
    python3-dev \
    libavformat-dev \
    libavcodec-dev \
    libavutil-dev \
    libswscale-dev \
    libavfilter-dev \
    librabbitmq-dev \
    libnanomsg-dev \
    libevent-dev \
    libglib2.0-dev \
    libgtk-3-dev \
    libwebsockets-dev \
    libsctp-dev \
    libusrsctp-dev \
    cmake \
    git \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /opt

# Clone dan build Janus WebRTC Server dengan optimasi RAM
RUN git clone --depth 1 --branch v1.2.0 https://github.com/meetecho/janus-gateway.git && \
    cd janus-gateway && \
    git checkout v1.2.0 && \
    sh autogen.sh && \
    ./configure --prefix=/opt/janus --enable-websockets --enable-data-channels --enable-plugin-videoroom --enable-plugin-audiobridge --enable-plugin-recordplay --enable-plugin-streaming && \
    make -j$(nproc) && \
    make install && \
    make configs

# Stage 2: Runtime stage
FROM ubuntu:22.04 AS runtime

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV JANUS_CONF=/opt/janus/etc/janus

# Install only runtime dependencies
RUN apt-get update && apt-get install -y \
    libmicrohttpd12 \
    libjansson4 \
    libnice10 \
    libssl3 \
    libsrtp2-1 \
    libsofia-sip-ua0 \
    libopus0 \
    libspeexdsp1 \
    libogg0 \
    libvpx6 \
    libcurl4 \
    liblua5.3-0 \
    libconfig9 \
    python3 \
    libavformat58 \
    libavcodec58 \
    libavutil56 \
    libswscale5 \
    libavfilter7 \
    librabbitmq4 \
    libnanomsg5 \
    libevent-2.1-7 \
    libglib2.0-0 \
    libgtk-3-0 \
    libwebsockets16 \
    libsctp2 \
    libusrsctp1 \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Create janus user
RUN groupadd -r janus && useradd -r -g janus janus

# Copy hasil build dari builder stage
COPY --from=builder /opt/janus /opt/janus

# Create necessary directories
RUN mkdir -p /opt/janus/share/janus/recordings \
    && mkdir -p /opt/janus/share/janus/demos \
    && mkdir -p /var/log/janus

# Copy configuration files
COPY config/ /opt/janus/etc/janus/

# Change ownership
RUN chown -R janus:janus /opt/janus && \
    chown -R janus:janus /var/log/janus

# Switch to janus user
USER janus

# Expose ports
# 8088: HTTP/HTTPS interface
# 8188: Admin/monitor interface
# 10000-20000: RTP/RTCP range
EXPOSE 8088 8188 10000-20000/udp

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:8088/janus/info || exit 1

# Start Janus dengan WebSocket support
CMD ["/opt/janus/bin/janus", "--stun-server=stun:stun.l.google.com:19302", "--configs-folder=/opt/janus/etc/janus"]