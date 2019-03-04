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

func NewSpecDB(initial GroupSpec) *mockGroupSpecDB {
	return &mockGroupSpecDB{
		current: initial,
	}
}

type mockGroupSpecDB struct {
	current GroupSpec
	history []GroupSpec
}

var _ GroupSpecDB = &mockGroupSpecDB{}

func (mgs *mockGroupSpecDB) UpdateSpec(newSpec GroupSpec) {
	mgs.history = append(mgs.history, mgs.current)
	mgs.current = newSpec
}

func (mgs *mockGroupSpecDB) GetCurrent() GroupSpec {
	return mgs.current
}

func (mgs *mockGroupSpecDB) GetHistory() []GroupSpec {
	return mgs.history
}
