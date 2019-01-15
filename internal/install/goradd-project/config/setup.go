package config

// We use init to setup some global variables used by the framework
// We can make sure this gets called first by importing config with an underscore

func init() {
	initDatabases()
	initGoradd()
	initApp()
}
