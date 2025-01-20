package storage

type Block struct {
	Name    string
	Sensors map[string]Sensor
}
