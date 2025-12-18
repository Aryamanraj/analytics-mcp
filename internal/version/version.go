package version

// Build-time variables. Override via -ldflags.
var (
	Version   = "dev"
	Commit    = "dev"
	BuildDate = "dev"
)

// Info describes build/version metadata.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
}

// Get returns version info, defaulting empty fields to "dev".
func Get() Info {
	return Info{
		Version:   defaultOr(Version, "dev"),
		Commit:    defaultOr(Commit, "dev"),
		BuildDate: defaultOr(BuildDate, "dev"),
	}
}

func defaultOr(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
