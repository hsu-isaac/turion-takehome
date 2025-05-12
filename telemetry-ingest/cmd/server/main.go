package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"
	"telemetry-ingest/internal/database"
	"telemetry-ingest/internal/models"
	"time"
)

const (
	UDP_PORT    = ":8089"
	BUFFER_SIZE = 1024
)

func main() {
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

	addr, err := net.ResolveUDPAddr("udp", UDP_PORT)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to start UDP server: %v", err)
	}
	defer conn.Close()

	log.Printf("Telemetry ingestion service listening on UDP port %s", UDP_PORT)

	buffer := make([]byte, BUFFER_SIZE)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP packet: %v", err)
			continue
		}

		if err := processPacket(buffer[:n], db); err != nil {
			log.Printf("Error processing packet: %v", err)
		}
	}
}

func processPacket(data []byte, db *database.Database) error {
	buf := bytes.NewReader(data)

	var primaryHeader models.CCSDSPrimaryHeader
	if err := binary.Read(buf, binary.BigEndian, &primaryHeader); err != nil {
		return err
	}

	var secondaryHeader models.CCSDSSecondaryHeader
	if err := binary.Read(buf, binary.BigEndian, &secondaryHeader); err != nil {
		return err
	}

	var payload models.TelemetryPayload
	if err := binary.Read(buf, binary.BigEndian, &payload); err != nil {
		return err
	}

	anomalies := payload.Validate()
	hasAnomaly := len(anomalies) > 0

	record := &models.TelemetryRecord{
		Timestamp:   time.Unix(int64(secondaryHeader.Timestamp), 0),
		SubsystemID: secondaryHeader.SubsystemID,
		Temperature: payload.Temperature,
		Battery:     payload.Battery,
		Altitude:    payload.Altitude,
		Signal:      payload.Signal,
		HasAnomaly:  hasAnomaly,
	}

	if err := db.StoreTelemetry(record); err != nil {
		return err
	}

	if hasAnomaly {
		for i := range anomalies {
			anomalies[i].SubsystemID = secondaryHeader.SubsystemID
		}

		if err := db.StoreAnomalies(anomalies); err != nil {
			return err
		}

		for _, anomaly := range anomalies {
			log.Printf("ALERT: %s detected - Value: %.2f (Expected Range: %s)",
				anomaly.AnomalyType,
				anomaly.Value,
				anomaly.ExpectedRange)
		}
	}

	return nil
}
