import { ulid } from 'ulid';
import * as kdljs from 'kdljs';
import { Time } from '@internationalized/date';
import { parseKdlValue, safeString } from './FormEditor';

// type definitions
export type QuestionType = 'input' | 'textarea' | 'radio' | 'checkbox' | 'file' | 'select' | 'date';

export type Option = {
	id: string;
	value: string;
	label: string;
	error?: string;
};

export type Question = {
	id: string;
	type: QuestionType;
	title: string;
	required: boolean;
	options?: Option[];
	placeholder?: string;
	validations?: {
		'max-chars'?: number;
		'min-chars'?: number;
		regex?: string;
	};
	'max-file-size'?: number;
	'max-files'?: number;
	'allowed-types'?: string[];
	error?: string;
};

export type FormData = {
	title: string;
	description: string;
	visibility: string;
	opens?: string;
	closes?: string;
	opensTime?: Time;
	closesTime?: Time;
	anonymous: boolean;
	max_responses?: number | null;
	individual_limit: number;
	editable_responses: boolean;
};

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

	function parseKdlForm(kdl: string): Question[] {
		const ast = kdljs.parse(kdl);
		if (!ast || !ast.output || !ast.output.length) throw new Error('Invalid KDL');
		const formNode = ast.output.find((n: any) => n.name === 'form');
		if (!formNode) throw new Error('No form node');
		const qs: Question[] = [];
		for (const child of formNode.children) {
			if (child.name === 'question') {
				const q: Question = {
					id: '',
					type: 'input',
					title: '',
					required: false
				};
				if (child.properties && typeof child.properties === 'object') {
					for (const key in child.properties) {
						const value = child.properties[key];
						if (key === 'id') q.id = safeString(value);
						else if (key === 'type') {
							let t = safeString(value);
							if (t === 'multiple_choice') t = 'radio';
							if (t === 'text') t = 'input';
							if (t === 'textarea') t = 'textarea';
							q.type = t as QuestionType;
						}
					}
				}
				if (Array.isArray(child.values) && child.values.includes('required')) {
					q.required = true;
				}
				if (child.children) {
					for (const c of child.children) {
						if (c.name === 'title') q.title = parseKdlValue(c.values[0]);
						else if (c.name === 'placeholder') q.placeholder = parseKdlValue(c.values[0]);
						else if (c.name === 'max-file-size') {
							q['max-file-size'] = Number(parseKdlValue(c.values[0]));
						} else if (c.name === 'max-files') {
							q['max-files'] = Number(parseKdlValue(c.values[0]));
						} else if (c.name === 'allowed-types') {
							q['allowed-types'] = c.values?.map((v: any) => parseKdlValue(v)) || [];
						} else if (c.name === 'option') {
							if (!q.options) q.options = [];
							let id = '',
								value = '',
								label = '';
							if (c.properties && typeof c.properties === 'object') {
								for (const key in c.properties) {
									const v = c.properties[key];
									if (key === 'value') value = safeString(v);
									else if (key === 'label') label = safeString(v);
									else if (key === 'id') id = safeString(v);
								}
							}
							const optionId =
								id || (value && !q.options.some((opt) => opt.id === value) ? value : ulid());
							q.options.push({ id: optionId, value, label });
						} else if (c.name === 'validations') {
							q.validations = {};
							for (const v of c.children || []) {
								if (v.name === 'regex') q.validations.regex = parseKdlValue(v.values[0]);
								else if (v.name === 'min-chars')
									q.validations['min-chars'] = Number(parseKdlValue(v.values[0]));
								else if (v.name === 'max-chars')
									q.validations['max-chars'] = Number(parseKdlValue(v.values[0]));
							}
						}
					}
				}
				qs.push(q);
			}
		}
		return qs;
	}

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
					questions = parseKdlForm(initialForm.structure);
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
