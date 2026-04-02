package mpv

type Mpv struct {
	isPlaying            bool
	additionalParameters []string
	watchLaterDir        string
}

func New() *Mpv {
	return &Mpv{}
}
