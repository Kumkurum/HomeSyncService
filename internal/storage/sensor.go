package storage

type Sensor struct {
	Name string
	Data []SensorData
}

func (s *Sensor) NewSensor(name string, size int) {
	s.Name = name
	s.Data = make([]SensorData, size)
}

func (s *Sensor) GetSensorData() []SensorData {}

func (s *Sensor) AddData(timeStamp uint64, value float64) {

}
