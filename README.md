http-echo
=========
HTTP Echo is a small go web server that serves the contents it was started with
as an HTML page.

The default port is 5678, but this is configurable via the `-listen` flag:

```
http-echo -listen=:8080 -text="hello world"
```

Then visit http://localhost:8080/ in your browser.

## Docker Usage

You can run http-echo using Docker from GitHub Container Registry:

```bash
# Run with default settings
docker run --rm -p 5678:5678 ghcr.io/kadoshita/http-echo:latest

# Run with custom text and port
docker run --rm -p 8080:8080 ghcr.io/kadoshita/http-echo:latest -listen=:8080 -text="Hello from Docker"

# Run with environment variable
docker run --rm -p 5678:5678 -e ECHO_TEXT="Hello World" ghcr.io/kadoshita/http-echo:latest
```

Available image tags:
- `latest` - Latest stable release
- `v{version}` - Specific version (e.g., `v1.2.3`)
- `dev` - Development builds from main branch
