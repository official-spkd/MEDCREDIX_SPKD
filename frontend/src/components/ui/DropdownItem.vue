<script setup lang="ts">
// DropdownItem — item menu dropdown; menutup menu saat diklik.
import { inject } from 'vue'
const close = inject<() => void>('dropdown-close', () => {})
const props = defineProps<{ danger?: boolean }>()
</script>
<template>
  <button
    type="button"
    :class="[
      'relative flex w-full cursor-pointer select-none items-center gap-2 rounded-md px-2 py-1.5 text-xs font-medium outline-none transition-colors',
      danger ? 'text-danger focus:bg-danger-light' : 'focus:bg-muted/70',
    ]"
    @click="close()"
  >
    <slot />
  </button>
</template>
