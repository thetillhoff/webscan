package types

// Option is a functional option that configures a value of type C.
type Option[C any] func(*C)

// ApplyOptions applies all options to cfg in order.
func ApplyOptions[C any](cfg *C, opts []Option[C]) {
	for _, opt := range opts {
		opt(cfg)
	}
}
