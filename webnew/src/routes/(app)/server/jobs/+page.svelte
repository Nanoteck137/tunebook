<script lang="ts">
	import { getApiClient, handleApiError } from "$lib";
	import { cn } from "$lib/utils";
	import { Badge, Button } from "$lib/components/ui";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { ListChecks, Play, RefreshCw } from "@lucide/svelte";
	import { onMount } from "svelte";
	import { toast } from "svelte-sonner";
	import type { JobTy, TaskSyncStateEventTaskTy } from "../events";
	import {
		connectServerSse,
		formatJobDate,
		JobSyncStateEvent,
		jobStatusVariant,
		TaskSyncStateEvent,
	} from "../events";
	import Spacer from "$lib/components/Spacer.svelte";

	const { data } = $props();
	const apiClient = getApiClient();

	let jobs = $state<JobTy[]>(data.jobs);
	let tasks = $state<TaskSyncStateEventTaskTy[]>([]);

	let filter = $state("all");

	const filters = [
		{ label: "All", value: "all" },
		{ label: "Pending", value: "pending" },
		{ label: "Running", value: "running" },
		{ label: "Completed", value: "completed" },
		{ label: "Failed", value: "failed" },
	];

	const counts = $derived({
		all: jobs.length,
		pending: jobs.filter((j) => j.status === "pending").length,
		running: jobs.filter((j) => j.status === "running").length,
		completed: jobs.filter((j) => j.status === "completed").length,
		failed: jobs.filter((j) => j.status === "failed").length,
	});

	const filteredJobs = $derived(
		filter === "all" ? jobs : jobs.filter((j) => j.status === filter),
	);

	let eventSource = $state<EventSource | null>(null);

	onMount(() => {
		connectServerSse({
			"task-sync-state": (d) => {
				tasks = TaskSyncStateEvent.parse(d).tasks;
			},
			"job-sync-state": (d) => {
				jobs = JobSyncStateEvent.parse(d).jobs;
			},
		}).then((e) => {
			eventSource = e;
		});

		return () => {
			eventSource?.close();
		};
	});
</script>

<div class="flex flex-col gap-4">
	<section>
		<SectionHeader>
			<RefreshCw />
			Tasks
		</SectionHeader>

		<Spacer />

		<div class="flex flex-col divide-y rounded-lg border bg-card px-4">
			{#each tasks as task (task.name)}
				<div class="flex items-center justify-between gap-4 py-3">
					<div class="flex min-w-0 flex-col">
						<span class="text-sm font-medium">{task.displayName}</span>
						<span class="text-xs text-muted-foreground">
							{task.isRunning ? "Running..." : "Idle"}
							{#if task.schedule}
								&middot; {task.schedule}
							{/if}
						</span>
					</div>
					{#if !task.isRunning}
						<Button
							variant="outline"
							size="sm"
							onclick={async () => {
								const res = await apiClient.runTask(task.name);
								if (!res.success) {
									return handleApiError(res.error);
								}

								toast.success("Dispatched task");
							}}
						>
							<Play size={14} />
							Run
						</Button>
					{/if}
				</div>
			{/each}
		</div>
	</section>

	<section>
		<SectionHeader>
			<ListChecks />
			Jobs
		</SectionHeader>

		<Spacer />

		<div class="flex flex-col gap-4 rounded-lg border bg-card p-4">
			<div class="flex flex-wrap items-center gap-2">
				{#each filters as f (f.value)}
					<Button
						variant="ghost"
						size="sm"
						onclick={() => (filter = f.value)}
						class={cn(
							"text-foreground hover:bg-accent hover:text-accent-foreground",
							filter === f.value
								? "bg-accent text-accent-foreground"
								: "text-muted-foreground",
						)}
					>
						{f.label}
						<span class="text-xs opacity-70"
							>{counts[f.value as keyof typeof counts]}</span
						>
					</Button>
				{/each}
			</div>

			<div class="overflow-x-auto">
				<table class="w-full text-left text-sm">
					<thead>
						<tr class="text-muted-foreground">
							<th class="px-3 py-2 font-medium">Name</th>
							<th class="px-3 py-2 font-medium">Status</th>
							<th class="px-3 py-2 font-medium">Attempts</th>
							<th class="px-3 py-2 font-medium">Created</th>
							<th class="px-3 py-2 font-medium">Updated</th>
						</tr>
					</thead>
					<tbody>
						{#if filteredJobs.length === 0}
							<tr>
								<td
									colspan={5}
									class="px-3 py-6 text-center text-sm text-muted-foreground"
								>
									No {filter === "all" ? "jobs" : `${filter} jobs`}.
								</td>
							</tr>
						{:else}
							{#each filteredJobs as job (job.id)}
								<tr class="border-t">
									<td class="px-3 py-2">{job.displayName}</td>
									<td class="px-3 py-2">
										<Badge variant={jobStatusVariant(job.status)}>
											{job.status}
										</Badge>
									</td>
									<td class="px-3 py-2">{job.attempts}/{job.maxAttempts}</td>
									<td class="px-3 py-2">{formatJobDate(job.created)}</td>
									<td class="px-3 py-2">{formatJobDate(job.updated)}</td>
								</tr>
								{#if job.error}
									<tr class="border-t">
										<td
											colspan={5}
											class="px-3 py-2 font-mono text-xs text-destructive"
										>
											{job.error}
										</td>
									</tr>
								{/if}
							{/each}
						{/if}
					</tbody>
				</table>
			</div>
		</div>
	</section>
</div>
