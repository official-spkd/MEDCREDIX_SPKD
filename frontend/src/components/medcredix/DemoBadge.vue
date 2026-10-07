<script setup lang="ts">
// DemoBadge — penanda data demo (fiktif, tanpa data pasien nyata).
import { FlaskConical } from 'lucide-vue-next'
defineProps<{ class?: string }>()
</script>
<template>
  <span :class="['inline-flex items-center gap-1 rounded-full border border-brand-2/40 bg-brand-2-light px-2.5 py-0.5 text-[10px] font-semibold text-brand-2-dark', $props.class]">
    <FlaskConical class="h-3 w-3" /> Data Demo
  </span>
</template>
