package export

import (
	"attendance/src/attendance"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func CSV(id string) string {
	filePath := fmt.Sprint("/tmp/datt_", time.Now().UnixNano(), ".csv")
	fmt.Println("Creating CSV file", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return ""
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"Clock in time", "Clock out time", "Total hours"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error writing header to CSV:", err)
		return ""
	}

	// Write data rows
	user := attendance.Data[id]
	var totalHours int64
	for i := 0; i < len(user.Clockin); i++ {
		clockin := user.Clockin[i]
		clockout := user.Clockout[i]
		hours := (clockout - clockin) / 3600 // Convert to hours

		// Calculate total hours including extra hours
		totalHours += hours

		// Convert Unix timestamps to human-readable format
		clockinTime := time.Unix(clockin, 0).Format("2006-01-02 15:04:05")
		clockoutTime := time.Unix(clockout, 0).Format("2006-01-02 15:04:05")

		// Convert total hours to string for CSV
		totalHoursStr := strconv.FormatInt(totalHours, 10)

		// Write data row to CSV
		row := []string{clockinTime, clockoutTime, totalHoursStr}
		err = writer.Write(row)
		if err != nil {
			fmt.Println("Error writing data to CSV:", err)
			return ""
		}
	}
	extraHoursRow := []string{"Extra Hours:", "", strconv.FormatInt(user.ExtraHours, 10)}
	totalHoursRow := []string{"Total Hours:", "", strconv.FormatInt(totalHours+user.ExtraHours, 10)}
	writer.Write(extraHoursRow)
	writer.Write(totalHoursRow)
	return filePath
}
