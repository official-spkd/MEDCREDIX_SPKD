<script setup lang="ts">
import { computed, inject } from 'vue'
const props = defineProps<{ value: string }>()
const tabs = inject<{ current: { value: string }; select: (v: string) => void }>('tabs')!
const active = computed(() => tabs.current.value === props.value)
</script>
<template>
  <button
    type="button"
    :class="[
      'inline-flex items-center justify-center gap-1.5 whitespace-nowrap rounded-md px-3 py-1 text-xs font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
      active ? 'bg-card text-foreground shadow-sm' : 'hover:text-foreground',
    ]"
    @click="tabs.select(value)"
  >
    <slot />
  </button>
</template>
