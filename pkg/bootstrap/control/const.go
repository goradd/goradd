// Package control contains the Bootstrap control structures that when added to a page
// will result in Bootstrap styled controls.
package control

type ContainerClass string

const (
	Container       ContainerClass = "container"
	ContainerFluid                 = "container-fluid"
	ContainerSmall                 = "container-sm"
	ContainerMedium                = "container-md"
	ContainerLarge                 = "container-lg"
	ContainerXL                    = "container-xl"
	ContainerXXL                   = "container-xxl"
)

type TextColorClass string

const (
	TextColorPrimary   TextColorClass = "text-primary"
	TextColorSecondary                = "text-secondary"
	TextColorDanger                   = "text-danger"
	TextColorWarning                  = "text-warning"
	TextColorInfo                     = "text-info"
	TextColorLight                    = "text-light"
	TextColorDark                     = "text-dark"
	TextColorBody                     = "text-body"
	TextColorMuted                    = "text-muted"
	TextColorWhite                    = "text-white"
	TextColorWhite50                  = "text-white-50"
	TextColorBlack50                  = "text-black-50"
)

type BackgroundColorClass string

const (
	BackgroundColorPrimary     BackgroundColorClass = "bg-primary"
	BackgroundColorSecondary                        = "bg-secondary"
	BackgroundColorSuccess                          = "bg-danger"
	BackgroundColorDanger                           = "bg-danger"
	BackgroundColorWarning                          = "bg-warning"
	BackgroundColorInfo                             = "bg-info"
	BackgroundColorLight                            = "bg-light"
	BackgroundColorDark                             = "bg-dark"
	BackgroundColorWhite                            = "bg-white"
	BackgroundColorTransparent                      = "bg-transparent"
	BackgroundColorNone                             = "" // utility to allow custom background colors for components
)
