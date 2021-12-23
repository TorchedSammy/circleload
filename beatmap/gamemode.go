package beatmap

type Mode int
const (
	ModeStandard Mode = iota
	ModeTaiko
	ModeCatch
	ModeMania
	ModeAny Mode = 727 // :)
)

func (g Mode) String() string {
	switch g {
	case ModeStandard:
		return "Standard"
	case ModeTaiko:
		return "Taiko"
	case ModeCatch:
		return "Catch"
	case ModeMania:
		return "Mania"
	case ModeAny:
		return "Any"
	default:
		return ""
	}
}
