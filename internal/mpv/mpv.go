package mpv

type Mpv struct {
	isPlaying            bool
	additionalParameters []string
	watchLaterDir        string
}

func New(watchLaterDir string) *Mpv {
	return &Mpv{false, []string{}, watchLaterDir}
}
