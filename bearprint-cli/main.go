package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ian-antking/bearprint/bearprint-cli/config"
	"github.com/ian-antking/bearprint/bearprint-cli/printservice"
	"github.com/ian-antking/bearprint/shared/printer"
)

var version = "dev"

func main() {
	host := flag.String("host", "", "Server host")
	port := flag.String("port", "", "Server port")
	q_flag := flag.Bool("q", false, "Treat input as a single QR code")
	qrcode_flag := flag.Bool("qrcode", false, "Treat input as a single QR code")
	nc_flag := flag.Bool("nc", false, "Do not cut the paper after printing")
	no_cut_flag := flag.Bool("no-cut", false, "Do not cut the paper after printing")
	v_flag := flag.Bool("v", false, "Print version information and exit")
	version_flag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *v_flag || *version_flag {
		fmt.Fprintf(os.Stderr, "bearprint-cli version %s\n", version)
		os.Exit(0)
	}

	cfg, err := config.NewConfig(*host, *port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	printerClient := printservice.NewClient(cfg.ServerHost, cfg.ServerPort)

	var items []printer.PrintItem

	if *q_flag || *qrcode_flag {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}

		content := strings.TrimSpace(string(stdin))
		if len(content) > 0 {
			items = append(items, printer.PrintItem{
				Type:    printer.QRCode,
				Content: content,
				Align:   printer.AlignCenter,
			})
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			items = append(items, printer.PrintItem{
				Type:    printer.Text,
				Content: scanner.Text(),
			})
		}
	}

	shouldCut := !(*nc_flag || *no_cut_flag)

	if len(items) > 0 && shouldCut {
		items = append(items, printer.PrintItem{Type: printer.Cut})
	}

	if err := printerClient.Print(items); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
