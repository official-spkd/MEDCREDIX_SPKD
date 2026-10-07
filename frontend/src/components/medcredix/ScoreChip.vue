<script setup lang="ts">
// ScoreChip — kartu skor domain ringkas, klik untuk drilldown.
import { ChevronRight } from 'lucide-vue-next'
import { computed } from 'vue'
const props = defineProps<{ label: string; score: number; group?: string; onClick?: () => void }>()
const t = computed(() => (props.score >= 90 ? 'success' : props.score >= 75 ? 'info' : props.score >= 60 ? 'warning' : 'danger'))
const toneHex: Record<string, string> = { success: '#059669', warning: '#D97706', danger: '#DC2626', info: '#0284C7' }
</script>

<template>
  <button
    type="button"
    :class="[
      'group flex w-full items-center justify-between gap-3 rounded-xl border bg-card p-4 text-left transition-all',
      'hover:border-brand-2/60 hover:shadow-sm focus-visible:ring-2 focus-visible:ring-ring',
    ]"
    @click="onClick?.()"
  >
    <div class="min-w-0">
      <div v-if="group" class="text-[10px] font-semibold uppercase tracking-wide text-muted-foreground">{{ group }}</div>
      <div class="truncate font-display text-sm font-semibold text-foreground">{{ label }}</div>
    </div>
    <div class="flex items-center gap-2">
      <span class="font-display text-xl font-bold" :style="{ color: toneHex[t] }">{{ Math.round(score) }}</span>
      <ChevronRight class="h-4 w-4 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5 group-hover:text-brand-2" />
    </div>
  </button>
</template>
