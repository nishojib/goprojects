package tordle

import "testing"

func TestFeedbackString(t *testing.T) {
	tt := map[string]struct {
		f    feedback
		want string
	}{
		"nominal feedback": {
			f: feedback{
				absentCharacter,
				wrongPosition,
				correctPosition,
				absentCharacter,
				correctPosition,
			},
			want: "⬜️🟡💚⬜️💚",
		},
		"no feedback": {
			f:    feedback{hint(4)},
			want: "💔",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := tc.f.String()

			if got != tc.want {
				t.Errorf("got: %q, want: %q", got, tc.want)
			}
		})
	}
}
