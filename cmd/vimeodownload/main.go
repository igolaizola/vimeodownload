package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/igolaizola/vimeodownload"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	// Create signal based context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Launch command
	cmd := newCommand()
	if err := cmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func newCommand() *ffcli.Command {
	fs := flag.NewFlagSet("vimeodownload", flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	id := fs.String("id", "", "video id")
	referer := fs.String("referer", "", "referer")
	out := fs.String("out", "", "output file")

	return &ffcli.Command{
		Name:       "vimeodownload",
		ShortUsage: "vimeodownload [flags] <key> <value data...>",
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("VIMEODOWNLOAD"),
		},
		ShortHelp: "vimeodownload action",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			if *id == "" {
				return errors.New("id is required")
			}
			return vimeodownload.Download(ctx, *id, *referer, *out)
		},
	}
}
