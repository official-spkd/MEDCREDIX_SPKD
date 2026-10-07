<script setup lang="ts">
// Badge — padanan shadcn/ui Badge.
const props = defineProps<{ variant?: 'default' | 'outline' | 'secondary' }>()
const cls =
  props.variant === 'outline'
    ? 'border border-border text-foreground'
    : props.variant === 'secondary'
      ? 'bg-muted text-foreground'
      : 'bg-brand-light text-brand'
</script>

<template>
  <span :class="['inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold', cls]">
    <slot />
  </span>
</template>
