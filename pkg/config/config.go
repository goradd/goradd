package config

import "fmt"

// Not used

//var service ConfigService

type serviceType struct {
	appStarted bool
	params map [string]interface{}
}

type ConfigService interface {
	Get(name string) interface{}
	Set(name string, value interface{})
	AppStarted()
}

func (s *serviceType) Get(name string) interface{} {
	return s.params[name]
}

func (s *serviceType) Set(name string, value interface{})  {
	if s.appStarted {
		panic(fmt.Errorf("you cannot set a config value after the application has started. Attempting to set %s to %v", name, value))
	}
	s.params[name] = value
}

func (s *serviceType) AppStarted() {
	s.appStarted = true
}

// Allow application to use a different service type.
func SetService(c ConfigService) {
	//service = c
}


func init() {
	//service = new (serviceType)	// default is to start our service at init time. This can be replaced.
}