package plugin

import (
	"pfi/sensorbee/fluentd"
	"pfi/sensorbee/sensorbee/bql"
)

func init() {
	if err := bql.RegisterGlobalSourceCreator("fluentd", bql.SourceCreatorFunc(fluentbee.NewSource)); err != nil {
		panic(err)
	}
	if err := bql.RegisterGlobalSinkCreator("fluentd", bql.SinkCreatorFunc(fluentbee.NewSink)); err != nil {
		panic(err)
	}
}
