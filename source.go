package fluentd

import (
	"github.com/sirupsen/logrus"
	"github.com/fluent/fluentd-forwarder"
	"github.com/op/go-logging"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"time"
)

type source struct {
	input *fluentd_forwarder.ForwardInput
	ctx   *core.Context
	w     core.Writer

	ioParams *bql.IOParams
	bind     string
	tagField string
}

func (s *source) Emit(rset []fluentd_forwarder.FluentRecordSet) error {
	now := time.Now().UTC()
	for _, rs := range rset {
		for _, r := range rs.Records {
			t := &core.Tuple{
				ProcTimestamp: now,
				Timestamp:     time.Unix(int64(r.Timestamp), 0),
			}

			m, err := data.NewMap(r.Data)
			if err != nil {
				s.ctx.ErrLog(err).WithFields(logrus.Fields{
					"source_type": s.ioParams.TypeName,
					"source_name": s.ioParams.Name,
					"data":        r.Data,
				}).Error("Cannot create a data.Map from the data")
				continue
			}
			m[s.tagField] = data.String(rs.Tag)

			t.Data = m
			if err := s.w.Write(s.ctx, t); err != nil {
				s.ctx.ErrLog(err).WithFields(logrus.Fields{
					"source_type": s.ioParams.TypeName,
					"source_name": s.ioParams.Name,
				}).Error("Cannot write a tuple")
			}
		}
	}
	return nil
}

func (s *source) GenerateStream(ctx *core.Context, w core.Writer) error {
	// Because the input isn't running yet, it's safe to set params here.
	s.ctx = ctx
	s.w = w

	s.input.Start()
	s.input.WaitForShutdown()
	return nil
}

func (s *source) Stop(ctx *core.Context) error {
	s.input.Stop()
	s.input.WaitForShutdown()
	return nil
}

type nullWriter struct {
}

func (n *nullWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

// NewSource create a new Source receiving data from fluentd's out_forward.
func NewSource(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Source, error) {
	s := &source{
		ioParams: ioParams,
		bind:     "127.0.0.1:24224",
		tagField: "tag",
	}

	if v, ok := params["bind"]; ok {
		b, err := data.AsString(v)
		if err != nil {
			return nil, err
		}
		s.bind = b
	}

	if v, ok := params["tag_field"]; ok {
		t, err := data.AsString(v)
		if err != nil {
			return nil, err
		}
		s.tagField = t
	}

	// TODO: optionally support logging
	logger, err := logging.GetLogger("fluentd-forwarder")
	if err != nil {
		return nil, err
	}
	logger.SetBackend(logging.AddModuleLevel(logging.NewLogBackend(&nullWriter{}, "", 0)))

	input, err := fluentd_forwarder.NewForwardInput(logger, s.bind, s)
	if err != nil {
		return nil, err
	}

	s.input = input
	return s, nil
}
