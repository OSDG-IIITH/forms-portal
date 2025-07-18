// types
export type QuestionType = 'input' | 'textarea' | 'radio' | 'checkbox' | 'file' | 'select' | 'date';

export type Option = {
  id: string;
  value: string;
  label: string;
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
    email?: boolean;
  };
  'max-file-size'?: number;
  'max-files'?: number;
  'allowed-types'?: string[];
};

export type FormData = {
  title: string;
  description: string;
  visibility: string;
  opens?: string;
  closes?: string;
  opensTime?: any;
  closesTime?: any;
  anonymous: boolean;
  max_responses?: number | null;
  individual_limit: number;
  editable_responses: boolean;
};

// util functions
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
      const typesString = q['allowed-types'].map(t => escapeKdlString(t)).join(' ');
      kdl += `    "allowed-types" ${typesString}\n`;
    }
  }
  if (q.validations && (q.validations['max-chars'] || q.validations['min-chars'] || q.validations.regex || q.validations.email)) {
    kdl += '    validations ';
    if(q.validations.email) {
      kdl += 'email ';
    }
    kdl +='{\n';
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
