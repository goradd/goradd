The scripts in this directory copy and process static files that your application will serve
through Go's embed system.

The scripts here will copy asset files from the goradd-project/web asset directories, and
the goradd installation directory, to the deploy directory.
These files get bundled into the application binary and served directly
from there, which offers easy deployment and some amount of additional security, as the files are
not easily changed.

The scripts will minify and compress files as needed so that they can be served
as fast as possible.

The scripts here will compress files into both gzip and brotli formats, which are generally supported
by all browsers. Brotli files are usually smaller and faster to decompress, so if the web browser supports
it, Goradd will prefer brotli.

You can make your application smaller by modifying these script to only stage gzip or brotli files.
If the client only supports a type of compression that is not found, Goradd will decompress the file
before sending it.

You can also modify the scripts to copy the originals without compression, and those will be served.

Files are staged in the goradd-project/deploy/embed directory
