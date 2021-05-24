package herbplugin

type Location struct {
	Path string
}

func NewLoaction() *Location {
	return &Location{}
}

type Trusted struct {
	Paths   []string
	Domains []string
}

func NewTrusted() *Trusted {
	return &Trusted{}
}

type Options struct {
	Location *Location
	Params   map[string]string
	Trusted  *Trusted
}

func NewOptions() *Options {
	return &Options{
		Location: NewLoaction(),
		Params:   map[string]string{},
		Trusted:  NewTrusted(),
	}
}
