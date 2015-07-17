package dns_router

import ()

type NullHealthCheck struct {
}

func NewNullHealthCheck() *NullHealthCheck {
	hc := &NullHealthCheck{}
	return hc
}

func (hc *NullHealthCheck) Check(state []CheckState) error {
	return nil
}

func (hc *NullHealthCheck) BackendAlive() bool {
	return true
}

func (hc *NullHealthCheck) SetAlive(idx int, alive CheckState) error {
	return nil
}
