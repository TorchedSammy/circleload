package beatmap

type Mode int
const (
	ModeStandard Mode = iota
	ModeTaiko
	ModeCatch
	ModeMania
	ModeAny Mode = 727 // :)
)

func (m Mode) String() string {
	switch m {
	case ModeStandard:
		return "standard"
	case ModeTaiko:
		return "taiko"
	case ModeCatch:
		return "catch"
	case ModeMania:
		return "mania"
	case ModeAny:
		return "any"
	default:
		return ""
	}
}
