package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/testcontainers"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func main() {
	cfg := &fiber.Config{}

	// Adding module containers
	testcontainers.AddModule(cfg, context.Background(), redis.Run, "redis:latest")
	testcontainers.AddModule(cfg, context.Background(), postgres.Run, "postgres:latest")

	// Adding a generic container
	testcontainers.Add(cfg, context.Background(), "postgres:latest")

	app := fiber.New(*cfg)

	log.Fatal(app.Listen(":3000"))
}
