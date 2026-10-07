<script setup lang="ts">
// Popover — panel mengambang kecil (padanan shadcn/ui Popover).
import { onClickOutside } from '@vueuse/core'
import { ref } from 'vue'
const open = ref(false)
const root = ref<HTMLElement | null>(null)
onClickOutside(root, () => (open.value = false))
</script>
<template>
  <div ref="root" class="relative inline-flex">
    <div class="inline-flex" @click="open = !open"><slot name="trigger" /></div>
    <div
      v-if="open"
      class="absolute bottom-full left-1/2 z-50 mb-2 w-72 -translate-x-1/2 rounded-lg border border-border bg-popover p-3 text-popover-foreground shadow-md"
    >
      <slot />
    </div>
  </div>
</template>
