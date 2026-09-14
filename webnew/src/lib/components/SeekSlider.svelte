<script lang="ts">
  import { cn } from "$lib/utils";

  interface Props {
    value: number;
    // eslint-disable-next-line no-unused-vars
    onValue: (value: number) => void;
    growOnHover?: boolean;
    buffered?: number;
    class?: string;
    ariaLabel?: string;
  }

  let {
    value,
    onValue,
    growOnHover = true,
    buffered = 0,
    class: className = "",
    ariaLabel = "Seek",
  }: Props = $props();

  let trackEl: HTMLDivElement | undefined = $state();
  let dragging = $state(false);
  let dragValue = $state(0);

  let settling = $state(false);
  let lastDragValue = 0;

  let displayValue = $derived.by(() => {
    if (dragging) return dragValue;
    if (settling) return lastDragValue;
    return value;
  });

  function clamp(v: number) {
    return Math.min(1, Math.max(0, v));
  }

  function percentFromEvent(e: PointerEvent) {
    if (!trackEl) return 0;
    const rect = trackEl.getBoundingClientRect();
    return clamp((e.clientX - rect.x) / rect.width);
  }

  function onPointerDown(e: PointerEvent) {
    dragging = true;
    dragValue = percentFromEvent(e);
    (e.target as HTMLElement).setPointerCapture(e.pointerId);
  }

  function onPointerMove(e: PointerEvent) {
    if (!dragging) return;
    dragValue = percentFromEvent(e);
  }

  // eslint-disable-next-line no-unused-vars, @typescript-eslint/no-unused-vars
  function onPointerUp(_e: PointerEvent) {
    if (!dragging) return;
    onValue(dragValue);
    lastDragValue = dragValue;
    settling = true;
    dragging = false;
  }

  function onKeyDown(e: KeyboardEvent) {
    const step = e.shiftKey ? 0.1 : 0.025;
    let newVal = value;

    if (e.key === "ArrowRight" || e.key === "ArrowUp") {
      newVal = clamp(value + step);
    } else if (e.key === "ArrowLeft" || e.key === "ArrowDown") {
      newVal = clamp(value - step);
    } else {
      return;
    }

    e.preventDefault();
    onValue(newVal);
    lastDragValue = newVal;
    settling = true;
  }

  $effect(() => {
    if (settling && Math.abs(value - lastDragValue) < 0.005) {
      settling = false;
    }
  });

  let thumbClasses = $derived(
    "absolute top-1/2 size-3 -translate-x-1/2 -translate-y-1/2 rounded-full border-2 border-primary bg-background shadow-sm transition-[width,height,opacity]" +
      (growOnHover ? " group-hover:size-4 group-hover:opacity-100" : "") +
      (dragging ? " size-4 opacity-100" : " opacity-0") +
      (!growOnHover && dragging ? " !opacity-100" : ""),
  );

  let trackClasses = $derived(
    "absolute left-0 right-0 top-1/2 h-1 -translate-y-1/2 rounded-full bg-muted" +
      (growOnHover ? " transition-[height] group-hover:h-1.5" : ""),
  );

  let fillClasses = $derived(
    "pointer-events-none absolute left-0 top-1/2 h-1 -translate-y-1/2 rounded-full bg-primary" +
      (growOnHover ? " transition-[height] group-hover:h-1.5" : ""),
  );
</script>

<div
  class={cn("group relative h-5 w-full touch-none", className)}
  bind:this={trackEl}
  onpointerdown={onPointerDown}
  onpointermove={onPointerMove}
  onpointerup={onPointerUp}
  onpointercancel={onPointerUp}
  onkeydown={onKeyDown}
  role="slider"
  tabindex="0"
  aria-valuenow={Math.round(displayValue * 100)}
  aria-valuemin="0"
  aria-valuemax="100"
  aria-label={ariaLabel}
>
  <!-- Track background -->
  <div class={trackClasses}></div>

  <!-- Buffered fill -->
  <div
    class="pointer-events-none absolute left-0 top-1/2 h-1 -translate-y-1/2 rounded-full bg-muted-foreground/20"
    style="width: {buffered * 100}%"
  ></div>

  <!-- Track fill -->
  <div class={fillClasses} style="width: {displayValue * 100}%"></div>

  <!-- Thumb -->
  <div class={thumbClasses} style="left: {displayValue * 100}%"></div>
</div>
