package entities

const (
	// NodeAvailabilityActive ACTIVE
	NodeAvailabilityActive = "active"
	// NodeAvailabilityPause PAUSE
	NodeAvailabilityPause = "pause"
	// NodeAvailabilityDrain DRAIN
	NodeAvailabilityDrain = "drain"

	// NodeStateUnknown UNKNOWN
	NodeStateUnknown = "unknown"
	// NodeStateDown DOWN
	NodeStateDown = "down"
	// NodeStateReady READY
	NodeStateReady = "ready"
	// NodeStateDisconnected DISCONNECTED
	NodeStateDisconnected = "disconnected"
)

// Node represents swarm node
type Node struct {
	ID            string `json:"id"`
	Addr          string `json:"addr"`
	Hostname      string `json:"hostname"`
	Availability  string `json:"availability"`
	State         string `json:"state"`
	Role          string `json:"role"`
	NanoCPUs      int64  `json:"nano_cpus"`
	MemoryBytes   int64  `json:"memory_bytes"`
	EngineVersion string `json:"engine_version"`
}
