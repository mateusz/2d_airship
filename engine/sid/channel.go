package sid

type Channel struct {
	src    SignalSource
	volume float64
	paused bool
}

func NewChannel(vol float64) *Channel {
	return &Channel{
		volume: vol,
	}
}
