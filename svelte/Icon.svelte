<script lang="ts">
import { getContext } from "svelte";
import { Context, iconKey } from "./icon-context";

const context = getContext<Context>(iconKey);

export let className = context.className;
export let width = context.width;
export let height = context.height;
export let x = context.x;
export let y = context.y;
</script>

<svg class={className} {width} {height} {x} {y}>
  [SVG_CONTENT]
</svg>
