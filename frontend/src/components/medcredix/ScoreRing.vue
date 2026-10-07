<script setup lang="ts">
// ScoreRing — cincin skor 0–100, dapat diklik untuk drilldown "Mengapa".
import { computed } from 'vue'
const props = defineProps<{
  value: number
  size?: number
  stroke?: number
  label?: string
  sub?: string
  tone?: string
  onClick?: () => void
}>()
const size = props.size ?? 132
const strokeWidth = props.stroke ?? 11
const t = computed(() => props.tone ?? (props.value >= 90 ? 'success' : props.value >= 75 ? 'info' : props.value >= 60 ? 'warning' : 'danger'))
const toneHex: Record<string, string> = {
  success: '#059669', warning: '#D97706', danger: '#DC2626', info: '#0284C7', muted: '#5B6B84',
}
const r = (size - strokeWidth) / 2
const circ = 2 * Math.PI * r
const off = computed(() => circ * (1 - Math.min(100, Math.max(0, props.value)) / 100))
</script>

<template>
  <button
    type="button"
    :disabled="!onClick"
    :class="['relative inline-flex items-center justify-center', onClick ? 'cursor-pointer group' : '']"
    :aria-label="label ? `${label}: ${value} dari 100` : `Skor ${value}`"
    @click="onClick?.()"
  >
    <svg :width="size" :height="size" class="-rotate-90">
      <circle :cx="size / 2" :cy="size / 2" :r="r" :stroke-width="strokeWidth" class="stroke-muted" fill="none" />
      <circle
        :cx="size / 2" :cy="size / 2" :r="r" :stroke-width="strokeWidth" :stroke="toneHex[t]"
        fill="none" :stroke-dasharray="circ" :stroke-dashoffset="off" stroke-linecap="round"
        style="transition: stroke-dashoffset 900ms cubic-bezier(0.22, 1, 0.36, 1)"
      />
    </svg>
    <span class="absolute inset-0 flex flex-col items-center justify-center">
      <span :class="['font-display font-bold leading-none', size > 100 ? 'text-4xl' : 'text-2xl']" :style="{ color: toneHex[t] }">
        {{ Math.round(value) }}<span class="text-base font-semibold text-muted-foreground">/100</span>
      </span>
      <span v-if="label" class="mt-1.5 max-w-[86%] text-center text-[11px] font-medium leading-tight text-muted-foreground">{{ label }}</span>
      <span v-if="sub" class="mt-0.5 text-[10px] text-muted-foreground/80">{{ sub }}</span>
    </span>
  </button>
</template>
