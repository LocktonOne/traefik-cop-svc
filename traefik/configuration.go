package traefik

type Configuration struct {
	HTTP HTTPConfiguration `json:"http,omitempty"`
}

type HTTPConfiguration struct {
	Routers  map[string]Router  `json:"routers,omitempty" toml:"routers,omitempty" yaml:"routers,omitempty"`
	Services map[string]Service `json:"services,omitempty" toml:"services,omitempty" yaml:"services,omitempty"`
}

type Service struct {
	LoadBalancer ServersLoadBalancer `json:"loadBalancer,omitempty" toml:"loadBalancer,omitempty" yaml:"loadBalancer,omitempty"`
}

type Router struct {
	Service string `json:"service,omitempty" toml:"service,omitempty" yaml:"service,omitempty"`
	Rule    string `json:"rule,omitempty" toml:"rule,omitempty" yaml:"rule,omitempty"`
}

type ServersLoadBalancer struct {
	Servers []Server `json:"servers,omitempty" toml:"servers,omitempty" yaml:"servers,omitempty" label-slice-as-struct:"server"`
}

type Server struct {
	URL    string `json:"url,omitempty" toml:"url,omitempty" yaml:"url,omitempty" label:"-"`
	Scheme string `toml:"-" json:"-" yaml:"-"`
	Port   string `toml:"-" json:"-" yaml:"-"`
}
