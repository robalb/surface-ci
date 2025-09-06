package pipeline

type Surface struct {
	Domains []string `yaml:"domains"`
	IPs     []string `yaml:"ips"`
	URLs    []string `yaml:"urls"`
}
