<script lang="ts">
	import { Github, Mail } from '@lucide/svelte';
</script>

<footer class="border-t w-full mt-auto bg-background">
	<div class="container mx-auto flex h-16 items-center justify-between px-4 md:px-6">
		<div class="flex items-center gap-2">
			<p class="text-sm text-muted-foreground">
				Made with ♥︎ at <a href="https://osdg.iiit.ac.in" target="_blank" class="hover:underline">OSDG</a>
			</p>
		</div>
		
		<div class="flex items-center gap-4">
			<div class="flex gap-6">
				<a
					href="https://github.eadcom/osdg-iiith/forms-portal"
					target="_blank"
					rel="noreferrer"
					class="text-muted-foreground hover:text-foreground transition-colors"
					aria-label="GitHub"
				>
					<Github class="h-5 w-5" />
				</a>
				<a
                    href="mailto:osdg@students.iiit.ac.in?subject=FormsIIIT"
					class="text-muted-foreground hover:text-foreground transition-colors"
					aria-label="Email"
				>
					<Mail class="h-5 w-5" />
				</a>
			</div>
		</div>
	</div>
</footer>
