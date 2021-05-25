package goutils

import (
	"fmt"
	"strconv"
)

const (
	B = 1
	K = 1024 * B
	M = 1024 * K
	G = 1024 * M
	T = 1024 * G
	P = 1024 * T
	E = 1024 * P
)

type DiskSize uint64

func (s DiskSize) String() string {
	n := float64(s)
	if s < K {
		return strconv.Itoa(int(s))
	} else if s < M {
		return fmt.Sprintf("%.2fK", n/K)
	} else if s < G {
		return fmt.Sprintf("%.2fM", n/M)
	} else if s < T {
		return fmt.Sprintf("%.2fG", n/G)
	} else if s < P {
		return fmt.Sprintf("%.2fT", n/T)
	} else if s < E {
		return fmt.Sprintf("%.2fP", n/P)
	} else {
		return fmt.Sprintf("%.2fE", n/E)
	}
}
