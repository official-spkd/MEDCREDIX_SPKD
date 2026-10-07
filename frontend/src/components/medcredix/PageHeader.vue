<script setup lang="ts">
// PageHeader — kepala halaman standar.
import DemoBadge from '~/components/medcredix/DemoBadge.vue'
defineProps<{ title: string; sub?: string; demo?: boolean }>()
</script>
<template>
  <div class="mb-5 flex flex-wrap items-end justify-between gap-3">
    <div>
      <div class="flex items-center gap-2.5">
        <h1 class="font-display text-xl font-bold tracking-tight text-foreground">{{ title }}</h1>
        <DemoBadge v-if="demo" />
      </div>
      <p v-if="sub" class="mt-1 text-sm text-muted-foreground">{{ sub }}</p>
    </div>
    <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2"><slot name="actions" /></div>
  </div>
</template>
