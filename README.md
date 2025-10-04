# Zebrash API

[![Docker Hub](https://img.shields.io/docker/pulls/deinticketshop/zebrash-api)](https://hub.docker.com/r/deinticketshop/zebrash-api)
[![GitHub](https://img.shields.io/github/stars/dein-ticket-shop/zebrash-api?style=social)](https://github.com/dein-ticket-shop/zebrash-api)

REST API for rendering ZPL via the [Zebrash library](https://github.com/ingridhq/zebrash)

**🐳 Docker Hub:** https://hub.docker.com/r/deinticketshop/zebrash-api  
**📂 GitHub:** https://github.com/dein-ticket-shop/zebrash-api

## Features

-   HTTP API endpoint for ZPL to PNG conversion
-   Configurable image dimensions in millimeters and DPI resolution
-   CORS support for cross-origin requests
-   Docker support for easy deployment
-   Health check endpoint
-   Proper error handling and logging

## API Endpoints

### Health Check

```
GET /health
```

Returns a JSON response indicating the server status.

### ZPL Rendering

```
POST /render/{x}/{y}/{dpi}
```

**Parameters:**

-   `x`: Image width in millimeters (integer)
-   `y`: Image height in millimeters (integer)
-   `dpi`: Dots per inch resolution (integer) - automatically converted to DPMM internally

**Request Body:**
Raw ZPL data as plain text

**Response:**
PNG image data with `Content-Type: image/png`

**Example:**

```bash
curl -X POST \
  "http://localhost:3009/render/101/152/203" \
  -H "Content-Type: text/plain" \
  -d "^XA^FO50,50^A0N,50,50^FDHello World^FS^XZ" \
  --output label.png
```

This example renders a label that is 101mm wide by 152mm high at 203 DPI resolution.

## Quick Start

### Using Docker (Recommended)

The fastest way to get started is using the pre-built Docker image:

```bash
docker run -p 3009:3009 deinticketshop/zebrash-api:latest
```

The server will be available at `http://localhost:3009`

### Using Go directly

1. Clone the repository:

```bash
git clone https://github.com/dein-ticket-shop/zebrash-api.git
cd zebrash-api
```

2. Install dependencies:

```bash
go mod download
```

3. Run the server:

```bash
go run main.go
```

The server will start on port 3009.

### Using Docker

Run the container using the pre-built image from Docker Hub:

```bash
docker run -p 3009:3009 deinticketshop/zebrash-api:latest
```

Or build your own image locally:

```bash
docker build -t zebrash-api .
docker run -p 3009:3009 zebrash-api
```

### Using Docker Compose

Create a `docker-compose.yml` file using the public image:

```yaml
version: "3.8"
services:
    zebrash-api:
        image: deinticketshop/zebrash-api:latest
        ports:
            - "3009:3009"
        restart: unless-stopped
```

Or build from source:

```yaml
version: "3.8"
services:
    zebrash-api:
        build: .
        ports:
            - "3009:3009"
        restart: unless-stopped
```

Then run:

```bash
docker-compose up -d
```

## Testing the API

### Test with sample ZPL

```bash
# Create a simple test ZPL file
echo "^XA^FO50,50^A0N,50,50^FDHello World^FS^XZ" > test.zpl

# Send to API and save as PNG (101mm x 152mm at 203 DPI)
curl -X POST \
  "http://localhost:3009/render/101/152/203" \
  -H "Content-Type: text/plain" \
  --data-binary @test.zpl \
  --output test-output.png
```

### Test health endpoint

```bash
curl http://localhost:3009/health
```

Expected response:

```json
{ "status": "healthy" }
```

## ZPL Examples

Here are some example ZPL commands you can test with:

### Simple Text Label

```zpl
^XA
^FO50,50^A0N,50,50^FDHello World^FS
^XZ
```

### Label with Barcode

```zpl
^XA
^FO50,50^BY3
^BCN,100,Y,N,N
^FD123456789^FS
^FO50,200^A0N,30,30^FDBarcode: 123456789^FS
^XZ
```

### Complex Label

```zpl
^XA
^CF0,30
^FO25,25^FDShipping Label^FS
^CF0,20
^FO25,70^FDFrom: Company Name^FS
^FO25,95^FD123 Business St^FS
^FO25,120^FDCity, ST 12345^FS
^FO25,170^FDTo: Customer Name^FS
^FO25,195^FD456 Customer Ave^FS
^FO25,220^FDCity, ST 67890^FS
^FO25,270^BY2
^BCN,80,Y,N,N
^FD1234567890^FS
^XZ
```

## Configuration

The server currently runs on port 3009. To change this, modify the `port` variable in `main.go`:

```go
port := ":3009"  // Change to your desired port
```

## CORS Support

The API includes CORS (Cross-Origin Resource Sharing) headers to allow requests from web applications running on different domains. This enables browser-based applications to make requests to the API without encountering CORS errors.

## Dependencies

-   [Gin Web Framework](https://github.com/gin-gonic/gin) - HTTP web framework with CORS support
-   [Zebrash](https://github.com/ingridhq/zebrash) - ZPL to PNG conversion library

## Error Handling

The API returns appropriate HTTP status codes:

-   `200 OK` - Successful PNG generation
-   `400 Bad Request` - Invalid parameters or empty ZPL data
-   `500 Internal Server Error` - ZPL rendering failed

Error responses include descriptive messages in the response body.

## Logging

The server logs important events including:

-   Server startup
-   Rendering requests with parameters
-   Successful renders with output size
-   Errors during rendering

## Docker Image

### Available Tags

-   `latest` - Latest stable version
-   `v1.0.0`, `v1.1.0`, etc. - Specific version tags

### Pulling the Image

```bash
docker pull deinticketshop/zebrash-api:latest
```

### Image Details

-   Based on Alpine Linux for minimal size (~20MB)
-   Runs as non-root user for security
-   Includes health check endpoint
-   Multi-stage build for optimized image size
-   Available on Docker Hub: https://hub.docker.com/r/deinticketshop/zebrash-api

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is open source. Please check the license file for details.
