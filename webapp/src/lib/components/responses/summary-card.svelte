<script lang="ts">
    import { Card } from "$lib/components/ui/card";
    import { Files, Users, Clock } from "@lucide/svelte";
    export let totalResponses: number;
    export let uniqueRespondents: number;
    export let lastResponse: string;
    export let opens: string;
    export let closes: string;

    let status = "";
    let statusColor = "";
    const now = new Date();
    if (opens && closes) {
        const openDate = new Date(opens);
        const closeDate = new Date(closes);
        if (now < openDate) {
            status = "Scheduled";
            statusColor = "bg-muted text-muted-foreground";
        } else if (now > closeDate) {
            status = "Closed";
            statusColor = "bg-destructive/10 text-destructive";
        } else {
            status = "Live";
            statusColor = "bg-primary/10 text-primary";
        }
    } else if (closes) {
        const closeDate = new Date(closes);
        if (now > closeDate) {
            status = "Closed";
            statusColor = "bg-destructive/10 text-destructive";
        } else {
            status = "Live";
            statusColor = "bg-primary/10 text-primary";
        }
    } else {
        status = "Live";
        statusColor = "bg-primary/10 text-primary";
    }

    console.log(lastResponse);

    function relativeTime(dateString: string): string {
        const date = new Date(dateString);
        if (isNaN(date.getTime())) return dateString;
        const now = new Date();
        const diff = (now.getTime() - date.getTime()) / 1000;
        if (diff < 60) return `${Math.floor(diff)}s`;
        if (diff < 3600) return `${Math.floor(diff / 60)}min`;
        if (diff < 86400) return `${Math.floor(diff / 3600)}hr`;
        if (diff < 172800) return `1 day`;
        return `${Math.floor(diff / 86400)} days`;
    }
</script>

<div class="flex items-center justify-between mb-3">
    <h2 class="text-2xl font-bold tracking-tight text-foreground">Responses</h2>
    <span class={`rounded-full px-3 py-1 text-xs font-medium ${statusColor}`}
        >{status}</span
    >
</div>
<Card class="p-0 overflow-hidden md:py-3">
    <div class="flex flex-col gap-0">
        <div
            class="flex flex-row md:flex-col lg:flex-row items-center justify-center gap-2 px-2 py-2"
        >
            <div class="flex-1 flex flex-col items-center justify-center py-4">
                <div class="flex items-center gap-2">
                    <Files class="w-5 h-5 text-primary" />
                    <span class="text-xl md:text-4xl font-bold text-primary"
                        >{totalResponses}</span
                    >
                </div>
                <span class="text-xs text-muted-foreground mt-1"
                    ><span class="hidden md:inline">Total</span> Responses</span
                >
            </div>
            <div class="flex-1 flex flex-col items-center justify-center py-4">
                <div class="flex items-center gap-2">
                    <Users class="w-5 h-5 text-primary" />
                    <span class="text-xl md:text-4xl font-bold text-primary"
                        >{uniqueRespondents}</span
                    >
                </div>
                <span class="text-xs text-muted-foreground mt-1">
                    <span class="hidden md:inline">Unique</span> Respondents
                </span>
            </div>
            <div class="flex-1 flex flex-col items-center justify-center py-4">
                <div class="flex items-center gap-2">
                    <Clock class="w-5 h-5 text-primary" />
                    <span class="text-xl md:text-4xl font-bold text-primary">
                        {relativeTime(lastResponse)}
                    </span>
                </div>
                <span class="text-xs text-muted-foreground mt-1"
                    >Last Response</span
                >
            </div>
        </div>
    </div>
</Card>