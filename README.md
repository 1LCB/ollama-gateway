# Ollama Gateway Documentation

## Overview

The Ollama Gateway is a lightweight HTTP reverse proxy with advanced features such as rate limiting, security filtering, API key validation, metrics collection, and hot-reloading of configurations. It enables secure and efficient routing to multiple Ollama servers.

---

## Features

1. **Reverse Proxy**: Distributes incoming requests to configured Ollama servers using a round-robin strategy.
2. **Rate Limiting**: Limits the number of requests per client within a specified time window.
3. **API Key Authentication**: Verifies requests based on API keys for access control.
4. **Security Filters**: Supports IP allowlist and denylist configurations.
5. **Metrics**: Provides Prometheus-compatible metrics for monitoring.
6. **Hot Configuration Reload**: Allows reloading of the configuration file without restarting the server.
7. **Key Generation**: Generates new API keys dynamically via an HTTP endpoint.

---

## Configuration

The gateway configuration is stored in a `config.json` file. Below is an example:

```json
{
    "ollamaAddresses": ["http://localhost:11434"],
    "gatewayAddress": "0.0.0.0:8080",
    "logging": true,
    "authHeaderName": "Authorization",
    "apiKeys": [
        "ahP7DieD7rNpTUa6No--iJpOaXY3TDK8dIjKg0cp-hI="
    ],
    "rateLimit": {
        "maxRequests": 25,
        "timeWindowSeconds": 30,
        "enabled": false
    },
    "security": {
        "allowIps": [],
        "denyIps": []
    },
    "metrics": {
        "enabled": true,
        "endpoint": "/metrics",
        "namespace": "gateway"
    }
}
```

### Key Configuration Fields

- **`ollamaAddresses`**: List of Ollama server addresses to forward requests to.
- **`gatewayAddress`**: The address and port on which the gateway listens.
- **`logging`**: Enables or disables logging.
- **`authHeaderName`**: The name of the HTTP header used for API key authentication.
- **`rateLimit`**: Configures the maximum requests allowed within a time window.
- **`security`**: Configures IP allowlists and denylists for access control.
- **`metrics`**: Configures Prometheus metrics endpoint and namespace.

---

## Run using Docker

You can run the Ollama Gateway easily using Docker, simply execute the following command:

```bash
docker run -d \
  -v /path/to/config.json:/config.json \
  -p 8080:8080 \
  1lcb/ollama-gateway:<tag>
```

---

## HTTP Endpoints

- **`/`**: Reverse proxy endpoint to forward requests to Ollama servers.
- **`/config/reload`**: Reloads the configuration file.
- **`/config/new-key`**: Generates a new API key.
- **`/metrics`**: Exposes Prometheus metrics (customizable endpoint).

---

## Metrics

The following Prometheus metrics are available:

1. **`gateway_requests_total`**: Total number of HTTP requests, labeled by method, endpoint, and status.
2. **`gateway_requests_duration_seconds`**: Latency histogram of HTTP requests, labeled by method and endpoint.
3. **`gateway_requests_by_ip_total`**: Total number of requests grouped by IP address, endpoint, and status.
4. **`gateway_requests_success_total`**: Total number of successful HTTP requests, labeled by method and endpoint.
5. **`gateway_requests_failure_total`**: Total number of failed HTTP requests, labeled by method and endpoint.
6. **`gateway_apikeys_total`**: Number of active API keys.
7. **`gateway_ollama_servers_total`**: Number of configured Ollama servers.

---

## Middleware

1. **RateLimitMiddleware**: Enforces rate limiting based on the configuration.
2. **IPAndAPIKeyMiddleware**: Validates incoming requests using IP filters and API keys.
3. **MetricsMiddleware**: Collects and exposes metrics for Prometheus.

---

## Example: Using the Ollama Gateway with the Default Library

With just a few simple changes, you can easily integrate the Ollama Gateway into your existing code. The gateway automatically handles authentication, rate limiting, IP filtering, and more. Here's how:

```python
from ollama import Client

# Set the Ollama Gateway address and authentication headers
ollama_gateway_address = "http://localhost:8080"
headers = {
    "Authorization": "ahP7DieD7rNpTUa6No--iJpOaXY3TDK8dIjKg0cp-hI"
}

# Create a client pointing to the Ollama Gateway
client = Client(
    host=ollama_gateway_address,
    headers=headers
)

# Send a chat message to the model
stream = client.chat(
    model="llama3.1:8b",
    messages=[{"role": "user", "content": "Hello, how are you today?"}],
    stream=True
)

# Print the response from the model
for chunk in stream:
    print(chunk.message.content, end="", flush=True)
print()
```


## Example: Hot Configuration Reload

To reload the configuration:

1. Update the `config.json` file.
2. Send a `GET` request to `/config/reload`.

On success, the gateway updates its settings dynamically without downtime.

---

## Example: Generating a New API Key

Send a `GET` request to `/config/new-key` to generate a new API key. The response will include the newly created key. The generated key will also be added in the `ApiKeys` list from `config.json` file

---

## Implementation Details

### Reverse Proxy Logic

The gateway uses a round-robin strategy to distribute requests across Ollama servers:

```go
func getAddressRoundRobin() string {
    mu.Lock()
    address := cfg.OllamaAddresses[currentIndex]
    if len(cfg.OllamaAddresses) > 1 {
        currentIndex = (currentIndex + 1) % len(cfg.OllamaAddresses)
    }
    mu.Unlock()

    return address
}
```

### Hot Configuration Reload

The gateway supports hot reloading via the `ReloadConfig` function. It updates metrics and logger settings dynamically:

```go
func ReloadConfig() error {
    utils.ReloadLogger()

    if err := config.ReloadConfig(); err != nil {
        return fmt.Errorf("failed to reload config: %v", err)
    }

	if config.Metrics.Enabled{
		metrics.APIKeysTotal.Set(float64(len(config.APIKeys)))
		metrics.OllamaServersTotal.Set(float64(len(config.OllamaAddresses)))
	}

    return nil
}
```

---

## Conclusion

The Ollama Gateway provides a robust and lightweight solution for managing traffic to Ollama servers. With its features like API key validation, rate limiting, and Prometheus metrics, it ensures security, scalability, and observability in a simple, configurable manner.
