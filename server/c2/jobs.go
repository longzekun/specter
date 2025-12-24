package c2

import (
	"fmt"

	"github.com/longzekun/specter/server/constants"
	"github.com/longzekun/specter/server/core"
)

func StartMTLSListenerJob(host string, port uint32) {
	bind := fmt.Sprintf("%s:%d", host, port)

	//	开始监听
	ln, err := StartMutualListener(host, uint16(port))
	if err != nil {
		return
	}

	job := &core.Job{
		ID:          core.NextJobID(),
		Name:        constants.MtlsStr,
		Description: fmt.Sprintf("mutual tls listener %s", bind),
		Protocol:    constants.TCPListenerStr,
		Port:        uint16(port),
		JobCtl:      make(chan bool),
	}

	go func() {
		<-job.JobCtl

		//	关闭监听
		ln.Close()

		core.Jobs.Remove(job)
	}()

	core.Jobs.Add(job)
}
