package taskpool

const (
	defaultThreshold = 1
)

type Conf struct {
	// 扩容的阈值，当 len(task chan) > Threshold时会新建goroutine
	Threshold int32
}

func NewConf() *Conf {
	return &Conf{
		Threshold: defaultThreshold,
	}
}
