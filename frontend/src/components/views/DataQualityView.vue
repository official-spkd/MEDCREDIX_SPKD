<script setup lang="ts">
// DataQualityView — port dari views-data.tsx (DataQualityView).
import { useMedCredix } from '~/stores/medcredix'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import { fmtDateTime, fmtID } from '~/utils/format'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const data = ref<any>(null)
onMounted(async () => {
  try { data.value = await api('/data?type=quality') } catch { data.value = null }
})
const issues = computed(() => data.value?.issues ?? [])
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Data Quality Center"
      sub="Masalah kualitas data dapat menjadi risiko → rekomendasi → aksi. Dengan sinyal anti-manipulasi (anomaly) yang ditandai khusus." demo
    />
    <div class="grid grid-cols-3 gap-3">
      <div v-for="s in [
        { label: 'Isu Terbuka', value: data?.summary?.terbuka ?? 0, tone: 'warning' },
        { label: 'Total Isu Terdeteksi', value: data?.summary?.total ?? 0, tone: 'info' },
        { label: 'Sinyal Anti-Manipulasi', value: data?.summary?.anomaly ?? 0, tone: 'danger' },
      ]" :key="s.label" class="rounded-xl border bg-card p-4">
        <div class="text-[10px] font-bold uppercase tracking-wide text-muted-foreground">{{ s.label }}</div>
        <div :class="['mt-1 font-display text-2xl font-bold', s.tone === 'danger' ? 'text-danger' : s.tone === 'warning' ? 'text-warning' : 'text-brand']">{{ s.value }}</div>
      </div>
    </div>

    <div class="grid gap-3 lg:grid-cols-2">
      <div v-for="i in issues" :key="i.id" :class="['rounded-xl border bg-card p-4', i.anomaly && 'border-warning/40 bg-warning-light/20']">
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-1.5">
              <span v-if="i.anomaly" class="rounded-full bg-warning px-2 py-0.5 text-[9px] font-bold text-white">ANOMALY FLAG</span>
              <Badge variant="outline" class="text-[10px]">{{ i.dimension }}</Badge>
              <StatusBadge :status="i.status" />
            </div>
            <h3 class="mt-2 font-display text-sm font-semibold">{{ i.title }}</h3>
            <p class="mt-1 text-xs leading-relaxed text-muted-foreground">{{ i.description }}</p>
          </div>
          <StatusBadge :status="i.severity" :tone="i.severity === 'Tinggi' ? 'danger' : i.severity === 'Sedang' ? 'warning' : 'info'" />
        </div>
        <div class="mt-3 flex flex-wrap items-center justify-between gap-2 border-t pt-3 text-[10px] text-muted-foreground">
          <span>{{ i.sourceSystem }} · {{ fmtID(i.affectedRecords) }} record terdampak · {{ fmtDateTime(i.detectedAt) }}</span>
          <Button
            size="sm" variant="ghost" class="h-6 gap-1 text-[10px] text-brand-2"
            @click="store.openAction({ title: `Selesaikan: ${i.title}`, description: i.description, sourceType: 'Inspeksi', sourceRef: i.dimension, priority: i.severity })"
          >Buat Tindakan</Button>
        </div>
      </div>
    </div>
  </div>
</template>
