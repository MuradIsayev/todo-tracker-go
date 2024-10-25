package status

// Describes the status of an item
type ItemStatus int

const (
	TODO ItemStatus = iota
	IN_PROGRESS
	DONE
)

// Returns the string representation of the ItemStatus
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
