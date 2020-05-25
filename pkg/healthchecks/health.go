package healthchecks

type HealthStatus byte

const (
	NotOk HealthStatus = iota
	Ok
)
