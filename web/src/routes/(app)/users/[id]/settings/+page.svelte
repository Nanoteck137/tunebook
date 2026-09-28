<script lang="ts">
	import { Button, Separator } from "$lib/components/ui";
	import ChangeDisplayName from "./ChangeDisplayName.svelte";
	import ChangeProfilePicture from "./ChangeProfilePicture.svelte";
	import ApiToken from "./ApiToken.svelte";
	import NewApiTokenModal from "./NewApiTokenModal.svelte";
	import QuickCode from "./QuickCode.svelte";
	import { KeyRound, Plus, Settings, ShieldCheck, User } from "@lucide/svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";

	let { data } = $props();

	let openNewApiTokenModal = $state(false);
</script>

<div class="flex flex-col gap-8">
	<SectionHeader>
		<Settings />
		Settings
	</SectionHeader>

	<section class="flex flex-col gap-4">
		<SectionHeader>
			<User />
			Profile
		</SectionHeader>

		<div class="flex flex-col gap-6 rounded-lg border p-4">
			<ChangeDisplayName currentName={data.userData.displayName} />

			<Separator />

			<ChangeProfilePicture currentPicture={data.userData.picture.large} />
		</div>
	</section>

	<section class="flex flex-col gap-4">
		<SectionHeader>
			<ShieldCheck />
			Security
		</SectionHeader>

		<div class="rounded-lg border p-4">
			<QuickCode />
		</div>
	</section>

	<section class="flex flex-col gap-4">
		<SectionHeader count={data.tokens.length}>
			<KeyRound />
			API Tokens
			{#snippet actions()}
				<Button
					size="sm"
					variant="outline"
					onclick={() => {
						openNewApiTokenModal = true;
					}}
				>
					<Plus size={14} />
					New Token
				</Button>
			{/snippet}
		</SectionHeader>

		{#if data.tokens.length === 0}
			<div class="flex flex-col items-center gap-2 rounded-lg border py-16">
				<KeyRound size={32} class="text-muted-foreground/40" />
				<p class="text-sm text-muted-foreground">
					No API tokens yet. Create one to access Tunebook programmatically.
				</p>
			</div>
		{:else}
			<div class="flex flex-col divide-y rounded-lg border px-4">
				{#each data.tokens as token (token.id)}
					<ApiToken {token} />
				{/each}
			</div>
		{/if}
	</section>
</div>

<NewApiTokenModal bind:open={openNewApiTokenModal} />
