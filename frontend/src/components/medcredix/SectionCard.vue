<script setup lang="ts">
// SectionCard — kartu bagian dengan judul, sub, ikon & slot action.
defineProps<{ title: any; sub?: any }>()
</script>

<template>
  <div class="rounded-xl border border-border/80 bg-card">
    <div class="flex flex-row items-start justify-between space-y-0 p-5 pb-3">
      <div class="flex items-start gap-2.5">
        <span v-if="$slots.icon" class="mt-0.5 rounded-lg bg-brand-light p-1.5 text-brand"><slot name="icon" /></span>
        <div>
          <div class="font-display text-[15px] font-semibold leading-tight">{{ title }}</div>
          <p v-if="sub" class="mt-0.5 text-xs text-muted-foreground">{{ sub }}</p>
        </div>
      </div>
      <div v-if="$slots.action"><slot name="action" /></div>
    </div>
    <div class="px-5 pb-5"><slot /></div>
  </div>
</template>
