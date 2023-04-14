package receiver

type Receiver interface {
	Name() string
	Receive(event Event) error
}

const (
	FormatSummary = "summary"
	FormatJSON    = "json"
	FormatYaml    = "yaml"
)

type Options struct {
	Format      string // short, full
	LogFileName string
}
