package gopb

// InitialismsMap returns a map[string]bool of o.Initialisms.
func (o *LintOptions) InitialismsMap() map[string]bool {
	fi := o.Initialisms
	initialisms := make(map[string]bool, len(fi))
	for _, i := range fi {
		initialisms[i] = true
	}
	return initialisms
}
