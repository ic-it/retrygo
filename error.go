package retrygo

// ErrRecovered is returned when a panic is recovered.
type ErrRecovered struct{ V any }

func (e ErrRecovered) Error() string {
	return "recovered"
}
