package cid

import (
	"context"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ColorizedMessage Used to display provider log color messages.
// Please set the environment variable
func ColorizedMessage(ctx context.Context, level string, msg string, args ...map[string]interface{}) {
	if _, set := os.LookupEnv("CID_LOG_COLOR"); set {
		color.NoColor = false
	}
	switch level {
	case "TRACE":
		tflog.Trace(ctx, color.GreenString(msg), args...)
	case "DEBUG":
		tflog.Debug(ctx, color.GreenString(msg), args...)
	case "INFO":
		tflog.Info(ctx, color.HiBlueString(msg), args...)
	case "WARN", "WARNING":
		tflog.Warn(ctx, color.HiRedString(msg), args...)
	case "ERR", "ERROR":
		tflog.Error(ctx, color.HiRedString(msg), args...)
	}
}
