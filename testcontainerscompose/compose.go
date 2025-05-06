package testcontainers

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

// stack represents a testcontainer that implements the fiber.RuntimeDependency interface.
// It manages the lifecycle of a compose.DockerCompose stack.
type stack struct {
	ctx   context.Context
	stack *compose.DockerCompose
}

// Start initializes and starts the compose stack. It implements the fiber.RuntimeDependency interface.
func (c *stack) Start(ctx context.Context) error {
	err := c.stack.Up(ctx)
	if err != nil {
		return fmt.Errorf("up stack: %w", err)
	}

	return nil
}

// String returns a human-readable representation of the compose stack's state.
// It implements the fiber.RuntimeDependency interface.
func (c *stack) String() string {
	if c.stack == nil {
		return "compose-stack (not started)"
	}

	return fmt.Sprintf("compose-stack (%s)", strings.Join(c.stack.Services(), ", "))
}

// State returns the state of the compose stack.
// It implements the fiber.RuntimeDependency interface.
func (c *stack) State(ctx context.Context) (string, error) {
	if c.stack == nil {
		return "compose-stack (not started)", nil
	}

	var errors []error
	var statuses []string
	for _, s := range c.stack.Services() {
		ctr, err := c.stack.ServiceContainer(ctx, s)
		if err != nil {
			errors = append(errors, fmt.Errorf("service container: %w", err))
			continue
		}

		st, err := ctr.State(ctx)
		if err != nil {
			errors = append(errors, fmt.Errorf("service state: %w", err))
			continue
		}

		statuses = append(statuses, s+": "+st.Status)
	}

	if len(errors) > 0 {
		return "", fmt.Errorf("failed to get status of compose stack: %w", errors)
	}

	return strings.Join(statuses, ", "), nil
}

// Terminate stops and removes the container. It implements the fiber.RuntimeDependency interface.
func (c *stack) Terminate(ctx context.Context) error {
	return c.stack.Down(
		ctx,
		compose.RemoveOrphans(true),
		compose.RemoveVolumes(true),
		compose.RemoveImagesLocal,
	)
}

// Add is a convenience function that adds a container to the Fiber app using the testcontainers.Run function.
// It's equivalent to calling Add with the testcontainers.Run function.
func AddStack(cfg *fiber.Config, ctx context.Context, r ...io.Reader) error {
	dc, err := compose.NewDockerComposeWith(compose.WithStackReaders(r...))
	if err != nil {
		return fmt.Errorf("new docker compose: %w", err)
	}

	s := &stack{
		ctx:   ctx,
		stack: dc,
	}

	cfg.DevTimeDependencies = append(cfg.DevTimeDependencies, s)

	return nil
}
