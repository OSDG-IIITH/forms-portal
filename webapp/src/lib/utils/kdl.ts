import * as kdljs from 'kdljs';
import { ulid } from 'ulid';
import { type Question, type QuestionType, type FormConfig, QuestionSchema, FormConfigSchema } from '$lib/types/form';

export function parseKdlValue(node: any): string {
    if (node === null || node === undefined) return '';
    if (typeof node === 'string') return node;
    if (typeof node.value === 'string') return node.value;
    if (typeof node.value === 'number' || typeof node.value === 'boolean') return String(node.value);
    return String(node);
}

export function safeString(val: any): string {
    if (val === null || val === undefined) return '';
    if (typeof val === 'string') return val;
    if (typeof val === 'number' || typeof val === 'boolean') return String(val);
    return '';
}

export function escapeKdlString(str: string): string {
    return '"' + str.replace(/\\/g, '\\\\').replace(/"/g, '\\"') + '"';
}

export function questionToKdl(q: Question): string {
    let kdl = `  question id=${escapeKdlString(q.id)} type=${escapeKdlString(q.type)}${q.required ? ' required' : ''} {\n`;
    kdl += `    title ${escapeKdlString(q.title)}\n`;
    if (q.description) {
        kdl += `    description ${escapeKdlString(q.description)}\n`;
    }
    if (q.placeholder) {
        kdl += `    placeholder ${escapeKdlString(q.placeholder)}\n`;
    }
    if (q.options && q.options.length > 0) {
        for (const opt of q.options) {
            kdl += `    option value=${escapeKdlString(opt.value)} label=${escapeKdlString(opt.label)}\n`;
        }
    }
    if (q['max-file-size'] || q['max-files'] || q['allowed-types']?.length) {
        if (q['max-file-size']) {
            kdl += `    "max-file-size" ${q['max-file-size']} mb\n`;
        }
        if (q['max-files'] && q['max-files'] !== -1) {
            kdl += `    "max-files" ${q['max-files']}\n`;
        }
        if (q['allowed-types'] && q['allowed-types'].length > 0) {
            const typesString = q['allowed-types'].map((t) => escapeKdlString(t)).join(' ');
            kdl += `    "allowed-types" ${typesString}\n`;
        }
    }
    if (
        q.validations &&
        (q.validations['max-chars'] || q.validations['min-chars'] || q.validations.regex)
    ) {
        kdl += '    validations {\n';
        if (q.validations.regex) {
            kdl += `      regex ${escapeKdlString(q.validations.regex)}\n`;
        }
        if (q.validations['min-chars']) {
            kdl += `      "min-chars" ${q.validations['min-chars']}\n`;
        }
        if (q.validations['max-chars']) {
            kdl += `      "max-chars" ${q.validations['max-chars']}\n`;
        }
        kdl += '    }\n';
    }
    kdl += `  }\n`;
    return kdl;
}

export function generateFormKdl(questions: Question[]): string {
    let kdl = 'form {\n';
    kdl += `  version 1\n`;
    for (const q of questions) {
        kdl += questionToKdl(q);
    }
    kdl += '}\n';
    return kdl;
}

export function parseKdlForm(kdl: string): { config: FormConfig; questions: Question[] } {
    const ast = kdljs.parse(kdl);
    if (!ast || !ast.output || !ast.output.length) throw new Error('Invalid KDL');
    const formNode = ast.output.find((n: any) => n.name === 'form');
    if (!formNode) throw new Error('No form node');

    const config: FormConfig = {
        title: '',
        description: '',
        visibility: 'public'
    };
    const questions: Question[] = [];

    for (const child of formNode.children) {
        if (child.name === 'title') config.title = parseKdlValue(child.values[0]);
        else if (child.name === 'description') config.description = parseKdlValue(child.values[0]);
        else if (child.name === 'visibility') config.visibility = parseKdlValue(child.values[0]);
        else if (child.name === 'question') {
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
                        const optionId = id || (value && !q.options.some((opt) => opt.id === value) ? value : ulid());
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
            questions.push(QuestionSchema.parse(q));
        }
    }
    FormConfigSchema.parse(config);
    return { config, questions };
}
