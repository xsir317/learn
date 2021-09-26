package configs

type Conf struct {
	Ip   string
	Port int
}

func NewConf(ip string, port int) *Conf {
	return &Conf{
		Ip:   ip,
		Port: port}
}
