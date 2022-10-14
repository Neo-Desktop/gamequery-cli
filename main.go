package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/wisp-gg/gamequery"
	"github.com/wisp-gg/gamequery/api"
)

var (
	IP      string
	PORT    uint16
	TIMEOUT int
	INDENT  bool
)

func init() {
	flagIP := flag.String("ip", "", "Required: IP address to query [env: GQ_IP]")
	flagPORT := flag.Int("port", 0, "Required: Port number to test [env: GQ_PORT]")
	flagTIMEOUT := flag.Int("timeout", 5, "Timeout value in milliseconds [env: GQ_TIMEOUT]")
	flagHELP := flag.Bool("help", false, "Displays this help message")
	flagINDENT := flag.Bool("indent", false, "Should the output be indented")

	envIP := os.Getenv("GQ_IP")
	envPORT := os.Getenv("GQ_PORT")
	envTIMEOUT := os.Getenv("GQ_TIMEOUT")

	flag.Parse()

	if flagIP != nil && *flagIP != "" {
		IP = *flagIP
	} else if envIP != "" {
		IP = envIP
	} else {
		fmt.Fprintf(os.Stderr, "please specify an IP\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if flagPORT != nil && *flagPORT > 0 {
		PORT = uint16(*flagPORT)
	} else if envPORT != "" {
		PORT16, err := strconv.ParseInt(envPORT, 10, 16)
		if err != nil || PORT16 <= 0 {
			fmt.Fprintf(os.Stderr, "unable to process GQ_PORT (%s)\n\n", envPORT)
			flag.Usage()
			os.Exit(1)
		}
		PORT = uint16(PORT16)
	} else {
		fmt.Fprintf(os.Stderr, "please specify a Port\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if flagTIMEOUT != nil && *flagTIMEOUT >= 0 {
		TIMEOUT = *flagTIMEOUT
	} else if envTIMEOUT != "" {
		TIMEOUT64, err := strconv.ParseInt(envTIMEOUT, 10, 32)
		if err != nil || TIMEOUT64 < 0 {
			fmt.Fprintf(os.Stderr, "unable to process GQ_TIMEOUT (%s)\n\n", envPORT)
			flag.Usage()
			os.Exit(1)
		}
		TIMEOUT = int(TIMEOUT64)
	} else {
		fmt.Fprintf(os.Stderr, "please specify a Timeout\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if *flagHELP {
		flag.Usage()
		os.Exit(1)
	}

	INDENT = *flagINDENT
}

type results struct {
	Protocol string
	*api.Response
}

func main() {
	timeout := time.Duration(TIMEOUT) * time.Millisecond

	res, protocol, err := gamequery.Detect(api.Request{
		IP:      IP,
		Port:    PORT,
		Timeout: &timeout,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query: %s\n\n", err)
		os.Exit(2)
	}

	out := results{
		Protocol: protocol,
		Response: &res,
	}

	var outBytes []byte

	if INDENT {
		outBytes, err = json.MarshalIndent(out, "", "   ")
	} else {
		outBytes, err = json.Marshal(out)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to package response: %s\n\n", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "%s\n\n", outBytes)
}
