# Security

As a framework, GoRADD provides some security benefits out of the box, and also provides
a security architecture that developers can use to easily secure their websites, APIs,
and inputs, and meet regulatory requirements.

## General Security Issues
### HTTPS and SSL
To fully secure any website that is gathering data from users, a website should serve its
pages using SSL, so that it can deliver its pages using an https:// schema.

There are many ways to do this. Since the main() function is located in the goradd-project
folder, which is an area under control of the developer, the developer has the ability to
adapt the application to the particular circumstances of the deployment.

For example, the developer can take advantage of Go's built-in SSL capabilities by providing 
TLS certificate and key files, either as command-line options or configuration variables. 
Another way is to run the GoRADD application as an http service on a particular port, and 
then use a reverse proxy through a web proxy server like Apache or Nginx to manage the
SSL certificate. This method is useful for servers running virtual servers, or running the
application using a container service like Docker.

Customize the app startup process in the RunWebServer() function in goradd-project/web/app/app.go.

### Cross-site Scripting Attack Protection
Built-in to GoRADD is a CSRF token that ensures that the page that submits information is
the same page that was collecting the information. It is important to secure
any web page that is collecting information by delivering it over https, so that the
CSRF token cannot be intercepted. Even if the information collected does not first require
a login, it should still be protected with SSL encryption.

If the CSRF token is not correctly matched, an error will be displayed to the user, and
GoRADD will prevent the input from being sent to the application.

### SQL Injection Protection
All inputs to SQL databases are sent using mechanisms provided by the respective database drivers
which encode all values that might be used as a SQL injection attack. Typically, these mechanisms
require placeholders in SQL INSERT statements where values will get encoded and inserted by the driver,
vs. using code that directly inserts values into strings of SQL code.

See the FormatArgument() function in pkg/orm/db/sql/mysql/db.go and pkg/orm/db/sql/pgsql/db.go for examples
of the placeholder code that is used by each of the respective database drivers.

### Javascript and HTML Injection Protection
As a framework, GoRADD must be able to process any kind of input, including input that
could be legitimately accepting JavaScript code or HTML. As such, GoRADD by default does
not attempt to strip out such code. Many of the popular open-source Go validation
libraries can easily find false positives that strip out legitimate text that is not HTML
or JavaScript. If the developer chooses to, the developer can set up a sanitizer for all
inputs by setting the GlobalSanitizer variable in the configuration settings.

GoRADD uses the approach of attempting to ensure that all untrusted text that is 
output to a browser is encoded so that it cannot be interpreted as HTML or JavaScript by the browser.
However, this takes a certain amount of knowledge and cooperation from the developer.

HTML controls that accept input are automatically encoded when they reproduce their
output. All text controls descend from the Textbox control, and so they inherit its validation
and encoding capabilities. Span and Panel controls (div controls), by default will HTML
encode its Text value. This encoding behavior can only be turned off on a per-control basis
if the developer decides to do this.

The main place for developers to be careful of is in the GoRADD templates. Whenever content
is delivered from an untrusted source, like the database, it should be encoded. Developers
can do this by using one of the output tags with an exclamation mark in it (!). For
example, in all the default panel templates that are output by the GoRADD code generator,
you will see {{!v tags, which will output any Go value and HTML encode the result.

## Secure Coding Practices

The following is from https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/

This is a description of how GoRADD follows these secure coding practices.

### Input Validation

- [x] Conduct all data validation on a trusted system (e.g., The server)  
 - All validation is server-side in Goradd. 
 - Client-side validation is optionally available, but not relied upon.

- [x] Identify all data sources and classify them into trusted and untrusted. Validate all data from untrusted sources (e.g., Databases, file streams, etc.)
  - Trusted:
    - Internal caches
    - Session storage
  - Untrusted:
    - Http input, including headers, hidden form values, page state, and cookies
      - Layered validation is provided by the page and control packages in the framework.
    - File uploads
      - Validation must be provided by the developer when using the file upload control.
    - Database data
      - Routines are provided to allow the developer to sanitize and/or encode output from database values as needed. 

- [x] There should be a centralized input validation routine for the application
  - Standard validation routines are funnelled through the functions in the validate.go file
  in the strings package. Type specific validations are handled by individual controls.

- [x] Specify proper character sets, such as UTF-8, for all sources of input
  - UTF-8 is used internally throughout

- [x] Encode data to a common character set before validating (Canonicalize)
  - UTF-8 is used internally throughout

- [x] All validation failures should result in input rejection
  - Internal framework values coming from the client are validated and on error, the entire input is rejected. See Context.fillApp().
  - Controls are validated server side, and errors are reported to the client, while also rejecting input. See ControlBase.Validate().

- [x] Determine if the system supports UTF-8 extended character sets and if so, validate after UTF-8 decoding is completed
  - Go supports UTF-8 extended and handles decoding in the standard library.

- [x] Validate all client provided data before processing, including all parameters, URLs and HTTP header content (e.g. Cookie names and values). Be sure to include automated post backs from JavaScript, Flash or other embedded code
 - Postback data is checked in Context.fillApp().
 - All REST API data should be checked by the developer.

- [x] Verify that header values in both requests and responses contain only ASCII characters
  - Incoming header values are checked in Application.validateHttpHandler().
  - Outgoing header values are checked when using buffered output. Due to how the Go http handlers work, it's up to the developer to check the outgoing headers of REST API calls.

- [x] Validate data from redirects (An attacker may submit malicious content directly to the target of the redirect, thus circumventing application logic and any validation performed before the redirect)  
  - GoRADD does not redirect in a way that would circumvent validation.

- [x] Validate for expected data types
  - Individual controls provide appropriate validation for each data type

- [x] Validate data range
  - Basic ranges are checked by default controls. Developers can write their own validations to further enhance range checking.

- [x]  Validate data length
  - Many controls have a data length check for incoming data. The codegen process by default will set this length to the length of the corresponding database field.

- [x] Validate all input against a "white" list of allowed characters, whenever possible
  - This can be done by the developer for specific situations in the Validate() function.
  - As a framework, GoRADD cannot implement whitelists.

- [x] If any potentially hazardous characters must be allowed as input, be sure that you implement additional controls like output encoding, secure task specific APIs and accounting for the utilization of that data throughout the application . Examples of common hazardous characters include:
< > " ' % ( ) & + \ \' \"
  - As described above, a variety of mechanisms are used to prevent malicious input from impacting the product.

- [x] If your standard validation routine cannot address the following inputs, then they should be checked discretely
  - Check for null bytes (%00)
    - The default sanitizer rejects strings that try to embed a null.
  - Check for new line characters (%0d, %0a, \r, \n)
    - The textbox sanitizer filters newlines for single-line textboxes.
  - Check for “dot-dot-slash" (../ or ..\) path alterations characters. In cases where UTF-8 extended character set encoding is supported, address alternate representation like: %c0%ae%c0%ae/
  (Utilize canonicalization to address double encoding or other forms of obfuscation attacks)
    - This should be handled on a case-by-case basis using filepath.Clean()

### Output Encoding
- [x] Conduct all encoding on a trusted system (e.g., The server)
- [x] Utilize a standard, tested routine for each type of outbound encoding
  - HTML encoding is handled by Go's html.EscapeString() function.
- [x] Contextually output encode all data returned to the client that originated outside the application's trust boundary. HTML entity encoding is one example, but does not work in all cases
  - HTML encoding is handled in the GoRADD templates as explained above.
  - Developers must handle encoding in other contexts.
- [x] Encode all characters unless they are known to be safe for the intended interpreter
  - HTML is encoded before output to the browser.
- [x] Contextually sanitize all output of un-trusted data to queries for SQL, XML, and LDAP
  - As mentioned above, GoRADD uses the standard capabilities of the SQL drivers to encode input.
- [x] Sanitize all output of un-trusted data to operating system commands
  - The GoRADD framework does not directly call operating system commands.

### Authentication and Password Management

Authentication and password management is a constantly changing field. GoRADD does
not have a built-in password management system. Password management is left to the developer.

### Session Management
- [x] Use the server or framework’s session management controls. The application should only recognize these session identifiers as valid
  - GoRADD uses alexedwards/scs by default, a popular sesssion management system
  - The session management system may be configured by removing the comments in goradd-project/web/app/app.go around the Application.SetupSessionManager() function.

- [x] Session identifier creation must always be done on a trusted system (e.g., The server)
  - GoRADD's session identifier is created by the session manager on the server

- [x] Session management controls should use well vetted algorithms that ensure sufficiently random session identifiers
  - The alexedwards/scs uses Go's rand.Read() function, a cryptographically secure random number generator.

- [x] Set the domain and path for cookies containing authenticated session identifiers to an appropriately restricted value for the site
  - By default, a cookie is created for the entire site.
  - The developer can change this behavior through calls to scs

- [ ] Logout functionality should fully terminate the associated session or connection
  - The developer should call session.Reset() on logout.

- [ ] Logout functionality should be available from all pages protected by authorization
  - This is the responsibility of the developer.

- [x] Establish a session inactivity timeout that is as short as possible, based on balancing risk and business functional requirements. In most cases it should be no more than several hours.
  - The default configuration has an idle timeout of 6 hours.
  - Custom timeouts can be set in the Application.SetupSessionManager() function mentioned above.

- [x] Disallow persistent logins and enforce periodic session terminations, even when the session is active. Especially for applications supporting rich network connections or connecting to critical systems. Termination times should support business requirements and the user should receive sufficient notification to mitigate negative impacts
  - The default configuration has a lifetime timeout of 24 hours.
  - Custom timeouts can be set in the Application.SetupSessionManager() function mentioned above.

- [ ] If a session was established before login, close that session and establish a new session after a successful login
  - The developer should call session.Reset() after a login.

- [ ] Generate a new session identifier on any re-authentication
  - The developer should call session.Reset() after a login.

- [ ] Do not allow concurrent logins with the same user ID
  - This is up to the developer

- [x] Do not expose session identifiers in URLs, error messages or logs. Session identifiers should only be located in the HTTP cookie header. For example, do not pass session identifiers as GET parameters
  - GoRADD does not reference the session identifier. 

- [x] Protect server side session data from unauthorized access, by other users of the server, by implementing appropriate access controls on the server
  - By default, session data is stored in memory.
  - The scs session manager allows a variety of alternate storage mechanisms to be used. It is up to the developer to secure whatever method is used.

- [x] Generate a new session identifier and deactivate the old one periodically. (This can mitigate certain session hijacking scenarios where the original identifier was compromised)
  - Sessions have a settable lifetime.

- [x] Generate a new session identifier if the connection security changes from HTTP to HTTPS, as can occur during authentication. Within an application, it is recommended to consistently utilize HTTPS rather than switching between HTTP to HTTPS.
  - The developer should call session.Reset() when the connection changes

- [x] Supplement standard session management for sensitive server-side operations, like account management, by utilizing per-session strong random tokens or parameters. This method can be used to prevent Cross Site Request Forgery attacks
- [x] Supplement standard session management for highly sensitive or critical operations by utilizing per-request, as opposed to per-session, strong random tokens or parameters
  - CSRF attacks are handled with the CSRF token checks
  
- [ ] Set the "secure" attribute for cookies transmitted over a TLS connection
  - The developer should set SessionManager.Cookie.Secure to true in Application.SetupSessionManager() if the application is being served with SSL.

- [x] Set cookies with the HttpOnly attribute, unless you specifically require client-side scripts within your application to read or set a cookie's value
  - By default, the HttpOnly attribute is set to true. 
  - Developers can control this attribute through the application.SetupSessionManager() function.

### Access Control
Access Control is under the control of the developer. 

### Cryptographic Practices
Cryptographic services are provided by Go's rand package.

### Error Handling and Logging

Validation errors are handled by each control by setting the control's Message value.
These messages are viewable on the browser only when a control is wrapped in a FormFieldWrapper
or derived control.

Serious application errors are logged and a generic error message is displayed to the user.
This error message can be customize in the config.setupErrorMessage() function.

See GoRADD's log package for details on the application logger. The logger can be customized
in the config.initLogs() function.

TBD:
- [ ] Do not disclose sensitive information in error responses, including system details, session identifiers or account information
- [ ] Use error handlers that do not display debugging or stack trace information
- [ ] Implement generic error messages and use custom error pages
- [ ] The application should handle application errors and not rely on the server configuration
- [ ] Properly free allocated memory when error conditions occur
- [ ] Error handling logic associated with security controls should deny access by default
- [ ] All logging controls should be implemented on a trusted system (e.g., The server)
- [ ] Logging controls should support both success and failure of specified security events
- [ ] Ensure logs contain important log event data
- [ ] Ensure log entries that include un-trusted data will not execute as code in the intended log viewing interface or software
- [ ] Restrict access to logs to only authorized individuals
- [ ] Utilize a master routine for all logging operations
- [ ] Do not store sensitive information in logs, including unnecessary system details, session identifiers or passwords
- [ ] Ensure that a mechanism exists to conduct log analysis
- [ ] Log all input validation failures
- [ ] Log all authentication attempts, especially failures
- [ ] Log all access control failures
- [ ] Log all apparent tampering events, including unexpected changes to state data
- [ ] Log attempts to connect with invalid or expired session tokens
- [ ] Log all system exceptions
- [ ] Log all administrative functions, including changes to the security configuration settings
- [ ] Log all backend TLS connection failures
- [ ] Log cryptographic module failures
- [ ] Use a cryptographic hash function to validate log entry integrity

### Data Protection
- [ ] Implement least privilege, restrict users to only the functionality, data and system information that is required to perform their tasks
  - To be managed by the developer

- [x] Protect all cached or temporary copies of sensitive data stored on the server from unauthorized access and purge those temporary working files a soon as they are no longer required.
  - All caches are kept in memory by default. 
  - The developer has the option of moving these to sotrage-based systems, but it is up to the developer to manage and secure that data.

- [x] Encrypt highly sensitive stored information, like authentication verification data, even on the server side. Always use well vetted algorithms, see "Cryptographic Practices" for additional guidance
  - All such data is kept by default in memory and is not stored.
  - The developer has the option of moving these to drive-based systems, but it is up to the developer to manage and secure that data.

- [x] Protect server-side source-code from being downloaded by a user
  - Since Go is a compiled language, there is no server-side source code.

- [x] Do not store passwords, connection strings or other sensitive information in clear text or in any non-cryptographically secure manner on the client side. This includes embedding in insecure formats like: MS viewstate, Adobe flash or compiled code
  - GoRADD does not utilize client-side storage. See the javascripts in web/assets/js.

- [x] Remove comments in user accessible production code that may reveal backend system or other sensitive information
  - In production, all comments in javascript and css are removed by the minimization process.
  - See the scripts in goradd-project/build/app which utilize the tdewolff/minify library to minify javascript, and sass to minify css.

- [x] Remove unnecessary application and system documentation as this can reveal useful information to attackers
  - GoRADD is compiled, and is the applicaiton binary is not delivered with documentation.

- [ ] Do not include sensitive information in HTTP GET request parameters
  - This is the responsibility of the developer

- [ ] Disable auto complete features on forms expected to contain sensitive information, including authentication
  - This can be accomplished by the developer by calling SetAttribute("autocomplete","off") on the control.

- [ ] Disable client side caching on pages containing sensitive information. Cache-Control: no-store, may be used in conjunction with the HTTP header control "Pragma: no-cache", which is less effective, but is HTTP/1.0 backward compatible
  - The developer can modify the response header for GoRADD pages in the Exit() function on the form control.
  - The developer has direct access to the response writer for header control when responding to API calls.

- [ ] The application should support the removal of sensitive data when that data is no longer required. (e.g. personal information or certain financial data)
  - This is the responsibility of the developer

- [x] Implement appropriate access controls for sensitive data stored on the server. This includes cached data, temporary files and data that should be accessible only by specific system users
  - All such data is kept by default in memory and is not stored.
  - The developer has the option of moving these to drive-based systems, but it is up to the developer to manage and secure that data.

### Communication Security
GoRADD does not implement any backdoor communication with other servers. If the developer
chooses to include this in the application, it is up to the developer to secure it.

### System Configuration
- [ ] Ensure servers, frameworks and system components are running the latest approved version
- [ ] Ensure servers, frameworks and system components have all patches issued for the version in use
  - These are the responsibility of the developer

- [x] Turn off directory listings
  - GoRADD does not perform directory listings by default.

- [ ] Restrict the web server, process and service accounts to the least privileges possible
  - This is the responsibility of the developer

- [x] When exceptions occur, fail securely
  - GoRADD has a default panic handler that will display a generic error message to the user that contains no sensitive information.
  - This error message can be customized in the config.setupErrorMessage() function.

- [ ] Remove all unnecessary functionality and files
  - This is the responsibility of the developer in the build process.

- [x] Remove test code or any functionality not intended for production, prior to deployment
  - The GoRADD build system will set the config.Release variable to true. Many parts of the framework respond to this by removing debug functionality from the application.

- [ ] Prevent disclosure of your directory structure in the robots.txt file by placing directories not intended for public indexing into an isolated parent directory. Then "Disallow" that entire parent directory in the robots.txt file rather than Disallowing each individual directory
  - The developer can create a robots.txt file and place it in the goradd-project/web/root directory to deploy it. GoRADD does not provide a default file.

- [x] Define which HTTP methods, Get or Post, the application will support and whether it will be handled differently in different pages in the application
  - This is the responsibility of the developer for REST APIs.
  - Standard GoRADD pages only support GET by default, and they only support POST requests that come from the same page.

- [x] Disable unnecessary HTTP methods, such as WebDAV extensions. If an extended HTTP method that supports file handling is required, utilize a well-vetted authentication mechanism
  - GoRADD does not respond HTTP methods that are not defined by the developer.

- [x] If the web server handles both HTTP 1.0 and 1.1, ensure that both are configured in a similar manor or insure that you understand any difference that may exist (e.g. handling of extended HTTP methods)
  - The Go server generally responds only to HTTP 1.1 and 2.0 requests.

- [x] Remove unnecessary information from HTTP response headers related to the OS, web-server version and application frameworks
  - GoRADD does not announce itself to the browser.

- [ ] The security configuration store for the application should be able to be output in human readable form to support auditing
  - GoRADD does not have a central security configuration store.

- [ ] Implement an asset management system and register system components and software in it
  - This is the responsibility of the developer for system components
  - GoRADD uses Go's standard go.mod file to register 3rd party OSS components that it uses.

- [ ] Isolate development environments from the production network and provide access only to authorized development and test groups. Development environments are often configured less securely than production environments and attackers may use this difference to discover shared weaknesses or as an avenue for exploitation
  - This is the responsibility of the developer

- [ ] Implement a software change control system to manage and record changes to the code both in development and production
  - This is the responsibility of the developer for application code. The goradd-project directory should be
    checked in to source control.
  - GoRADD uses GitHub for change control for the framework itself.


### Database Security
- [x] Use strongly typed parameterized queries
  - GoRADD uses the paramaterized query capabilities of its database drivers.

- [x] Utilize input validation and output encoding and be sure to address meta characters. If these fail, do not run the database command
  - These are covered above in the Input Validation and Output Encoding sections

- [x] Ensure that variables are strongly typed
  - Go is a strongly typed language. All database variables use type equivalents when 
    transferring them from Go to the database.

- [ ] The application should use the lowest possible level of privilege when accessing the database
  - This is the responsibility of the developer to setup the database user with appropriate privileges.
  - See the database configuration process in the goradd-project/config/db.go file.

- [ ] Use secure credentials for database access
  - This is the responsibility of the developer.
  - Use the ability to send the database credentials into the program through the db.cfg file.

- [ ] Connection strings should not be hard coded within the application. Connection strings should be stored in a separate configuration file on a trusted system and they should be encrypted.

- [ ] Use stored procedures to abstract data access and allow for the removal of permissions to the base tables in the database
  - This is the responsibility of the developer, but GoRADD is designed to directly access databases with SQL.

- [ ] Close the connection as soon as possible
  - Go's database drivers use a database pool for efficiently reusing connections

- [ ] Remove or change all default database administrative passwords. Utilize strong passwords/phrases or implement multi-factor authentication
  - This is the responsibility of the developer

- [ ] Turn off all unnecessary database functionality (e.g., unnecessary stored procedures or services, utility packages, install only the minimum set of features and options required (surface area reduction))
  - This is the responsibility of the developer

- [ ] Remove unnecessary default vendor content (e.g., sample schemas)
  - This is the responsibility of the developer

- [ ] Disable any default accounts that are not required to support business requirements
  - This is the responsibility of the developer

- [ ] The application should connect to the database with different credentials for every trust distinction (e.g., user, read-only user, guest, administrators)
  - This is the responsibility of the developer


### File Management
GoRADD provides a file upload control in the pkg/page/control/button/file_select.go file.
It is up to the developer to care for the file upload process. See that file for more details.


Memory Management:
- [] Utilize input and output control for un-trusted data

- [ ] Double check that the buffer is as large as specified
  - Go buffers are flexible and generally do not have a fixed size.

- [x] When using functions that accept a number of bytes to copy, such as strncpy(), be aware that if the destination buffer size is equal to the source buffer size, it may not NULL-terminate the string
  - Go strings are not internally null terminated.

- [x] Check buffer boundaries if calling the function in a loop and make sure there is no danger of writing past the allocated space
- [ ] Truncate all input strings to a reasonable length before passing them to the copy and concatenation functions
  - As a framework, it is not practical for GoRADD to predict the lengths.
  
- [x] Specifically close resources, don’t rely on garbage collection. (e.g., connection objects, file handles, etc.)
  - This is the responsibility of the devloper
  
- [x] Use non-executable stacks when available
  - This is the responsibility of the devloper
  
- [ ] Avoid the use of known vulnerable functions (e.g., printf, strcat, strcpy etc.)
  - Go's memory model makes these function not vulnerable.

- [x] Properly free allocated memory upon the completion of functions and at all exit points
  - Go is a garbage collected language and will free allocated memory when it sees fit.


