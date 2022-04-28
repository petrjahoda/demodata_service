package main

import (
	"database/sql"
	"github.com/kardianos/service"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const version = "2022.2.1.28"
const programName = "Demodata Service"
const programDescription = "Created demodata life it comes from Zapsi devices"
const downloadInSeconds = 10
const config = "user=postgres password=pj79.. dbname=system host=database port=5432 sslmode=disable application_name=zapsi_demodata_service"
const numberOfDevicesToCreate = 20

var serviceRunning = false

var (
	activeDevices  []database.Device
	runningDevices []database.Device
	deviceSync     sync.Mutex
)

type program struct{}

func (p *program) Start(s service.Service) error {
	logInfo("MAIN", "Starting "+programName+" on "+s.Platform())
	logInfo("MAIN", "© "+strconv.Itoa(time.Now().Year())+" Petr Jahoda")
	go p.run()
	serviceRunning = true
	return nil
}

func (p *program) run() {
	logInfo("MAIN", "Program version "+version+" started")
	db, _ := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	writeProgramVersionIntoSettings(db)
	createDevicesAndWorkplaces(db)
	logInfo("MAIN", "Devices checked")
	createTerminals(db)
	logInfo("MAIN", "Terminals checked")
	createWorkshiftsForWorkplaces(db)
	logInfo("MAIN", "Workshifts checked")
	initialDataOK := createInitialData(db)
	logInfo("MAIN", "Initial data checked")
	for {
		start := time.Now()
		logInfo("MAIN", "Program running")
		logInfo("MAIN", "Active devices: "+strconv.Itoa(len(activeDevices))+", running devices: "+strconv.Itoa(len(runningDevices)))
		if initialDataOK {
			for _, activeDevice := range activeDevices {
				activeDeviceIsRunning := checkDevice(activeDevice)
				if !activeDeviceIsRunning {
					go runDevice(activeDevice, db)
				}
			}
			logInfo("MAIN", "Updating default terminal records")
			updateDefaultTerminalRecords(db)
		}
		if time.Since(start) < (downloadInSeconds * time.Second) {
			sleeptime := downloadInSeconds*time.Second - time.Since(start)
			logInfo("MAIN", "Sleeping for "+sleeptime.String())
			time.Sleep(sleeptime)
		}
	}
}

func createInitialData(db *gorm.DB) bool {
	end := time.Now()
	beginning := end.AddDate(0, 0, -7)
	var analogData database.DevicePortAnalogRecord
	db.Last(&analogData)
	if analogData.ID > 0 {
		logInfo("MAIN", "Initial data already created")
		return true
	}
	var devices []database.Device
	db.Where("device_type_id = 1").Find(&devices)
	logInfo("MAIN", "Creating initial data for "+strconv.Itoa(len(devices))+" devices")
	for _, device := range devices {
		logInfo("MAIN", "Creating initial data for "+device.Name)
		insertTime := beginning
		var analogPort database.DevicePort
		db.Where("device_id = ? and device_port_type_id = 2", device.ID).Find(&analogPort)
		var digitalPort database.DevicePort
		db.Where("device_id = ? and device_port_type_id = 1", device.ID).Find(&digitalPort)
		var analogRecordsToInsert []database.DevicePortAnalogRecord
		var digitalRecordsToInsert []database.DevicePortDigitalRecord
		for insertTime.Before(end) {
			generatedState := rand.Intn(3)
			generatedDuration := rand.Intn(180)
			finalTime := insertTime.Add(time.Duration(generatedDuration) * time.Minute)
			if generatedState == 1 { // Production
				for insertTime.Before(finalTime) {
					var digitalData database.DevicePortDigitalRecord
					digitalData.DateTime = insertTime
					digitalData.DevicePortID = int(digitalPort.ID)
					digitalData.Data = 1
					digitalRecordsToInsert = append(digitalRecordsToInsert, digitalData)
					var digitalData2 database.DevicePortDigitalRecord
					digitalData2.DateTime = insertTime.Add(1 * time.Second)
					digitalData2.DevicePortID = int(digitalPort.ID)
					digitalData2.Data = 0
					digitalRecordsToInsert = append(digitalRecordsToInsert, digitalData2)
					var analogData database.DevicePortAnalogRecord
					analogData.DateTime = insertTime
					analogData.DevicePortID = int(analogPort.ID)
					min := 80
					max := 100
					randomNumber := rand.Intn(max-min) + min
					analogData.Data = float32(randomNumber)
					analogRecordsToInsert = append(analogRecordsToInsert, analogData)
					insertTime = insertTime.Add(10 * time.Second)
				}
			} else if generatedState == 2 { // Downtime
				for insertTime.Before(finalTime) {
					var analogData database.DevicePortAnalogRecord
					analogData.DateTime = insertTime
					analogData.DevicePortID = int(analogPort.ID)
					min := 2
					max := 20
					randomNumber := rand.Intn(max-min) + min
					analogData.Data = float32(randomNumber)
					analogRecordsToInsert = append(analogRecordsToInsert, analogData)
					insertTime = insertTime.Add(10 * time.Second)
				}

			} else {
				insertTime = insertTime.Add(time.Duration(generatedDuration) * time.Minute)
			}
			if len(analogRecordsToInsert) > 1000 {
				db.Clauses(clause.OnConflict{DoNothing: true}).Create(&analogRecordsToInsert)
				analogRecordsToInsert = nil
			}
			if len(digitalRecordsToInsert) > 1000 {
				db.Clauses(clause.OnConflict{DoNothing: true}).Create(&digitalRecordsToInsert)
				digitalRecordsToInsert = nil
			}
		}
		logInfo("MAIN", "Initial data created for "+device.Name)
	}
	logInfo("MAIN", "Initial data created")
	return true
}

func updateDefaultTerminalRecords(db *gorm.DB) {
	if serviceRunning {
		var orderRecords []database.OrderRecord
		db.Where("order_id = 1").Find(&orderRecords)
		for _, orderRecord := range orderRecords {
			db.Exec("UPDATE order_records SET order_id=(SELECT id FROM orders ORDER BY RANDOM() LIMIT 1), user_id=(SELECT id FROM users where user_type_id=1 ORDER BY RANDOM() LIMIT 1), updated_at=? WHERE id = ?", time.Now(), orderRecord.ID)
		}
		db.Exec("UPDATE downtime_records SET downtime_id=(SELECT id FROM downtimes ORDER BY RANDOM() LIMIT 1),updated_at=? WHERE id = (select id from downtime_records where date_time_end is null AND downtime_id = 1 limit 1)", time.Now())

		db.Exec("UPDATE user_records SET user_id=(SELECT id FROM users where user_type_id=1 ORDER BY RANDOM() LIMIT 1),updated_at=? WHERE id = (select id from user_records where date_time_end is null AND user_id = 1 limit 1)", time.Now())
	}
}

func createWorkshiftsForWorkplaces(db *gorm.DB) {
	var workplaceWorkShifts []database.WorkplaceWorkshift
	var workplaces []database.Workplace
	db.Find(&workplaceWorkShifts)
	db.Find(&workplaces)
	if len(workplaceWorkShifts) == 0 {
		logInfo("MAIN", "Creating workplace workshifts")
		for _, workplace := range workplaces {
			createWorkshiftsForWorkplace(workplace.ID, db)
		}
	}
}

func createWorkshiftsForWorkplace(workplaceId uint, db *gorm.DB) {
	for i := 1; i <= 3; i++ {
		newWorkplaceWorkshift := database.WorkplaceWorkshift{
			WorkplaceID: int(workplaceId),
			WorkshiftID: i,
		}
		db.Create(&newWorkplaceWorkshift)
	}
}

func createTerminals(db *gorm.DB) {
	var deviceType database.DeviceType
	db.Where("name=?", "Touchscreen Terminal").Find(&deviceType)
	var activeTerminals []database.Device
	db.Where("device_type_id=?", deviceType.ID).Where("activated = ?", "1").Find(&activeTerminals)
	if len(activeTerminals) == 0 {
		logInfo("MAIN", "Creating terminals")
		for i := 0; i < numberOfDevicesToCreate; i++ {
			addTerminalWithWorkplace("CNC Terminal "+strconv.Itoa(i), "192.168.1."+strconv.Itoa(i), "CNC "+strconv.Itoa(i), db)
		}
	}
}

func addTerminalWithWorkplace(workplaceName string, ipAddress string, terminalName string, db *gorm.DB) {
	var deviceType database.DeviceType
	db.Where("name=?", "Touchscreen Terminal").Find(&deviceType)
	db.Where("name=?", "Touchscreen Terminal").Find(&deviceType)
	newTerminal := database.Device{Name: workplaceName, DeviceTypeID: int(deviceType.ID), IpAddress: ipAddress, TypeName: "Touchscreen Terminal", Activated: true}
	db.Create(&newTerminal)
	var workplace database.Workplace
	db.Where("name = ?", terminalName).Find(&workplace)
	newRecord := database.DeviceWorkplaceRecord{
		DeviceID:    int(newTerminal.ID),
		WorkplaceID: int(workplace.ID),
	}
	db.Create(&newRecord)

}
func (p *program) Stop(s service.Service) error {
	serviceRunning = false
	for len(runningDevices) != 0 {
		logInfo("MAIN", "Stopping, still running devices: "+strconv.Itoa(len(runningDevices)))
		time.Sleep(1 * time.Second)
	}
	logInfo("MAIN", "Stopped on platform "+s.Platform())
	return nil
}

func main() {
	serviceConfig := &service.Config{
		Name:        programName,
		DisplayName: programName,
		Description: programDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		logError("MAIN", err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("MAIN", "Problem starting "+serviceConfig.Name)
	}
}

func createDevicesAndWorkplaces(db *gorm.DB) {
	var deviceType database.DeviceType
	db.Where("name=?", "Zapsi").Find(&deviceType)
	db.Where("device_type_id=?", deviceType.ID).Where("activated = ?", "1").Find(&activeDevices)
	if len(activeDevices) == 0 {
		logInfo("MAIN", "Creating devices")
		for i := 0; i < numberOfDevicesToCreate; i++ {
			addDeviceWithWorkplace("CNC "+strconv.Itoa(i), "192.168.0."+strconv.Itoa(i), db)
		}
	}
}

func addDeviceWithWorkplace(workplaceName string, ipAddress string, db *gorm.DB) {
	var deviceType database.DeviceType
	db.Where("name=?", "Zapsi").Find(&deviceType)
	newDevice := database.Device{Name: workplaceName, DeviceTypeID: int(deviceType.ID), IpAddress: ipAddress, TypeName: "Zapsi", Activated: true}
	db.Create(&newDevice)
	var device database.Device
	db.Where("name=?", workplaceName).Find(&device)
	deviceDigitalPort := database.DevicePort{Name: "Production", Unit: "ks", PortNumber: 1, DevicePortTypeID: 1, DeviceID: int(device.ID)}
	deviceAnalogPort := database.DevicePort{Name: "Amperage", Unit: "A", PortNumber: 3, DevicePortTypeID: 2, DeviceID: int(device.ID)}
	db.Create(&deviceDigitalPort)
	db.Create(&deviceAnalogPort)
	var state database.State
	db.Where("name=?", "Poweroff").Find(&state)
	var mode database.WorkplaceMode
	db.Where("name=?", "Production").Find(&mode)
	newWorkplace := database.Workplace{Name: workplaceName, Code: workplaceName, WorkplaceModeID: 1, Voltage: 230}
	db.Create(&newWorkplace)
	var workplace database.Workplace
	db.Where("name=?", workplaceName).Find(&workplace)
	var devicePortDigital database.DevicePort
	db.Where("name=?", "Production").Where("device_id=?", device.ID).Find(&devicePortDigital)
	var productionState database.State
	db.Where("name=?", "Production").Find(&productionState)
	digitalPort := database.WorkplacePort{Name: "Production", DevicePortID: int(devicePortDigital.ID), WorkplaceID: int(workplace.ID), StateID: sql.NullInt32{Int32: int32(productionState.ID), Valid: true}, CounterOK: true}
	db.Create(&digitalPort)
	var devicePortAnalog database.DevicePort
	db.Where("name=?", "Amperage").Where("device_id=?", device.ID).Find(&devicePortAnalog)
	var poweroffState database.State
	db.Where("name=?", "Poweroff").Find(&poweroffState)
	analogPort := database.WorkplacePort{Name: "Amperage", DevicePortID: int(devicePortAnalog.ID), WorkplaceID: int(workplace.ID), StateID: sql.NullInt32{Int32: int32(poweroffState.ID), Valid: true}}
	db.Create(&analogPort)
}

func checkDevice(device database.Device) bool {
	for _, runningDevice := range runningDevices {
		if runningDevice.Name == device.Name {
			return true
		}
	}
	return false
}

func runDevice(device database.Device, db *gorm.DB) {
	logInfo(device.Name, "Device started running")
	deviceSync.Lock()
	runningDevices = append(runningDevices, device)
	deviceSync.Unlock()
	deviceIsActive := true
	actualCycle := 0
	totalCycles := 0
	actualState := "poweroff"
	var digitalPort database.DevicePort
	db.Where("device_id=?", device.ID).Where("device_port_type_id=1").Where("port_number=1").Find(&digitalPort)
	var analogPort database.DevicePort
	db.Where("device_id=?", device.ID).Where("device_port_type_id=2").Where("port_number=3").Find(&analogPort)
	logInfo(device.Name, "Digital port id: "+strconv.Itoa(int(digitalPort.ID))+", analog port id: "+strconv.Itoa(int(analogPort.ID)))
	pieceInserted := false
	for deviceIsActive && serviceRunning {
		start := time.Now()
		if actualCycle >= totalCycles {
			actualCycle, actualState, totalCycles = generateNewState()
		}
		switch actualState {
		case "production":
			logInfo(device.Name, "Production -> "+strconv.Itoa(actualCycle)+" of "+strconv.Itoa(totalCycles))
			pieceInserted = generateProductionData(db, digitalPort, analogPort, pieceInserted)
		case "downtime":
			logInfo(device.Name, "Downtime -> "+strconv.Itoa(actualCycle)+" of "+strconv.Itoa(totalCycles))
			generateDowntimeData(db, analogPort)
		case "poweroff":
			logInfo(device.Name, "Poweroff -> "+strconv.Itoa(actualCycle)+" of "+strconv.Itoa(totalCycles))
		}
		logInfo(device.Name, "Processing takes "+time.Since(start).String())
		sleep(device, start)
		deviceIsActive = checkActive(device)
		actualCycle++
	}
	removeDeviceFromRunningDevices(device)
	logInfo(device.Name, "Device not active, stopped running")

}

func sleep(device database.Device, start time.Time) {
	if time.Since(start) < (downloadInSeconds * time.Second) {
		sleepTime := downloadInSeconds*time.Second - time.Since(start)
		logInfo(device.Name, "Sleeping for "+sleepTime.String())
		time.Sleep(sleepTime)
	}
}

func checkActive(device database.Device) bool {
	for _, activeDevice := range activeDevices {
		if activeDevice.Name == device.Name {
			logInfo(device.Name, "Device still active")
			return true
		}
	}
	logInfo(device.Name, "Device not active")
	return false
}

func removeDeviceFromRunningDevices(device database.Device) {
	deviceSync.Lock()
	for idx, runningDevice := range runningDevices {
		if device.Name == runningDevice.Name {
			runningDevices = append(runningDevices[0:idx], runningDevices[idx+1:]...)
		}
	}
	deviceSync.Unlock()
}

func writeProgramVersionIntoSettings(db *gorm.DB) {
	var settings database.Setting
	db.Where("name=?", programName).Find(&settings)
	settings.Name = programName
	settings.Value = version
	db.Save(&settings)
	logInfo("MAIN", "Updated version in database for "+programName)
}
