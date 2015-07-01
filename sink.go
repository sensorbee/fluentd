package fluentd

import (
	"github.com/fluent/fluentd-forwarder"
	"github.com/op/go-logging"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type sink struct {
	output *fluentd_forwarder.ForwardOutput

	tagField       string
	defaultTag     string
	removeTagField bool
}

func (s *sink) Write(ctx *core.Context, t *core.Tuple) error {
	var tag string
	if v, ok := t.Data[s.tagField]; ok {
		tag, _ = data.ToString(v)
	}
	if tag == "" {
		tag = s.defaultTag
	}

	// TODO: this seems very inefficient
	m := map[string]interface{}{}
	for k, v := range t.Data {
		m[k] = v // assuming any data.Type can be serialized as msgpack
	}

	if s.removeTagField {
		delete(m, s.tagField)
	}

	// TODO: this seems very inefficient
	return s.output.Emit([]fluentd_forwarder.FluentRecordSet{
		{
			Tag: tag,
			Records: []fluentd_forwarder.TinyFluentRecord{
				{
					Timestamp: uint64(t.Timestamp.Unix()),
					Data:      m,
				},
			},
		},
	})
}

func (s *sink) Close(ctx *core.Context) error {
	s.output.Stop()
	s.output.WaitForShutdown()
	return nil
}

// NewSink creates a new Sink for SensorBee which sends tuples to fluentd.
func NewSink(ctx *core.Context, params data.Map) (core.Sink, error) {
	conf := struct {
		// fluentd parameters
		forwardTo           string
		retryInterval       time.Duration
		connectionTimeout   time.Duration
		writeTimeout        time.Duration
		flushInterval       time.Duration
		journalGroupPath    string
		maxJournalChunkSize int64

		// plugin parameters
		tagField       string
		defaultTag     string
		removeTagField bool
	}{
		forwardTo:           "127.0.0.1:24224",
		retryInterval:       5 * time.Second,
		connectionTimeout:   10 * time.Second,
		writeTimeout:        10 * time.Second,
		flushInterval:       5 * time.Second,
		journalGroupPath:    "*", // TODO: check the meaning of this option carefully
		maxJournalChunkSize: 16777216,

		tagField:       "tag",
		defaultTag:     "sensorbee",
		removeTagField: true,
	}

	if v, ok := params["forward_to"]; ok {
		f, err := data.AsString(v)
		if err != nil {
			return nil, err
		}
		conf.forwardTo = f
	}

	if err := getDuration(params, "retry_interval", &conf.retryInterval); err != nil {
		return nil, err
	}
	if err := getDuration(params, "connection_timeout", &conf.connectionTimeout); err != nil {
		return nil, err
	}
	if err := getDuration(params, "write_timeout", &conf.writeTimeout); err != nil {
		return nil, err
	}
	if err := getDuration(params, "flush_interval", &conf.flushInterval); err != nil {
		return nil, err
	}

	if v, ok := params["journal_group_path"]; ok {
		p, err := data.AsString(v)
		if err != nil {
			return nil, err
		}
		conf.journalGroupPath = p
	}

	if v, ok := params["max_journal_chunk_size"]; ok {
		s, err := data.AsInt(v) // TODO: support something like '100MB'
		if err != nil {
			return nil, err
		}
		conf.maxJournalChunkSize = s
	}

	if v, ok := params["tag_field"]; ok {
		f, err := data.AsString(v)
		if err != nil {
			return nil, err
		}
		conf.tagField = f
	}

	if v, ok := params["default_tag"]; ok {
		t, err := data.AsString(v)
		if err != nil {
			return nil, err
		}
		conf.defaultTag = t
	}

	if v, ok := params["remove_tag_field"]; ok {
		r, err := data.ToBool(v)
		if err != nil {
			return nil, err
		}
		conf.removeTagField = r
	}

	// TODO: optionally support logging
	logger, err := logging.GetLogger("fluentd-forwarder")
	if err != nil {
		return nil, err
	}
	logger.SetBackend(logging.AddModuleLevel(logging.NewLogBackend(&nullWriter{}, "", 0)))

	out, err := fluentd_forwarder.NewForwardOutput(
		logger,
		conf.forwardTo,
		conf.retryInterval,
		conf.connectionTimeout,
		conf.writeTimeout,
		conf.flushInterval,
		conf.journalGroupPath,
		conf.maxJournalChunkSize)
	if err != nil {
		return nil, err
	}

	s := &sink{
		output:         out,
		tagField:       conf.tagField,
		defaultTag:     conf.defaultTag,
		removeTagField: conf.removeTagField,
	}
	out.Start()
	return s, nil
}

func getDuration(params data.Map, field string, d *time.Duration) error {
	v, ok := params[field]
	if !ok {
		return nil
	}

	// TODO: support other types
	str, err := data.AsString(v)
	if err != nil {
		return err
	}

	*d, err = time.ParseDuration(str)
	return err
}
