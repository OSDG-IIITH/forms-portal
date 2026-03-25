import { z } from 'zod';
import type { Time } from '@internationalized/date';

export const QuestionTypeSchema = z.enum([
	'input',
	'textarea',
	'radio',
	'checkbox',
	'file',
	'select',
	'date',
	'section-header'
]);

export type QuestionType = z.infer<typeof QuestionTypeSchema>;

export const OptionSchema = z.object({
	id: z.string(),
	value: z.string(),
	label: z.string().min(1, "Option text cannot be empty."),
	error: z.string().optional()
});

export type Option = z.infer<typeof OptionSchema>;

export const QuestionSchema = z.object({
	id: z.string(),
	type: QuestionTypeSchema,
	title: z.string().min(1, "Question text cannot be empty."),
	required: z.boolean(),
	description: z.string().optional(),
	options: z.array(OptionSchema).optional(),
	placeholder: z.string().optional(),
	validations: z
		.object({
			'max-chars': z.number().optional(),
			'min-chars': z.number().optional(),
			regex: z.string().optional()
		})
		.optional(),
	'max-file-size': z.number().optional(),
	'max-files': z.number().optional(),
	'allowed-types': z.array(z.string()).optional(),
	error: z.string().optional()
}).superRefine((data, ctx) => {
	if (data.options && (data.type === 'radio' || data.type === 'checkbox' || data.type === 'select')) {
		const seenLabels = new Map<string, number[]>();

		data.options.forEach((opt, idx) => {
			const label = opt.label.trim().toLowerCase();
			if (label) {
				if (!seenLabels.has(label)) {
					seenLabels.set(label, []);
				}
				seenLabels.get(label)!.push(idx);
			}
		});

		seenLabels.forEach((indices) => {
			if (indices.length > 1) {
				indices.forEach((idx) => {
					ctx.addIssue({
						code: z.ZodIssueCode.custom,
						message: "Duplicate option found.",
						path: ['options', idx, 'error'] // We want to target the error field of the specific option
					});
				});
			}
		});
	}
});

export type Question = z.infer<typeof QuestionSchema>;

export const FormDataSchema = z.object({
	title: z.string(),
	description: z.string(),
	visibility: z.string(),
	opens: z.string().optional(),
	closes: z.string().optional(),
	opensTime: z.custom<Time>((val) => val !== undefined && val !== null, "Invalid Time").optional(), // Zod doesn't hold class instances natively well without custom, but loose check is fine here as it's runtime arg usually
	closesTime: z.custom<Time>((val) => val !== undefined && val !== null, "Invalid Time").optional(),
	anonymous: z.boolean(),
	max_responses: z.number().nullable().optional(),
	individual_limit: z.number(),
	editable_responses: z.boolean()
});

export type FormData = z.infer<typeof FormDataSchema>;

export type EditorForm = {
	id?: string;
	title?: string;
	description?: string | null;
	structure?: string;
	opens?: string | null;
	closes?: string | null;
	anonymous?: boolean | null;
	max_responses?: number | null;
	individual_limit?: number | null;
	editable_responses?: boolean | null;
};

export const FormConfigSchema = z.object({
	title: z.string(),
	description: z.string(),
	visibility: z.string()
});

export type FormConfig = z.infer<typeof FormConfigSchema>;
