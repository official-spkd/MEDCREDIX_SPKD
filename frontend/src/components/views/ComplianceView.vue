<script setup lang="ts">
// ComplianceView — port dari views-assurance.tsx (ComplianceView).
import { Scale, Database, Download } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import { fmtDate } from '~/utils/format'
import Badge from '~/components/ui/Badge.vue'
import Dialog from '~/components/ui/Dialog.vue'
import Select from '~/components/ui/Select.vue'
import SelectOption from '~/components/ui/SelectOption.vue'
import DemoBadge from '~/components/medcredix/DemoBadge.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const status = ref('all')
const selected = ref<any>(null)
const data = ref<any>(null)
onMounted(async () => {
  try { data.value = await api('/compliance') } catch { data.value = null }
})
const mappings = computed(() => (data.value?.mappings ?? []).filter((m: any) => status.value === 'all' || m.status === status.value))
const regs = computed(() => data.value?.regulations ?? [])
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Regulatory Compliance"
      sub="Regulasi → versi → requirement → pemetaan modul/fitur/data/validasi → bukti → audit. Tidak pernah klaim '100% patuh'." demo
    />
    <div class="grid grid-cols-2 gap-3 md:grid-cols-4">
      <div v-for="s in [
        { label: 'Patuh', value: data?.summary?.patuh ?? 0, tone: 'success' },
        { label: 'Sebagian', value: data?.summary?.sebagian ?? 0, tone: 'warning' },
        { label: 'Belum', value: data?.summary?.belum ?? 0, tone: 'danger' },
        { label: 'Sedang Ditinjau', value: data?.summary?.ditinjau ?? 0, tone: 'info' },
      ]" :key="s.label" class="rounded-xl border bg-card p-4">
        <div class="text-[10px] font-bold uppercase tracking-wide text-muted-foreground">{{ s.label }}</div>
        <div :class="['mt-1 font-display text-2xl font-bold', s.tone === 'success' ? 'text-success' : s.tone === 'warning' ? 'text-warning' : s.tone === 'danger' ? 'text-danger' : 'text-info']">{{ s.value }}</div>
      </div>
    </div>

    <SectionCard title="Basis Regulasi (dengan Versioning)" sub="Regulasi berubah → versi lama superseded → analisis dampak → pembaruan — tidak ada asumsi regulasi yang di-hard-code">
      <template #icon><Scale class="h-4 w-4" /></template>
      <div class="grid gap-2 md:grid-cols-2">
        <button
          v-for="r in regs" :key="r.id"
          class="rounded-xl border bg-card p-3.5 text-left transition-all hover:border-brand-2/50 hover:shadow-sm"
          @click="selected = r"
        >
          <div class="flex items-start justify-between gap-2">
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-1.5">
                <Badge variant="outline" class="text-[9px] font-bold">{{ r.category }}</Badge>
                <StatusBadge :status="r.status" />
              </div>
              <div class="mt-1.5 truncate text-xs font-semibold">{{ r.title }}</div>
              <div class="mt-0.5 text-[10px] text-muted-foreground">{{ r.number }} · {{ r.authority }}</div>
            </div>
            <span class="font-display text-sm font-bold text-muted-foreground">{{ r.versions?.length ?? 1 }}v</span>
          </div>
        </button>
      </div>
    </SectionCard>

    <SectionCard title="Matriks Kepatuhan (Compliance Matrix)" :sub="`${mappings.length} pemetaan requirement → fitur → validasi → bukti`">
      <template #action>
        <div class="w-40">
          <Select v-model="status">
            <SelectOption value="all">Semua Status</SelectOption>
            <SelectOption v-for="s in ['Patuh', 'Sebagian', 'Belum', 'Sedang Ditinjau']" :key="s" :value="s">{{ s }}</SelectOption>
          </Select>
        </div>
      </template>
      <div class="overflow-x-auto rounded-xl border">
        <table class="w-full text-xs">
          <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
            <tr>
              <th class="px-3 py-2.5">Regulasi</th><th class="px-3 py-2.5">Pasal</th><th class="px-3 py-2.5">Requirement</th>
              <th class="px-3 py-2.5">Modul / Fitur</th><th class="px-3 py-2.5">Status</th><th class="px-3 py-2.5">Review</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="m in mappings" :key="m.id" class="border-t">
              <td class="px-3 py-2.5"><div class="font-mono text-[10px] font-semibold">{{ m.regulation?.number }}</div></td>
              <td class="px-3 py-2.5 text-muted-foreground">{{ m.article }}</td>
              <td class="max-w-64 px-3 py-2.5"><div class="truncate font-medium" :title="m.requirement">{{ m.requirement }}</div></td>
              <td class="px-3 py-2.5 text-[10px] text-muted-foreground">{{ m.module }} / {{ m.feature }}</td>
              <td class="px-3 py-2.5"><StatusBadge :status="m.status" /></td>
              <td class="px-3 py-2.5 text-[10px] text-muted-foreground">{{ fmtDate(m.lastReviewedAt) }}<br />{{ m.reviewer }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p class="mt-3 rounded-lg bg-info-light px-3 py-2 text-[11px] leading-relaxed text-info">
        Dirancang untuk mendukung kepatuhan terhadap regulasi kesehatan Indonesia yang berlaku — tunduk pada verifikasi regulasi terkini dan implementasi spesifik rumah sakit.
      </p>
    </SectionCard>

    <!-- Detail drawer -->
    <Dialog :open="!!selected" @update:open="(o) => !o && (selected = null)">
      <template v-if="selected">
        <div class="flex items-center gap-2 text-[11px] uppercase tracking-wider text-muted-foreground">
          {{ selected.number }} · {{ selected.authority }} <DemoBadge />
        </div>
        <h2 class="mt-1 text-left text-lg font-bold leading-snug">{{ selected.title }}</h2>
        <div class="mt-4 space-y-4">
          <p class="text-sm leading-relaxed text-foreground/90">{{ selected.summary }}</p>
          <div class="rounded-xl border p-4">
            <div class="mb-2 text-xs font-semibold">Riwayat Versi</div>
            <div class="space-y-2">
              <div v-for="v in selected.versions ?? []" :key="v.id" :class="['rounded-lg border p-3', v.status === 'Berlaku' ? 'border-success/40 bg-success-light/30' : 'border-border bg-muted/40']">
                <div class="flex items-center justify-between gap-2">
                  <span class="text-xs font-semibold">{{ v.version }}</span>
                  <StatusBadge :status="v.status" :tone="v.status === 'Berlaku' ? 'success' : v.status === 'Superseded' ? 'muted' : 'info'" />
                </div>
                <div class="mt-1 text-[11px] text-muted-foreground">{{ v.changeSummary }}</div>
                <div class="mt-1 flex flex-wrap gap-2 text-[10px] text-muted-foreground">
                  <span>Ditetapkan {{ fmtDate(v.issuedAt) }}</span>
                  <span v-if="v.impactLevel" :class="['font-semibold', v.impactLevel === 'Tinggi' ? 'text-danger' : 'text-warning']">· Dampak {{ v.impactLevel }}</span>
                  <span v-if="v.supersededBy">· Digantikan: {{ v.supersededBy }}</span>
                </div>
              </div>
            </div>
          </div>
          <a :href="selected.source" target="_blank" rel="noreferrer" class="inline-flex items-center gap-1.5 text-xs font-semibold text-brand-2 hover:underline">
            <Database class="h-3.5 w-3.5" /> Sumber resmi (JDIH) <Download class="h-3 w-3" />
          </a>
        </div>
      </template>
    </Dialog>
  </div>
</template>
