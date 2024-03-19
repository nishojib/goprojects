package tordle

import "strings"

// hint describes the validity of a character in a word.
type hint byte

const (
	absentCharacter hint = iota
	wrongPosition
	correctPosition
)

// String implements the Stringer interface.
func (h hint) String() string {
	switch h {
	case absentCharacter:
		return "⬜️" // grey square
	case wrongPosition:
		return "🟡" // yellow circle
	case correctPosition:
		return "💚" // green heart
	default:
		// This should never happen.
		return "💔" // red broken heart
	}
}

// Feedback is a list of hints, one per character of the word.
type Feedback []hint

// String implements the Stringer interface for a slice of hints.
func (fb Feedback) String() string {
	sb := strings.Builder{}
	for _, h := range fb {
		sb.WriteString(h.String())
	}
	return sb.String()
}

// Equal determines equality of two feedbacks
func (fb Feedback) Equal(other Feedback) bool {
	if len(fb) != len(other) {
		return false
	}

	for index, value := range fb {
		if value != other[index] {
			return false
		}
	}

	return true
}

// GameWon returns whether a feedback indicates a player has found all characters.
func (fb Feedback) GameWon() bool {
	for _, c := range fb {
		if c != correctPosition {
			return false
		}
	}

	return true
}
