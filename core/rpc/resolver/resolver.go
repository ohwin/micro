package resolver

import "google.golang.org/grpc/resolver"

func Register(builder resolver.Builder) {
	resolver.Register(builder)
}
