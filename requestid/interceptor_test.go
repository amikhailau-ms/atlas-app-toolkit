package requestid

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type testRequest struct{}

type testResponse struct{}

type mockServerTransportStream struct {
	ctx context.Context
}

func (*mockServerTransportStream) Method() string {
	return "unimplemente"
}

func (*mockServerTransportStream) SetHeader(metadata.MD) error {
	return nil
}

func (*mockServerTransportStream) SendHeader(metadata.MD) error {
	return nil
}

func (*mockServerTransportStream) SetTrailer(metadata.MD) error { return nil }

func (m *mockServerTransportStream) Context() context.Context {
	return m.ctx
}

func (*mockServerTransportStream) SendMsg(m interface{}) error {
	return nil
}

func (*mockServerTransportStream) RecvMsg(m interface{}) error {
	return nil
}

type mockServerStream struct {
	ctx context.Context
}

func (*mockServerStream) Method() string {
	return "unimplemente"
}

func (*mockServerStream) SetHeader(metadata.MD) error {
	return nil
}

func (*mockServerStream) SendHeader(metadata.MD) error {
	return nil
}

func (*mockServerStream) SetTrailer(metadata.MD) {}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

func (*mockServerStream) SendMsg(m interface{}) error {
	return nil
}

func (*mockServerStream) RecvMsg(m interface{}) error {
	return nil
}

func DummyContextWithServerTransportStream() context.Context {
	expectedStream := &mockServerTransportStream{}
	return grpc.NewContextWithServerTransportStream(context.Background(), expectedStream)
}

func TestUnaryServerInterceptorWithoutRequestId(t *testing.T) {
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		reqID, exists := FromContext(ctx)
		if exists && reqID == "" {
			t.Errorf("requestId must be generated by interceptor")
		}
		return &testResponse{}, nil
	}
	ctx := DummyContextWithServerTransportStream()
	_, err := UnaryServerInterceptor()(ctx, testRequest{}, nil, handler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnaryServerInterceptorWithDummyRequestId(t *testing.T) {
	dummyRequestID := newRequestID()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		reqID, exists := FromContext(ctx)
		if !exists || reqID != dummyRequestID {
			t.Errorf("expected requestID: %q, returned requestId: %q", dummyRequestID, reqID)
		}
		return &testResponse{}, nil
	}
	ctx := DummyContextWithServerTransportStream()
	md := metadata.Pairs(DefaultRequestIDKey, dummyRequestID)
	newCtx := metadata.NewIncomingContext(ctx, md)
	_, err := UnaryServerInterceptor()(newCtx, testRequest{}, nil, handler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStreamServerInterceptorWithoutRequestId(t *testing.T) {
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		reqID, exists := FromContext(stream.Context())
		if exists && reqID == "" {
			t.Errorf("requestId must be generated by interceptor")
		}
		return nil
	}
	ctx := DummyContextWithServerTransportStream()
	ms := &mockServerStream{ctx}
	err := StreamServerInterceptor()(testRequest{}, ms, nil, handler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStreamServerInterceptorWithDummyRequestId(t *testing.T) {
	dummyRequestID := newRequestID()
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		reqID, exists := FromContext(stream.Context())
		if !exists || reqID != dummyRequestID {
			t.Errorf("expected requestID: %q, returned requestId: %q", dummyRequestID, reqID)
		}
		return nil
	}
	ctx := DummyContextWithServerTransportStream()
	md := metadata.Pairs(DefaultRequestIDKey, dummyRequestID)
	newCtx := metadata.NewIncomingContext(ctx, md)
	streamInterceptor := StreamServerInterceptor()
	if err := streamInterceptor(testRequest{}, &mockServerStream{newCtx}, nil, handler); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
