package plugin

import (
	"pfi/sensorbee/fluentd"
	"pfi/sensorbee/sensorbee/bql"
)

func init() {
	bql.MustRegisterGlobalSourceCreator("fluentd", bql.SourceCreatorFunc(fluentd.NewSource))
	bql.MustRegisterGlobalSinkCreator("fluentd", bql.SinkCreatorFunc(fluentd.NewSink))
}
