<script lang="ts">
	import { Button, Separator } from "$lib/components/ui";
	import ChangeDisplayName from "./ChangeDisplayName.svelte";
	import ChangeProfilePicture from "./ChangeProfilePicture.svelte";
	import ApiToken from "./ApiToken.svelte";
	import NewApiTokenModal from "./NewApiTokenModal.svelte";
	import { Camera, KeyRound, Plus, QrCode, User } from "@lucide/svelte";
	import QuickCodeModal from "./QuickCodeModal.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";

	let { data } = $props();

	let openNewApiTokenModal = $state(false);
	let openQuickCodeModal = $state(false);
</script>

<div class="flex flex-col gap-6">
	<div
		class="flex flex-col gap-6 rounded-lg border bg-linear-to-b from-[oklch(0.93_0.045_75)] to-background p-4 shadow-sm sm:p-6 md:flex-row md:items-end md:gap-8 dark:from-[oklch(0.24_0.03_80)] dark:to-background"
	>
		<div class="flex min-w-0 flex-col gap-2">
			<p
				class="text-xs font-semibold uppercase tracking-wider text-muted-foreground"
			>
				Settings
			</p>

			<h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
				{data.userData.displayName}
			</h1>

			<p class="text-sm text-muted-foreground">
				Manage your profile, security, and API access.
			</p>
		</div>
	</div>

	<div class="flex flex-col gap-10">
		<section>
			<SectionHeader>
				<Camera />
				Profile Picture
			</SectionHeader>
			<p class="mb-3 px-2 text-sm text-muted-foreground">
				Change your profile picture.
			</p>
			<ChangeProfilePicture currentPicture={data.userData.picture.large} />
		</section>

		<section>
			<SectionHeader>
				<User />
				Display Name
			</SectionHeader>
			<p class="mb-3 px-2 text-sm text-muted-foreground">
				Change how your name appears across Tunebook.
			</p>
			<ChangeDisplayName />
		</section>

		<Separator />

		<section>
			<SectionHeader>
				<QrCode />
				Quick Code
				{#snippet actions()}
					<Button
						size="sm"
						variant="outline"
						onclick={() => {
							openQuickCodeModal = true;
						}}
					>
						<QrCode size={14} />
						Enter Code
					</Button>
				{/snippet}
			</SectionHeader>
			<p class="mb-3 px-2 text-sm text-muted-foreground">
				Claim a quick connect code to log in on another device.
			</p>
		</section>

		<section>
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
			<p class="mb-3 px-2 text-sm text-muted-foreground">
				Manage API tokens for programmatic access.
			</p>
			<div class="flex flex-col">
				{#each data.tokens as token (token.id)}
					<ApiToken {token} />
				{/each}
			</div>
		</section>
	</div>
</div>

<NewApiTokenModal bind:open={openNewApiTokenModal} />
<QuickCodeModal bind:open={openQuickCodeModal} />
