package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
	bible_lib "constellation/src/lib"
)

const quizScript = `(function () {
  var difficultyPanel = document.getElementById('quiz-difficulty-panel');
  var gameArea = document.getElementById('quiz-game-area');
  var difficultyButtons = document.querySelectorAll('[data-difficulty]');
  var changeDifficulty = document.getElementById('quiz-change-difficulty');
  var question = document.getElementById('quiz-question');
  var choices = document.getElementById('quiz-choices');
  var submit = document.getElementById('quiz-submit');
  var nextRound = document.getElementById('quiz-next-round');
  var feedback = document.getElementById('quiz-feedback');
  var guideBubble = document.getElementById('quiz-guide-bubble');
  var guideSprite = document.getElementById('quiz-guide-sprite');
  var guideName = document.getElementById('quiz-guide-name');
  var guideRole = document.getElementById('quiz-guide-role');
  var guideNote = document.getElementById('quiz-guide-note');
  var guideTemplates = document.querySelectorAll('#quiz-character-templates template');
  if (!difficultyPanel || !gameArea || !changeDifficulty || !question || !choices || !submit || !nextRound || !feedback || !guideBubble || !guideSprite || !guideName || !guideRole || !guideNote) return;

  var roundValue = document.getElementById('quiz-round-value');
  var scoreValue = document.getElementById('quiz-score-value');
  var reference = document.getElementById('quiz-reference');

  var current = null;
  var selectedText = '';
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
    selectedText = '';
    current = null;
    nextRound.hidden = true;
    submit.disabled = true;
    choices.innerHTML = '';
    choices.classList.remove('is-locked');
    question.textContent = 'Loading question...';
    reference.textContent = '';
    guideSprite.innerHTML = '';
    guideName.textContent = 'Your guide';
    guideRole.textContent = '';
    guideNote.textContent = 'Finding a guide for this round...';
    setFeedback('Fetching a question from the archive...');

    fetch('/api/quiz?difficulty=' + encodeURIComponent(difficulty), { cache: 'no-store' })
      .then(function (response) {
        if (!response.ok) throw new Error('request failed');
        return response.json();
      })
      .then(function (data) {
        current = data;
        round += 1;
        roundValue.textContent = String(round).padStart(2, '0');
        question.textContent = data.question || '';
        reference.textContent = data.reference || '';
        setGuide(data.character, data.character ? data.character.verse : '');

        (data.choices || []).forEach(function (choice, index) {
          var button = document.createElement('button');
          button.type = 'button';
          button.className = 'quiz-choice';
          button.textContent = choice;
          button.setAttribute('data-index', String(index));
          button.setAttribute('aria-pressed', 'false');
          bindTap(button, function () {
            selectChoice(index, button);
          });
          choices.appendChild(button);
        });

        locked = false;
        setFeedback('Pick the best answer.');
      })
      .catch(function () {
        setFeedback('Could not load the question. Refresh to try again.', 'is-wrong');
      });
  }

  function selectChoice(index, button) {
    if (locked) return;

    selectedText = button.textContent;
    choices.querySelectorAll('.quiz-choice').forEach(function (choice) {
      choice.classList.remove('is-selected');
      choice.setAttribute('aria-pressed', 'false');
    });
    button.classList.add('is-selected');
    button.setAttribute('aria-pressed', 'true');
    submit.disabled = false;
    setFeedback('Answer selected. Submit when ready.');
  }

  function checkAnswer() {
    if (locked || !current || !selectedText) return;

    locked = true;
    choices.classList.add('is-locked');
    submit.disabled = true;
    setFeedback('Checking your answer...');

    fetch('/api/quiz', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: current.id, answer: selectedText })
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

        choices.querySelectorAll('.quiz-choice').forEach(function (button) {
          button.disabled = true;
          if (button.textContent === result.answerText) {
            button.classList.add('is-answer');
          } else if (button.classList.contains('is-selected')) {
            button.classList.add('is-wrong-choice');
          }
        });

        reference.textContent = result.reference || '';
        if (current.character) {
          setGuide(current.character, current.character.name + ' says: “' + result.answerText + ' — ' + result.reference + '.”');
        }
        nextRound.hidden = false;
      })
      .catch(function () {
        setFeedback('Could not check that answer. Try the next question.', 'is-wrong');
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

  bindTap(submit, checkAnswer);
  bindTap(nextRound, loadRound);
  bindTap(changeDifficulty, showDifficulty);
  difficultyButtons.forEach(function (button) {
    bindTap(button, function () {
      startGame(button.getAttribute('data-difficulty'));
    });
  });
})();`

func Quiz() string {
	details := map[bible_lib.Difficulty]string{
		bible_lib.Easy:   "Well-known stories and facts from the Bible.",
		bible_lib.Medium: "Details you know once you slow down and think.",
		bible_lib.Hard:   "Specific people, places, and moments to recall.",
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
								elements.Span(constellation.WebComponent{Class: "eyebrow", Text: "BIBLE QUIZ"}),
								elements.H1(constellation.WebComponent{Class: "section-title", Text: "Answer the call."}),
								elements.P(constellation.WebComponent{
									Class: "section-lead",
									Text:  "Pick a difficulty, choose one answer from four, and see how much of the story you remember. Every question reveals its Bible reference after you answer.",
								}),
							},
						}),
						elements.Div(constellation.WebComponent{
							Class:      "difficulty-panel",
							Attributes: map[string]string{"id": "quiz-difficulty-panel"},
							Children: []any{
								elements.H2(constellation.WebComponent{Class: "difficulty-heading", Text: "Choose your difficulty"}),
								elements.P(constellation.WebComponent{
									Class: "difficulty-lead",
									Text:  "The questions get more specific as you level up. Pick a level to enter the quiz.",
								}),
								elements.Div(constellation.WebComponent{
									Class:    "difficulty-options",
									Children: gameDifficultyButtons(details),
								}),
							},
						}),
						elements.Div(constellation.WebComponent{
							Attributes: map[string]string{"id": "quiz-game-area", "hidden": "hidden"},
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
														elements.Span(constellation.WebComponent{Class: "round-value", Attributes: map[string]string{"id": "quiz-round-value"}, Text: "01"}),
														elements.Span(constellation.WebComponent{Class: "score", Children: []any{
															elements.Span(constellation.WebComponent{Text: "SCORE"}),
															elements.Span(constellation.WebComponent{Class: "score-value", Attributes: map[string]string{"id": "quiz-score-value"}, Text: "0"}),
														}}),
														elements.Button(constellation.WebComponent{
															Class: "pixel-button difficulty-change",
															Attributes: map[string]string{
																"id":   "quiz-change-difficulty",
																"type": "button",
															},
															Text: "Difficulty",
														}),
														elements.Span(constellation.WebComponent{Class: "reference", Attributes: map[string]string{"id": "quiz-reference"}, Text: ""}),
													},
												}),
												elements.H2(constellation.WebComponent{
													Class:      "quiz-prompt",
													Attributes: map[string]string{"id": "quiz-question"},
													Text:       "Loading question...",
												}),
												elements.Div(constellation.WebComponent{
													Class:      "quiz-choices",
													Attributes: map[string]string{"id": "quiz-choices"},
													Children: []any{
														elements.Span(constellation.WebComponent{Class: "loading", Text: "Loading question..."}),
													},
												}),
												elements.Div(constellation.WebComponent{
													Class: "game-actions",
													Children: []any{
														elements.Button(constellation.WebComponent{
															Class: "pixel-button pixel-button-primary",
															Attributes: map[string]string{
																"id":       "quiz-submit",
																"type":     "button",
																"disabled": "disabled",
															},
															Text: "Submit answer",
														}),
													},
												}),
												elements.Div(constellation.WebComponent{
													Class: "feedback",
													Attributes: map[string]string{
														"id":        "quiz-feedback",
														"role":      "status",
														"aria-live": "polite",
													},
													Text: "Pick a difficulty to begin.",
												}),
												elements.Button(constellation.WebComponent{
													Class: "pixel-button next-round",
													Attributes: map[string]string{
														"id":     "quiz-next-round",
														"type":   "button",
														"hidden": "hidden",
													},
													Text: "Next question",
												}),
											},
										}),
										elements.Aside(constellation.WebComponent{
											Class: "game-guide",
											Children: append(characterGuideIntro("quiz"),
												elements.H3(constellation.WebComponent{Class: "guide-title", Text: "How to play"}),
												elements.Ol(constellation.WebComponent{
													Class: "guide-list",
													Children: []any{
														elements.Li(constellation.WebComponent{Text: "Pick a difficulty to begin."}),
														elements.Li(constellation.WebComponent{Text: "Read the question carefully."}),
														elements.Li(constellation.WebComponent{Text: "Choose one of the four answers."}),
														elements.Li(constellation.WebComponent{Text: "Submit and read the Bible reference."}),
														elements.Li(constellation.WebComponent{Text: "Move on to the next question."}),
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
		elements.Script(constellation.WebComponent{Children: []any{quizScript}}),
		Footer(),
	})
}
