# Security
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
  - The framework implements a layered approach in the page and control packages for validating http input, which can be additionally customized by the developer.
  - Http inputs specific to the framework are validated by the framework.
  - See the "Validate()" function provided by the ControlBase struct.


- [x] Specify proper character sets, such as UTF-8, for all sources of input
  - UTF-8 is used internally throughout

- [x] Encode data to a common character set before validating (Canonicalize)
  - UTF-8 is used internally throughout

- [x] All validation failures should result in input rejection
  - Internal framework values coming from the client are validated and on error, the entire input is rejected. See Context.fillApp().
  - Controls are validated server side, and errors are reported to the client, while also rejecting input. See ControlBase.Validate().

- [x] Determine if the system supports UTF-8 extended character sets and if so, validate after UTF-8 decoding is completed
  - Go supports UTF-8 extended

- [x] Validate all client provided data before processing, including all parameters, URLs and HTTP header content (e.g. Cookie names and values). Be sure to include automated post backs from JavaScript, Flash or other embedded code
 - Postback data is checked in Context.fillApp.
 - All REST API data should be checked by the developer.

- [x] Verify that header values in both requests and responses contain only ASCII characters
  - Incoming header values are checked in Application.validateHttpHandler().
  - Outgoing header values are checked when using buffered output. Due to how the Go http handlers work, it's up to the developer to check the outgoing headers of REST API calls.


- [x] Validate data from redirects (An attacker may submit malicious content directly to the target of the redirect, thus circumventing application logic and any validation performed before the redirect)  
  - Goradd does not redirect in a way that would circumvent validation.

- [x] Validate for expected data types
  - Individual controls provide appropriate validation for each data type


- [x] Validate data range
  - Basic ranges are checked by default controls. Developers can write their own validations to further enhance range checking.

- [x]  Validate data length
  - Many controls have a data length check for incoming data. The codegen process by default will set this length to the length of the corresponding database field.

- [x] Validate all input against a "white" list of allowed characters, whenever possible
  - This can be done by the developer for specific situation in the Validate() function.


- [x] If any potentially hazardous characters must be allowed as input, be sure that you implement additional controls like output encoding, secure task specific APIs and accounting for the utilization of that data throughout the application . Examples of common hazardous characters include:
< > " ' % ( ) & + \ \' \"
  - BlueMonday is used to validate input text by default, and handles special characters.

	[16] If your standard validation routine cannot address the following inputs, then they should be checked discretely
o	Check for null bytes (%00)
o	Check for new line characters (%0d, %0a, \r, \n)
o	Check for “dot-dot-slash" (../ or ..\) path alterations characters. In cases where UTF-8 extended character set encoding is supported, address alternate representation like: %c0%ae%c0%ae/
(Utilize canonicalization to address double encoding or other forms of obfuscation attacks)


Output Encoding:
	[17] Conduct all encoding on a trusted system (e.g., The server)
	[18] Utilize a standard, tested routine for each type of outbound encoding
	[19] Contextually output encode all data returned to the client that originated outside the application's trust boundary. HTML entity encoding is one example, but does not work in all cases
	[20] Encode all characters unless they are known to be safe for the intended interpreter
	[21] Contextually sanitize all output of un-trusted data to queries for SQL, XML, and LDAP
	[22] Sanitize all output of un-trusted data to operating system commands



Authentication and Password Management:
	[23] Require authentication for all pages and resources, except those specifically intended to be public
	[24] All authentication controls must be enforced on a trusted system (e.g., The server)  
	[25] Establish and utilize standard, tested, authentication services whenever possible
	[26] Use a centralized implementation for all authentication controls, including libraries that call external authentication services
	[27] Segregate authentication logic from the resource being requested and use redirection to and from the centralized authentication control
	[28] All authentication controls should fail securely
	[29] All administrative and account management functions must be at least as secure as the primary authentication mechanism
	[30] If your application manages a credential store, it should ensure that only cryptographically strong one-way salted hashes of passwords are stored and that the table/file that stores the passwords and keys is write-able only by the application. (Do not use the MD5 algorithm if it can be avoided)
	[31] Password hashing must be implemented on a trusted system (e.g., The server).
	[32] Validate the authentication data only on completion of all data input, especially for sequential authentication implementations
	[33] Authentication failure responses should not indicate which part of the authentication data was incorrect. For example, instead of "Invalid username" or "Invalid password", just use "Invalid username and/or password" for both. Error responses must be truly identical in both display and source code
	[34] Utilize authentication for connections to external systems that involve sensitive information or functions
	[35] Authentication credentials for accessing services external to the application should be encrypted and stored in a protected location on a trusted system (e.g., The server). The source code is NOT a secure location
	[36] Use only HTTP POST requests to transmit authentication credentials
	[37] Only send non-temporary passwords over an encrypted connection or as encrypted data, such as in an encrypted email. Temporary passwords associated with email resets may be an exception
	[38] Enforce password complexity requirements established by policy or regulation. Authentication credentials should be sufficient to withstand attacks that are typical of the threats in the deployed environment. (e.g., requiring the use of alphabetic as well as numeric and/or special characters)
	[39] Enforce password length requirements established by policy or regulation. Eight characters is commonly used, but 16 is better or consider the use of multi-word pass phrases
	[40] Password entry should be obscured on the user's screen. (e.g., on web forms use the input type "password")
	[41] Enforce account disabling after an established number of invalid login attempts (e.g., five attempts is common).  The account must be disabled for a period of time sufficient to discourage brute force guessing of credentials, but not so long as to allow for a denial-of-service attack to be performed
	[42] Password reset and changing operations require the same level of controls as account creation and authentication.
	[43] Password reset questions should support sufficiently random answers. (e.g., "favorite book" is a bad question because “The Bible” is a very common answer)
	[44] If using email based resets, only send email to a pre-registered address with a temporary link/password
	[45] Temporary passwords and links should have a short expiration time
	[46] Enforce the changing of temporary passwords on the next use
	[47] Notify users when a password reset occurs
	[48] Prevent password re-use
	[49] Passwords should be at least one day old before they can be changed, to prevent attacks on password re-use
	[50] Enforce password changes based on requirements established in policy or regulation. Critical systems may require more frequent changes. The time between resets must be administratively controlled
	[51] Disable "remember me" functionality for password fields
	[52] The last use (successful or unsuccessful) of a user account should be reported to the user at their next successful login
	[53] Implement monitoring to identify attacks against multiple user accounts, utilizing the same password. This attack pattern is used to bypass standard lockouts, when user IDs can be harvested or guessed
	[54] Change all vendor-supplied default passwords and user IDs or disable the associated accounts
	[55] Re-authenticate users prior to performing critical operations
	[56] Use Multi-Factor Authentication for highly sensitive or high value transactional accounts
	[57] If using third party code for authentication, inspect the code carefully to ensure it is not affected by any malicious code


Session Management:
	[58] Use the server or framework’s session management controls. The application should only recognize these session identifiers as valid
	[59] Session identifier creation must always be done on a trusted system (e.g., The server)
	[60] Session management controls should use well vetted algorithms that ensure sufficiently random session identifiers
	[61] Set the domain and path for cookies containing authenticated session identifiers to an appropriately restricted value for the site
	[62] Logout functionality should fully terminate the associated session or connection
	[63] Logout functionality should be available from all pages protected by authorization
	[64] Establish a session inactivity timeout that is as short as possible, based on balancing risk and business functional requirements. In most cases it should be no more than several hours
	[65] Disallow persistent logins and enforce periodic session terminations, even when the session is active. Especially for applications supporting rich network connections or connecting to critical systems. Termination times should support business requirements and the user should receive sufficient notification to mitigate negative impacts
	[66] If a session was established before login, close that session and establish a new session after a successful login
	[67] Generate a new session identifier on any re-authentication
	[68] Do not allow concurrent logins with the same user ID
	[69] Do not expose session identifiers in URLs, error messages or logs. Session identifiers should only be located in the HTTP cookie header. For example, do not pass session identifiers as GET parameters
	[70] Protect server side session data from unauthorized access, by other users of the server, by implementing appropriate access controls on the server
	[71] Generate a new session identifier and deactivate the old one periodically. (This can mitigate certain session hijacking scenarios where the original identifier was compromised)
	[72] Generate a new session identifier if the connection security changes from HTTP to HTTPS, as can occur during authentication. Within an application, it is recommended to consistently utilize HTTPS rather than switching between HTTP to HTTPS.
	[73] Supplement standard session management for sensitive server-side operations, like account management, by utilizing per-session strong random tokens or parameters. This method can be used to prevent Cross Site Request Forgery attacks
	[74] Supplement standard session management for highly sensitive or critical operations by utilizing per-request, as opposed to per-session, strong random tokens or parameters
	[75] Set the "secure" attribute for cookies transmitted over an TLS connection
	[76] Set cookies with the HttpOnly attribute, unless you specifically require client-side scripts within your application to read or set a cookie's value


Access Control:
	[77] Use only trusted system objects, e.g. server side session objects, for making access authorization decisions
	[78] Use a single site-wide component to check access authorization. This includes libraries that call external authorization services
	[79] Access controls should fail securely
	[80] Deny all access if the application cannot access its security configuration information
	[81] Enforce authorization controls on every request, including those made by server side scripts, "includes" and requests from rich client-side technologies like AJAX and Flash
	[82] Segregate privileged logic from other application code
	[83] Restrict access to files or other resources, including those outside the application's direct control, to only authorized users
	[84] Restrict access to protected URLs to only authorized users
	[85] Restrict access to protected functions to only authorized users
	[86] Restrict direct object references to only authorized users
	[87] Restrict access  to services to only authorized users
	[88] Restrict access  to application data to only authorized users
	[89] Restrict access to user and data attributes and policy information used by access controls
	[90] Restrict access security-relevant configuration information to only authorized users
	[91] Server side implementation and presentation layer representations of access control rules must match
	[92] If state data must be stored on the client, use encryption and integrity checking on the server side to catch state tampering.
	[93] Enforce application logic flows to comply with business rules
	[94] Limit the number of transactions a single user or device can perform in a given period of time. The transactions/time should be above the actual business requirement, but low enough to deter automated attacks
	[95] Use the "referer" header as a supplemental check only, it should never be the sole authorization check, as it is can be spoofed
	[96] If long authenticated sessions are allowed, periodically re-validate a user’s authorization to ensure that their privileges have not changed and if they have, log the user out and force them to re-authenticate
	[97] Implement account auditing and enforce the disabling of unused accounts (e.g., After no more than 30 days from the expiration of an account’s password.)
	[98] The application must support disabling of accounts and terminating sessions when authorization ceases (e.g., Changes to role, employment status, business process, etc.)
	[99] Service accounts or accounts supporting connections to or from external systems should have the least privilege possible
	[100] Create an Access Control Policy to document an application's business rules, data types and access authorization criteria and/or processes so that access can be properly provisioned and controlled. This includes identifying access requirements for both the data and system resources


Cryptographic Practices:
	[101] All cryptographic functions used to protect secrets from the application user must be implemented on a trusted system (e.g., The server)
	[102] Protect master secrets from unauthorized access
	[103] Cryptographic modules should fail securely
	[104] All random numbers, random file names, random GUIDs, and random strings should be generated using the cryptographic module’s approved random number generator when these random values are intended to be un-guessable
	[105] Cryptographic modules used by the application should be compliant to FIPS 140-2 or an equivalent standard. (See http://csrc.nist.gov/groups/STM/cmvp/validation.html)
	[106] Establish and utilize a policy and process for how cryptographic keys will be managed


Error Handling and Logging:
	[107] Do not disclose sensitive information in error responses, including system details, session identifiers or account information
	[108] Use error handlers that do not display debugging or stack trace information
	[109] Implement generic error messages and use custom error pages
	[110] The application should handle application errors and not rely on the server configuration
	[111] Properly free allocated memory when error conditions occur
	[112] Error handling logic associated with security controls should deny access by default
	[113] All logging controls should be implemented on a trusted system (e.g., The server)
	[114] Logging controls should support both success and failure of specified security events
	[115] Ensure logs contain important log event data
	[116] Ensure log entries that include un-trusted data will not execute as code in the intended log viewing interface or software
	[117] Restrict access to logs to only authorized individuals
	[118] Utilize a master routine for all logging operations
	[119] Do not store sensitive information in logs, including unnecessary system details, session identifiers or passwords
	[120] Ensure that a mechanism exists to conduct log analysis
	[121] Log all input validation failures
	[122] Log all authentication attempts, especially failures
	[123] Log all access control failures
	[124] Log all apparent tampering events, including unexpected changes to state data
	[125] Log attempts to connect with invalid or expired session tokens
	[126] Log all system exceptions
	[127] Log all administrative functions, including changes to the security configuration settings
	[128] Log all backend TLS connection failures
	[129] Log cryptographic module failures
	[130] Use a cryptographic hash function to validate log entry integrity
Data Protection:
	[131] Implement least privilege, restrict users to only the functionality, data and system information that is required to perform their tasks
	[132] Protect all cached or temporary copies of sensitive data stored on the server from unauthorized access and purge those temporary working files a soon as they are no longer required.
	[133] Encrypt highly sensitive stored information, like authentication verification data, even on the server side. Always use well vetted algorithms, see "Cryptographic Practices" for additional guidance
	[134] Protect server-side source-code from being downloaded by a user
	[135] Do not store passwords, connection strings or other sensitive information in clear text or in any non-cryptographically secure manner on the client side. This includes embedding in insecure formats like: MS viewstate, Adobe flash or compiled code
	[136] Remove comments in user accessible production code that may reveal backend system or other sensitive information
	[137] Remove unnecessary application and system documentation as this can reveal useful information to attackers
	[138] Do not include sensitive information in HTTP GET request parameters
	[139] Disable auto complete features on forms expected to contain sensitive information, including authentication  
	[140] Disable client side caching on pages containing sensitive information. Cache-Control: no-store, may be used in conjunction with the HTTP header control "Pragma: no-cache", which is less effective, but is HTTP/1.0 backward compatible
	[141] The application should support the removal of sensitive data when that data is no longer required. (e.g. personal information or certain financial data)
	[142] Implement appropriate access controls for sensitive data stored on the server. This includes cached data, temporary files and data that should be accessible only by specific system users


Communication Security:
	[143] Implement encryption for the transmission of all sensitive information. This should include TLS for protecting the connection and may be supplemented by discrete encryption of sensitive files or non-HTTP based connections
	[144] TLS certificates should be valid and have the correct domain name, not be expired, and be installed with intermediate certificates when required
	[145] Failed TLS connections should not fall back to an insecure connection
	[146] Utilize TLS connections for all content requiring authenticated access and for all other sensitive information
	[147] Utilize TLS for connections to external systems that involve sensitive information or functions
	[148] Utilize a single standard TLS implementation that is configured appropriately
	[149] Specify character encodings for all connections
	[150] Filter parameters containing sensitive information from the HTTP referer, when linking to external sites



System Configuration:
	[151] Ensure servers, frameworks and system components are running the latest approved version
	[152] Ensure servers, frameworks and system components have all patches issued for the version in use
	[153] Turn off directory listings
	[154] Restrict the web server, process and service accounts to the least privileges possible
	[155] When exceptions occur, fail securely
	[156] Remove all unnecessary functionality and files
	[157] Remove test code or any functionality not intended for production, prior to deployment
	[158] Prevent disclosure of your directory structure in the robots.txt file by placing directories not intended for public indexing into an isolated parent directory. Then "Disallow" that entire parent directory in the robots.txt file rather than Disallowing each individual directory
	[159] Define which HTTP methods, Get or Post, the application will support and whether it will be handled differently in different pages in the application
	[160] Disable unnecessary HTTP methods, such as WebDAV extensions. If an extended HTTP method that supports file handling is required, utilize a well-vetted authentication mechanism
	[161] If the web server handles both HTTP 1.0 and 1.1, ensure that both are configured in a similar manor or insure that you understand any difference that may exist (e.g. handling of extended HTTP methods)
	[162] Remove unnecessary information from HTTP response headers related to the OS, web-server version and application frameworks
	[163] The security configuration store for the application should be able to be output in human readable form to support auditing
	[164] Implement an asset management system and register system components and software in it
	[165] Isolate development environments from the production network and provide access only to authorized development and test groups. Development environments are often configured less securely than production environments and attackers may use this difference to discover shared weaknesses or as an avenue for exploitation
	[166] Implement a software change control system to manage and record changes to the code both in development and production


Database Security:
	[167] Use strongly typed parameterized queries
	[168] Utilize input validation and output encoding and be sure to address meta characters. If these fail, do not run the database command
	[169] Ensure that variables are strongly typed
	[170] The application should use the lowest possible level of privilege when accessing the database
	[171] Use secure credentials for database access
	[172] Connection strings should not be hard coded within the application. Connection strings should be stored in a separate configuration file on a trusted system and they should be encrypted.
	[173] Use stored procedures to abstract data access and allow for the removal of permissions to the base tables in the database
	[174] Close the connection as soon as possible
	[175] Remove or change all default database administrative passwords. Utilize strong passwords/phrases or implement multi-factor authentication
	[176] Turn off all unnecessary database functionality (e.g., unnecessary stored procedures or services, utility packages, install only the minimum set of features and options required (surface area reduction))
	[177] Remove unnecessary default vendor content (e.g., sample schemas)
	[178] Disable any default accounts that are not required to support business requirements
	[179] The application should connect to the database with different credentials for every trust distinction (e.g., user, read-only user, guest, administrators)


File Management:
	[180] Do not pass user supplied data directly to any dynamic include function
	[181] Require authentication before allowing a file to be uploaded
	[182] Limit the type of files that can be uploaded to only those types that are needed for business purposes
	[183] Validate uploaded files are the expected type by checking file headers. Checking for file type by extension alone is not sufficient
	[184] Do not save files in the same web context as the application. Files should either go to the content server or in the database.
	[185] Prevent or restrict the uploading of any file that may be interpreted by the web server.
	[186] Turn off execution privileges on file upload directories
	[187] Implement safe uploading in UNIX by mounting the targeted file directory as a logical drive using the associated path or the chrooted environment
	[188] When referencing existing files, use a white list of allowed file names and types. Validate the value of the parameter being passed and if it does not match one of the expected values, either reject it or use a hard coded default file value for the content instead
	[189] Do not pass user supplied data into a dynamic redirect. If this must be allowed, then the redirect should accept only validated, relative path URLs
	[190] Do not pass directory or file paths, use index values mapped to pre-defined list of paths
	[191] Never send the absolute file path to the client
	[192] Ensure application files and resources are read-only
	[193] Scan user uploaded files for viruses and malware


Memory Management:
	[194] Utilize input and output control for un-trusted data
	[195] Double check that the buffer is as large as specified
	[196] When using functions that accept a number of bytes to copy, such as strncpy(), be aware that if the destination buffer size is equal to the source buffer size, it may not NULL-terminate the string
	[197] Check buffer boundaries if calling the function in a loop and make sure there is no danger of writing past the allocated space
	[198] Truncate all input strings to a reasonable length before passing them to the copy and concatenation functions
	[199] Specifically close resources, don’t rely on garbage collection. (e.g., connection objects, file handles, etc.)
	[200] Use non-executable stacks when available
	[201] Avoid the use of known vulnerable functions (e.g., printf, strcat, strcpy etc.)
	[202] Properly free allocated memory upon the completion of functions and at all exit points





Generally, input validation from web pages is handled by the individual control
structures that each page uses. Validation is handled in layers, with basic
UTF-8 checking handled by Go and in the framework, validation specific to a control
type handled by the validation routine of that control, 

