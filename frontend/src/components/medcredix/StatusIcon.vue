<script setup lang="ts">
// StatusIcon — ikon status semantik (selalu ikon, bukan warna saja — WCAG).
import { CheckCircle2, AlertTriangle, XCircle, Info, CircleDashed } from 'lucide-vue-next'
defineProps<{ tone: string }>()
</script>
<template>
  <CheckCircle2 v-if="tone === 'success'" class="h-4 w-4 text-success" aria-label="Baik" />
  <AlertTriangle v-else-if="tone === 'warning'" class="h-4 w-4 text-warning" aria-label="Perlu perhatian" />
  <XCircle v-else-if="tone === 'danger'" class="h-4 w-4 text-danger" aria-label="Kritis" />
  <Info v-else-if="tone === 'info'" class="h-4 w-4 text-info" aria-label="Informasi" />
  <CircleDashed v-else class="h-4 w-4 text-muted-foreground" aria-label="Tidak ada data" />
</template>
