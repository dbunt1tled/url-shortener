package enum

type Status int64

const (
	Inactive Status = iota
	Active
)

func (s Status) String() string {
	switch s {
	case Inactive:
		return "inactive"
	case Active:
		return "active"
	default:
		return "unknown"
	}
}
