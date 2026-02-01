import { parse } from 'kdljs';

export interface Template {
    slug: string;
    title: string;
    description: string;
    kdl: string;
}

export function getTemplates(): Template[] {
    const modules = import.meta.glob('$lib/data/templates/*.kdl', { query: '?raw', eager: true });
    const templates: Template[] = [];

    for (const path in modules) {
        const kdl = (modules[path] as { default: string }).default;
        const slug = path.split('/').pop()?.replace('.kdl', '') || '';

        let title = slug;
        let description = '';

        try {
            const doc = parse(kdl) as any;
            const formNode = doc?.output?.find((node: any) => node.name === 'form');

            if (formNode?.children) {
                const titleNode = formNode.children.find((node: any) => node.name === 'title');
                const descNode = formNode.children.find((node: any) => node.name === 'description');

                if (titleNode?.values?.length > 0) {
                    title = String(titleNode.values[0]);
                }

                if (descNode?.values?.length > 0) {
                    description = String(descNode.values[0]);
                }
            }
        } catch (e) {
            console.error(`Failed to parse template: ${path}`, e);
        }

        templates.push({
            slug,
            title,
            description,
            kdl
        });
    }

    return templates.sort((a, b) => a.slug.localeCompare(b.slug));
}

export function getTemplate(slug: string): Template | undefined {
    const templates = getTemplates();
    return templates.find((t) => t.slug === slug);
}
