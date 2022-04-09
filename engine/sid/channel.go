package sid

const (
	SID_FADE_IN  = iota
	SID_FADE_OUT = iota
)

type Channel struct {
	src         SignalSource
	volume      float64
	paused      bool
	fadeSamples int
	// 0 = faded completely, fadeSamples = unfaded
	fadeCurrent   int
	fadeDirection int
}

func NewChannel(vol float64) *Channel {
	return &Channel{
		volume:        vol,
		fadeSamples:   44100.0 * 0.100,
		fadeCurrent:   0,
		fadeDirection: SID_FADE_IN,
	}
}
