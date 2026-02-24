package cmd

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nanoteck137/dwebble/apis"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/core"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		app := core.NewBaseApp(&config.LoadedConfig)

		err := app.Bootstrap()
		if err != nil {
			slog.Error("Failed to bootstrap app", "err", err)
			os.Exit(-1)
		}

		e, err := apis.Server(app)
		if err != nil {
			slog.Error("Failed to create server", "err", err)
			os.Exit(-1)
		}

		slog.Info("Starting server...")

		done := make(chan bool, 1)

		// listen for interrupt signal to gracefully shutdown the application
		go func() {
			sigch := make(chan os.Signal, 1)
			signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
			<-sigch

			done <- true
		}()

		go func() {
			err = e.Start(app.Config().ListenAddr)
			if err != nil {
				slog.Error("Failed to start server", "err", err)
				os.Exit(-1)
			}

			done <- true
		}()

		<-done

		slog.Info("Stopping server...")

		err = app.Shutdown()
		if err != nil {
			slog.Error("Failed to shutdown app", "err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
