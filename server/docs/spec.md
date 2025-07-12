# Form Specification

A form specification file is written using [KDL](https://kdl.dev). It is used to describe a form's various questions and sections in a way that is easy for humans to write and programs to parse.

An example form spec is as follows:

```kdl
form {
	version 1

	title "C Programming - Course Feedback"
	description "Thank you for taking the time to provide feedback. Your responses are anonymous and will help us improve the course for future students."

	question id="semester" type="input" required {
		title "Which semester are you taking this course in?"

		placeholder "Enter the 3 character code. For example, M25 or S24."
		validations {
			regex "^[mMsS][0-9][0-9]$"
			max-chars 3
		}
	}

	section id="content-and-lectures" {
		title "Part 1: Course Content & Lectures"
		description "Please reflect on the lectures, topics covered, and pace of the course."

		question id="lecture-rating" type="likert" required {
			title "The lectures were clear and easy to understand."

			icon "emoji"
			steps 5

			min-label "Strongly Disagree"
			max-label "Strongly Agree"
		}

		question id="course-pace" type="radio" required {
			title "How did you find the pace of the course?"

			option value="slow" label="Too Slow"
			option value="good" label="Just Right"
			option value="fast" label="Too Fast"
		}

		question id="overall-rating" type="likert" required {
			title "How would you rate the overall quality of the course content?"

			icon "star"
			steps 5

			min-label "Poor"
			max-label "Excellent"
		}
	}

	section id="assignments-and-labs" {
		title "Part 2: Assignments and Labs"
		description "Please reflect on the assignments and labs conducted throughout."

		question id="assignment-ratings" type="matrix" required {
			title "Please rate the following aspects of the weekly assignments."

			category value="relevance" label="Relevance to lecture topics"
			category value="difficulty" label="Difficulty level"
			category value="time" label="Time commitment required"

			option value="low" label="Too Low"
			option value="appropriate" label="Appropriate"
			option value="high" label="Too High"
		}

		question id="favourite-assignments" type="checkbox" {
			title "Which assignment(s) did you find most engaging?"

			option value=1 label="Assignment 1: Variables & Control Flow"
			option value=2 label="Assignment 2: Functions & Scope"
			option value=3 label="Assignment 3: Data Structures"
			option value=4 label="Assignment 4: Final Project"
		}
	}

	question id="suggestions" type="textarea" {
		title "Do you have any other comments or suggestions for improving the course?"

		placeholder "Feel free to mention anything, from specific topics you enjoyed to suggestions for the teaching staff."
		validations {
			max-words 500
		}
	}
}
```
