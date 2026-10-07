<script setup lang="ts">
// EmptyState — kondisi kosong yang jujur.
import { Info } from 'lucide-vue-next'
defineProps<{ title: string; desc?: string }>()
</script>
<template>
  <div class="flex flex-col items-center justify-center gap-2 rounded-xl border border-dashed bg-muted/30 px-6 py-10 text-center">
    <span class="rounded-full bg-muted p-3 text-muted-foreground"><slot name="icon"><Info class="h-5 w-5" /></slot></span>
    <div class="font-display text-sm font-semibold">{{ title }}</div>
    <p v-if="desc" class="max-w-sm text-xs leading-relaxed text-muted-foreground">{{ desc }}</p>
    <div v-if="$slots.action"><slot name="action" /></div>
  </div>
</template>
