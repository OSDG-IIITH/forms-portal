import { getTemplate } from '$lib/server/templates';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url }) => {
    const templateSlug = url.searchParams.get('template');
    let templateKdl = '';

    let templateTitle = '';
    let templateDescription = '';

    if (templateSlug) {
        const template = getTemplate(templateSlug);
        if (template) {
            templateKdl = template.kdl;
            templateTitle = template.title;
            templateDescription = template.description;
        }
    }

    return {
        templateKdl,
        templateTitle,
        templateDescription
    };
};
