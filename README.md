# Combine Go contexts

This package provides the `contexts.Combine` function, which combines multiple [contexts](https://pkg.go.dev/context) into a single context. The combined context is Done if any of its contexts are Done. It returns the Err of the first context to become Done. The combined Deadline is the oldest Deadline returned by the contexts. The combined Value method calls Value on each of the contexts and returns the first non-nil result.

# Install

```
go get github.com/jbarnette/contexts
```

## Example

Here's how to inject a combined context into a Go gRPC server:

```go
// The gRPC server derives its stream contexts from context.Background, so it's difficult
// to make shared context values available to the RPC handlers.
base := context.WithValue(context.Background(), "sharedKey", "sharedValue")

// This interceptor combines a base context with the stream context, which makes the
// base context's values available to the RPC handler without disturbing any values in
// the original stream context.
interceptor := func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
  return handler(contexts.Combine(base, ctx), req)
}

// The server's unary handlers will receive a combined context.
grpc.NewServer(grpc.UnaryInterceptor(interceptor))
```
