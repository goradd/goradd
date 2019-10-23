package action

import "encoding/gob"

// ActionGroup groups multiple actions as a single action. To use it, call Group() and pass a list of actions.
type ActionGroup struct {
	Actions []ActionI
}

// ΩRenderScript renders the group of actions as a single action.
func (g ActionGroup) ΩRenderScript(params ΩrenderParams) (s string) {
	for _,a := range g.Actions {
		s += a.ΩRenderScript(params)
	}
	return
}

// Call Group to join multiple actions into a single action.
func Group(actions ...ActionI) ActionGroup {
	var foundCallback bool
	for _,a := range actions {
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

func (g ActionGroup) HasServerAction() bool {
	if a := g.GetCallbackAction(); a != nil {
		return a.IsServerAction()
	}
	return false
}

func (g ActionGroup) HasCallbackAction() bool {
	if a := g.GetCallbackAction(); a != nil {
		return true
	}
	return false
}


// GetCallbackAction returns the embedded callback action in the group, if one exists. Note that
// you can only have at most one callback action in a group
func (g ActionGroup) GetCallbackAction() CallbackActionI {
	if g.Actions == nil || len(g.Actions) == 0 {
		return nil
	}
	a := g.Actions[len(g.Actions) - 1]
	if a2, ok := a.(CallbackActionI); ok {
		return a2
	}
	return nil
}

func init() {
	gob.Register(ActionGroup{})
}