import { ulid } from 'ulid';
import * as kdljs from 'kdljs';
import { Time } from '@internationalized/date';
import { parseKdlForm, safeString } from '$lib/utils/kdl';
import type { Question, QuestionType, Option, FormData } from '$lib/types/form';

// store creation function
export function createFormStore(initialForm: any) {
	const formData = $state<FormData>({
		title: 'Untitled Form',
		description: 'Add a description',
		visibility: 'private',
		opens: undefined,
		closes: undefined,
		opensTime: undefined,
		closesTime: undefined,
		anonymous: false,
		max_responses: null,
		individual_limit: 1,
		editable_responses: false
	});

	let questions = $state<Question[]>([]);

	const isFormValid = $derived.by(() => {
		for (const question of questions) {
			if (question.error) return false;

			if (question.options) {
				for (const option of question.options) {
					if (option.error) return false;
				}
			}
		}
		return true;
	});

	function loadForm() {
		if (initialForm) {
			formData.title = initialForm.title || 'Untitled Form';
			formData.description = initialForm.description || 'Add a description';

			if (initialForm.opens) {
				const [dateStr, timeStr] = initialForm.opens.split('T');
				formData.opens = dateStr;
				if (timeStr) {
					const [hours, minutes] = timeStr.split(':').map(Number);
					formData.opensTime = new Time(hours || 0, minutes || 0);
				}
			} else {
				formData.opens = undefined;
				formData.opensTime = undefined;
			}

			if (initialForm.closes) {
				const [dateStr, timeStr] = initialForm.closes.split('T');
				formData.closes = dateStr;
				if (timeStr) {
					const [hours, minutes] = timeStr.split(':').map(Number);
					formData.closesTime = new Time(hours || 0, minutes || 0);
				}
			} else {
				formData.closes = undefined;
				formData.closesTime = undefined;
			}

			formData.anonymous = initialForm.anonymous || false;
			formData.max_responses = initialForm.max_responses ?? null;
			formData.individual_limit = initialForm.individual_limit ?? 1;
			formData.editable_responses = initialForm.editable_responses || false;

			if (initialForm.structure) {
				try {
					const result = parseKdlForm(initialForm.structure);
					questions = result.questions;
				} catch (e) {
					console.error('Error parsing form KDL:', e);
					questions = [];
				}
			} else {
				questions = [];
			}
		}
	}

	// validation

	function validateQuestions() {
		for (const question of questions) {
			question.error = undefined;

			if (!question.title || question.title.trim() === '') {
				question.error = "Question text cannot be empty.";
			}

			if (question.options && (question.type === 'radio' || question.type === 'checkbox' || question.type === 'select')) {
				const seenLabels = new Set<string>();

				for (const option of question.options) {
					option.error = undefined;

					if (!option.label || option.label.trim() === '') {
						option.error = "Option text cannot be empty.";
						continue;
					}

					const trimmedLabel = option.label.trim().toLowerCase();
					if (seenLabels.has(trimmedLabel)) {
						option.error = "Duplicate option found.";
					} else {
						seenLabels.add(trimmedLabel);
					}
				}

				const labelCounts = new Map<string, Option[]>();
				for (const option of question.options) {
					if (option.label && option.label.trim() !== '') {
						const trimmedLabel = option.label.trim().toLowerCase();
						if (!labelCounts.has(trimmedLabel)) {
							labelCounts.set(trimmedLabel, []);
						}
						labelCounts.get(trimmedLabel)!.push(option);
					}
				}

				for (const [label, options] of labelCounts) {
					if (options.length > 1) {
						for (const option of options) {
							option.error = "Duplicate option found.";
						}
					}
				}
			}
		}
	}

	// actions

	function addQuestion(type: QuestionType) {
		const newQuestion: Question = {
			id: ulid(),
			type,
			title: '',
			required: false
		};
		if (type === 'radio' || type === 'checkbox' || type === 'select') {
			newQuestion.options = [
				{ id: ulid(), value: '', label: '' },
				{ id: ulid(), value: '', label: '' }
			];
		}
		if (type === 'input' || type === 'textarea') {
			newQuestion.placeholder = '';
			newQuestion.validations = {};
		}
		if (type === 'file') {
			newQuestion['max-file-size'] = 10;
			newQuestion['max-files'] = -1;
			newQuestion['allowed-types'] = [];
		}
		if (type === 'section-header') {
			newQuestion.description = '';
			newQuestion.required = false; // Section headers are never required
		}
		questions.push(newQuestion);
	}

	function updateQuestion(id: string, updatedQuestion: Partial<Question>) {
		const index = questions.findIndex((q) => q.id === id);
		if (index !== -1) {
			questions[index] = { ...questions[index], ...updatedQuestion };

			// Only clear question title error if title is actually being updated and is valid
			if (updatedQuestion.title !== undefined && questions[index].error) {
				if (updatedQuestion.title.trim() !== '') {
					questions[index].error = undefined;
				}
			}
		}
	}

	function removeQuestion(id: string) {
		const index = questions.findIndex((q) => q.id === id);
		if (index > -1) {
			questions.splice(index, 1);
		}
	}

	function addOption(questionId: string) {
		const question = questions.find((q) => q.id === questionId);
		if (!question?.options) return;

		const newOption: Option = { id: ulid(), value: '', label: '' };
		const updatedOptions = [...question.options, newOption];
		updateQuestion(questionId, { options: updatedOptions });
	}

	function removeOption(questionId: string, optionId: string) {
		const question = questions.find((q) => q.id === questionId);
		if (!question?.options) return;

		const updatedOptions = question.options.filter((option) => option.id !== optionId);
		updateQuestion(questionId, { options: updatedOptions });
	}

	function updateOption(questionId: string, optionId: string, updatedOption: Partial<Option>) {
		const question = questions.find((q) => q.id === questionId);
		if (!question?.options) return;

		const updatedOptions = question.options.map(opt => {
			if (opt.id === optionId) {
				const newOption = { ...opt, ...updatedOption };
				if (updatedOption.label !== undefined && opt.error && updatedOption.label.trim() !== '') {
					newOption.error = undefined;
				}
				return newOption;
			}
			return opt;
		});
		updateQuestion(questionId, { options: updatedOptions });
	}

	function moveQuestionUp(index: number) {
		if (index > 0) {
			const q = questions[index];
			questions.splice(index, 1);
			questions.splice(index - 1, 0, q);
		}
	}

	function moveQuestionDown(index: number) {
		if (index < questions.length - 1) {
			const q = questions[index];
			questions.splice(index, 1);
			questions.splice(index + 1, 0, q);
		}
	}

	function updateFormData(data: Partial<FormData>) {
		Object.assign(formData, data);
	}

	// init
	loadForm();

	return {
		get formData() {
			return formData;
		},
		get questions() {
			return questions;
		},
		get isFormValid() {
			return isFormValid;
		},
		addQuestion,
		updateQuestion,
		removeQuestion,
		addOption,
		removeOption,
		updateOption,
		moveQuestionUp,
		moveQuestionDown,
		updateFormData,
		validateQuestions
	};
}

export type FormStore = ReturnType<typeof createFormStore>;
