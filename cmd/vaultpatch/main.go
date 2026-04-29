// Package main is the entry point for the vaultpatch CLI tool.
// It wires together configuration loading, Vault client initialization,
// diff computation, patch application, rendering, and audit logging.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/your-org/vaultpatch/internal/audit"
	"github.com/your-org/vaultpatch/internal/config"
	"github.com/your-org/vaultpatch/internal/diff"
	"github.com/your-org/vaultpatch/internal/env"
	"github.com/your-org/vaultpatch/internal/patch"
	"github.com/your-org/vaultpatch/internal/render"
	"github.com/your-org/vaultpatch/internal/snapshot"
	"github.com/your-org/vaultpatch/internal/vault"
)

const usage = `vaultpatch — diff and apply HashiCorp Vault secret changes across environments.

Usage:
  vaultpatch diff   <source-env> <target-env> [flags]
  vaultpatch apply  <source-env> <target-env> [flags]
  vaultpatch snap   <env>                     [flags]

Flags:
`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
		return fmt.Errorf("no subcommand provided")
	}

	// Global flags
	fs := flag.NewFlagSet("vaultpatch", flag.ContinueOnError)
	cfgPath := fs.String("config", "", "path to config file (default: vaultpatch.yaml)")
	auditPath := fs.String("audit-log", "", "path to append audit log entries (optional)")
	dryRun := fs.Bool("dry-run", false, "preview changes without writing to Vault")
	color := fs.Bool("color", true, "enable colored diff output")

	subcmd := args[0]
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	positional := fs.Args()

	// Load configuration
	var cfg *config.Config
	var err error
	if *cfgPath != "" {
		cfg, err = config.Load(*cfgPath)
	} else {
		cfg, err = config.LoadDefault()
	}
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Set up audit logger
	var logger *audit.Logger
	if *auditPath != "" {
		logger, err = audit.NewFileLogger(*auditPath)
		if err != nil {
			return fmt.Errorf("opening audit log: %w", err)
		}
	} else {
		logger = audit.NewLogger(nil)
	}

	// Build Vault client
	client, err := vault.NewClient(cfg.VaultAddress, cfg.Token)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	// Build environment resolver
	resolver := env.NewResolver(cfg.Environments)

	ctx := context.Background()

	switch subcmd {
	case "diff":
		if len(positional) < 2 {
			return fmt.Errorf("diff requires <source-env> and <target-env>")
		}
		return runDiff(ctx, client, resolver, logger, positional[0], positional[1], *color)

	case "apply":
		if len(positional) < 2 {
			return fmt.Errorf("apply requires <source-env> and <target-env>")
		}
		return runApply(ctx, client, resolver, logger, positional[0], positional[1], *dryRun)

	case "snap":
		if len(positional) < 1 {
			return fmt.Errorf("snap requires <env>")
		}
		return runSnap(ctx, client, resolver, positional[0])

	default:
		return fmt.Errorf("unknown subcommand %q", subcmd)
	}
}

func runDiff(ctx context.Context, client vault.Reader, resolver *env.Resolver, logger *audit.Logger, srcEnv, dstEnv string, colorEnabled bool) error {
	srcPath, err := resolver.Resolve(srcEnv)
	if err != nil {
		return err
	}
	dstPath, err := resolver.Resolve(dstEnv)
	if err != nil {
		return err
	}

	srcSecrets, err := client.ReadSecrets(ctx, srcPath)
	if err != nil {
		return fmt.Errorf("reading source secrets: %w", err)
	}
	dstSecrets, err := client.ReadSecrets(ctx, dstPath)
	if err != nil {
		return fmt.Errorf("reading target secrets: %w", err)
	}

	changes := diff.Diff(srcSecrets, dstSecrets)
	renderer := render.New(os.Stdout, colorEnabled)
	renderer.RenderDiff(srcEnv, dstEnv, changes)

	summary := diff.Summary(changes)
	fmt.Fprintf(os.Stdout, "\nSummary: %d added, %d removed, %d modified\n",
		summary.Added, summary.Removed, summary.Modified)

	_ = logger.LogDiff(srcEnv, dstEnv, changes)
	return nil
}

func runApply(ctx context.Context, client vault.ReadWriter, resolver *env.Resolver, logger *audit.Logger, srcEnv, dstEnv string, dryRun bool) error {
	srcPath, err := resolver.Resolve(srcEnv)
	if err != nil {
		return err
	}
	dstPath, err := resolver.Resolve(dstEnv)
	if err != nil {
		return err
	}

	srcSecrets, err := client.ReadSecrets(ctx, srcPath)
	if err != nil {
		return fmt.Errorf("reading source secrets: %w", err)
	}
	dstSecrets, err := client.ReadSecrets(ctx, dstPath)
	if err != nil {
		return fmt.Errorf("reading target secrets: %w", err)
	}

	changes := diff.Diff(srcSecrets, dstSecrets)
	applyErr := patch.Apply(ctx, client, dstPath, changes, dryRun)

	summary := patch.Summary(changes, applyErr)
	fmt.Fprintf(os.Stdout, "Apply result: %d written, %d skipped, errors: %v\n",
		summary.Written, summary.Skipped, summary.Errors)

	_ = logger.LogApply(srcEnv, dstEnv, changes, applyErr)
	return applyErr
}

func runSnap(ctx context.Context, client vault.Reader, resolver *env.Resolver, envName string) error {
	path, err := resolver.Resolve(envName)
	if err != nil {
		return err
	}

	store, err := snapshot.NewStore(".vaultpatch/snapshots")
	if err != nil {
		return fmt.Errorf("creating snapshot store: %w", err)
	}

	snap, err := snapshot.Capture(ctx, client, envName, path)
	if err != nil {
		return fmt.Errorf("capturing snapshot: %w", err)
	}

	if err := store.Save(snap); err != nil {
		return fmt.Errorf("saving snapshot: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Snapshot saved for environment %q (%d keys)\n", envName, len(snap.Secrets))
	return nil
}
