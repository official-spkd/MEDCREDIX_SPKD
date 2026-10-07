<script setup lang="ts">
// EvidenceView — port dari views-assurance.tsx (EvidenceView).
import { Search } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import { fmtDateTime, fmtID } from '~/utils/format'
import Badge from '~/components/ui/Badge.vue'
import Input from '~/components/ui/Input.vue'
import Select from '~/components/ui/Select.vue'
import SelectOption from '~/components/ui/SelectOption.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const q = ref('')
const refType = ref('all')
const data = ref<any>(null)

async function load() {
  try {
    const params = new URLSearchParams()
    if (q.value) params.set('q', q.value)
    if (refType.value !== 'all') params.set('refType', refType.value)
    data.value = await api(`/evidence?${params}`)
  } catch { data.value = null }
}
onMounted(load)
watch(refType, load)
let t: ReturnType<typeof setTimeout> | null = null
watch(q, () => { if (t) clearTimeout(t); t = setTimeout(load, 300) })
const items = computed(() => data.value?.items ?? [])
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Evidence Center"
      sub="Tempat semua bukti hidup: sumber, timestamp, jumlah record, resource FHIR, status validasi. Tidak pernah ada bukti yang dipalsukan." demo
    />
    <div class="grid grid-cols-2 gap-3 md:grid-cols-5">
      <div v-for="[t, c] in Object.entries(data?.byType ?? {}).slice(0, 5)" :key="t" class="rounded-xl border bg-card p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wide text-muted-foreground">{{ t }}</div>
        <div class="mt-0.5 font-display text-xl font-bold text-brand">{{ c }}</div>
      </div>
    </div>
    <div class="flex flex-wrap items-center gap-2 rounded-xl border bg-card p-3">
      <div class="relative">
        <Search class="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
        <Input v-model="q" placeholder="Cari judul, metrik, atau sistem sumber…" class="h-8 w-64 pl-8 text-xs" />
      </div>
      <div class="w-44">
        <Select v-model="refType">
          <SelectOption value="all">Semua Jenis Referensi</SelectOption>
          <SelectOption v-for="t in ['Skor', 'Risiko', 'Klaim', 'Peringatan', 'Indikator', 'Regulasi', 'Aksi', 'Kredensial']" :key="t" :value="t">{{ t }}</SelectOption>
        </Select>
      </div>
      <span class="ml-auto text-xs text-muted-foreground">{{ items.length }} item bukti</span>
    </div>
    <div class="grid gap-3 lg:grid-cols-2">
      <div v-for="e in items" :key="e.id" class="rounded-xl border bg-card p-4">
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <Badge variant="outline" class="text-[9px] font-bold">{{ e.refType }}</Badge>
            <h3 class="mt-1.5 text-xs font-semibold leading-snug">{{ e.title }}</h3>
            <div class="mt-1 font-display text-lg font-bold text-brand">{{ e.value }}{{ e.unit === '%' ? '%' : ` ${e.unit ?? ''}` }}</div>
          </div>
          <StatusBadge :status="e.validationStatus" />
        </div>
        <div class="mt-2.5 flex flex-wrap gap-x-3 gap-y-1 border-t pt-2.5 text-[10px] text-muted-foreground">
          <span>{{ e.metric }}</span>
          <span>· {{ fmtID(e.recordCount ?? 0) }} record</span>
          <span>· Sumber: {{ e.sourceSystem }}</span>
          <span v-if="e.fhirResource">· FHIR: {{ e.fhirResource }}</span>
          <span>· Sync: {{ fmtDateTime(e.lastSyncAt) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
