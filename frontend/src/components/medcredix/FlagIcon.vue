<script setup lang="ts">
// FlagIcon — ikon flag konsistensi klaim (ok/warn/bad).
import { CheckCircle2, AlertTriangle, XCircle } from 'lucide-vue-next'
defineProps<{ flag: string }>()
</script>
<template>
  <CheckCircle2 v-if="flag === 'ok'" class="h-4 w-4 text-success" aria-label="Konsisten" />
  <AlertTriangle v-else-if="flag === 'warn'" class="h-4 w-4 text-warning" aria-label="Perlu ditinjau" />
  <XCircle v-else class="h-4 w-4 text-danger" aria-label="Sinyal tinggi" />
</template>
