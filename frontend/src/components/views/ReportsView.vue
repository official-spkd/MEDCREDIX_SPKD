<script setup lang="ts">
// ReportsView — port dari views-assurance.tsx (ReportsView).
import { FileBarChart, Search, Download, FileText, Printer, PlusCircle, Loader2 } from 'lucide-vue-next'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref } from 'vue'
import { api, downloadCsv } from '~/composables/useApi'
import { fmtDateTime } from '~/utils/format'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import Dialog from '~/components/ui/Dialog.vue'
import DemoBadge from '~/components/medcredix/DemoBadge.vue'
import EmptyState from '~/components/medcredix/EmptyState.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'

const generating = ref<string | null>(null)
const selected = ref<any>(null)
const data = ref<any>(null)

async function load() {
  try { data.value = await api('/reports') } catch { data.value = null }
}
onMounted(load)

async function generate(type: string) {
  generating.value = type
  try {
    await api('/reports', { method: 'POST', body: { type } })
    toast({ title: 'Laporan dibuat', description: 'Snapshot tersimpan dengan referensi audit.' })
    load()
  } catch (e: any) {
    toast({ title: 'Gagal membuat laporan', description: e?.message, variant: 'destructive' })
  } finally {
    generating.value = null
  }
}

async function openReport(r: any) {
  try {
    const res = await api<any>(`/reports?id=${r.id}`)
    selected.value = { ...r, content: res.content }
  } catch (e: any) {
    toast({ title: 'Gagal membuka laporan', description: e?.message, variant: 'destructive' })
  }
}

const reports = computed(() => data.value?.reports ?? [])

function renderValue(v: any): string {
  if (v == null) return '—'
  if (typeof v === 'object') return JSON.stringify(v).slice(0, 120)
  return String(v)
}
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="MedCredix Reports"
      sub="Setiap laporan memuat: tanggal generate, periode, sumber data, metodologi, versi, pembuat, dan referensi audit." demo
    />
    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="t in data?.types ?? []" :key="t.type" class="rounded-xl border bg-card p-4">
        <div class="flex items-start justify-between gap-2">
          <div>
            <div class="text-xs font-bold">{{ t.title }}</div>
            <div class="mt-0.5 text-[10px] leading-snug text-muted-foreground">{{ t.desc }}</div>
          </div>
          <FileBarChart class="h-4 w-4 shrink-0 text-brand-2" />
        </div>
        <Button size="sm" class="mt-3 w-full gap-1.5 bg-brand text-xs hover:bg-brand-hover" :disabled="generating === t.type" @click="generate(t.type)">
          <Loader2 v-if="generating === t.type" class="h-3.5 w-3.5 animate-spin" /><PlusCircle v-else class="h-3.5 w-3.5" /> Generate
        </Button>
      </div>
    </div>

    <SectionCard title="Riwayat Laporan" :sub="`${reports.length} snapshot tersimpan`">
      <EmptyState
        v-if="reports.length === 0" title="Belum ada laporan"
        desc="Generate laporan pertama Anda — snapshot tersimpan permanen dengan referensi audit."
      >
        <template #icon><FileText class="h-5 w-5" /></template>
      </EmptyState>
      <div v-else class="space-y-2">
        <div v-for="r in reports" :key="r.id" class="flex flex-wrap items-center justify-between gap-3 rounded-xl border p-3.5">
          <div class="min-w-0">
            <div class="text-xs font-semibold">
              {{ r.title }} <Badge variant="outline" class="ml-1 text-[9px]">{{ r.version }}</Badge>
            </div>
            <div class="mt-0.5 text-[10px] text-muted-foreground">
              {{ fmtDateTime(r.generatedAt) }} · oleh {{ r.generatedBy }} · Periode {{ r.period }} · Audit {{ r.auditRef }}
            </div>
          </div>
          <div class="flex gap-1.5">
            <Button size="sm" variant="outline" class="h-7 gap-1 text-[11px]" @click="openReport(r)"><Search class="h-3 w-3" /> Lihat</Button>
            <Button size="sm" variant="outline" class="h-7 gap-1 text-[11px]" @click="downloadCsv(`/reports?export=${r.id}`)"><Download class="h-3 w-3" /> CSV</Button>
            <Button size="sm" variant="outline" class="h-7 gap-1 text-[11px]" @click="openReport(r)"><Printer class="h-3 w-3" /> Cetak</Button>
          </div>
        </div>
      </div>
    </SectionCard>

    <!-- Detail drawer -->
    <Dialog :open="!!selected" @update:open="(o) => !o && (selected = null)">
      <template v-if="selected">
        <div class="flex items-center gap-2 text-[11px] uppercase tracking-wider text-muted-foreground">
          {{ selected.version }} · Periode {{ selected.period }} <DemoBadge />
        </div>
        <h2 class="mt-1 text-left text-lg font-bold">{{ selected.title }}</h2>
        <div class="print-area mt-4 space-y-4 text-sm">
          <div class="rounded-xl border bg-muted/40 p-4 text-[11px] leading-relaxed text-muted-foreground">
            <div class="mb-1 font-semibold text-foreground">Metodologi</div>
            {{ selected.methodology }}
            <div class="mt-2 grid grid-cols-2 gap-1 text-[10px]">
              <span>Dibuat: {{ fmtDateTime(selected.generatedAt) }}</span>
              <span>Oleh: {{ selected.generatedBy }}</span>
              <span>Ref. Audit: {{ selected.auditRef }}</span>
              <span>Versi: {{ selected.version }}</span>
            </div>
          </div>
          <div v-for="[key, val] in Object.entries(selected.content ?? {})" :key="key" class="rounded-xl border">
            <div class="border-b bg-muted/50 px-3.5 py-2 font-display text-xs font-semibold capitalize">
              {{ key.replace(/([A-Z])/g, ' $1') }}
            </div>
            <div v-if="Array.isArray(val)" class="divide-y">
              <div v-for="(row, i) in val" :key="i" class="grid gap-1 px-3.5 py-2 text-[11px] sm:grid-cols-2">
                <div v-for="[k, v] in Object.entries(row)" :key="k">
                  <span class="text-muted-foreground">{{ k }}:</span> <strong>{{ renderValue(v) }}</strong>
                </div>
              </div>
            </div>
            <div v-else-if="typeof val === 'object' && val !== null" class="grid gap-1.5 px-3.5 py-3 text-[11px] sm:grid-cols-2">
              <div v-for="[k, v] in Object.entries(val as any)" :key="k">
                <span class="text-muted-foreground">{{ k }}:</span> <strong>{{ renderValue(v) }}</strong>
              </div>
            </div>
            <div v-else class="px-3.5 py-3 text-xs">{{ renderValue(val) }}</div>
          </div>
        </div>
      </template>
    </Dialog>
  </div>
</template>
