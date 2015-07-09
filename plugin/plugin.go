package plugin

import (
	"pfi/sensorbee/fluentd"
	"pfi/sensorbee/sensorbee/bql"
)

func init() {
	if err := bql.RegisterGlobalSourceCreator("fluentd", bql.SourceCreatorFunc(fluentd.NewSource)); err != nil {
		panic(err)
	}
	if err := bql.RegisterGlobalSinkCreator("fluentd", bql.SinkCreatorFunc(fluentd.NewSink)); err != nil {
		panic(err)
	}
}
