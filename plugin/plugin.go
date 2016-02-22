package plugin

import (
	"github.com/sensorbee/fluentd"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
)

func init() {
	bql.MustRegisterGlobalSourceCreator("fluentd", bql.SourceCreatorFunc(fluentd.NewSource))
	bql.MustRegisterGlobalSinkCreator("fluentd", bql.SinkCreatorFunc(fluentd.NewSink))
}
