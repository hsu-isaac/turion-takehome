package models

import (
	"time"
)

// CCSDS Primary Header (6 bytes)
type CCSDSPrimaryHeader struct {
	PacketID      uint16 // Version(3 bits), Type(1 bit), SecHdrFlag(1 bit), APID(11 bits)
	PacketSeqCtrl uint16 // SeqFlags(2 bits), SeqCount(14 bits)
	PacketLength  uint16 // Total packet length minus 7
}

// CCSDS Secondary Header (10 bytes)
type CCSDSSecondaryHeader struct {
	Timestamp   uint64 // Unix timestamp
	SubsystemID uint16 // Identifies the subsystem
}

// TelemetryPayload
type TelemetryPayload struct {
	Temperature float32 // Temperature in Celsius
	Battery     float32 // Battery percentage
	Altitude    float32 // Altitude in kilometers
	Signal      float32 // Signal strength in dB
}

type TelemetryRecord struct {
	Timestamp   time.Time
	SubsystemID uint16
	Temperature float32
	Battery     float32
	Altitude    float32
	Signal      float32
	HasAnomaly  bool
}

type Anomaly struct {
	Timestamp     time.Time
	SubsystemID   uint16
	AnomalyType   string
	Value         float32
	ExpectedRange string
}

func (t *TelemetryPayload) ValidateTemperature() (bool, string) {
	if t.Temperature > 35.0 {
		return false, "High temperature anomaly"
	}
	if t.Temperature < 20.0 || t.Temperature > 30.0 {
		return false, "Temperature out of normal range"
	}
	return true, ""
}

func (t *TelemetryPayload) ValidateBattery() (bool, string) {
	if t.Battery < 40.0 {
		return false, "Low battery anomaly"
	}
	if t.Battery < 70.0 || t.Battery > 100.0 {
		return false, "Battery out of normal range"
	}
	return true, ""
}

func (t *TelemetryPayload) ValidateAltitude() (bool, string) {
	if t.Altitude < 400.0 {
		return false, "Low altitude anomaly"
	}
	if t.Altitude < 500.0 || t.Altitude > 550.0 {
		return false, "Altitude out of normal range"
	}
	return true, ""
}

func (t *TelemetryPayload) ValidateSignal() (bool, string) {
	if t.Signal < -80.0 {
		return false, "Weak signal anomaly"
	}
	if t.Signal < -60.0 || t.Signal > -40.0 {
		return false, "Signal strength out of normal range"
	}
	return true, ""
}

func (t *TelemetryPayload) Validate() []Anomaly {
	var anomalies []Anomaly
	now := time.Now()

	if ok, msg := t.ValidateTemperature(); !ok {
		anomalies = append(anomalies, Anomaly{
			Timestamp:     now,
			AnomalyType:   msg,
			Value:         t.Temperature,
			ExpectedRange: "20.0°C - 30.0°C",
		})
	}

	if ok, msg := t.ValidateBattery(); !ok {
		anomalies = append(anomalies, Anomaly{
			Timestamp:     now,
			AnomalyType:   msg,
			Value:         t.Battery,
			ExpectedRange: "70% - 100%",
		})
	}

	if ok, msg := t.ValidateAltitude(); !ok {
		anomalies = append(anomalies, Anomaly{
			Timestamp:     now,
			AnomalyType:   msg,
			Value:         t.Altitude,
			ExpectedRange: "500km - 550km",
		})
	}

	if ok, msg := t.ValidateSignal(); !ok {
		anomalies = append(anomalies, Anomaly{
			Timestamp:     now,
			AnomalyType:   msg,
			Value:         t.Signal,
			ExpectedRange: "-60dB to -40dB",
		})
	}

	return anomalies
}
