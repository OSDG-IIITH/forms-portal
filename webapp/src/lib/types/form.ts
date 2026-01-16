import type { Time } from '@internationalized/date';

export type QuestionType =
	| 'input'
	| 'textarea'
	| 'radio'
	| 'checkbox'
	| 'file'
	| 'select'
	| 'date'
	| 'section-header';

export interface Option {
	id: string;
	value: string;
	label: string;
	error?: string;
}

export interface Question {
	id: string;
	type: QuestionType;
	title: string;
	required: boolean;
	description?: string;
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
}

export interface FormData {
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
}

export interface FormConfig {
  title: string;
  description: string;
  visibility: string;
}
