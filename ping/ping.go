package ping

import (
	probing "github.com/prometheus-community/pro-bing"
)

func PingAdress(addr string, errChan chan error) (func(), func() error) {
	pinger, err := probing.NewPinger(addr)
	if err != nil {
		panic(err)
	}

	pinger.OnSendError = func(_ *probing.Packet, err error) {
		errChan <- err
	}

	pinger.SetLogger(nil)

	return pinger.Stop, pinger.Run
}
