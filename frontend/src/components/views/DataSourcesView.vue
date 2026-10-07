<script setup lang="ts">
// DataSourcesView — port dari views-data.tsx (DataSourcesView).
import { Cable, CheckCircle2, ShieldAlert, RefreshCw, Loader2, PlayCircle, Activity } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import { fmtDateTime, fmtID } from '~/utils/format'
import Button from '~/components/ui/Button.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const data = ref<any>(null)
const testing = ref(false)
const testResult = ref<any>(null)

async function load() {
  try { data.value = await api('/data?type=sources') } catch { data.value = null }
}
onMounted(load)

async function testConnection() {
  testing.value = true
  try {
    const res = await api<any>('/data', { method: 'POST', body: { action: 'test-connection' } })
    testResult.value = res
    toast({ title: `Test Koneksi selesai: ${res.summary.ok} OK · ${res.summary.perlu} perlu perbaikan` })
  } catch (e: any) {
    toast({ title: 'Test koneksi gagal', description: e?.message, variant: 'destructive' })
  } finally {
    testing.value = false
  }
}
const sources = computed(() => data.value?.sources ?? [])
const jobs = computed(() => data.value?.jobs ?? [])
const canConfigure = computed(() => store.user?.permissions.includes('CONFIGURE'))
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Data Sources & Integration Health"
      sub="10 sistem sumber — jadwal pull 4×/hari (02.00, 08.00, 14.00, 20.00 WIB), cross-reference anti-manipulasi 02.45/08.45, scoring inkremental tiap 15 menit." demo
    >
      <template #actions>
        <Button v-if="canConfigure" size="sm" class="gap-1.5 bg-brand-2 hover:bg-brand-2-hover" :disabled="testing" @click="testConnection">
          <Loader2 v-if="testing" class="h-3.5 w-3.5 animate-spin" /><PlayCircle v-else class="h-3.5 w-3.5" /> Test Koneksi (24 Resource)
        </Button>
      </template>
    </PageHeader>

    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
      <div v-for="s in sources" :key="s.id" class="rounded-xl border bg-card p-4">
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <div class="truncate text-xs font-semibold">{{ s.name }}</div>
            <div class="mt-0.5 text-[9px] uppercase tracking-wide text-muted-foreground">{{ s.type }}</div>
          </div>
          <StatusBadge :status="s.status" />
        </div>
        <div class="mt-3 grid grid-cols-2 gap-2 text-center">
          <div class="rounded-lg bg-muted/50 py-1.5">
            <div :class="['font-display text-sm font-bold', s.successRate >= 95 ? 'text-success' : 'text-warning']">{{ s.successRate }}%</div>
            <div class="text-[8px] text-muted-foreground">SUCCESS RATE</div>
          </div>
          <div class="rounded-lg bg-muted/50 py-1.5">
            <div :class="['font-display text-sm font-bold', s.avgLatencyMs <= 300 ? 'text-success' : 'text-warning']">{{ s.avgLatencyMs }}ms</div>
            <div class="text-[8px] text-muted-foreground">LATENSI</div>
          </div>
        </div>
        <div class="mt-2.5 space-y-1 text-[9px] text-muted-foreground">
          <div>Sync terakhir: {{ fmtDateTime(s.lastSyncAt) }}</div>
          <div>FHIR resources: {{ s.fhirResources }} · Coverage {{ s.coveragePct }}% · Uptime {{ s.uptime30d }}%</div>
        </div>
        <p v-if="s.notes" class="mt-2 rounded bg-warning-light/70 px-2 py-1.5 text-[9px] leading-snug text-warning">{{ s.notes }}</p>
      </div>
    </div>

    <SectionCard
      v-if="testResult"
      title="Hasil Test Koneksi — 24 Resource FHIR (TDS v2.0.0)"
      :sub="`Dijalankan ${fmtDateTime(testResult.testedAt)} · ${testResult.summary.ok} OK · ${testResult.summary.perlu} perlu perbaikan · ${testResult.summary.gagal} gagal`"
    >
      <template #icon><Activity class="h-4 w-4" /></template>
      <template #action>
        <Button variant="outline" size="sm" class="h-7 gap-1 text-xs" @click="testResult = null"><RefreshCw class="h-3 w-3" /> Tutup</Button>
      </template>
      <div class="grid gap-1.5 sm:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="r in testResult.results" :key="r.resource"
          :class="['flex items-center justify-between rounded-lg border px-3 py-2 text-xs', r.status === 'OK' ? 'border-success/30 bg-success-light/30' : 'border-warning/40 bg-warning-light/40']"
        >
          <div>
            <span class="font-mono text-[10px] font-semibold">{{ r.resource }}</span>
            <div class="text-[9px] text-muted-foreground">{{ r.message }}</div>
          </div>
          <CheckCircle2 v-if="r.status === 'OK'" class="h-4 w-4 text-success" />
          <ShieldAlert v-else class="h-4 w-4 text-warning" />
        </div>
      </div>
    </SectionCard>

    <SectionCard title="Integration Job Log" sub="Metadata saja — isi respons (body) tidak pernah disimpan (kebijakan privasi)">
      <template #icon><Cable class="h-4 w-4" /></template>
      <div class="max-h-80 overflow-auto rounded-xl border scrollbar-thin">
        <table class="w-full text-xs">
          <thead class="sticky top-0 bg-muted/95 text-left text-[10px] uppercase tracking-wide text-muted-foreground backdrop-blur">
            <tr>
              <th class="px-3 py-2.5">Waktu</th><th class="px-3 py-2.5">Job</th><th class="px-3 py-2.5 text-right">Record</th>
              <th class="px-3 py-2.5 text-right">Durasi</th><th class="px-3 py-2.5">Status & Pesan</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="j in jobs" :key="j.id" class="border-t">
              <td class="px-3 py-2 font-mono text-[10px] text-muted-foreground">{{ fmtDateTime(j.startedAt) }}</td>
              <td class="px-3 py-2">{{ j.jobType }}</td>
              <td class="px-3 py-2 text-right font-mono text-[10px]">{{ fmtID(j.recordsPulled) }}</td>
              <td class="px-3 py-2 text-right font-mono text-[10px]">{{ j.durationMs ? `${(j.durationMs / 1000).toFixed(1)}s` : '—' }}</td>
              <td class="px-3 py-2">
                <div class="flex items-center gap-2">
                  <StatusBadge :status="j.status" />
                  <span v-if="j.message" class="truncate text-[10px] text-muted-foreground" :title="j.message">{{ j.message }}</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </SectionCard>
  </div>
</template>
