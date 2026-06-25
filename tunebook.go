package tunebook

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

var AppName = "tunebook"
var CliAppName = AppName + "-cli"

var Version = "no-version"
var Commit = "no-commit"

func VersionTemplate(appName string) string {
	return fmt.Sprintf(
		"%s: %s (%s)\n",
		appName, Version, Commit,
	)
}

func DefaultLogger() *slog.Logger {
	// TODO(patrik): Don't use github.com/golang-cz/devslog for prod
	logger := slog.New(devslog.NewHandler(os.Stdout, nil))
	return logger 
}

func init() {
	slog.SetDefault(DefaultLogger())
}
