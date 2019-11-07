package main

import (
	"net"
	"os"
	"time"

	"github.com/siyka-au/gofins/fins"
	"gopkg.in/alecthomas/kingpin.v2"
)

type finscContext struct {
	client        *fins.Client
	localIP       net.IP
	localPort     int
	localNetwork  byte
	localNode     byte
	localUnit     byte
	remoteIP      net.IP
	remotePort    int
	remoteNetwork byte
	remoteNode    byte
	remoteUnit    byte
	timeout       time.Duration
}

func (finsc *finscContext) makeFinsConnection(c *kingpin.ParseContext) error {

	clientAddr := fins.NewAddress(finsc.localIP, finsc.localPort, finsc.localNetwork, finsc.localNode, finsc.localUnit)
	plcAddr := fins.NewAddress(finsc.remoteIP, finsc.remotePort, finsc.remoteNetwork, finsc.remoteNode, finsc.remoteUnit)

	var err error
	finsc.client, err = fins.NewClient(clientAddr, plcAddr)
	if err != nil {
		panic(err)
	}

	finsc.client.SetTimeout(finsc.timeout)
	return nil
}

func main() {
	finsc := finscContext{}
	app := kingpin.New("finsc", "help").DefaultEnvars().PreAction(finsc.makeFinsConnection)

	app.Flag("local-ip", "local IP address").Default("0.0.0.0").IPVar(&finsc.localIP)
	app.Flag("local-port", "local port").Default("9600").IntVar(&finsc.localPort)
	app.Flag("local-network", "local FINS network").Default("0").Uint8Var(&finsc.localNetwork)
	app.Flag("local-node", "local FINS node").Default("2").Uint8Var(&finsc.localNode)
	app.Flag("local-unit", "local FINS unit").Default("0").Uint8Var(&finsc.localUnit)

	app.Flag("remote-ip", "remote IP address").Required().IPVar(&finsc.remoteIP)
	app.Flag("remote-port", "remote port").Default("9600").IntVar(&finsc.remotePort)
	app.Flag("remote-network", "remote FINS network").Default("0").Uint8Var(&finsc.remoteNetwork)
	app.Flag("remote-node", "remote FINS node").Default("1").Uint8Var(&finsc.remoteNode)
	app.Flag("remote-unit", "remote FINS unit").Default("0").Uint8Var(&finsc.remoteUnit)

	app.Flag("timeout", "timeout for commands").Default("20ms").DurationVar(&finsc.timeout)

	configureClock(app, &finsc)
	configureCpu(app, &finsc)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
