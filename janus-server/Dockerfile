# Multi-stage build untuk Janus Gateway - OPTIMIZED VERSION
# Stage 1: Build stage
FROM ubuntu:22.04 AS builder

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV JANUS_VERSION=1.2.0

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    curl \
    wget \
    cmake \
    meson \
    ninja-build \
    pkg-config \
    gengetopt \
    libtool \
    automake \
    autoconf \
    # Janus dependencies
    libmicrohttpd-dev \
    libjansson-dev \
    libssl-dev \
    libglib2.0-dev \
    libopus-dev \
    libogg-dev \
    libcurl4-openssl-dev \
    liblua5.3-dev \
    libconfig-dev \
    libavutil-dev \
    libavcodec-dev \
    libavformat-dev \
    libspeexdsp-dev \
    && rm -rf /var/lib/apt/lists/*

# Build libnice (required for ICE)
WORKDIR /tmp/libnice
RUN git clone --depth 1 --branch 0.1.21 https://gitlab.freedesktop.org/libnice/libnice.git . && \
    meson setup build --prefix=/usr --buildtype=release && \
    ninja -C build && \
    ninja -C build install && \
    ldconfig

# Build libsrtp (required for SRTP)
WORKDIR /tmp/libsrtp
RUN wget https://github.com/cisco/libsrtp/archive/v2.5.0.tar.gz && \
    tar xfv v2.5.0.tar.gz && \
    cd libsrtp-2.5.0 && \
    ./configure --prefix=/usr --enable-openssl && \
    make shared_library -j$(nproc) && \
    make install && \
    ldconfig

# Build usrsctp (for DataChannels)
WORKDIR /tmp/usrsctp
RUN git clone --depth 1 https://github.com/sctplab/usrsctp.git . && \
    ./bootstrap && \
    ./configure --prefix=/usr --disable-programs --disable-inet --disable-inet6 && \
    make -j$(nproc) && \
    make install && \
    ldconfig

# Build libwebsockets (for WebSocket support)
WORKDIR /tmp/libwebsockets
RUN git clone --depth 1 --branch v4.3-stable https://libwebsockets.org/repo/libwebsockets.git . && \
    mkdir build && cd build && \
    cmake -DLWS_MAX_SMP=1 \
          -DLWS_WITHOUT_EXTENSIONS=0 \
          -DCMAKE_INSTALL_PREFIX:PATH=/usr \
          -DCMAKE_C_FLAGS="-fpic" \
          -DCMAKE_BUILD_TYPE=Release .. && \
    make -j$(nproc) && \
    make install && \
    ldconfig

# Build Sofia-SIP (for SIP plugin)
WORKDIR /tmp/sofia-sip
RUN git clone --depth 1 --branch v1.13.16 https://github.com/freeswitch/sofia-sip.git . && \
    ./bootstrap.sh && \
    ./configure --prefix=/usr && \
    make -j$(nproc) && \
    make install && \
    ldconfig

# Clone and build Janus Gateway
WORKDIR /tmp/janus
RUN git clone --depth 1 --branch v${JANUS_VERSION} https://github.com/meetecho/janus-gateway.git . && \
    sh autogen.sh && \
    ./configure --prefix=/opt/janus \
        --enable-websockets \
        --enable-data-channels \
        --enable-post-processing \
        --enable-plugin-videoroom \
        --enable-plugin-audiobridge \
        --enable-plugin-recordplay \
        --enable-plugin-streaming \
        --enable-plugin-textroom \
        --enable-rest \
        --disable-rabbitmq \
        --disable-mqtt \
        --disable-unix-sockets && \
    make -j$(nproc) && \
    make install && \
    make configs

# Stage 2: Runtime stage
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# Install ONLY runtime dependencies (no -dev packages)
RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    # Runtime libraries
    libmicrohttpd12 \
    libjansson4 \
    libssl3 \
    libglib2.0-0 \
    libopus0 \
    libogg0 \
    libcurl4 \
    liblua5.3-0 \
    libconfig9 \
    libavutil56 \
    libavcodec58 \
    libavformat58 \
    libspeexdsp1 \
    && rm -rf /var/lib/apt/lists/*

# Copy built libraries from builder
COPY --from=builder /usr/lib/x86_64-linux-gnu/libnice.so* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libgio-*.so* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/libsrtp2.so* /usr/lib/
COPY --from=builder /usr/lib/libusrsctp.so* /usr/lib/
COPY --from=builder /usr/lib/libwebsockets.so* /usr/lib/
COPY --from=builder /usr/lib/libsofia-sip-ua.so* /usr/lib/

# Copy Janus installation
COPY --from=builder /opt/janus /opt/janus

# Update library cache
RUN ldconfig

# Create janus user and necessary directories
RUN groupadd -r janus && useradd -r -g janus janus && \
    mkdir -p /opt/janus/var/log/janus \
             /opt/janus/var/recordings \
             /opt/janus/share/janus/streams \
             /opt/janus/etc/janus \
             /var/log/janus && \
    chown -R janus:janus /opt/janus /var/log/janus

# Copy configuration files if exists
COPY --chown=janus:janus config/* /opt/janus/etc/janus/ 2>/dev/null || true

# Switch to janus user
USER janus

# Expose ports
# 8088: HTTP REST API
# 8188: Admin/Monitor HTTP
# 8989: WebSocket
# 7088: Admin WebSocket
# 10000-20000: RTP/RTCP
EXPOSE 8088 8188 8989 7088
EXPOSE 10000-20000/udp

WORKDIR /opt/janus

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=5 \
    CMD curl -f http://localhost:8088/janus/info || exit 1

# Start Janus
CMD ["/opt/janus/bin/janus", \
     "--stun-server=stun:stun.l.google.com:19302", \
     "--configs-folder=/opt/janus/etc/janus"]