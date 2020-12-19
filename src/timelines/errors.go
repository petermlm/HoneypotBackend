package timelines

import "fmt"

type InvalidIP7 struct {
	IP string
}

func newInvalidIP7(ip string) error {
	return &InvalidIP7{ip}
}

func (e *InvalidIP7) Error() string {
	return fmt.Sprintf("Invalid IP: %s", e.IP)
}

type CantMarshalConnAttemp struct {
	Err error
}

func newCantMarshalConnAttemp(err error) error {
	return &CantMarshalConnAttemp{err}
}

func (e *CantMarshalConnAttemp) Error() string {
	return e.Err.Error()
}

type CantUnmarshalConnAttemp struct {
	Err error
}

func newCantUnarshalConnAttemp(err error) error {
	return &CantUnmarshalConnAttemp{err}
}

func (e *CantUnmarshalConnAttemp) Error() string {
	return e.Err.Error()
}

type CantSeparateAddrError struct {
	addr string
}

func newCantSeparateAddrError(addr string) error {
	return &CantSeparateAddrError{addr: addr}
}

func (e *CantSeparateAddrError) Error() string {
	return fmt.Sprintf("Addr can't be separated by ':', %s", e.addr)
}

type InvalidRange struct {
	RangeValue string
}

func newInvalidRange(rangeValue string) error {
	return &InvalidRange{rangeValue}
}

func (e *InvalidRange) Error() string {
	return fmt.Sprintf("Invalid Flux Range: %s", e.RangeValue)
}
