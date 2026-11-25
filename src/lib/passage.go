package bible_lib

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

type PassageMistakes map[Difficulty]map[string]string

type PassageTexts map[Difficulty]string

var Versions = map[Difficulty]string{
	Easy:   "NIV",
	Medium: "NLT",
	Hard:   "KJV",
}

type Passage struct {
	ID        string
	Reference string
	Texts     PassageTexts
	Mistakes  PassageMistakes
}

var Passages = []Passage{
	{
		ID:        "john-3-16",
		Reference: "John 3:16",
		Texts: PassageTexts{
			Easy:   "For God so loved the world that he gave his one and only Son, that whoever believes in him shall not perish but have eternal life.",
			Medium: "For this is how God loved the world: He gave his one and only Son, so that everyone who believes in him will not perish but have eternal life.",
			Hard:   "For God so loved the world, that he gave his only begotten Son, that whosoever believeth in him should not perish, but have everlasting life.",
		},
		Mistakes: PassageMistakes{
			Easy:   {"loved": "love"},
			Medium: {"will": "would", "not": "be"},
			Hard:   {"should not": "shall", "but have": "and not"},
		},
	},
	{
		ID:        "genesis-1-1",
		Reference: "Genesis 1:1",
		Texts: PassageTexts{
			Easy:   "In the beginning God created the heavens and the earth.",
			Medium: "In the beginning God created the heavens and the earth.",
			Hard:   "In the beginning God created the heaven and the earth.",
		},
		Mistakes: PassageMistakes{
			Easy:   {"In": "From", "beginning": "start"},
			Medium: {"created": "made"},
			Hard:   {"beginning": "begining"},
		},
	},
	{
		ID:        "psalm-23-1",
		Reference: "Psalm 23:1",
		Texts: PassageTexts{
			Easy:   "The Lord is my shepherd; I have all that I need.",
			Medium: "The Lord is my shepherd, I lack nothing.",
			Hard:   "The Lord is my shepherd; I shall not want.",
		},
		Mistakes: PassageMistakes{
			Easy:   {"my": "the", "have all that I": "don't"},
			Medium: nil,
			Hard:   {"shepherd": "shephard"},
		},
	},
	{
		ID:        "proverbs-3-5",
		Reference: "Proverbs 3:5",
		Texts: PassageTexts{
			Easy:   "Trust in the Lord with all your heart; do not depend on your own understanding.",
			Medium: "Trust in the Lord with all your heart and lean not on your own understanding;",
			Hard:   "Trust in the Lord with all thine heart; and lean not unto thine own understanding.",
		},
		Mistakes: PassageMistakes{
			Easy:   nil,
			Medium: {"in the Lord with all your": "your", "lean not": "and lean"},
			Hard:   {"Trust in the Lord with all thine": "The Lord trusted your", "lean not": "lean"},
		},
	},
	{
		ID:        "luke-6-31",
		Reference: "Luke 6:31",
		Texts: PassageTexts{
			Easy:   "Do to others as you would like them to do to you.",
			Medium: "Do to others as you would have them do to you.",
			Hard:   "And as ye would that men should do to you, do ye also to them likewise.",
		},
		Mistakes: PassageMistakes{
			Easy:   nil,
			Medium: nil,
			Hard:   {"men": "they"},
		},
	},
	{
		ID:        "philippians-4-13",
		Reference: "Philippians 4:13",
		Texts: PassageTexts{
			Easy:   "For I can do everything through Christ, who gives me strength.",
			Medium: "I can do all this through him who gives me strength.",
			Hard:   "I can do all things through Christ which strengtheneth me.",
		},
		Mistakes: PassageMistakes{
			Easy:   {"I can": "you can't", "though": "without"},
			Medium: {"him": "his"},
			Hard:   {"through": "though"},
		},
	},
	{
		ID:        "romans-8-28",
		Reference: "Romans 8:28",
		Texts: PassageTexts{
			Easy:   "And we know that God causes everything to work together for the good of those who love God and are called according to his purpose for them.",
			Medium: "And we know that in all things God works for the good of those who love him, who have been called according to his purpose.",
			Hard:   "And we know that all things work together for good to them that love God, to them who are the called according to his purpose.",
		},
		Mistakes: PassageMistakes{
			Easy:   {"the good": "better,"},
			Medium: {"God": "god", "purpose": "propose"},
			Hard:   {"we": "you", "good": "god", "God": "god"},
		},
	},
	{
		ID:        "joshua-1-8",
		Reference: "Joshua 1:8",
		Texts: PassageTexts{
			Easy:   "Study this Book of Instruction continually. Meditate on it day and night so you will be sure to obey everything written in it. Only then will you prosper and succeed in all you do.",
			Medium: "Keep this Book of the Law always on your lips; meditate on it day and night, so that you may be careful to do everything written in it. Then you will be prosperous and successful.",
			Hard:   "This book of the law shall not depart out of thy mouth; but thou shalt meditate therein day and night, that thou mayest observe to do according to all that is written therein: for then thou shalt make thy way prosperous, and then thou shalt have good success.",
		},
		Mistakes: PassageMistakes{
			Easy:   nil,
			Medium: {"meditate": "mediate", "prosperous": "generous"},
			Hard:   {"shall not": "shall", "meditate": "medidate"},
		},
	},
}
