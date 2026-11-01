package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nanoteck137/tunebook/apis"
	"github.com/nanoteck137/tunebook/app"
	"github.com/nanoteck137/tunebook/config"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the api server",
	Run: func(cmd *cobra.Command, args []string) {
		cfgFile, _ := cmd.Flags().GetString("config")

		config, err := config.Load(cfgFile)
		if err != nil {
			slog.Error("failed to load config", "err", err)
			os.Exit(1)
		}

		a := app.NewApp(config)

		err = a.Bootstrap()
		if err != nil {
			slog.Error("failed to bootstrap app", "err", err)
			os.Exit(1)
		}

		a.Start()

		e, err := apis.Server(a)
		if err != nil {
			slog.Error("failed to create server", "err", err)
			os.Exit(1)
		}

		slog.Info("starting server")

		done := make(chan bool, 1)

		// listen for interrupt signal to gracefully shutdown the application
		go func() {
			sigch := make(chan os.Signal, 1)
			signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
			<-sigch

			done <- true
		}()

		go func() {
			err = e.Start(a.Config().ListenAddr)
			if err != nil {
				slog.Error("failed to start server", "err", err)
				os.Exit(1)
			}

			done <- true
		}()

		<-done

		slog.Info("stopping server")

		err = a.Shutdown()
		if err != nil {
			slog.Error("failed to shutdown app", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
