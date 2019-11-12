package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/siyka-au/gofins/fins"
	"gopkg.in/alecthomas/kingpin.v2"
)

var iso8601DateRegex = regexp.MustCompile(`^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])[Tt](2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])(.([0-9]+))?([Zz])?$`)

type clockCommand struct {
	fins   **fins.Client
	setArg string
}

func configureClock(app *kingpin.Application, finsc *finscContext) {
	clock := app.Command("clock", "Remote clock")
	cc := &clockCommand{}
	cc.fins = &finsc.client
	clock.Command("read", "Read remote clock").Action(cc.readClock)
	clockSet := clock.Command("set", "Set remote clock").Action(cc.setClock)
	clockSet.Arg("date-time", "Date and time to set").Required().StringVar(&cc.setArg)
}

func (cc *clockCommand) readClock(c *kingpin.ParseContext) error {
	t, e := (**cc.fins).ReadClock()
	if e != nil {
		return e
	}
	fmt.Println(t)
	return nil
}

func (cc *clockCommand) setClock(c *kingpin.ParseContext) error {
	var t time.Time

	if cc.setArg == "now" {
		t = time.Now()
	} else {
		matches := iso8601DateRegex.FindStringSubmatch(cc.setArg)

		if len(matches) == 10 {
			year, _ := strconv.Atoi(matches[1])
			month, _ := strconv.Atoi(matches[2])
			day, _ := strconv.Atoi(matches[3])
			hours, _ := strconv.Atoi(matches[4])
			minutes, _ := strconv.Atoi(matches[5])
			seconds, _ := strconv.Atoi(matches[6])
			ms, _ := strconv.Atoi(matches[8])

			t = time.Date(year, time.Month(month), day, hours, minutes, seconds, ms*1000000, time.Local)
		} else {
			return errors.New("Date and time to be in ISO 8601 format of YYYY-MM-DDThh:mm:ss.mmmZ, or 'now'")
		}
	}

	fmt.Printf("Setting remote clock to %s", t)
	e := (**cc.fins).WriteClock(t)
	if e != nil {
		panic(e)
	}

	return nil
}
