import { getTemplates } from '$lib/server/templates';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
    return {
        templates: getTemplates()
    };
};
