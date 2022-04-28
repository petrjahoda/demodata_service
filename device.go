package main

import (
	"github.com/petrjahoda/database"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func generateDowntimeData(db *gorm.DB, analogPort database.DevicePort) {
	timeToInsert := time.Now()
	min := 1
	max := 20
	randomNumber := rand.Intn(max-min) + min
	db.Exec("INSERT INTO device_port_analog_records (date_time,device_port_id,data) VALUES (?,?,?)", timeToInsert, int(analogPort.ID), float32(randomNumber))
}

func generateProductionData(db *gorm.DB, digitalPort database.DevicePort, analogPort database.DevicePort, pieceInserted bool) bool {
	timeToInsert := time.Now()
	timeToInsertForZero := timeToInsert.Add(1 * time.Second)
	if !pieceInserted {
		db.Exec("INSERT INTO device_port_digital_records (date_time,device_port_id,data) VALUES (?,?,1),(?,?,0)", timeToInsert, int(digitalPort.ID), timeToInsertForZero, int(digitalPort.ID))
		pieceInserted = true
	} else {
		pieceInserted = false
	}
	min := 20
	max := 100
	randomNumber := rand.Intn(max-min) + min
	if randomNumber%5 == 0 {
		min = 80
		max = 100
	}
	randomNumber = rand.Intn(max-min) + min
	db.Exec("INSERT INTO device_port_analog_records (date_time,device_port_id,data) VALUES (?,?,?)", timeToInsert, int(analogPort.ID), float32(randomNumber))
	return pieceInserted
}

func generateNewState() (actualCycle int, actualState string, totalCycles int) {
	min := 1
	max := 4
	randomNumber := rand.Intn(max-min) + min
	switch randomNumber {
	case 1:
		return 0, "poweroff", 150
	case 2:
		return 0, "downtime", 150
	default:
		return 0, "production", 300
	}
}
