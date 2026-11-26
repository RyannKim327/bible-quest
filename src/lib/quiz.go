package bible_lib

type QuizQuestion struct {
	ID         string
	Difficulty Difficulty
	Prompt     string
	Choices    []string
	Answer     int
	Reference  string
}

var QuizQuestions = []QuizQuestion{
	{
		ID:         "ark-builder",
		Difficulty: Easy,
		Prompt:     "Who built the ark?",
		Choices:    []string{"Abraham", "Moses", "Noah", "David"},
		Answer:     2,
		Reference:  "Genesis 6:14",
	},
	{
		ID:         "first-book",
		Difficulty: Easy,
		Prompt:     "What is the first book of the Bible?",
		Choices:    []string{"Genesis", "Exodus", "Matthew", "Psalms"},
		Answer:     0,
		Reference:  "Genesis 1:1",
	},
	{
		ID:         "garden-serpent",
		Difficulty: Easy,
		Prompt:     "Which animal spoke to Eve in the garden?",
		Choices:    []string{"Donkey", "Camel", "Dove", "Serpent"},
		Answer:     3,
		Reference:  "Genesis 3:1",
	},
	{
		ID:         "first-king",
		Difficulty: Easy,
		Prompt:     "Who was the first king of Israel?",
		Choices:    []string{"David", "Saul", "Solomon", "Samuel"},
		Answer:     1,
		Reference:  "1 Samuel 10:24",
	},
	{
		ID:         "jesus-birthplace",
		Difficulty: Easy,
		Prompt:     "Where was Jesus born?",
		Choices:    []string{"Nazareth", "Jerusalem", "Bethlehem", "Capernaum"},
		Answer:     2,
		Reference:  "Luke 2:7",
	},
	{
		ID:         "creation-days",
		Difficulty: Easy,
		Prompt:     "How many days did God take to create the world?",
		Choices:    []string{"Seven", "Ten", "Forty", "Six"},
		Answer:     3,
		Reference:  "Genesis 1:31",
	},
	{
		ID:         "joseph-sold",
		Difficulty: Medium,
		Prompt:     "Which brother sold Joseph into slavery?",
		Choices:    []string{"Reuben", "Levi", "Benjamin", "Judah"},
		Answer:     3,
		Reference:  "Genesis 37:26",
	},
	{
		ID:         "goliath",
		Difficulty: Medium,
		Prompt:     "What did David use to defeat Goliath?",
		Choices:    []string{"A spear", "A sword", "A sling and a stone", "A bow and arrow"},
		Answer:     2,
		Reference:  "1 Samuel 17:50",
	},
	{
		ID:         "peter-denial",
		Difficulty: Medium,
		Prompt:     "Which apostle denied Jesus three times?",
		Choices:    []string{"Andrew", "Peter", "James", "John"},
		Answer:     1,
		Reference:  "Luke 22:61",
	},
	{
		ID:         "birth-angel",
		Difficulty: Medium,
		Prompt:     "Which angel announced the birth of Jesus?",
		Choices:    []string{"Michael", "Raphael", "Gabriel", "Uriel"},
		Answer:     2,
		Reference:  "Luke 1:26",
	},
	{
		ID:         "cana-wine",
		Difficulty: Medium,
		Prompt:     "What did Jesus turn water into at Cana?",
		Choices:    []string{"Bread", "Wine", "Oil", "Vinegar"},
		Answer:     1,
		Reference:  "John 2:1",
	},
	{
		ID:         "after-john",
		Difficulty: Medium,
		Prompt:     "Which book comes immediately after John?",
		Choices:    []string{"Acts", "Romans", "James", "Revelation"},
		Answer:     0,
		Reference:  "Acts 1:1",
	},
	{
		ID:         "chariot-fire",
		Difficulty: Hard,
		Prompt:     "Which prophet was taken up in a chariot of fire?",
		Choices:    []string{"Elisha", "Elijah", "Enoch", "Isaiah"},
		Answer:     1,
		Reference:  "2 Kings 2:11",
	},
	{
		ID:         "first-martyr",
		Difficulty: Hard,
		Prompt:     "Who was the first Christian martyr?",
		Choices:    []string{"Stephen", "Peter", "Paul", "James"},
		Answer:     0,
		Reference:  "Acts 7:59",
	},
	{
		ID:         "twelve-spies",
		Difficulty: Hard,
		Prompt:     "Which two spies brought back grapes from Canaan?",
		Choices:    []string{"Moses and Aaron", "Gideon and Samson", "Joshua and Caleb", "David and Jonathan"},
		Answer:     2,
		Reference:  "Numbers 13:23",
	},
	{
		ID:         "statue-dream",
		Difficulty: Hard,
		Prompt:     "Which king saw a great statue in a dream?",
		Choices:    []string{"Belshazzar", "Nebuchadnezzar", "Cyrus", "Darius"},
		Answer:     1,
		Reference:  "Daniel 2:1",
	},
	{
		ID:         "patmos-author",
		Difficulty: Hard,
		Prompt:     "Who wrote Revelation while exiled on Patmos?",
		Choices:    []string{"Peter", "Paul", "James", "John"},
		Answer:     3,
		Reference:  "Revelation 1:9",
	},
	{
		ID:         "roman-governor",
		Difficulty: Hard,
		Prompt:     "Which Roman governor sentenced Jesus to be crucified?",
		Choices:    []string{"Herod", "Pontius Pilate", "Felix", "Caesar Augustus"},
		Answer:     1,
		Reference:  "Matthew 27:26",
	},
}
