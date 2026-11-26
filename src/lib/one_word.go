package bible_lib

const (
	OneWordCharacter = "character"
	OneWordBook      = "book"
)

type OneWordClue struct {
	ID         string
	Difficulty Difficulty
	Word       string
	Kind       string
	Answer     string
	Aliases    []string
	Reference  string
}

var OneWordClues = []OneWordClue{
	{
		ID:         "ark-noah",
		Difficulty: Easy,
		Word:       "ark",
		Kind:       OneWordCharacter,
		Answer:     "Noah",
		Reference:  "Genesis 6:14",
	},
	{
		ID:         "lions-daniel",
		Difficulty: Easy,
		Word:       "lions",
		Kind:       OneWordCharacter,
		Answer:     "Daniel",
		Reference:  "Daniel 6:16",
	},
	{
		ID:         "coat-joseph",
		Difficulty: Easy,
		Word:       "coat",
		Kind:       OneWordCharacter,
		Answer:     "Joseph",
		Reference:  "Genesis 37:3",
	},
	{
		ID:         "beginning-genesis",
		Difficulty: Easy,
		Word:       "beginning",
		Kind:       OneWordBook,
		Answer:     "Genesis",
		Reference:  "Genesis 1:1",
	},
	{
		ID:         "songs-psalms",
		Difficulty: Easy,
		Word:       "songs",
		Kind:       OneWordBook,
		Answer:     "Psalms",
		Aliases:    []string{"Psalm"},
		Reference:  "Psalm 1:1",
	},
	{
		ID:         "plagues-exodus",
		Difficulty: Easy,
		Word:       "plagues",
		Kind:       OneWordBook,
		Answer:     "Exodus",
		Reference:  "Exodus 9:14",
	},
	{
		ID:         "bush-moses",
		Difficulty: Medium,
		Word:       "bush",
		Kind:       OneWordCharacter,
		Answer:     "Moses",
		Reference:  "Exodus 3:2",
	},
	{
		ID:         "ram-abraham",
		Difficulty: Medium,
		Word:       "ram",
		Kind:       OneWordCharacter,
		Answer:     "Abraham",
		Reference:  "Genesis 22:13",
	},
	{
		ID:         "wisdom-solomon",
		Difficulty: Medium,
		Word:       "wisdom",
		Kind:       OneWordCharacter,
		Answer:     "Solomon",
		Reference:  "1 Kings 3:12",
	},
	{
		ID:         "fish-jonah",
		Difficulty: Medium,
		Word:       "fish",
		Kind:       OneWordCharacter,
		Answer:     "Jonah",
		Reference:  "Jonah 1:17",
	},
	{
		ID:         "suffering-job",
		Difficulty: Medium,
		Word:       "suffering",
		Kind:       OneWordBook,
		Answer:     "Job",
		Reference:  "Job 1:1",
	},
	{
		ID:         "sayings-proverbs",
		Difficulty: Medium,
		Word:       "sayings",
		Kind:       OneWordBook,
		Answer:     "Proverbs",
		Reference:  "Proverbs 1:1",
	},
	{
		ID:         "chariot-elijah",
		Difficulty: Hard,
		Word:       "chariot",
		Kind:       OneWordCharacter,
		Answer:     "Elijah",
		Reference:  "2 Kings 2:11",
	},
	{
		ID:         "thorn-paul",
		Difficulty: Hard,
		Word:       "thorn",
		Kind:       OneWordCharacter,
		Answer:     "Paul",
		Reference:  "2 Corinthians 12:7",
	},
	{
		ID:         "purple-lydia",
		Difficulty: Hard,
		Word:       "purple",
		Kind:       OneWordCharacter,
		Answer:     "Lydia",
		Reference:  "Acts 16:14",
	},
	{
		ID:         "potter-jeremiah",
		Difficulty: Hard,
		Word:       "potter",
		Kind:       OneWordCharacter,
		Answer:     "Jeremiah",
		Reference:  "Jeremiah 18:1",
	},
	{
		ID:         "apocalypse-revelation",
		Difficulty: Hard,
		Word:       "apocalypse",
		Kind:       OneWordBook,
		Answer:     "Revelation",
		Reference:  "Revelation 1:1",
	},
	{
		ID:         "camel-john",
		Difficulty: Hard,
		Word:       "camel",
		Kind:       OneWordCharacter,
		Answer:     "John the Baptist",
		Aliases:    []string{"John"},
		Reference:  "Mark 1:6",
	},
}
