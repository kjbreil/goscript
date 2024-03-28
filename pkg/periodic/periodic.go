package periodic

// Periodic is the list of cron expressions to run periodically
type Periodic []string

// Periodics is a helper function to add multiple strings without needing a []string{}
func Periodics(times ...string) []string {
	rtn := make([]string, len(times))
	// TODO: Validate the cron strings
	copy(rtn, times)
	return rtn
}
