// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package api

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"

	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/MetalBlockchain/metalgo/trace"
)

var _ http.Handler = (*tracedHandler)(nil)

type tracedHandler struct {
	h            http.Handler
	serveHTTPTag string
	tracer       trace.Tracer
}

func TraceHandler(h http.Handler, name string, tracer trace.Tracer) http.Handler {
	return &tracedHandler{
		h:            h,
		serveHTTPTag: fmt.Sprintf("%s.ServeHTTP", name),
		tracer:       tracer,
	}
}

func (h *tracedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := h.tracer.Start(ctx, h.serveHTTPTag, oteltrace.WithAttributes(
		attribute.String("method", r.Method),
		attribute.String("url", r.URL.Redacted()),
		attribute.String("proto", r.Proto),
		attribute.String("host", r.Host),
		attribute.String("remoteAddr", r.RemoteAddr),
		attribute.String("requestURI", r.RequestURI),
	))
	defer span.End()

	r = r.WithContext(ctx)
	h.h.ServeHTTP(w, r)
}
