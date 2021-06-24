package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/quickfixgo/examples/cmd/tradeclient/internal"
	"github.com/quickfixgo/quickfix"
)

//TradeClient implements the quickfix.Application interface
type TradeClient struct {
}

//OnCreate implemented as part of Application interface
func (e TradeClient) OnCreate(sessionID quickfix.SessionID) {}

//OnLogon implemented as part of Application interface
func (e TradeClient) OnLogon(sessionID quickfix.SessionID) {}

//OnLogout implemented as part of Application interface
func (e TradeClient) OnLogout(sessionID quickfix.SessionID) {}

//FromAdmin implemented as part of Application interface
func (e TradeClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return nil
}

//ToAdmin implemented as part of Application interface
func (e TradeClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {}

//ToApp implemented as part of Application interface
func (e TradeClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	fmt.Printf("Sending %s\n", msg)
	return
}

//FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e TradeClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	fmt.Printf("FromApp: %s\n", msg.String())
	return
}

func main() {
	flag.Parse()

	cfgFileName := path.Join("config", "tradeclient.cfg")
	if flag.NArg() > 0 {
		cfgFileName = flag.Arg(0)
	}

	cfg, err := os.Open(cfgFileName)
	if err != nil {
		fmt.Printf("Error opening %v, %v\n", cfgFileName, err)
		return
	}
	defer cfg.Close()
	stringData, readErr := ioutil.ReadAll(cfg)
	if readErr != nil {
		fmt.Println("Error reading cfg,", readErr)
		return
	}

	appSettings, err := quickfix.ParseSettings(bytes.NewReader(stringData))
	if err != nil {
		fmt.Println("Error reading cfg,", err)
		return
	}

	app := TradeClient{}
	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

	if err != nil {
		fmt.Println("Error creating file log factory,", err)
		return
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	if err != nil {
		fmt.Printf("Unable to create Initiator: %s\n", err)
		return
	}

	err = initiator.Start()
	if err != nil {
		fmt.Printf("Unable to start Initiator: %s\n", err)
		return
	}

	printConfig(bytes.NewReader(stringData))

Loop:
	for {
		action, err := internal.QueryAction()
		if err != nil {
			break
		}

		switch action {
		case "1":
			err = internal.QueryEnterOrder()

		case "2":
			err = internal.QueryCancelOrder()

		case "3":
			err = internal.QueryMarketDataRequest()

		case "4":
			//quit
			break Loop

		default:
			err = fmt.Errorf("unknown action: '%v'", action)
		}

		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	initiator.Stop()
}

func printConfig(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	color.Set(color.Bold)
	fmt.Println("Started FIX initiator with config:")
	color.Unset()

	color.Set(color.FgHiMagenta)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	color.Unset()
}
