package auth

import "time"

// ApiEndpoint is the url for the authorization api. Feel free to change the default during startup.
var ApiEndpoint = "/api/auth"

// HackerDelay is the amount of time to delay a response if we detect a hacker. Feel free to change the default during startup.
var HackerDelay = 20 * time.Second

// Number of seconds to permit between login attempts. Attempts to login faster than this will be rejected.
var LoginRateLimit int64 = 3
