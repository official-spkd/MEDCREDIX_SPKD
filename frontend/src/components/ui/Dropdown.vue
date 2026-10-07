<script setup lang="ts">
// Dropdown — menu dropdown (padanan shadcn/ui DropdownMenu).
// Slot: trigger; konten ditampilkan sebagai menu mengambang.
import { onClickOutside } from '@vueuse/core'
import { provide, ref } from 'vue'
const open = ref(false)
const root = ref<HTMLElement | null>(null)
onClickOutside(root, () => (open.value = false))
provide('dropdown-close', () => (open.value = false))
</script>
<template>
  <div ref="root" class="relative inline-flex">
    <div class="inline-flex" @click="open = !open"><slot name="trigger" /></div>
    <div
      v-if="open"
      class="absolute right-0 top-full z-50 mt-1 min-w-40 overflow-hidden rounded-lg border border-border bg-popover p-1 text-popover-foreground shadow-md"
    >
      <slot />
    </div>
  </div>
</template>
