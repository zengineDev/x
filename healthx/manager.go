package healthx

type HealthCheckResult map[string]bool

type HealthCheckManager struct {
	checks []HealthCheckContract
	failed bool
}

func NewManager() HealthCheckManager {
	return HealthCheckManager{failed: false}
}

func (c *HealthCheckManager) Use(check HealthCheckContract) {
	c.checks = append(c.checks, check)
}

func (c *HealthCheckManager) Run() HealthCheckResult {
	result := HealthCheckResult{}

	for _, check := range c.checks {
		r := check.Run()
		result[r.CheckID()] = r.IsFailed()
		c.failed = r.IsFailed()
	}

	return result

}

func (c *HealthCheckManager) SomethingFailed() bool {
	return c.failed
}

type HealthCheckContract interface {
	Run() HealthCheckStatusContract
}

type HealthCheckStatusContract interface {
	Status() string
	IsFailed() bool
	CheckID() string
}

type HealthCheckStatus struct {
	StatusText string
	Failed     bool
	ID         string
}

func (s HealthCheckStatus) Status() string {
	return s.StatusText
}

func (s HealthCheckStatus) IsFailed() bool {
	return s.Failed
}

func (s HealthCheckStatus) CheckID() string {
	return s.ID
}
