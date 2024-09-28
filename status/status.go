package status

type ItemStatus int

const (
	TODO ItemStatus = iota
	IN_PROGRESS
	DONE
)

func (itemStatus ItemStatus) String() string {
	switch itemStatus {
	case TODO:
		return "TODO"
	case IN_PROGRESS:
		return "IN_PROGRESS"
	case DONE:
		return "DONE"
	default:
		return "UNKNOWN"
	}
}
