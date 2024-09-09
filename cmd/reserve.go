/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lib_reserve/lib"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// reserveCmd represents the reserve command
var reserveCmd = &cobra.Command{
	Use:   "reserve",
	Short: "Reserve your seat",
	Long:  fmt.Sprintf("Default time table:\n\t%s", timeTable),
	Run: func(cmd *cobra.Command, args []string) {
		if single && len(timeDuration) != 2 {
			log.Fatal(errors.New("time duration must have two --time or -t"))
		}
		err := lib.Get_cookies()
		if err != nil {
			log.Fatal(err)
		}
		reserve()
	},
}

func init() {
	rootCmd.AddCommand(reserveCmd)

	reserveCmd.Flags().StringVarP(&date, "date", "d", date, "date to reserve")
	reserveCmd.Flags().StringArrayVarP(&timeDuration, "time", "t", timeDuration, "time duration to reserve")
	reserveCmd.Flags().StringVarP(&room, "room", "r", room, "room to reserve(only 101, 206, 3A supported)")
	reserveCmd.Flags().StringVar(&lib.UserAgent, "user-agent", lib.UserAgent, "define user-agent")
	reserveCmd.Flags().IntVarP(&seat, "seat", "s", seat, "seat to reserve")
	reserveCmd.Flags().BoolVarP(&single, "single", "1", false, "reserve only one seat")
}

var (
	timeDuration        = []string{"9:00:00", "12:00:00"}
	date         string = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	room         string = "101"
	seat         int    = 129
	single       bool   = false
	// time table
	timeTable = [...][2]string{
		{"9:00:00", "12:00:00"},
		{"14:00:00", "18:00:00"},
		{"19:00:00", "22:15:00"}}
)

func reserve() {
	// small func to send request
	reserve_request :=
		func(date string, time [2]string, seatno string) {
			client := &http.Client{}
			url := "http://libbooking.gzhu.edu.cn/ic-web/reserve"
			payload := fmt.Sprintf(`{"sysKind":8, "appAccNo":101586823, "memberKind":1, "resvMember":[101586823], "resvBeginTime":"%s %s", "resvEndTime":"%s %s", "testName":"", "captcha":"", "resvProperty":0, "resvDev":["%s"], "memo":""}`,
				date, time[0], date, time[1], seatno)
			//log.Println(payload)
			req, _ := http.NewRequest("POST",
				url, bytes.NewBuffer([]byte(payload)))
			req.Header.Set("Cookie", lib.Cookies)
			req.Header.Add("User-Agent", lib.UserAgent)
			req.Header.Add("Content-Type", "application/json")
			//log.Println(req.Header)
			resp, err := client.Do(req)
			if err != nil {
				log.Panic(err)
			}
			defer resp.Body.Close()
			ret, _ := io.ReadAll(resp.Body)
			var dat map[string]interface{}
			if err := json.Unmarshal(ret, &dat); err != nil {
				log.Panic(err)
			}
			log.Println(string(ret))
			log.Printf("%s %s %s-%s %s", date, time,
				room, strconv.Itoa(seat),
				dat["message"])
		}
	// calc seatno
	seatno, err := lib.Calc_seat(room, seat)
	if err != nil {
		log.Panic(err)
	}
	// send requests
	if single {
		reserve_request(date, [2]string(timeDuration), seatno)
	} else {
		for _, ts := range timeTable {
			reserve_request(date, ts, seatno)
		}
	}
}
