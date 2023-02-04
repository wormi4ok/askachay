package youtube

import "testing"

func Test_isSongName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			"Valid input",
			"Gorillaz - Clint Eastwood",
			true,
		},
		{
			"Invalid input",
			"Hello",
			false,
		},
		{
			"Missing separator",
			"Queen Bohemian Rhapsody",
			false,
		},
		{
			"Special chars",
			"Guns N' Roses - Knockin' On Heaven's Door",
			true,
		},
		{
			"Cyrillic input",
			"Сплин - Когда пройдёт 100 лет",
			true,
		},
		{
			"Complex case",
			"LinkinPark - Figure.09",
			true,
		},
		{
			"Missing spaces around separator",
			"Linkin Park-Numb",
			false,
		},
		{
			"Dashes in song title",
			"Rammstein - Links 2-3-4",
			true,
		},
		{
			"Multiline string",
			`Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget - dolor.
Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSongName(tt.input); got != tt.want {
				t.Errorf("isSongName() = %v, want %v", got, tt.want)
			}
		})
	}
}
