package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	up := flag.Bool("up", false, "Run migrations up")
	down := flag.Bool("down", false, "Run migrations down")
	steps := flag.Int("steps", 0, "Number of migrations to run (0 for all)")
	targetVersion := flag.Uint("version", 0, "Migrate to specific version")
	flag.Parse()

	if *up && *down {
		log.Fatal("Cannot specify both -up and -down")
	}
	if !*up && !*down && *targetVersion == 0 {
		log.Fatal("Must specify either -up, -down, or -version")
	}

	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPass := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbName := getEnvOrDefault("DB_NAME", "telemetry")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error creating migration driver: %v", err)
	}

	migrationsDir := filepath.Join("migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Fatalf("Migrations directory not found: %s", migrationsDir)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"postgres", driver)
	if err != nil {
		log.Fatalf("Error creating migration instance: %v", err)
	}

	if *targetVersion > 0 {
		if err := m.Migrate(*targetVersion); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error migrating to version %d: %v", *targetVersion, err)
		}
		log.Printf("Successfully migrated to version %d", *targetVersion)
	} else if *up {
		if *steps > 0 {
			if err := m.Steps(*steps); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Error running %d migrations up: %v", *steps, err)
			}
			log.Printf("Successfully ran %d migrations up", *steps)
		} else {
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Error running migrations up: %v", err)
			}
			log.Println("Successfully ran all migrations up")
		}
	} else if *down {
		if *steps > 0 {
			if err := m.Steps(-*steps); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Error running %d migrations down: %v", *steps, err)
			}
			log.Printf("Successfully ran %d migrations down", *steps)
		} else {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Error running migrations down: %v", err)
			}
			log.Println("Successfully ran all migrations down")
		}
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Error getting migration version: %v", err)
	}
	if err == migrate.ErrNilVersion {
		log.Println("No migrations have been run")
	} else {
		log.Printf("Current migration version: %d (dirty: %v)", version, dirty)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
