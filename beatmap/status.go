package beatmap

type Status int
const (
	StatusGraveyard Status = iota - 2
	StatusWIP
	StatusPending
	StatusRanked
	StatusApproved
	StatusQualified
	StatusLoved
	StatusAny Status = 727
)

func (s Status) String() string {
	switch s {
	case StatusGraveyard:
		return "graveyard"
	case StatusWIP:
		return "wip"
	case StatusPending:
		return "pending"
	case StatusRanked:
		return "ranked"
	case StatusApproved:
		return "approved"
	case StatusQualified:
		return "qualified"
	case StatusLoved:
		return "loved"
	case StatusAny:
		return "any"
	default:
		return "unknown"
	}
}

func StatusFromString(status string) Status {
	switch status {
	case "graveyard":
		return StatusGraveyard
	case "wip":
		return StatusWIP
	case "pending":
		return StatusPending
	case "ranked":
		return StatusRanked
	case "approved":
		return StatusApproved
	case "qualified":
		return StatusQualified
	case "loved":
		return StatusLoved
	case "any":
		return StatusAny
	default:
		return StatusAny
	}
}
