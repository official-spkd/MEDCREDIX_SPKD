<script setup lang="ts">
// InfoPopover — perbaikan QA produk lama: ikon status wajib bisa diklik untuk penjelasan.
import { Info } from 'lucide-vue-next'
import Popover from '~/components/ui/Popover.vue'
defineProps<{ title: string }>()
</script>
<template>
  <Popover>
    <template #trigger>
      <button type="button" class="inline-flex align-middle text-muted-foreground transition-colors hover:text-brand-2" :aria-label="`Penjelasan: ${title}`">
        <slot name="icon"><Info class="h-3.5 w-3.5" /></slot>
      </button>
    </template>
    <div class="text-xs leading-relaxed">
      <div class="mb-1 font-display text-xs font-semibold text-foreground">{{ title }}</div>
      <div class="text-muted-foreground"><slot /></div>
    </div>
  </Popover>
</template>
