package telegram

import "time"

func formatFullTime(now time.Time, d time.Duration) string {
	var wib = func() *time.Location {
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			return time.FixedZone("WIB", 7*60*60)
		}
		return loc
	}()
	
	return now.Add(d).In(wib).Format("15:04 WIB")
}
