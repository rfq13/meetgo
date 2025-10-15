# Janus Server Production Configuration

This document outlines the necessary steps to configure the Janus WebRTC server for a secure production environment. The primary focus is on enabling HTTPS for the admin and API interfaces.

## 1. Obtain SSL/TLS Certificates

For a production environment, you must use valid SSL/TLS certificates. You can obtain them from a Certificate Authority (CA) like Let's Encrypt. You will need the following files:

*   `fullchain.pem`: Your full certificate chain.
*   `privkey.pem`: Your private key.

## 2. Configure Janus for HTTPS

To enable HTTPS, you need to update the `janus.jcfg` configuration file to specify the paths to your SSL/TLS certificates.

1.  **Place Certificates:** Copy your `fullchain.pem` and `privkey.pem` files to a secure location accessible by the Janus container. A good practice is to create a new `certs` directory inside `janus-server` and mount it into the container.

2.  **Update `janus.jcfg`:** Uncomment the `https_port` and add the paths to your certificate and private key in the `[http]` section of `janus-server/config/janus.jcfg`:

    ```ini
    # HTTP transport settings
    http: {
        # HTTP server port
        http_port = 8088

        # HTTPS server port
        https_port = 7889

        # ... other settings

        # Path to certificates
        https_certs = {
            fullchain = "/path/to/your/fullchain.pem"
            privkey = "/path/to/your/privkey.pem"
        }
    }
    ```

    **Note:** Replace `/path/to/your/` with the actual path inside the container where the certificates are located.

## 3. Update Docker Compose

You need to mount the directory containing your certificates into the Janus container.

1.  **Mount the Certificates:** In your `docker-compose.yml` file, add a volume mount for the `janus` service:

    ```yaml
    services:
      janus:
        # ... other service configurations
        volumes:
          - ./janus-server/certs:/path/to/your/certs:ro # Mount certs read-only
          # ... other volumes
    ```

2.  **Update Admin URL:** Ensure that the `websocket` service (and any other service that communicates with the Janus admin interface) is configured to use the HTTPS URL:

    ```yaml
    services:
      websocket:
        environment:
          # ...
          JANUS_ADMIN_URL: https://janus:7889/admin
          # ...
    ```

## 4. (Recommended) Use a Reverse Proxy

For a more robust and secure production setup, it is highly recommended to use a reverse proxy like Nginx. The reverse proxy would handle all incoming public traffic over HTTPS (TLS Termination) and then communicate with the Janus server internally.

In this setup:

1.  **Nginx handles HTTPS:** Your SSL/TLS certificates would be configured in Nginx.
2.  **Janus runs on HTTP:** You can keep the Janus configuration as it is for local development (with HTTPS disabled).
3.  **Nginx proxies requests:** Nginx would be configured to proxy requests to the appropriate backend services (API, WebSocket, Janus) over plain HTTP within the Docker network.

This approach centralizes your security management and simplifies the configuration of individual services. The provided `docker-compose.yml` already includes a service for Nginx under the `production` profile that can be adapted for this purpose.