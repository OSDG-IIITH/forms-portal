export type Group = {
    id: string;
    owner: string;
    name: string;
    description?: string;
    type: 'domain' | 'list';
    domain?: string;
    members?: string[];
};
