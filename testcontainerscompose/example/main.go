package main

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	fibercompose "github.com/gofiber/contrib/testcontainerscompose"
)

const composeFile = `services:
  nginx:
    image: nginx:stable-alpine
    environment:
      bar: ${bar}
      foo: ${foo}
    ports:
      - "8081:80"
  mysql:
    image: mysql:8.0.36
    environment:
      - MYSQL_DATABASE=db
      - MYSQL_ROOT_PASSWORD=my-secret-pw
    ports:
     - "3307:3306"
`

func main() {
	cfg := &fiber.Config{}

	// Adding compose stack
	fibercompose.AddStack(cfg, context.Background(), strings.NewReader(composeFile))

	app := fiber.New(*cfg)

	log.Fatal(app.Listen(":3000"))
}
