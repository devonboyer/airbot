package shows

import "testing"
import "github.com/stretchr/testify/require"

func TestProcessList(t *testing.T) {
	shows := &ShowList{
		Records: []Show{
			Show{Fields: Fields{Name: "Game of Thrones"}},
			Show{Fields: Fields{Name: "Silicon Valley"}},
			Show{Fields: Fields{Name: "The Walking Dead"}},
			Show{Fields: Fields{Name: "Rick and Morty"}},
		},
	}
	bot := Bot{}
	expected := "Shows on tonight:\nGame of Thrones\nSilicon Valley\nThe Walking Dead\nRick and Morty\n"
	require.Equal(t, expected, bot.processList(shows))
}
