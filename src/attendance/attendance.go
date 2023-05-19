package attendance

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

const persistent = "db/attendance.gob"

var Data = make(map[string]User)

type User struct {
	Clockin    []int64 //Array of Unix time
	Clockout   []int64 //Array of Unix time
	ExtraHours int64   //Extra hours
}

func ClockIn(id string) {
	if _, ok := Data[id]; !ok {
		Data[id] = User{Clockin: []int64{time.Now().Unix()}}
		return
	}
	Data[id] = User{Clockin: append(Data[id].Clockin, time.Now().Unix()), Clockout: Data[id].Clockout, ExtraHours: Data[id].ExtraHours}
	fmt.Println(Data)
}

func ClockOut(id string) {
	if _, ok := Data[id]; !ok {
		Data[id] = User{Clockout: []int64{time.Now().Unix()}}
		return
	}
	Data[id] = User{Clockout: append(Data[id].Clockout, time.Now().Unix()), Clockin: Data[id].Clockin, ExtraHours: Data[id].ExtraHours}
	fmt.Println(Data)
}

func ChangeHours(id string, hours int64) {
	if _, ok := Data[id]; !ok {
		Data[id] = User{ExtraHours: hours}
		return
	}
	Data[id] = User{ExtraHours: Data[id].ExtraHours + hours, Clockin: Data[id].Clockin, Clockout: Data[id].Clockout}
	fmt.Println(Data)
}

func NewPeriod() {
	Data = make(map[string]User)
	Save()
}

func Save() {
	file, _ := os.Create(persistent)
	defer file.Close()
	encoder := gob.NewEncoder(file)
	encoder.Encode(Data)
}

func Load() {
	file, err := os.Open(persistent)
	if err != nil {
		return
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	decoder.Decode(&Data)
}

func ClockedIn(id string) bool {
	if _, ok := Data[id]; !ok {
		return false
	}
	if len(Data[id].Clockin) > len(Data[id].Clockout) {
		return true
	}
	return false
}
