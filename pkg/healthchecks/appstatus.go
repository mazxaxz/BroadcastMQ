package healthchecks

type AppStatus struct {
	Host       string                `json:"host"`
	Components []*AppStatusComponent `json:"components"`
}

type AppStatusComponent struct {
	Name   string       `json:"name"`
	Status HealthStatus `json:"status"`
}

// AddComponent adds new component to hc registry
func (as *AppStatus) AddComponent(name string, status HealthStatus) {
	if as.Components == nil {
		as.Components = make([]*AppStatusComponent, 0)
	}

	as.Components = append(as.Components, &AppStatusComponent{name, status})
}

// IsAnyUnhealthy returns true if any component is unhealthy
func (as *AppStatus) IsAnyUnhealthy() bool {
	if as.Components == nil {
		return true
	}

	for _, component := range as.Components {
		if component.Status == NotOk {
			return true
		}
	}

	return false
}
