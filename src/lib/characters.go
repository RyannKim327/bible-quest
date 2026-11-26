package bible_lib

type BibleCharacter struct {
	Key   string
	Name  string
	Role  string
	Verse string
}

var Characters = []BibleCharacter{
	{"jesus", "Jesus", "The Savior", "John 3:16"},
	{"abraham", "Abraham", "Father of faith", "Genesis 15:6"},
	{"moses", "Moses", "The lawgiver", "Exodus 3:10"},
	{"david", "David", "The shepherd king", "Psalm 23:1"},
	{"peter", "Peter", "The rock", "1 Peter 5:7"},
	{"paul", "Paul", "The missionary", "Philippians 4:13"},
	{"mary", "Mary", "The mother", "Luke 1:46"},
	{"daniel", "Daniel", "The faithful one", "Daniel 6:23"},
}

var CharacterPools = map[Difficulty][]BibleCharacter{
	Easy:   charactersByKeys("jesus", "abraham", "moses", "david", "mary", "daniel"),
	Medium: charactersByKeys("moses", "david", "peter", "paul", "mary", "daniel"),
	Hard:   charactersByKeys("peter", "paul", "daniel", "jesus", "moses", "abraham"),
}

func charactersByKeys(keys ...string) []BibleCharacter {
	byKey := make(map[string]BibleCharacter, len(Characters))
	for _, character := range Characters {
		byKey[character.Key] = character
	}

	characters := make([]BibleCharacter, 0, len(keys))
	for _, key := range keys {
		if character, ok := byKey[key]; ok {
			characters = append(characters, character)
		}
	}
	return characters
}
