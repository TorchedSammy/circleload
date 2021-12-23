package beatmap

type Gamemode int
const (
	Standard Gamemode = iota
	Taiko
	Catch
	Mania
	AnyMode Gamemode = 727 // :)
)

func (g Gamemode) String() string {
	switch g {
	case Standard:
		return "Standard"
	case Taiko:
		return "Taiko"
	case Catch:
		return "Catch"
	case Mania:
		return "Mania"
	case AnyMode:
		return "Any"
	default:
		return ""
	}
}
