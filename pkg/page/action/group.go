package action

import "encoding/gob"

// actionGroup groups multiple actions as a single action. To use it, call Group() and pass a list of actions.
type actionGroup struct {
	Actions []ActionI
}

type GroupI interface {
	GetCallbackAction() CallbackActionI
}

// RenderScript renders the group of actions as a single action.
func (g actionGroup) RenderScript(params RenderParams) (s string) {
	for _, a := range g.Actions {
		s += a.RenderScript(params)
	}
	return
}

// Group joins multiple actions into a single action.
// Any number of javascript actions can be included, but only one callback action can be included in the group.
func Group(actions ...ActionI) ActionI {
	var foundCallback bool
	for _, a := range actions {
		switch a.(type) {
		case actionGroup:
			panic("You cannot put an actionGroup into another actionGroup")
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

	return actionGroup{actions}
}

// GetCallbackAction returns the embedded callback action in the group, if one exists. Note that
// you can only have at most one callback action in a group
func (g actionGroup) GetCallbackAction() CallbackActionI {
	if g.Actions == nil || len(g.Actions) == 0 {
		return nil
	}
	a := g.Actions[len(g.Actions)-1]
	if a2, ok := a.(CallbackActionI); ok {
		return a2
	}
	return nil
}

func init() {
	gob.Register(actionGroup{})
}
