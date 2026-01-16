import { z } from 'zod';
import type { Question } from '$lib/types/form';

export function createResponseSchema(question: Question): z.ZodType<any> {
    let schema = z.string();

    if (question.required) {
        schema = schema.min(1, 'This field is required.');
    } else {
        // If not required, empty string is valid (and effectively "optional")
        // But we want to allow empty string without triggering other validations if it's not required?
        // Actually, if it's not required, an empty string should be valid.
        // If it HAS content, then min-chars/regex should apply.
        // Zod's .or(z.literal('')) works for this pattern usually, or .optional() if we treat empty as undefined, 
        // but here our model is likely string.

        // Let's stick to the base schema and handle the "optional" logic by making sure validations only run if value is present,
        // OR by using .or(z.literal('')) at the end if not required. 
        // However, in the input component, we usually bind to a string.
    }

    // Refinement for validations that should only apply if the string is NOT empty (unless min-chars is used to enforce non-empty)
    // Actually, if it is NOT required, we probably want to allow empty string.
    // If it IS required, we already added .min(1).

    // Let's build the chain.

    if (question.validations) {
        if (question.validations['min-chars'] !== undefined) {
            schema = schema.min(question.validations['min-chars'], `Minimum ${question.validations['min-chars']} characters required.`);
        }

        if (question.validations['max-chars'] !== undefined) {
            schema = schema.max(question.validations['max-chars'], `Maximum ${question.validations['max-chars']} characters allowed.`);
        }

        if (question.validations.regex) {
            try {
                // Validate that the regex is valid
                new RegExp(question.validations.regex);
                schema = schema.regex(new RegExp(question.validations.regex), 'Invalid format.');
            } catch (e) {
                // Ignore invalid regex definitions safely or log warn
                console.warn(`Invalid regex for question ${question.id}: ${question.validations.regex}`);
            }
        }
    }

    if (!question.required) {
        // If not required, allow empty string to pass even if it violates min/regex (common pattern: empty OR valid)
        // correct way in Zod for "string that is valid OR empty":
        return schema.or(z.literal(''));
    }

    return schema;
}
