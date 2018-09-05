package log

import (
	"github.com/pivotal-golang/lager"
	"os"
)

const ServcieName = "service_plan"

var Logger lager.Logger

func init() {
	Logger = lager.NewLogger(ServcieName)
	Logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))
}
