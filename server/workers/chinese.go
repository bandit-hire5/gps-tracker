package workers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gps/gps-tracker/conf"

	"github.com/gps/gps-tracker/models"
)

func NewChinese(config conf.Config, tracker *models.SupportedTracker) *Chinese {
	return &Chinese{config, tracker}
}

type Chinese struct {
	config conf.Config
	*models.SupportedTracker
}

func (t Chinese) handshake(imei string) string {
	return fmt.Sprintf("(%s%s)", imei, "AP01HSO")
}

func (t Chinese) Work(input string, w *bufio.Writer) {
	if w == nil {
		return
	}

	dataTrim := strings.Trim(input, " ")
	if dataTrim == "" {
		return
	}

	splitData := strings.Split(dataTrim, ")(")
	if len(splitData) == 0 {
		return
	}

	var sendHandshake bool
	db := t.config.DB()

	for _, splitItem := range splitData {
		clearStr := strings.TrimLeft(strings.TrimRight(splitItem, ")"), "(")

		imei := clearStr[0:12]
		command := clearStr[12:16]
		if command == "BP00" {
			if sendHandshake {
				continue
			}

			_, _ = w.WriteString(t.handshake(clearStr[0:12]))
			_ = w.Flush()

			sendHandshake = true

			continue
		}

		if command != "BR00" {
			continue
		}

		re := regexp.MustCompile(`^([\d]{12})([\w]{2}[\d]{2})([\d]{6})A([^N]+)N([^E]+)E([\d]{3}[\.][\d])([\d]{6})`)
		matchSlice := re.FindStringSubmatch(clearStr)

		if matchSlice == nil {
			continue
		}

		if len(matchSlice) != 8 {
			continue
		}

		registredTracker, err := db.FindRegistredTrackerByIMEI(imei)
		if err != nil {
			continue
		}

		if registredTracker == nil {
			continue
		}

		dataListItem := models.Tracker{
			IMEI:    imei,
			Command: command,
		}

		dateString := matchSlice[3]
		timeString := matchSlice[7]

		date, err := time.Parse(
			"06-01-02 15:04:05",
			fmt.Sprintf(
				"%s-%s-%s %s:%s:%s",
				dateString[0:2],
				dateString[2:4],
				dateString[4:6],
				timeString[0:2],
				timeString[2:4],
				timeString[4:6],
			),
		)
		if err != nil {
			continue
		}

		dataListItem.Date = date

		lat, err := strconv.ParseFloat(matchSlice[4], 64)
		if err != nil {
			continue
		}

		dataListItem.Lat = lat

		lon, err := strconv.ParseFloat(matchSlice[5], 64)
		if err != nil {
			continue
		}

		dataListItem.Lon = lon

		speed, err := strconv.ParseFloat(matchSlice[6], 64)
		if err != nil {
			continue
		}

		dataListItem.Speed = speed

		_ = db.AddNewTrackerLog(&dataListItem)
		t.logToFile(dataListItem)
	}
}

func (t Chinese) logToFile(input interface{}) {
	bytes, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("logs/chinese.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s\n", string(bytes))); err != nil {
		panic(err)
	}
}
