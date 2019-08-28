package action

type ActionGroup struct {
	Actions []ActionI
}

func (g ActionGroup) ΩRenderScript(params ΩrenderParams) (s string) {
	for _,a := range g.Actions {
		s += a.ΩRenderScript(params)
	}
	return
}

func Group(actions ...ActionI) ActionGroup {
	return ActionGroup{actions}
}

// TODO: Serialize