package postgresqlcnpgio


// Sleep represents a duration that the container should sleep.
type PoolerSpecTemplateSpecEphemeralContainersLifecyclePreStopSleep struct {
	// Seconds is the number of seconds to sleep.
	Seconds *float64 `field:"required" json:"seconds" yaml:"seconds"`
}

