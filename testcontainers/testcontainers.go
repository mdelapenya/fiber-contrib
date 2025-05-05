package testcontainers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/testcontainers/testcontainers-go"
)

// Container represents a testcontainer that implements the fiber.RuntimeDependency interface.
// It manages the lifecycle of a testcontainers.Container instance.
type container[T testcontainers.Container] struct {
	ctr         T
	img         string
	opts        []testcontainers.ContainerCustomizer
	runFn       func(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (T, error)
	initialized bool
}

// Start initializes and starts the container. It implements the fiber.RuntimeDependency interface.
func (c *container[T]) Start(ctx context.Context) error {
	ctr, err := c.runFn(ctx, c.img, c.opts...)
	if err != nil {
		return fmt.Errorf("run container: %w", err)
	}

	c.ctr = ctr
	c.initialized = true

	return nil
}

// String returns a human-readable representation of the container's state.
// It implements the fiber.RuntimeDependency interface.
func (c *container[T]) String() string {
	if !c.initialized {
		return fmt.Sprintf("%s (not started)", c.img)
	}

	if c.ctr.IsRunning() {
		return fmt.Sprintf("%s (running through testcontainers-go)", c.img)
	}

	return fmt.Sprintf("%s (not running)", c.img)
}

func (c *container[T]) State(ctx context.Context) (string, error) {
	if !c.initialized {
		return "", fmt.Errorf("container not initialized")
	}

	st, err := c.ctr.State(ctx)
	if err != nil {
		return "", fmt.Errorf("get container state: %w", err)
	}

	return st.Status, nil
}

// Terminate stops and removes the container. It implements the fiber.RuntimeDependency interface.
func (c *container[T]) Terminate(ctx context.Context) error {
	return c.ctr.Terminate(ctx)
}

// AddModule adds a Testcontainers module container as a runtime dependency to the Fiber app.
// The module should be a function like redis.Run or postgres.Run that returns a container type
// which embeds testcontainers.Container.
func AddModule[T testcontainers.Container](cfg *fiber.Config, ctx context.Context, runFn func(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (T, error), img string, opts ...testcontainers.ContainerCustomizer) {
	if cfg == nil {
		return
	}

	c := &container[T]{
		img:   img,
		opts:  opts,
		runFn: runFn,
	}

	cfg.DevTimeDependencies = append(cfg.DevTimeDependencies, c)
}

// Add is a convenience function that adds a container to the Fiber app using the testcontainers.Run function.
// It's equivalent to calling Add with the testcontainers.Run function.
func Add(cfg *fiber.Config, ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) {
	AddModule(cfg, ctx, testcontainers.Run, img, opts...)
}
