package db

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/vilterp/instance-manager-toy/proto"
)

type GroupSpecDB interface {
	UpdateSpec(spec *proto.GroupSpec)
	GetCurrent() *proto.GroupSpecInfo
	GetHistory() []*proto.GroupSpecInfo
}

func NewSpecDB(initial *proto.GroupSpec) *mockGroupSpecDB {
	return &mockGroupSpecDB{
		current: &proto.GroupSpecInfo{
			Spec:      initial,
			CreatedAt: ptypes.TimestampNow(),
		},
	}
}

type mockGroupSpecDB struct {
	current *proto.GroupSpecInfo
	history []*proto.GroupSpecInfo
}

var _ GroupSpecDB = &mockGroupSpecDB{}

func (mgs *mockGroupSpecDB) UpdateSpec(newSpec *proto.GroupSpec) {
	mgs.history = append(mgs.history, mgs.current)
	mgs.current = &proto.GroupSpecInfo{
		Spec:      newSpec,
		CreatedAt: ptypes.TimestampNow(),
	}
}

func (mgs *mockGroupSpecDB) GetCurrent() *proto.GroupSpecInfo {
	return mgs.current
}

func (mgs *mockGroupSpecDB) GetHistory() []*proto.GroupSpecInfo {
	return mgs.history
}
