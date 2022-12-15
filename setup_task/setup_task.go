package setuptask

type SetupTask interface {
	ToAtomic() []AtomicSetupTask
}

type AtomicSetupTask interface {
	GetTerraform() string
	GetType() string
}
