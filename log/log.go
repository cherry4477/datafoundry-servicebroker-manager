package log

import (
	"github.com/pivotal-golang/lager"
	"os"
)

const ServcieBrokerName = "openshift"

var Logger lager.Logger

func init() {
	Logger = lager.NewLogger(ServcieBrokerName)
	Logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))
}
