package errormessages

const (
	// EmptyZoneID No Zone ID specified
	EmptyZoneID = "Invalid config: Zone id must not be empty"
	// InvalidZoneID Zone ID is invalid
	InvalidZoneID = "Invalid config: Zone id is invalid"
	// EmptyCredentials No Credentials specified
	EmptyCredentials = "Invalid credentials: Key & email must not be empty"
	// InvalidEmail API Email is invalid
	InvalidEmail = "Invalid credentials: Email is invalid"
	// InvalidAPIKey API key is invalid
	InvalidAPIKey = "Invalid credentials: Key is invalid"
	// EmptyScript No Script details specified
	EmptyScript = "Invalid config: Script name & script file must not be empty"
	// InvalidCredentials Cannot login
	InvalidCredentials = "Invalid credentials: Key and/or email are invalid"
)
