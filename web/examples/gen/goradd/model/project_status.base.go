// Code generated by GoRADD. DO NOT EDIT.

package model

import (
	//"log"
	//"github.com/goradd/goradd/pkg/orm/query"
	"strconv"
)

const (
	ProjectStatusOpen      ProjectStatus = 1
	ProjectStatusCancelled ProjectStatus = 2
	ProjectStatusCompleted ProjectStatus = 3
	ProjectStatusPlanned   ProjectStatus = 4
)

// ProjectStatusMaxValue is the maximum enumerated value of ProjectStatus
// doc: type=ProjectStatus
const ProjectStatusMaxValue = 4

type ProjectStatus int

// String returns the name value of the type and satisfies the fmt.Stringer interface
func (p ProjectStatus) String() string {
	switch p {
	case 0:
		return ""
	case 1:
		return "Open"
	case 2:
		return "Cancelled"
	case 3:
		return "Completed"
	case 4:
		return "Planned"
	default:
		panic("index out of range")
	}
	return "" // prevent warning
}

// ID returns a string representation of the id and satisfies the IDer interface
func (p ProjectStatus) ID() string {
	return strconv.Itoa(int(p))
}

// ProjectStatusFromID converts a ProjectStatus ID to a ProjectStatus
func ProjectStatusFromID(id string) ProjectStatus {
	switch id {
	case "1":
		return ProjectStatus(1)
	case "2":
		return ProjectStatus(2)
	case "3":
		return ProjectStatus(3)
	case "4":
		return ProjectStatus(4)
	}
	return ProjectStatus(0)
}

// ProjectStatusesFromIDs converts a slice of ProjectStatus IDs to a slice of ProjectStatus
func ProjectStatusesFromIDs(ids []string) (values []ProjectStatus) {
	values = make([]ProjectStatus, 0, len(ids))
	for _, id := range ids {
		values = append(values, ProjectStatusFromID(id))
	}
	return
}

// ProjectStatusFromName converts a ProjectStatus name to a ProjectStatus
func ProjectStatusFromName(name string) ProjectStatus {
	switch name {
	case "Open":
		return ProjectStatus(1)
	case "Cancelled":
		return ProjectStatus(2)
	case "Completed":
		return ProjectStatus(3)
	case "Planned":
		return ProjectStatus(4)
	}
	return ProjectStatus(0)
}

// AllProjectStatuses returns a slice of all the ProjectStatus values.
func AllProjectStatuses() (values []ProjectStatus) {
	values = append(values, 1)
	values = append(values, 2)
	values = append(values, 3)
	values = append(values, 4)
	return
}

// AllProjectStatusesI returns a slice of all the ProjectStatus values as generic interfaces.
// doc: type=ProjectStatus
func AllProjectStatusesI() (values []any) {
	values = make([]interface{}, 4, 4)
	values[0] = ProjectStatus(1)
	values[1] = ProjectStatus(2)
	values[2] = ProjectStatus(3)
	values[3] = ProjectStatus(4)
	return
}

// Label returns the string that will be displayed to a user for this item. Together with
// the Value function, it satisfies the ItemLister interface that makes it easy
// to create a dropdown list of items.
func (p ProjectStatus) Label() string {
	return p.String()
}

// Value returns the value that will be used in dropdown lists and satisfies the
// Valuer and ItemLister interfaces.
func (p ProjectStatus) Value() interface{} {
	return p.ID()
}

func (p ProjectStatus) Name() string {
	switch p {
	case 0:
		return ""
	case 1:
		return "Open"
	case 2:
		return "Cancelled"
	case 3:
		return "Completed"
	case 4:
		return "Planned"
	default:
		panic("Index out of range")
	}
	return "" // prevent warning
}

func (p ProjectStatus) Description() string {
	switch p {
	case 0:
		return ""
	case 1:
		return "The project is currently active"
	case 2:
		return "The project has been canned"
	case 3:
		return "The project has been completed successfully"
	case 4:
		return "Project is in the planning stages and has not been assigned a manager"
	default:
		panic("Index out of range")
	}
	return "" // prevent warning
}

func (p ProjectStatus) Guidelines() string {
	switch p {
	case 0:
		return ""
	case 1:
		return "All projects that we are working on should be in this state"
	case 2:
		return ""
	case 3:
		return "Celebrate successes!"
	case 4:
		return "Get ready"
	default:
		panic("Index out of range")
	}
	return "" // prevent warning
}

func (p ProjectStatus) IsActive() bool {
	switch p {
	case 0:
		return false
	case 1:
		return true
	case 2:
		return true
	case 3:
		return true
	case 4:
		return false
	default:
		panic("Index out of range")
	}
	return false // prevent warning
}

// ProjectStatusNames returns a slice of all the Names associated with ProjectStatus values.
// doc: type=ProjectStatus
func ProjectStatusNames() []string {
	return []string{
		// 0 item will be a blank
		"",
		"Open",
		"Cancelled",
		"Completed",
		"Planned",
	}
}

// ProjectStatusDescriptions returns a slice of all the Descriptions associated with ProjectStatus values.
// doc: type=ProjectStatus
func ProjectStatusDescriptions() []string {
	return []string{
		// 0 item will be a blank
		"",
		"The project is currently active",
		"The project has been canned",
		"The project has been completed successfully",
		"Project is in the planning stages and has not been assigned a manager",
	}
}

// ProjectStatusGuidelines returns a slice of all the Guidelines associated with ProjectStatus values.
// doc: type=ProjectStatus
func ProjectStatusGuidelines() []string {
	return []string{
		// 0 item will be a blank
		"",
		"All projects that we are working on should be in this state",
		"",
		"Celebrate successes!",
		"Get ready",
	}
}

// ProjectStatusIsActives returns a slice of all the IsActives associated with ProjectStatus values.
// doc: type=ProjectStatus
func ProjectStatusIsActives() []bool {
	return []bool{
		// 0 item will be a blank
		false,
		true,
		true,
		true,
		false,
	}
}

func (p ProjectStatus) Get(key string) interface{} {

	switch key {
	case "Name":
		return p.Name()
	case "Description":
		return p.Description()
	case "Guidelines":
		return p.Guidelines()
	case "IsActive":
		return p.IsActive()
	default:
		return nil
	}
}
