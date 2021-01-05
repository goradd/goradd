This docker deployment does the basics of pulling a minimal alpine container, and then
building your app and copying it to the container. Additional containers you will need will depend on what database
you are using, how you are handling your session and pagestate stores, etc.

To do the build, run

go generate buildAll.go

The db.cfg file is the place to put your database credentials and customizations.
The default uses host.docker.internal as the address so that the docker image can work on
a Mac or Windows. For deployment, you should instead use one that defines your
dbname, user, and password for your server. Make sure you set the permissions on this file
to only allow the process that starts the file to be able to read the sensitive data in the file.
