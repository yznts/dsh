package ddb

// ResolveScheme resolves the scheme to a standard name.
func ResolveScheme(scheme string) string {
	switch scheme {
	case "sqlite", "sqlite3":
		return "sqlite"
	case "postgres", "postgresql":
		return "postgres"
	default:
		return ""
	}
}
