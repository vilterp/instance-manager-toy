package pure_manager

// the dumbest possible strategy
func WipeAndRestart(instances []*Instance, desiredSpec GroupSpec) ActionNode {
	return Ser([]ActionNode{
		Wipe(instances),
		StartFromScratch(desiredSpec),
	})
}

func Wipe(instances []*Instance) ActionNode {
	var out []ActionNode
	for _, i := range instances {
		out = append(out, Unit(ShutDownInstance{i.ID}))
	}
	return Par(out)
}

func StartFromScratch(spec GroupSpec) ActionNode {
	var out []ActionNode
	for i := 0; i < spec.NumInstances; i++ {
		out = append(out, Unit(StartInstance{
			Spec: InstanceSpec{
				Version: spec.Version,
			},
		}))
	}
	return Par(out)
}

func RollingRestart(instances []*Instance, newVersion Version) ActionNode {
	var out []ActionNode
	for _, i := range instances {
		if i.Version == newVersion {
			continue
		}
		out = append(out, Unit(&RestartInstance{
			ID:         i.ID,
			NewVersion: newVersion,
		}))
	}
	return Ser(out)
}

func ReplaceInstance(id InstanceID, spec GroupSpec) ActionNode {
	return Par([]ActionNode{
		Unit(ShutDownInstance{id}),
		Unit(StartInstance{InstanceSpec{spec.Version}}),
	})
}
