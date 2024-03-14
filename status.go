package main

type Status string

const (
	UP       Status = "UP"
	STARTING Status = "STARTING"
	DOWN     Status = "DOWN"
)

func (s Status) String() string {
	return string(s)
}
