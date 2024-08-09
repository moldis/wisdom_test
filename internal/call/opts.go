package call

type OptFn func(*Caller)

func WithServerAddr(address string) OptFn {
	return func(c *Caller) {
		c.addr = address
	}
}

func WithProtocol(protocol string) OptFn {
	return func(c *Caller) {
		c.protocol = protocol
	}
}
