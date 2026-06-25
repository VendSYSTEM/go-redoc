package iris

import (
	"github.com/kataras/iris/v12"
	"github.com/VendSYSTEM/go-redoc/v2"
	irisdoc "github.com/VendSYSTEM/go-redoc/iris"
)

func main() {
	doc := redoc.Redoc{
		Title:       "Example API",
		Description: "Example API Description",
		SpecFile:    "./openapi.json",
		SpecPath:    "/openapi.json",
		DocsPath:    "/docs",
	}

	app := iris.New()
	app.Use(irisdoc.New(doc))
	println("Documentation served at http://127.0.0.1:8000/docs")
	panic(app.Listen(":8000"))
}
