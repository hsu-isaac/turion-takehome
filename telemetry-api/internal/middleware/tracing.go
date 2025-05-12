package middleware

import (
	"context"
	"time"

	"telemetry-api/internal/observability"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TracingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := observability.GetTracer().Start(
			context.Background(),
			c.Path(),
			trace.WithAttributes(
				attribute.String("method", c.Method()),
				attribute.String("path", c.Path()),
			),
		)
		defer span.End()

		c.Locals("ctx", ctx)

		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		observability.RecordRequest(ctx, c.Method(), c.Path(), c.Response().StatusCode(), duration)

		span.SetAttributes(
			attribute.Float64("duration", duration.Seconds()),
			attribute.Int("status", c.Response().StatusCode()),
		)

		return err
	}
}

func GetContext(c *fiber.Ctx) context.Context {
	if ctx, ok := c.Locals("ctx").(context.Context); ok {
		return ctx
	}
	return context.Background()
}
