# go-redoc

`go-redoc` serves OpenAPI documentation with an embedded [ReDoc](https://github.com/ReDocly/redoc) UI. It embeds the HTML template and `redoc.standalone.js` in the Go binary, so the documentation page does not need a CDN at runtime.

The package supports `net/http` directly and includes adapters for `gin`, `echo`, `fiber`, and `iris`. It serves an existing OpenAPI JSON or YAML file. It does not generate the spec for you.

The embedded ReDoc bundle is pinned to `redoc@2.5.3` in `Taskfile.yml`.

## Install

```sh
go get github.com/VendSYSTEM/go-redoc/v2
```

Install the adapter module you use:

```sh
go get github.com/VendSYSTEM/go-redoc/gin
go get github.com/VendSYSTEM/go-redoc/echo
go get github.com/VendSYSTEM/go-redoc/fiber
go get github.com/VendSYSTEM/go-redoc/iris
```

## Quick Start

```go
package main

import (
	"net/http"

	"github.com/VendSYSTEM/go-redoc/v2"
)

func main() {
	doc := redoc.Redoc{
		Title:       "Example API",
		Description: "Example API documentation",
		SpecFile:    "./openapi.json",
		SpecPath:    "/openapi.json",
		DocsPath:    "/docs",
	}

	handler := doc.Handler()
	mux := http.NewServeMux()
	mux.Handle("/docs", handler)
	mux.Handle("/openapi.json", handler)

	panic(http.ListenAndServe(":8000", mux))
}
```

Open `http://127.0.0.1:8000/docs`.

## Configuration

`redoc.Redoc` accepts these fields:

| Field         | Required    | Description                                                                                                                                                     |
| ------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Title`       | No          | HTML page title and ReDoc document title.                                                                                                                       |
| `Description` | No          | HTML meta description.                                                                                                                                          |
| `SpecFile`    | Yes         | File path to the OpenAPI document. With `SpecFS`, this path is relative to the embedded filesystem. Without `SpecFS`, this path is read from the OS filesystem. |
| `SpecFS`      | No          | Optional `*embed.FS` used to read `SpecFile` from embedded files.                                                                                               |
| `SpecPath`    | Recommended | Public URL path where the handler serves the OpenAPI document. Use `/openapi.json`, `/openapi.yaml`, `/swagger.json`, or the route you expose.                  |
| `DocsPath`    | Recommended | Public URL path where the handler serves the ReDoc HTML page. Use `/docs` or another route you mount.                                                           |

Set `SpecPath` explicitly. `Handler` falls back to `/openapi.json` for the spec route when `SpecPath` is empty, but the rendered HTML uses the value present when `Body` runs.

Set `DocsPath` when you use the handler as middleware. If `DocsPath` is empty, the handler serves the docs page for any `GET` or `HEAD` request that does not match the spec path.

## Handler Behavior

`doc.Handler()` prepares the HTML and reads the spec once when you create the handler. If `SpecFile` is empty, or if the spec file cannot be read, the handler panics during setup.

The handler writes responses for `GET` and `HEAD` requests. It serves:

| Request path | Response                                                    |
| ------------ | ----------------------------------------------------------- |
| `SpecPath`   | The OpenAPI document with `Content-Type: application/json`. |
| `DocsPath`   | The ReDoc HTML page with `Content-Type: text/html`.         |

The handler does not implement application-level authentication, authorization, CORS, caching, compression, or OpenAPI generation. Add those concerns in your router, middleware, reverse proxy, or application code.

When you mount `doc.Handler()` directly as the only `net/http` handler, non-matching paths can receive an empty `200` response because the handler leaves unmatched requests untouched. Mount it on exact routes or wrap it with your router if you need `404` or `405` behavior.

## Embed the OpenAPI Spec

Use `SpecFS` when you want a single binary that contains both ReDoc and your OpenAPI file.

```go
package main

import (
	"embed"
	"net/http"

	"github.com/VendSYSTEM/go-redoc/v2"
)

//go:embed docs/openapi.yaml
var spec embed.FS

func main() {
	doc := redoc.Redoc{
		Title:    "Example API",
		SpecFile: "docs/openapi.yaml",
		SpecFS:   &spec,
		SpecPath: "/openapi.yaml",
		DocsPath: "/docs",
	}

	handler := doc.Handler()
	mux := http.NewServeMux()
	mux.Handle("/docs", handler)
	mux.Handle("/openapi.yaml", handler)

	panic(http.ListenAndServe(":8000", mux))
}
```

## Framework Adapters

All adapters wrap the same `doc.Handler()` behavior.

### Gin

```go
import (
	"github.com/gin-gonic/gin"
	"github.com/VendSYSTEM/go-redoc/v2"
	ginredoc "github.com/VendSYSTEM/go-redoc/gin"
)

doc := redoc.Redoc{SpecFile: "./openapi.json", SpecPath: "/openapi.json", DocsPath: "/docs"}

r := gin.New()
r.Use(ginredoc.New(doc))
```

### Echo

```go
import (
	"github.com/labstack/echo/v4"
	"github.com/VendSYSTEM/go-redoc/v2"
	echoredoc "github.com/VendSYSTEM/go-redoc/echo"
)

doc := redoc.Redoc{SpecFile: "./openapi.json", SpecPath: "/openapi.json", DocsPath: "/docs"}

r := echo.New()
r.Use(echoredoc.New(doc))
```

### Fiber

```go
import (
	"github.com/gofiber/fiber/v2"
	"github.com/VendSYSTEM/go-redoc/v2"
	fiberredoc "github.com/VendSYSTEM/go-redoc/fiber"
)

doc := redoc.Redoc{SpecFile: "./openapi.json", SpecPath: "/openapi.json", DocsPath: "/docs"}

r := fiber.New()
r.Use(fiberredoc.New(doc))
```

### Iris

```go
import (
	"github.com/kataras/iris/v12"
	"github.com/VendSYSTEM/go-redoc/v2"
	irisdoc "github.com/VendSYSTEM/go-redoc/iris"
)

doc := redoc.Redoc{SpecFile: "./openapi.json", SpecPath: "/openapi.json", DocsPath: "/docs"}

app := iris.New()
app.Use(irisdoc.New(doc))
```

## Generate an OpenAPI Spec

`go-redoc` only serves a spec. Use another tool to generate or maintain the OpenAPI file.

The [`_examples/gen`](_examples/gen) example uses `swag` with `go generate`:

```go
//go:generate swag init
```

Then it serves the generated file:

```go
doc := redoc.Redoc{
	Title:    "Example API",
	SpecFile: "./docs/swagger.json",
	SpecPath: "/swagger.json",
	DocsPath: "/docs",
}
```

## Examples

Runnable examples live in [`_examples`](_examples):

| Example                                  | Description                       |
| ---------------------------------------- | --------------------------------- |
| [`_examples/http`](_examples/http)       | Plain `net/http`.                 |
| [`_examples/gin`](_examples/gin)         | Gin middleware.                   |
| [`_examples/echo`](_examples/echo)       | Echo middleware.                  |
| [`_examples/fiber`](_examples/fiber)     | Fiber middleware.                 |
| [`_examples/iris`](_examples/iris)       | Iris middleware.                  |
| [`_examples/gorilla`](_examples/gorilla) | Gorilla mux path-prefix mounting. |
| [`_examples/gen`](_examples/gen)         | Generated spec with `swag`.       |

Most examples serve docs at `http://127.0.0.1:8000/docs`.

## Update the Embedded ReDoc Bundle

The ReDoc JavaScript bundle lives at `assets/redoc.standalone.js`. `Taskfile.yml` downloads it from jsDelivr through `REDOC_URL`.

To update ReDoc:

```sh
npm view redoc version
```

Update `REDOC_VERSION` in `Taskfile.yml`, then regenerate the embedded asset:

```sh
task redoc:update
```

Commit both `Taskfile.yml` and `assets/redoc.standalone.js` so consumers do not need network access at runtime.

## Development

Run the main test suite:

```sh
task test
```

Run formatting, vet, golangci-lint, and tests:

```sh
task all
```

Install pinned local tools with `mise`:

```sh
mise install
```

Install the lint dependency without `mise`:

```sh
task deps
```

CI runs `go test -race ./...` on Go `1.17` and `1.21`.

## Project Layout

| Path                           | Purpose                                                      |
| ------------------------------ | ------------------------------------------------------------ |
| `redoc.go`                     | Main `Redoc` type, embedded assets, and `net/http` handler.  |
| `assets/index.html`            | HTML template used by `Body`.                                |
| `assets/redoc.standalone.js`   | Embedded ReDoc standalone bundle.                            |
| `Taskfile.yml`                 | Development tasks and ReDoc bundle update command.           |
| `.mise.toml`                   | Local tool versions for Go, Task, and golangci-lint.         |
| `gin`, `echo`, `fiber`, `iris` | Framework adapter modules.                                   |
| `_examples`                    | Runnable examples for supported routers and generation flow. |
| `testdata`                     | Test OpenAPI document.                                       |

## License

MIT. See [`LICENSE`](LICENSE).
