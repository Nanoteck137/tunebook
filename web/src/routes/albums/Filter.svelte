<script lang="ts">
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import {
    Button,
    Card,
    Input,
    Label,
    Select,
    Tabs,
  } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { defaultSort, FullFilter, sortTypes } from "./types";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { FilterX } from "lucide-svelte";

  export type Props = {
    fullFilter: FullFilter;
  };

  const { fullFilter }: Props = $props();

  function submit(data: FullFilter) {
    setTimeout(() => {
      const query = $page.url.searchParams;
      query.delete("query");
      query.delete("sort");

      query.delete("filterType");
      query.delete("excludeType");

      query.set("query", data.query);
      query.set("sort", data.sort);

      // if (data.filters.type.length > 0) {
      //   query.set("filterType", data.filters.type.join(","));
      // }

      // if (data.excludes.type.length > 0) {
      //   query.set("excludeType", data.excludes.type.join(","));
      // }

      goto("?" + query.toString(), { invalidateAll: true });
    }, 0);
  }

  const { form, errors, enhance, reset } = superForm(
    defaults(fullFilter, zod(FullFilter)),
    {
      id: "filter",
      SPA: true,
      validators: zod(FullFilter),
      dataType: "json",
      resetForm: false,
      onUpdate({ form }) {
        if (form.valid) {
          submit(form.data);
        }
      },
    },
  );
</script>

<form action="GET" use:enhance>
  <Tabs.Root value="normal" class="">
    <Tabs.List class="w-full justify-start overflow-y-scroll">
      <Tabs.Trigger value="normal">Normal</Tabs.Trigger>
      <Tabs.Trigger value="search">Search</Tabs.Trigger>
      <Tabs.Trigger value="database">Database</Tabs.Trigger>
    </Tabs.List>
    <Tabs.Content value="normal">
      <Card.Root>
        <Card.Content class="flex flex-col gap-4">
          <!-- <FormItem>
            <Label for="query">Search</Label>
            <Input
              id="query"
              name="query"
              type="text"
              bind:value={$form.query}
            />
            <Errors errors={$errors.query} />
          </FormItem> -->

          <div class="flex flex-col gap-4">
            <!-- <div class="flex items-center gap-4">
              <p class="min-w-20">Filters</p>

              <Select.Root type="multiple" bind:value={$form.filters.type}>
                <Select.Trigger>Type</Select.Trigger>
                <Select.Content>
                  {#each showTypes as ty (ty.value)}
                    <Select.Item value={ty.value} label={ty.label} />
                  {/each}
                </Select.Content>
              </Select.Root>
            </div>

            <div class="flex items-center gap-4">
              <p class="min-w-20">Excludes</p>

              <Select.Root type="multiple" bind:value={$form.excludes.type}>
                <Select.Trigger>Type</Select.Trigger>
                <Select.Content>
                  {#each showTypes as ty (ty.value)}
                    <Select.Item value={ty.value} label={ty.label} />
                  {/each}
                </Select.Content>
              </Select.Root>
            </div> -->

            <div class="flex items-center gap-4">
              <p class="min-w-20">Sort</p>

              <Select.Root
                type="single"
                allowDeselect={false}
                bind:value={$form.sort}
              >
                <Select.Trigger>
                  {sortTypes.find((i) => i.value === $form.sort)?.label ??
                    "Sort"}
                </Select.Trigger>
                <Select.Content>
                  {#each sortTypes as ty (ty.value)}
                    <Select.Item value={ty.value} label={ty.label} />
                  {/each}
                </Select.Content>
              </Select.Root>
            </div>
          </div>
        </Card.Content>
        <Card.Footer class="flex gap-2">
          <Button
            variant="outline"
            onclick={() => {
              reset({
                data: {
                  query: "",
                  filters: { type: [] },
                  excludes: { type: [] },
                  sort: defaultSort,
                },
              });
            }}
          >
            <FilterX />
            Reset
          </Button>

          <Button type="submit">Filter</Button>
        </Card.Footer>
      </Card.Root>
    </Tabs.Content>
    <Tabs.Content value="password">Change your password here.</Tabs.Content>
  </Tabs.Root>

  <div class="h-[1000px]"></div>

  <div class="">
    <div class="flex overflow-x-scroll">
      <p class="rounded-t-xl border-l border-r border-t px-4 py-2">Search</p>
      <p
        class="ml-[1px] mt-[1px] border-b px-4 py-2 hover:ml-0 hover:mt-0 hover:rounded-t-xl hover:border-b-0 hover:border-l hover:border-r hover:border-t"
      >
        Filter
      </p>
      <p
        class="ml-[1px] mt-[1px] border-b px-4 py-2 hover:ml-0 hover:mt-0 hover:rounded-t-xl hover:border-b-0 hover:border-l hover:border-r hover:border-t"
      >
        Database
      </p>
      <div class="w-full border-b"></div>
    </div>

    <div
      class="rounded-b-xl border-b border-l border-r bg-card text-card-foreground shadow"
    >
      <Card.Content class="flex flex-col gap-4">
        <FormItem>
          <Label for="query">Search</Label>
          <Input
            id="query"
            name="query"
            type="text"
            bind:value={$form.query}
          />
          <Errors errors={$errors.query} />
        </FormItem>

        <div class="flex flex-col gap-4">
          <div class="flex items-center gap-4">
            <p class="min-w-20">Filters</p>

            <!-- <Select.Root type="multiple" bind:value={$form.filters.type}>
            <Select.Trigger>Type</Select.Trigger>
            <Select.Content>
              {#each showTypes as ty (ty.value)}
                <Select.Item value={ty.value} label={ty.label} />
              {/each}
            </Select.Content>
          </Select.Root> -->
          </div>

          <div class="flex items-center gap-4">
            <p class="min-w-20">Excludes</p>

            <!-- <Select.Root type="multiple" bind:value={$form.excludes.type}>
            <Select.Trigger>Type</Select.Trigger>
            <Select.Content>
              {#each showTypes as ty (ty.value)}
                <Select.Item value={ty.value} label={ty.label} />
              {/each}
            </Select.Content>
          </Select.Root> -->
          </div>

          <div class="flex items-center gap-4">
            <p class="min-w-20">Sort</p>

            <Select.Root
              type="single"
              allowDeselect={false}
              bind:value={$form.sort}
            >
              <Select.Trigger>
                {sortTypes.find((i) => i.value === $form.sort)?.label ??
                  "Sort"}
              </Select.Trigger>
              <Select.Content>
                {#each sortTypes as ty (ty.value)}
                  <Select.Item value={ty.value} label={ty.label} />
                {/each}
              </Select.Content>
            </Select.Root>
          </div>
        </div>
      </Card.Content>

      <Card.Footer class="flex gap-2">
        <Button
          variant="outline"
          onclick={() => {
            reset({
              data: {
                query: "",
                filters: { type: [] },
                excludes: { type: [] },
                sort: defaultSort,
              },
            });
          }}
        >
          <FilterX />
          Reset
        </Button>

        <Button type="submit">Filter</Button>
      </Card.Footer>
    </div>
  </div>

  <div class="h-10"></div>

  <Card.Root>
    <Card.Header>
      <Card.Title>Card Title</Card.Title>
      <Card.Description>Card Description</Card.Description>
    </Card.Header>

    <Card.Content class="flex flex-col gap-4">
      <FormItem>
        <Label for="query">Search</Label>
        <Input id="query" name="query" type="text" bind:value={$form.query} />
        <Errors errors={$errors.query} />
      </FormItem>

      <div class="flex flex-col gap-4">
        <div class="flex items-center gap-4">
          <p class="min-w-20">Filters</p>

          <!-- <Select.Root type="multiple" bind:value={$form.filters.type}>
            <Select.Trigger>Type</Select.Trigger>
            <Select.Content>
              {#each showTypes as ty (ty.value)}
                <Select.Item value={ty.value} label={ty.label} />
              {/each}
            </Select.Content>
          </Select.Root> -->
        </div>

        <div class="flex items-center gap-4">
          <p class="min-w-20">Excludes</p>

          <!-- <Select.Root type="multiple" bind:value={$form.excludes.type}>
            <Select.Trigger>Type</Select.Trigger>
            <Select.Content>
              {#each showTypes as ty (ty.value)}
                <Select.Item value={ty.value} label={ty.label} />
              {/each}
            </Select.Content>
          </Select.Root> -->
        </div>

        <div class="flex items-center gap-4">
          <p class="min-w-20">Sort</p>

          <Select.Root
            type="single"
            allowDeselect={false}
            bind:value={$form.sort}
          >
            <Select.Trigger>
              {sortTypes.find((i) => i.value === $form.sort)?.label ?? "Sort"}
            </Select.Trigger>
            <Select.Content>
              {#each sortTypes as ty (ty.value)}
                <Select.Item value={ty.value} label={ty.label} />
              {/each}
            </Select.Content>
          </Select.Root>
        </div>
      </div>
    </Card.Content>
    <Card.Footer class="flex gap-2">
      <Button
        variant="outline"
        onclick={() => {
          reset({
            data: {
              query: "",
              filters: { type: [] },
              excludes: { type: [] },
              sort: defaultSort,
            },
          });
        }}
      >
        <FilterX />
        Reset
      </Button>

      <Button type="submit">Filter</Button>
    </Card.Footer>
  </Card.Root>
</form>
