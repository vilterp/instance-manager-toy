package pure_manager

type GroupSpec struct {
	NumInstances int
	Version      Version
}

type GroupSpecDB interface {
	UpdateSpec(spec GroupSpec)
	GetCurrent() GroupSpec
	GetHistory() []GroupSpec
}
