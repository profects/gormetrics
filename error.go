package gormetrics

const ErrDbIsNil gormetricsErr = "db is nil"

// A simple type for constant errors. Makes errors easy to match.
type gormetricsErr string

func (e gormetricsErr) Error() string {
	return string(e)
}
