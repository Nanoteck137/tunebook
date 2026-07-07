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
	if Version == "no-version" || Commit == "no-commit" || os.Getenv("TUNEBOOK_DEV") != "" {
		return slog.New(devslog.NewHandler(os.Stdout, nil))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func init() {
	slog.SetDefault(DefaultLogger())
}
