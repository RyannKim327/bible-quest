package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
	bible_lib "constellation/src/lib"
)

const guessScript = `(function () {
  var difficultyPanel = document.getElementById('difficulty-panel');
  var gameArea = document.getElementById('game-area');
  var difficultyButtons = document.querySelectorAll('[data-difficulty]');
  var changeDifficulty = document.getElementById('change-difficulty');
  var passage = document.getElementById('passage');
  if (!passage || !difficultyPanel || !gameArea || !changeDifficulty) return;

  var roundValue = document.getElementById('round-value');
  var scoreValue = document.getElementById('score-value');
  var reference = document.getElementById('reference');
  var submitWord = document.getElementById('submit-word');
  var submitNone = document.getElementById('submit-none');
  var nextRound = document.getElementById('next-round');
  var feedback = document.getElementById('feedback');

  var current = null;
  var selected = [];
  var score = 0;
  var round = 0;
  var locked = true;
  var difficulty = '';

  function setFeedback(message, state) {
    feedback.textContent = message;
    feedback.className = 'feedback' + (state ? ' ' + state : '');
  }

  function loadRound() {
    if (!difficulty) return;

    locked = true;
    selected = [];
    current = null;
    nextRound.hidden = true;
    submitWord.disabled = true;
    submitNone.disabled = false;
    submitNone.textContent = 'No mistake';
    passage.innerHTML = '';
    passage.classList.remove('is-locked');
    reference.textContent = 'Loading passage...';
    setFeedback('Fetching a passage from the archive...');

    fetch('/api/guess-the-mistake?difficulty=' + encodeURIComponent(difficulty), { cache: 'no-store' })
      .then(function (response) {
        if (!response.ok) throw new Error('request failed');
        return response.json();
      })
      .then(function (data) {
        current = data;
        round += 1;
        roundValue.textContent = String(round).padStart(2, '0');
        reference.textContent = (data.reference || 'Unknown passage') + (data.version ? ' · ' + data.version : '');

        data.words.forEach(function (word, index) {
          var button = document.createElement('button');
          button.type = 'button';
          button.className = 'word';
          button.textContent = word;
          button.setAttribute('data-index', String(index));
          button.setAttribute('aria-pressed', 'false');
          button.addEventListener('click', function () {
            toggleWord(index, button);
          });
          passage.appendChild(button);
        });

        locked = false;
        setFeedback('Tap every word that looks wrong. Tap again to deselect.');
      })
      .catch(function () {
        setFeedback('Could not load the passage. Refresh to try again.', 'is-wrong');
      });
  }

  function toggleWord(index, button) {
    if (locked) return;

    var position = selected.indexOf(index);
    if (position === -1) {
      selected.push(index);
      button.classList.add('is-selected');
      button.setAttribute('aria-pressed', 'true');
    } else {
      selected.splice(position, 1);
      button.classList.remove('is-selected');
      button.setAttribute('aria-pressed', 'false');
    }

    submitWord.disabled = selected.length === 0;
    submitNone.textContent = selected.length === 0 ? 'No mistake' : 'Submit';

    if (selected.length === 0) {
      setFeedback('No words selected. Tap a word to toggle it on.');
    } else {
      setFeedback(selected.length + ' word' + (selected.length === 1 ? '' : 's') + ' selected. Tap again to deselect, then choose Missing word.');
    }
  }

  function checkAnswer(wordIndexes) {
    if (locked || !current) return;

    locked = true;
    passage.classList.add('is-locked');
    submitWord.disabled = true;
    submitNone.disabled = true;
    setFeedback('Checking your answer...');

    fetch('/api/guess-the-mistake', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: current.id, words: wordIndexes })
    })
      .then(function (response) {
        if (!response.ok) throw new Error('request failed');
        return response.json();
      })
      .then(function (result) {
        if (result.correct) {
          score += 1;
          scoreValue.textContent = String(score);
        }

        setFeedback(result.message, result.correct ? 'is-correct' : 'is-wrong');

        var words = passage.querySelectorAll('.word');
        (result.mistakeIndices || []).forEach(function (index) {
          if (words[index]) {
            words[index].classList.add('is-answer');
          }
        });

        nextRound.hidden = false;
      })
      .catch(function () {
        setFeedback('Could not check that answer. Try the next passage.', 'is-wrong');
        nextRound.hidden = false;
      });
  }

  submitWord.addEventListener('click', function () {
    if (selected.length === 0) {
      setFeedback('Pick at least one word first.', 'is-wrong');
      return;
    }
    checkAnswer(selected.slice());
  });

  submitNone.addEventListener('click', function () {
    checkAnswer([]);
  });

  function startGame(level) {
    difficulty = level;
    score = 0;
    round = 0;
    scoreValue.textContent = '0';
    roundValue.textContent = '01';
    changeDifficulty.textContent = level.charAt(0).toUpperCase() + level.slice(1);
    difficultyButtons.forEach(function (button) {
      button.classList.toggle('is-selected', button.getAttribute('data-difficulty') === level);
    });
    difficultyPanel.hidden = true;
    gameArea.hidden = false;
    loadRound();
  }

  function showDifficulty() {
    locked = true;
    difficulty = '';
    current = null;
    difficultyPanel.hidden = false;
    gameArea.hidden = true;
    difficultyButtons.forEach(function (button) {
      button.classList.remove('is-selected');
    });
    setFeedback('Choose a difficulty to begin.');
  }

  difficultyButtons.forEach(function (button) {
    button.addEventListener('click', function () {
      startGame(button.getAttribute('data-difficulty'));
    });
  });

  nextRound.addEventListener('click', loadRound);
  changeDifficulty.addEventListener('click', showDifficulty);
})();`

func guessDifficultyButtons() []any {
	options := []any{}
	for _, option := range []struct{ key, label, detail string }{
		{string(bible_lib.Easy), "Easy", bible_lib.Versions[bible_lib.Easy] + " · Obvious swaps and misspellings."},
		{string(bible_lib.Medium), "Medium", bible_lib.Versions[bible_lib.Medium] + " · Real words that still don't belong."},
		{string(bible_lib.Hard), "Hard", bible_lib.Versions[bible_lib.Hard] + " · Tiny changes that are easy to miss."},
	} {
		options = append(options, elements.Button(constellation.WebComponent{
			Class: "difficulty-option",
			Attributes: map[string]string{
				"type":            "button",
				"data-difficulty": option.key,
			},
			Children: []any{
				elements.Strong(constellation.WebComponent{Text: option.label}),
				elements.Span(constellation.WebComponent{Class: "difficulty-detail", Text: option.detail}),
			},
		}))
	}
	return options
}

func GuessTheMistake() string {
	return Parent([]any{
		Header(),
		elements.Section(constellation.WebComponent{
			Class: "guess-page",
			Children: []any{
				elements.Div(constellation.WebComponent{
					Class: "guess-wrap",
					Children: []any{
						elements.A(constellation.WebComponent{
							Class:      "back-link",
							Attributes: map[string]string{"href": "/"},
							Text:       "Back to quests",
						}),
						elements.Div(constellation.WebComponent{
							Class: "section-heading",
							Children: []any{
								elements.Span(constellation.WebComponent{Class: "eyebrow", Text: "GUESS THE MISTAKE"}),
								elements.H1(constellation.WebComponent{Class: "section-title", Text: "Find the word that doesn't belong."}),
								elements.P(constellation.WebComponent{
									Class: "section-lead",
									Text:  "Pick a difficulty first. A passage may hide one or more swapped words. Tap every word that looks wrong, then choose Missing word — or choose No mistake if the passage is right.",
								}),
							},
						}),
						elements.Div(constellation.WebComponent{
							Class:      "difficulty-panel",
							Attributes: map[string]string{"id": "difficulty-panel"},
							Children: []any{
								elements.H2(constellation.WebComponent{Class: "difficulty-heading", Text: "Choose your difficulty"}),
								elements.P(constellation.WebComponent{
									Class: "difficulty-lead",
									Text:  "The swap gets harder to spot as you level up. Pick a level to enter the quest.",
								}),
								elements.Div(constellation.WebComponent{
									Class:    "difficulty-options",
									Children: guessDifficultyButtons(),
								}),
							},
						}),
						elements.Div(constellation.WebComponent{
							Attributes: map[string]string{"id": "game-area", "hidden": "hidden"},
							Children: []any{
								elements.Div(constellation.WebComponent{
									Class: "guess-layout",
									Children: []any{
										elements.Div(constellation.WebComponent{
											Class: "guess-board",
											Children: []any{
												elements.Div(constellation.WebComponent{
													Class: "round-bar",
													Children: []any{
														elements.Span(constellation.WebComponent{Text: "ROUND"}),
														elements.Span(constellation.WebComponent{Class: "round-value", Attributes: map[string]string{"id": "round-value"}, Text: "01"}),
														elements.Span(constellation.WebComponent{Class: "score", Children: []any{
															elements.Span(constellation.WebComponent{Text: "SCORE"}),
															elements.Span(constellation.WebComponent{Class: "score-value", Attributes: map[string]string{"id": "score-value"}, Text: "0"}),
														}}),
														elements.Button(constellation.WebComponent{
															Class: "pixel-button difficulty-change",
															Attributes: map[string]string{
																"id":   "change-difficulty",
																"type": "button",
															},
															Text: "Difficulty",
														}),
														elements.Span(constellation.WebComponent{Class: "reference", Attributes: map[string]string{"id": "reference"}, Text: "Loading passage..."}),
													},
												}),
												elements.Div(constellation.WebComponent{
													Class:      "passage",
													Attributes: map[string]string{"id": "passage"},
													Children: []any{
														elements.Span(constellation.WebComponent{Class: "loading", Text: "Loading passage..."}),
													},
												}),
												elements.Div(constellation.WebComponent{
													Class: "guess-actions",
													Children: []any{
														elements.Button(constellation.WebComponent{
															Class: "pixel-button pixel-button-primary",
															Attributes: map[string]string{
																"id":       "submit-word",
																"type":     "button",
																"disabled": "disabled",
															},
															Text: "Missing word",
														}),
														elements.Button(constellation.WebComponent{
															Class: "pixel-button pixel-button-secondary",
															Attributes: map[string]string{
																"id":   "submit-none",
																"type": "button",
															},
															Text: "No mistake",
														}),
													},
												}),
												elements.Div(constellation.WebComponent{
													Class: "feedback",
													Attributes: map[string]string{
														"id":        "feedback",
														"role":      "status",
														"aria-live": "polite",
													},
													Text: "Tap a word to begin.",
												}),
												elements.Button(constellation.WebComponent{
													Class: "pixel-button next-round",
													Attributes: map[string]string{
														"id":     "next-round",
														"type":   "button",
														"hidden": "hidden",
													},
													Text: "Next passage",
												}),
											},
										}),
										elements.Aside(constellation.WebComponent{
											Class: "guess-guide",
											Children: []any{
												elements.Div(constellation.WebComponent{
													Class:    "guide-sprite",
													Children: []any{CharacterSprite("paul", "sprite sprite-guide")},
												}),
												elements.H3(constellation.WebComponent{Class: "guide-title", Text: "How to play"}),
												elements.Ol(constellation.WebComponent{
													Class: "guide-list",
													Children: []any{
														elements.Li(constellation.WebComponent{Text: "Pick a difficulty to begin."}),
														elements.Li(constellation.WebComponent{Text: "Read the passage carefully."}),
														elements.Li(constellation.WebComponent{Text: "Tap every word you think was swapped."}),
														elements.Li(constellation.WebComponent{Text: "Tap a word again to turn it off."}),
														elements.Li(constellation.WebComponent{Text: "Choose Missing word or No mistake."}),
													},
												}),
											},
										}),
									},
								}),
							},
						}),
					},
				}),
			},
		}),
		elements.Script(constellation.WebComponent{Children: []any{guessScript}}),
		Footer(),
	})
}
