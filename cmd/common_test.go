package cmd

// resetFlags reset all flags back to the default values
func resetFlags() {
	rootFlags.json = false
	rootFlags.quiet = false
}
