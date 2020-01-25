package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"go-template/common"
	"go-template/common/registry"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE\n")
	fmt.Fprintf(os.Stderr, "  serve <mode> [flags]\n")
}

func main() {
	var (
		port       = flag.Int("port", 5000, "The server port")
		jaegeraddr = flag.String("jaeger_addr", "jaeger:6831", "Jaeger address")
		consuladdr = flag.String("consul_addr", "consul:8500", "Consul address")
	)
	flag.Parse()

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var run func(port int, consul *registry.Client, jaegeraddr string) error

	switch strings.ToLower(os.Args[1]) {
	case common.ServiceTasker:
		run = runTasker
	default:
		usage()
		os.Exit(1)
	}

	consul, err := registry.NewClient(*consuladdr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := run(*port, consul, *jaegeraddr); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
