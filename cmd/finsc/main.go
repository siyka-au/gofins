package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/siyka-au/gofins/fins"
	"gopkg.in/alecthomas/kingpin.v2"
)

var iso8601DateRegex = regexp.MustCompile(`^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])[Tt](2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])(.([0-9]+))?([Zz])?$`)

var (
	app           = kingpin.New("fins", "help").DefaultEnvars()
	remoteIP      = app.Flag("remote-ip", "remote IP address").Required().IP()
	remotePort    = app.Flag("remote-port", "remote port").Default("9600").Int()
	remoteNetwork = app.Flag("remote-network", "remote FINS network").Default("0").Uint8()
	remoteNode    = app.Flag("remote-node", "remote FINS node").Default("1").Uint8()
	remoteUnit    = app.Flag("remote-unit", "remote FINS unit").Default("0").Uint8()

	localIP      = app.Flag("local-ip", "local IP address").Default("0.0.0.0").IP()
	localPort    = app.Flag("local-port", "local port").Default("9600").Int()
	localNetwork = app.Flag("local-network", "local FINS network").Default("0").Uint8()
	localNode    = app.Flag("local-node", "local FINS node").Default("2").Uint8()
	localUnit    = app.Flag("local-unit", "local FINS unit").Default("0").Uint8()

	timeout = app.Flag("timeout", "timeout for commands").Default("20ms").Duration()

	clock              = app.Command("clock", "Remote clock")
	clockRead          = clock.Command("read", "Read remote clock")
	clockWrite         = clock.Command("write", "Write remote clock")
	clockWriteDateTime = clockWrite.Arg("date-time", "Date and time to set").Required().String()
)

func main() {
	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	// context, e := app.ParseContext(os.Args[1:])
	// app.FatalIfError(e, "Unknown error")

	clientAddr := fins.NewAddress(*localIP, *localPort, *localNetwork, *localNode, *localUnit)
	plcAddr := fins.NewAddress(*remoteIP, *remotePort, *remoteNetwork, *remoteNode, *remoteUnit)

	c, err := fins.NewClient(clientAddr, plcAddr)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.SetTimeout(*timeout)

	switch command {
	case clockRead.FullCommand():
		t, e := c.ReadClock()
		app.FatalIfError(e, "")
		fmt.Println(t)

	case clockWrite.FullCommand():
		var t time.Time
		matches := iso8601DateRegex.FindStringSubmatch(*clockWriteDateTime)

		if *clockWriteDateTime == "now" {
			t = time.Now()
		} else if len(matches) == 10 {
			year, _ := strconv.Atoi(matches[1])
			month, _ := strconv.Atoi(matches[2])
			day, _ := strconv.Atoi(matches[3])
			hours, _ := strconv.Atoi(matches[4])
			minutes, _ := strconv.Atoi(matches[5])
			seconds, _ := strconv.Atoi(matches[6])
			ms, _ := strconv.Atoi(matches[8])

			t = time.Date(year, time.Month(month), day, hours, minutes, seconds, ms*1000, time.Local)

			fmt.Println(t)
			e := c.WriteClock(t)
			if e != nil {
				panic(e)
			}
		} else {
			app.Fatalf("Date and time to be in ISO 8601 format of YYYY-MM-DDThh:mm:ss.mmmZ, or 'now'")
		}

	}
}
