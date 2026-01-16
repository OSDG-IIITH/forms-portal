export interface DashboardForm {
    id: number | string;
    title: string;
    slug?: string;
    owner?: string;
    createdAt: string;
    modified?: string;
    status?: string;
    responses?: number;
    anonymous?: boolean;
    editable_responses?: boolean;
    _deleted?: boolean;
}
