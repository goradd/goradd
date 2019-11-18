This docker deployment does the basics of pulling a minimal alpine container, and then
building your app and copying it to the container. Additional containers you will need will depend on what database
you are using, how you are handling your session and pagestate stores, etc.

To do the build, run

go generate buildAll.go
