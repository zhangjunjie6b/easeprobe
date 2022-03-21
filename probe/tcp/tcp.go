package tcp

import (
	"fmt"
	"net"
	"time"

	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/probe"
	log "github.com/sirupsen/logrus"
)

// Kind is the type
const Kind string = "tcp"

// TCP implements a config for TCP
type TCP struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`

	//Control Option
	Timeout      time.Duration `yaml:"timeout,omitempty"`
	TimeInterval time.Duration `yaml:"interval,omitempty"`

	result *probe.Result `yaml:"-"`
}

// Kind return the HTTP kind
func (t *TCP) Kind() string {
	return Kind
}

// Interval get the interval
func (t *TCP) Interval() time.Duration {
	return t.TimeInterval
}

// Result get the probe result
func (t *TCP) Result() *probe.Result {
	return t.result
}

// Config HTTP Config Object
func (t *TCP) Config(gConf global.ProbeSettings) error {

	t.Timeout = gConf.NormalizeTimeOut(t.Timeout)
	t.TimeInterval = gConf.NormalizeInterval(t.TimeInterval)

	t.result = probe.NewResult()
	t.result.Endpoint = t.Host
	t.result.Name = t.Name
	t.result.PreStatus = probe.StatusInit
	t.result.TimeFormat = gConf.TimeFormat

	log.Debugf("[%s] configuration: %+v, %+v", t.Kind(), t, t.Result())
	return nil
}

// Probe return the checking result
func (t *TCP) Probe() probe.Result {

	now := time.Now()
	t.result.StartTime = now
	t.result.StartTimestamp = now.UnixMilli()

	conn, err := net.DialTimeout("tcp", t.Host, t.Timeout)
	t.result.RoundTripTime.Duration = time.Since(now)
	status := probe.StatusUp
	if err != nil {
		t.result.Message = fmt.Sprintf("Error: %v", err)
		log.Errorf("error: %v", err)
		status = probe.StatusDown
	} else {
		t.result.Message = "TCP Connection Established Successfully!"
		conn.Close()
	}
	t.result.PreStatus = t.result.Status
	t.result.Status = status

	t.result.DoStat(t.TimeInterval)

	return *t.result
}
