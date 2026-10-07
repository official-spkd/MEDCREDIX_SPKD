<script setup lang="ts">
// Hint — tooltip penjelasan singkat.
import { Info } from 'lucide-vue-next'
import Tooltip from '~/components/ui/Tooltip.vue'
</script>
<template>
  <Tooltip>
    <template #trigger>
      <button type="button" class="inline-flex text-muted-foreground/70 transition-colors hover:text-brand-2" aria-label="Penjelasan">
        <Info class="h-3.5 w-3.5" />
      </button>
    </template>
    <slot />
  </Tooltip>
</template>
