<script setup lang="ts">
// KpiCard — kartu KPI dengan ikon & sub-teks.
const props = defineProps<{
  label: string
  value: any
  unit?: string
  tone?: string
  onClick?: () => void
}>()
const toneCls: Record<string, string> = {
  success: 'bg-success-light text-success',
  warning: 'bg-warning-light text-warning',
  danger: 'bg-danger-light text-danger',
  info: 'bg-info-light text-info',
  muted: 'bg-muted text-muted-foreground',
}
</script>

<template>
  <div :class="['rounded-xl border border-border/80 bg-card', onClick ? 'cursor-pointer transition-all hover:border-brand-2/50 hover:shadow-sm' : '']">
    <button type="button" :disabled="!onClick" class="w-full p-4 text-left" @click="onClick?.()">
      <div class="flex items-start justify-between gap-2">
        <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground">{{ label }}</div>
        <span v-if="$slots.icon" :class="['rounded-lg p-1.5', toneCls[tone ?? 'info']]"><slot name="icon" /></span>
      </div>
      <div class="mt-2 font-display text-2xl font-bold leading-tight text-foreground">
        {{ value }}<span v-if="unit" class="ml-1 text-sm font-semibold text-muted-foreground">{{ unit }}</span>
      </div>
      <div v-if="$slots.sub" class="mt-1 flex flex-wrap items-center gap-1.5 text-xs text-muted-foreground"><slot name="sub" /></div>
    </button>
  </div>
</template>
