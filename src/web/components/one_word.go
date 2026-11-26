package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
	bible_lib "constellation/src/lib"
)

const oneWordScript = `(function () {
  var difficultyPanel = document.getElementById('oneword-difficulty-panel');
  var gameArea = document.getElementById('oneword-game-area');
  var difficultyButtons = document.querySelectorAll('[data-difficulty]');
  var changeDifficulty = document.getElementById('oneword-change-difficulty');
  var clue = document.getElementById('oneword-clue');
  var kind = document.getElementById('oneword-kind');
  var form = document.getElementById('oneword-form');
  var input = document.getElementById('oneword-input');
  var submit = document.getElementById('oneword-submit');
  var reveal = document.getElementById('oneword-reveal');
  var nextRound = document.getElementById('oneword-next-round');
  var feedback = document.getElementById('oneword-feedback');
  var guideBubble = document.getElementById('oneword-guide-bubble');
  var guideSprite = document.getElementById('oneword-guide-sprite');
  var guideName = document.getElementById('oneword-guide-name');
  var guideRole = document.getElementById('oneword-guide-role');
  var guideNote = document.getElementById('oneword-guide-note');
  var guideTemplates = document.querySelectorAll('#oneword-character-templates template');
  if (!difficultyPanel || !gameArea || !changeDifficulty || !clue || !kind || !form || !input || !submit || !reveal || !nextRound || !feedback || !guideBubble || !guideSprite || !guideName || !guideRole || !guideNote) return;

  var roundValue = document.getElementById('oneword-round-value');
  var scoreValue = document.getElementById('oneword-score-value');

  var current = null;
  var score = 0;
  var round = 0;
  var locked = true;
  var difficulty = '';

  function bindTap(element, handler) {
    var startX = 0;
    var startY = 0;
    var touched = false;
    element.addEventListener('touchstart', function (event) {
      startX = event.touches[0].clientX;
      startY = event.touches[0].clientY;
    }, { passive: true });
    element.addEventListener('touchend', function (event) {
      var touch = event.changedTouches[0];
      if (touched || Math.abs(touch.clientX - startX) > 10 || Math.abs(touch.clientY - startY) > 10) return;
      touched = true;
      handler(event);
      window.setTimeout(function () { touched = false; }, 700);
    }, { passive: true });
    element.addEventListener('click', function (event) {
      if (touched) return;
      handler(event);
    });
  }

  function setFeedback(message, state) {
    feedback.textContent = message;
    feedback.className = 'feedback' + (state ? ' ' + state : '');
  }

  function popGuide() {
    guideBubble.classList.remove('is-popping');
    void guideBubble.offsetWidth;
    guideBubble.classList.add('is-popping');
  }

  function setGuide(character, note) {
    if (!character) return;
    guideSprite.innerHTML = '';
    guideTemplates.forEach(function (template) {
      if (template.getAttribute('data-character-key') === character.key) {
        guideSprite.appendChild(template.content.cloneNode(true));
      }
    });
    guideName.textContent = character.name;
    guideRole.textContent = character.role;
    guideNote.textContent = note || character.verse;
    popGuide();
  }

  function loadRound() {
    if (!difficulty) return;

    locked = true;
    current = null;
    nextRound.hidden = true;
    reveal.hidden = true;
    reveal.textContent = '';
    input.value = '';
    input.disabled = false;
    submit.disabled = false;
    clue.textContent = '...';
    kind.textContent = '';
    guideSprite.innerHTML = '';
    guideName.textContent = 'Your guide';
    guideRole.textContent = '';
    guideNote.textContent = 'Finding a guide for this round...';
    setFeedback('Fetching a clue from the archive...');

    fetch('/api/one-word?difficulty=' + encodeURIComponent(difficulty), { cache: 'no-store' })
      .then(function (response) {
        if (!response.ok) throw new Error('request failed');
        return response.json();
      })
      .then(function (data) {
        current = data;
        round += 1;
        roundValue.textContent = String(round).padStart(2, '0');
        clue.textContent = data.word || '';
        kind.textContent = data.kind === 'book' ? 'BOOK' : 'CHARACTER';
        setGuide(data.character, data.character ? data.character.verse : '');
        locked = false;
        input.focus();
        setFeedback('Name the ' + (data.kind === 'book' ? 'book' : 'character') + ' this word points to.');
      })
      .catch(function () {
        setFeedback('Could not load the clue. Refresh to try again.', 'is-wrong');
      });
  }

  function checkAnswer() {
    if (locked || !current) return;

    var answer = input.value.trim();
    if (!answer) {
      setFeedback('Type a name before submitting.', 'is-wrong');
      input.focus();
      return;
    }

    locked = true;
    input.disabled = true;
    submit.disabled = true;
    setFeedback('Checking your answer...');

    fetch('/api/one-word', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: current.id, answer: answer })
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
        reveal.textContent = 'Answer: ' + result.answer + ' · ' + result.reference;
        reveal.hidden = false;
        if (current.character) {
          setGuide(current.character, current.character.name + ' says: “' + result.answer + ' — ' + result.reference + '.”');
        }
        nextRound.hidden = false;
      })
      .catch(function () {
        setFeedback('Could not check that answer. Try the next clue.', 'is-wrong');
        nextRound.hidden = false;
      });
  }

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

  form.addEventListener('submit', function (event) {
    event.preventDefault();
    checkAnswer();
  });
  bindTap(nextRound, loadRound);
  bindTap(changeDifficulty, showDifficulty);
  difficultyButtons.forEach(function (button) {
    bindTap(button, function () {
      startGame(button.getAttribute('data-difficulty'));
    });
  });
})();`

func OneWord() string {
	details := map[bible_lib.Difficulty]string{
		bible_lib.Easy:   "Famous people and well-known books.",
		bible_lib.Medium: "Less obvious links from the story.",
		bible_lib.Hard:   "Small details that take a closer look.",
	}

	return Parent([]any{
		Header(),
		elements.Section(constellation.WebComponent{
			Class: "game-page",
			Children: []any{
				elements.Div(constellation.WebComponent{
					Class: "game-wrap",
					Children: []any{
						elements.A(constellation.WebComponent{
							Class:      "back-link",
							Attributes: map[string]string{"href": "/"},
							Text:       "Back to quests",
						}),
						elements.Div(constellation.WebComponent{
							Class: "section-heading",
							Children: []any{
								elements.Span(constellation.WebComponent{Class: "eyebrow", Text: "ONE WORD"}),
								elements.H1(constellation.WebComponent{Class: "section-title", Text: "One word. One answer."}),
								elements.P(constellation.WebComponent{
									Class: "section-lead",
									Text:  "A single word hints at a person or a book from the Bible. Read the clue, type your answer, and find out if you got it right.",
								}),
							},
						}),
						elements.Div(constellation.WebComponent{
							Class:      "difficulty-panel",
							Attributes: map[string]string{"id": "oneword-difficulty-panel"},
							Children: []any{
								elements.H2(constellation.WebComponent{Class: "difficulty-heading", Text: "Choose your difficulty"}),
								elements.P(constellation.WebComponent{
									Class: "difficulty-lead",
									Text:  "The links get less obvious as you level up. Pick a level to enter the game.",
								}),
								elements.Div(constellation.WebComponent{
									Class:    "difficulty-options",
									Children: gameDifficultyButtons(details),
								}),
							},
						}),
						elements.Div(constellation.WebComponent{
							Attributes: map[string]string{"id": "oneword-game-area", "hidden": "hidden"},
							Children: []any{
								elements.Div(constellation.WebComponent{
									Class: "game-layout",
									Children: []any{
										elements.Div(constellation.WebComponent{
											Class: "game-board",
											Children: []any{
												elements.Div(constellation.WebComponent{
													Class: "round-bar",
													Children: []any{
														elements.Span(constellation.WebComponent{Text: "ROUND"}),
														elements.Span(constellation.WebComponent{Class: "round-value", Attributes: map[string]string{"id": "oneword-round-value"}, Text: "01"}),
														elements.Span(constellation.WebComponent{Class: "score", Children: []any{
															elements.Span(constellation.WebComponent{Text: "SCORE"}),
															elements.Span(constellation.WebComponent{Class: "score-value", Attributes: map[string]string{"id": "oneword-score-value"}, Text: "0"}),
														}}),
														elements.Button(constellation.WebComponent{
															Class: "pixel-button difficulty-change",
															Attributes: map[string]string{
																"id":   "oneword-change-difficulty",
																"type": "button",
															},
															Text: "Difficulty",
														}),
													},
												}),
												elements.Div(constellation.WebComponent{
													Class: "oneword-clue-wrap",
													Children: []any{
														elements.Span(constellation.WebComponent{Class: "answer-kind", Attributes: map[string]string{"id": "oneword-kind"}, Text: ""}),
														elements.P(constellation.WebComponent{
															Class:      "oneword-clue",
															Attributes: map[string]string{"id": "oneword-clue"},
															Text:       "...",
														}),
													},
												}),
												elements.Form(constellation.WebComponent{
													Class:      "oneword-form",
													Attributes: map[string]string{"id": "oneword-form"},
													Children: []any{
														elements.Label(constellation.WebComponent{
															Class:      "oneword-label",
															Attributes: map[string]string{"for": "oneword-input"},
															Text:       "Your answer",
														}),
														elements.Input(constellation.WebComponent{
															Class: "oneword-input",
															Attributes: map[string]string{
																"id":           "oneword-input",
																"type":         "text",
																"name":         "answer",
																"autocomplete": "off",
																"placeholder":  "Type a name...",
															},
														}),
														elements.Button(constellation.WebComponent{
															Class: "pixel-button pixel-button-primary",
															Attributes: map[string]string{
																"id":   "oneword-submit",
																"type": "submit",
															},
															Text: "Check answer",
														}),
													},
												}),
												elements.Div(constellation.WebComponent{
													Class:      "oneword-reveal",
													Attributes: map[string]string{"id": "oneword-reveal", "hidden": "hidden"},
												}),
												elements.Div(constellation.WebComponent{
													Class: "feedback",
													Attributes: map[string]string{
														"id":        "oneword-feedback",
														"role":      "status",
														"aria-live": "polite",
													},
													Text: "Pick a difficulty to begin.",
												}),
												elements.Button(constellation.WebComponent{
													Class: "pixel-button next-round",
													Attributes: map[string]string{
														"id":     "oneword-next-round",
														"type":   "button",
														"hidden": "hidden",
													},
													Text: "Next clue",
												}),
											},
										}),
										elements.Aside(constellation.WebComponent{
											Class: "game-guide",
											Children: append(characterGuideIntro("oneword"),
												elements.H3(constellation.WebComponent{Class: "guide-title", Text: "How to play"}),
												elements.Ol(constellation.WebComponent{
													Class: "guide-list",
													Children: []any{
														elements.Li(constellation.WebComponent{Text: "Pick a difficulty to begin."}),
														elements.Li(constellation.WebComponent{Text: "Read the single-word clue."}),
														elements.Li(constellation.WebComponent{Text: "Check whether it points to a character or a book."}),
														elements.Li(constellation.WebComponent{Text: "Type your answer and submit."}),
														elements.Li(constellation.WebComponent{Text: "Learn the Bible reference behind it."}),
													},
												}),
											),
										}),
									},
								}),
							},
						}),
					},
				}),
			},
		}),
		elements.Script(constellation.WebComponent{Children: []any{oneWordScript}}),
		Footer(),
	})
}
