package server

import (
	"fmt"

	"github.com/vilterp/instance-manager-toy/proto"
	"github.com/vilterp/instance-manager-toy/taskgraph"
)

type builder struct {
	b *taskgraph.SpecBuilder
}

func newBuilder() *builder {
	return &builder{
		b: taskgraph.NewSpecBuilder(),
	}
}

func (b *builder) Build() *proto.TaskGraphSpec {
	return b.b.Build()
}

// the dumbest possible strategy
func (b *builder) WipeAndRestart(nodes []*proto.Node, desiredSpec *proto.GroupSpec) proto.TaskID {
	w := b.Wipe(nodes)
	s := b.StartFromScratch(desiredSpec)
	return b.b.SerIDs("WipeAndRestart", []proto.TaskID{
		w,
		s,
	}).Tail
}

func (b *builder) KillSome(nodes []*proto.Node, n int64) (proto.TaskID, error) {
	if int64(len(nodes)) < n {
		return "", fmt.Errorf("not enough nodes to shut down")
	}
	var out []proto.TaskID
	for i := int64(0); i < n; i++ {
		node := nodes[i]
		out = append(out, b.b.Unit(proto.MkShutDown(proto.NodeID(node.Id))))
	}
	return b.b.ParIDs("KillSome", out), nil
}

func (b *builder) Wipe(nodes []*proto.Node) proto.TaskID {
	var out []proto.TaskID
	for _, i := range nodes {
		out = append(out, b.b.Unit(proto.MkShutDown(proto.NodeID(i.Id))))
	}
	return b.b.ParIDs("Wipe", out)
}

func (b *builder) StartFromScratch(spec *proto.GroupSpec) proto.TaskID {
	return b.StartNodes(spec.NumInstances, spec.Version)
}

func (b *builder) StartNodes(n int64, v int64) proto.TaskID {
	var out []proto.TaskID
	for i := int64(0); i < n; i++ {
		out = append(out, b.b.Unit(proto.MkStartNode(&proto.NodeSpec{
			Version: int64(v),
		})))
	}
	return b.b.ParIDs("StartNodes", out)
}

// TODO: max surge, max unavailable, etc
func (b *builder) RollingUpgrade(nodes []*proto.Node, newSpec *proto.NodeSpec) proto.TaskID {
	var out []taskgraph.TaskChain
	for _, node := range nodes {
		shutDown := b.b.Unit(proto.MkShutDown(proto.NodeID(node.Id)))
		start := b.b.Unit(proto.MkStartNode(newSpec))
		restart := b.b.SerIDs("Restart", []proto.TaskID{shutDown, start})
		out = append(out, restart)
	}
	return b.b.SerChains(out).Tail
}
