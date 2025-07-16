<script lang="ts">
	import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs';
	import FormItem from './form-item.svelte';
	import FormsGridView from './forms-grid-view.svelte';
	import FormsListView from './forms-list-view.svelte';
	import { IconSearch, IconFilter, IconLayoutGrid, IconList, IconChevronLeft, IconChevronRight, IconClock, IconEdit, IconFileFilled, IconUsers } from '@tabler/icons-svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Pagination } from 'bits-ui';
	import recentForms from '$lib/data/recent-forms.json';
	import { onMount } from 'svelte';

	const USE_LOCAL_FORMS_DATA = false;

	let forms = $state<typeof recentForms>([]);
	let totalCount = $state(0);
	let userHandle = $state<string>('nouser');
	let _loading = false;
	let _error = false;
	let searchQuery = $state('');
	let activeTab = $state('recent');
	let viewMode = $state('grid');
	let currentPage = $state(1);
	let searchTimeout: ReturnType<typeof setTimeout>;

	const loading = $derived(_loading);
	const error = $derived(_error);
	const perPage = 12;

	async function fetchForms(page = 1) {
		_loading = true;
		_error = false;
		currentPage = page;
		if (USE_LOCAL_FORMS_DATA) {
			forms = recentForms;
			totalCount = recentForms.length;
			_loading = false;
			return;
		}
		try {
			const params = new URLSearchParams({
				limit: perPage.toString(),
				offset: ((page - 1) * perPage).toString(),
			});
			const res = await fetch(`/api/forms?${params.toString()}`, { credentials: 'include' });
			if (!res.ok) throw new Error();
			const data = await res.json();
			if (data && Array.isArray(data.data) && data.pagination && typeof data.pagination.total === 'number') {
				forms = data.data;
				totalCount = data.pagination.total;
			} else if (Array.isArray(data)) {
				forms = data;
				totalCount = data.length;
			} else {
				throw new Error();
			}
			_loading = false;
		} catch {
			_error = true;
			forms = [];
			totalCount = 0;
			_loading = false;
		}
	}

	function getFilteredForms() {
		const q = searchQuery.trim().toLowerCase();
		if (!q) return forms;
		return forms.filter(f => f.title && f.title.toLowerCase().includes(q));
	}

	$effect(() => {
		fetchForms(currentPage);
	});

	$effect(() => {
		if (searchTimeout) clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			fetchForms(1);
		}, 300);
	});

	function handlePageChange(page: number) {
		fetchForms(page);
	}

	async function fetchCurrentUser() {
		try {
			const res = await fetch('/api/auth/info', { credentials: 'include' });
			if (res.ok) {
				const user = await res.json();
				if (user?.handle) {
					userHandle = user.handle;
				}
			}
		} catch (e) {
			console.warn('Failed to fetch user info, using `nouser` handle');
		}
	}

	onMount(() => {
		const storedViewMode = localStorage.getItem('forms-view-mode');
		if (storedViewMode === 'grid' || storedViewMode === 'list') viewMode = storedViewMode;
		fetchCurrentUser();
		fetchForms(1);
	});

	$effect(() => {
		if (viewMode === 'grid' || viewMode === 'list') localStorage.setItem('forms-view-mode', viewMode);
	});
</script>

<Tabs bind:value={activeTab}>
	<div class="space-y-6">
		<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
			<TabsList class="grid w-full sm:w-fit grid-cols-4 gap-1 order-1 sm:order-none">
				<TabsTrigger value="recent" class="flex flex-col items-center justify-center px-0 sm:px-6 py-2">
					<span class="sm:hidden"><IconClock class="h-5 w-5" aria-label="Recent" /></span>
					<span class="hidden sm:inline">Recent</span>
				</TabsTrigger>
				<TabsTrigger value="created" class="flex flex-col items-center justify-center px-0 sm:px-6 py-2">
					<span class="sm:hidden"><IconEdit class="h-5 w-5" aria-label="Created" /></span>
					<span class="hidden sm:inline">Created</span>
				</TabsTrigger>
				<TabsTrigger value="filled" class="flex flex-col items-center justify-center px-0 sm:px-6 py-2">
					<span class="sm:hidden"><IconFileFilled class="h-5 w-5" aria-label="Filled" /></span>
					<span class="hidden sm:inline">Filled</span>
				</TabsTrigger>
				<TabsTrigger value="shared" class="flex flex-col items-center justify-center px-0 sm:px-6 py-2">
					<span class="sm:hidden"><IconUsers class="h-5 w-5" aria-label="Shared" /></span>
					<span class="hidden sm:inline">Shared</span>
				</TabsTrigger>
			</TabsList>
			<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-2 order-2 sm:order-none w-full sm:w-auto">
				<div class="flex flex-row gap-2 w-full">
					<div class="relative flex-1">
						<IconSearch class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
						<Input
							placeholder="Search forms"
							class="pl-10 w-full"
							bind:value={searchQuery}
						/>
					</div>
					<Button variant="outline" size="sm" class="flex-shrink-0 h-9 w-9 p-0">
						<IconFilter class="h-4 w-4" />
					</Button>
					<div class="flex border rounded-md bg-background flex-shrink-0">
						<Button
							variant={viewMode === 'grid' ? 'default' : 'ghost'}
							size="sm"
							class="rounded-r-none border-0 px-3"
							onclick={() => viewMode = 'grid'}
						>
							<IconLayoutGrid class="h-4 w-4" />
						</Button>
						<Button
							variant={viewMode === 'list' ? 'default' : 'ghost'}
							size="sm"
							class="rounded-l-none border-0 px-3"
							onclick={() => viewMode = 'list'}
						>
							<IconList class="h-4 w-4" />
						</Button>
					</div>
				</div>
			</div>
		</div>

		<TabsContent value="recent" class="animate-in fade-in-0 slide-in-from-bottom-1 duration-300">
			{#if loading}
				{#if viewMode === 'grid'}
					<FormsGridView forms={[]} {loading} {userHandle} />
				{:else}
					<FormsListView forms={[]} {loading} {userHandle} />
				{/if}
			{:else if error}
				<div class="text-center py-8 text-muted-foreground">Failed to load forms</div>
			{:else if forms.length}
				<div class="space-y-6">
					{#if viewMode === 'grid'}
						<FormsGridView forms={getFilteredForms()} {loading} {userHandle} />
					{:else}
						<FormsListView forms={getFilteredForms()} {loading} {userHandle} />
					{/if}
					{#if totalCount > perPage || forms.length === perPage}
						<Pagination.Root count={totalCount} {perPage} bind:page={currentPage} onPageChange={handlePageChange}>
							{#snippet children({ pages, range })}
								<div class="flex items-center justify-center">
									<Pagination.PrevButton class="mr-4 inline-flex h-8 w-8 items-center justify-center rounded-md bg-transparent hover:bg-accent disabled:cursor-not-allowed disabled:opacity-50">
										<IconChevronLeft class="h-4 w-4" />
									</Pagination.PrevButton>
									<div class="flex items-center gap-1">
										{#each pages as page (page.key)}
											{#if page.type === "ellipsis"}
												<div class="flex h-8 w-8 items-center justify-center text-sm text-muted-foreground">...</div>
											{:else}
												<Pagination.Page
													{page}
													class="inline-flex h-8 w-8 items-center justify-center rounded-md bg-transparent text-sm font-medium hover:bg-accent data-selected:bg-primary data-selected:text-primary-foreground"
												>
													{page.value}
												</Pagination.Page>
											{/if}
										{/each}
									</div>
									<Pagination.NextButton class="ml-4 inline-flex h-8 w-8 items-center justify-center rounded-md bg-transparent hover:bg-accent disabled:cursor-not-allowed disabled:opacity-50">
										<IconChevronRight class="h-4 w-4" />
									</Pagination.NextButton>
								</div>
								<p class="text-center text-xs text-muted-foreground mt-2">
									Showing {range.start} - {range.end} of {totalCount}
								</p>
							{/snippet}
						</Pagination.Root>
					{/if}
				</div>
			{:else}
				<div class="text-center py-8 text-muted-foreground">No forms found</div>
			{/if}
		</TabsContent>

		<TabsContent value="created" class="animate-in fade-in-0 slide-in-from-bottom-1 duration-300">
			{#if loading}
				{#if viewMode === 'grid'}
					<FormsGridView forms={[]} {loading} {userHandle} />
				{:else}
					<FormsListView forms={[]} {loading} {userHandle} />
				{/if}
			{:else if error}
				<div class="text-center py-8 text-muted-foreground">Failed to load forms</div>
			{:else if forms.length}
				<div class="space-y-6">
					{#if viewMode === 'grid'}
						<FormsGridView forms={getFilteredForms()} {loading} {userHandle} />
					{:else}
						<FormsListView forms={getFilteredForms()} {loading} {userHandle} />
					{/if}
					{#if totalCount > perPage || forms.length === perPage}
						<Pagination.Root count={totalCount} {perPage} bind:page={currentPage} onPageChange={handlePageChange}>
							{#snippet children({ pages, range })}
								<div class="flex items-center justify-center">
									<Pagination.PrevButton class="mr-4 inline-flex h-8 w-8 items-center justify-center rounded-md bg-transparent hover:bg-accent disabled:cursor-not-allowed disabled:opacity-50">
										<IconChevronLeft class="h-4 w-4" />
									</Pagination.PrevButton>
									<div class="flex items-center gap-1">
										{#each pages as page (page.key)}
											{#if page.type === "ellipsis"}
												<div class="flex h-8 w-8 items-center justify-center text-sm text-muted-foreground">...</div>
											{:else}
												<Pagination.Page
													{page}
													class="inline-flex h-8 w-8 items-center justify-center rounded-md bg-transparent text-sm font-medium hover:bg-accent data-selected:bg-primary data-selected:text-primary-foreground"
												>
													{page.value}
												</Pagination.Page>
											{/if}
										{/each}
									</div>
									<Pagination.NextButton class="ml-4 inline-flex h-8 w-8 items-center justify-center rounded-md bg-transparent hover:bg-accent disabled:cursor-not-allowed disabled:opacity-50">
										<IconChevronRight class="h-4 w-4" />
									</Pagination.NextButton>
								</div>
								<p class="text-center text-xs text-muted-foreground mt-2">
									Showing {range.start} - {range.end} of {totalCount}
								</p>
							{/snippet}
						</Pagination.Root>
					{/if}
				</div>
			{:else}
				<div class="text-center py-8 text-muted-foreground">No forms found</div>
			{/if}
		</TabsContent>

		<TabsContent value="filled" class="animate-in fade-in-0 slide-in-from-bottom-1 duration-300">
			<FormsGridView forms={[]} loading={true} {userHandle} animationDelay={0} />
		</TabsContent>

		<TabsContent value="shared" class="animate-in fade-in-0 slide-in-from-bottom-1 duration-300">
			<div class="text-center py-8 text-muted-foreground">No shared forms</div>
		</TabsContent>
	</div>
</Tabs>
