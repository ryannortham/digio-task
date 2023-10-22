package log

type LogEntry struct {
	IP         string
	Identity   string
	UserID     string
	Time       string
	Method     string
	Path       string
	Protocol   string
	StatusCode int
	Size       int
	Referrer   string
	UserAgent  string
}
