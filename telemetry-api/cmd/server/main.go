package main

import (
	"encoding/json"
	"log"
	"os"
	"telemetry-api/internal/database"
	"telemetry-api/internal/handlers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Println("Starting telemetry API service...")

	db, err := database.NewDatabase(
		os.Getenv("DB_HOST"),
		5432,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	h := handlers.NewHandlers(db)

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Error handling request: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${ip} ${method} ${path} ${status} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     os.Stdout,
	}))

	api := app.Group("/api/v1")

	api.Get("/telemetry", func(c *fiber.Ctx) error {
		log.Printf("Received telemetry request: %s", c.OriginalURL())
		return h.GetTelemetry(c)
	})
	api.Get("/telemetry/current", h.GetCurrentTelemetry)
	api.Get("/telemetry/anomalies", h.GetAnomalies)

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			record, err := db.GetCurrentTelemetry()
			if err != nil {
				log.Printf("Error getting current telemetry: %v", err)
				continue
			}

			if err := c.WriteJSON(record); err != nil {
				log.Printf("Error writing to websocket: %v", err)
				break
			}

			time.Sleep(1 * time.Second)
		}
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
