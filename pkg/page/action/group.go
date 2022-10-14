package action

import "encoding/gob"

// ActionGroup groups multiple actions as a single action. To use it, call Group() and pass a list of actions.
type ActionGroup struct {
	Actions []ActionI
}

// RenderScript renders the group of actions as a single action.
func (g ActionGroup) RenderScript(params RenderParams) (s string) {
	for _, a := range g.Actions {
		s += a.RenderScript(params)
	}
	return
}

// Group joins multiple actions into a single action.
func Group(actions ...ActionI) ActionGroup {
	var foundCallback bool
	for _, a := range actions {
		switch a.(type) {
		case ActionGroup:
			panic("You cannot put an ActionGroup into another ActionGroup")
		case CallbackActionI:
			if foundCallback {
				panic("You can only associate one callback action with an event, and it must be the last action.")
			}
			foundCallback = true
		default:
			if foundCallback {
				panic("You can only associate one callback action with an event, and it must be the last action.")
			}
		}
	}
	// Note, the above could be more robust and allow multiple callback actions, but it would get quite tricky if different
	// kinds of actions were interleaved. We will wait until someone presents a compelling need for something like that.

	return ActionGroup{actions}
}

// GetCallbackAction returns the embedded callback action in the group, if one exists. Note that
// you can only have at most one callback action in a group
func (g ActionGroup) GetCallbackAction() G_CallbackActionI {
	if g.Actions == nil || len(g.Actions) == 0 {
		return nil
	}
	a := g.Actions[len(g.Actions)-1]
	if a2, ok := a.(CallbackActionI); ok {
		return a2.(G_CallbackActionI)
	}
	return nil
}

func init() {
	gob.Register(ActionGroup{})
}
