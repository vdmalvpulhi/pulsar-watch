// Package main is the entry point for the pulsar-watch CLI tool.
// It wires together all internal packages and exposes a command-line
// interface for monitoring, filtering, replaying, and exporting
// Apache Pulsar topic messages.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/yourorg/pulsar-watch/internal/config"
	"github.com/yourorg/pulsar-watch/internal/consumer"
	"github.com/yourorg/pulsar-watch/internal/exporter"
	"github.com/yourorg/pulsar-watch/internal/filter"
	"github.com/yourorg/pulsar-watch/internal/output"
	"github.com/yourorg/pulsar-watch/internal/stats"
	"github.com/yourorg/pulsar-watch/internal/watcher"
)

var version = "dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		cfgFile  string
		verbose  bool
		broker   string
		topic    string
		sub      string
		fmt_     string
		outFile  string
		keyPat   string
		payPat   string
		maxMsgs  int
	)

	root := &cobra.Command{
		Use:   "pulsar-watch",
		Short: "Monitor and replay Apache Pulsar topic messages",
		Long: `pulsar-watch is a lightweight CLI tool for monitoring Apache Pulsar topics.

It supports real-time message filtering, export to JSON or text, replay,
and live statistics — all from a single binary.`,
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatch(cmd.Context(), cfgFile, watchOptions{
				verbose:  verbose,
				broker:   broker,
				topic:    topic,
				sub:      sub,
				format:   fmt_,
				outFile:  outFile,
				keyPat:   keyPat,
				payPat:   payPat,
				maxMsgs:  maxMsgs,
			})
		},
	}

	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (YAML)")
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose/debug output")

	root.Flags().StringVarP(&broker, "broker", "b", "", "Pulsar broker URL (e.g. pulsar://localhost:6650)")
	root.Flags().StringVarP(&topic, "topic", "t", "", "Pulsar topic to watch")
	root.Flags().StringVarP(&sub, "subscription", "s", "pulsar-watch", "subscription name")
	root.Flags().StringVarP(&fmt_, "format", "f", "text", "export format: text, json")
	root.Flags().StringVarP(&outFile, "output", "o", "", "write exported messages to file (default: stdout)")
	root.Flags().StringVar(&keyPat, "key-pattern", "", "regex to filter by message key")
	root.Flags().StringVar(&payPat, "payload-pattern", "", "regex to filter by message payload")
	root.Flags().IntVar(&maxMsgs, "max-messages", 0, "stop after receiving this many matched messages (0 = unlimited)")

	return root
}

type watchOptions struct {
	verbose  bool
	broker   string
	topic    string
	sub      string
	format   string
	outFile  string
	keyPat   string
	payPat   string
	maxMsgs  int
}

func runWatch(ctx context.Context, cfgFile string, opts watchOptions) error {
	// Load configuration, then apply CLI flag overrides.
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	if opts.broker != "" {
		cfg.BrokerURL = opts.broker
	}
	if opts.topic != "" {
		cfg.Topic = opts.topic
	}
	if opts.sub != "" {
		cfg.Subscription = opts.sub
	}
	if opts.format != "" {
		cfg.ExportFormat = opts.format
	}
	if opts.maxMsgs > 0 {
		cfg.MaxMessages = opts.maxMsgs
	}

	out := output.New(os.Stdout, opts.verbose)

	// Build message filter.
	f, err := filter.New(opts.keyPat, opts.payPat)
	if err != nil {
		return fmt.Errorf("filter: %w", err)
	}

	// Build exporter (file or stdout).
	var exp *exporter.Exporter
	if opts.outFile != "" {
		f_, err := os.Create(opts.outFile)
		if err != nil {
			return fmt.Errorf("output file: %w", err)
		}
		defer f_.Close()
		exp, err = exporter.NewWithWriter(cfg.ExportFormat, f_)
	} else {
		exp, err = exporter.New(cfg.ExportFormat)
	}
	if err != nil {
		return fmt.Errorf("exporter: %w", err)
	}

	// Build consumer.
	c, err := consumer.New(cfg.BrokerURL, cfg.Topic, cfg.Subscription)
	if err != nil {
		return fmt.Errorf("consumer: %w", err)
	}
	defer c.Close()

	// Build stats tracker.
	st := stats.New()

	// Build and run watcher.
	w, err := watcher.New(c, f, exp, st, out, cfg.MaxMessages)
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}

	// Handle OS signals for graceful shutdown.
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	out.Info("watching topic %s on %s", cfg.Topic, cfg.BrokerURL)
	return w.Run(ctx)
}
