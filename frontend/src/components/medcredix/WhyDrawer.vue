<script setup lang="ts">
// WhyDrawer — "Mengapa Skor Ini?" drilldown universal:
// Kalkulasi → Kontributor → Bukti → Risiko → Rekomendasi → Tren.
import {
  Calculator, CheckCircle2, AlertTriangle, Fingerprint, Radar, Lightbulb, TrendingUp, ArrowRight, FileSearch,
} from 'lucide-vue-next'
import { AI_DISCLAIMER, fmtDate, fmtDateTime } from '~/utils/format'
import { useMedCredix } from '~/stores/medcredix'
import { computed, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import Button from '~/components/ui/Button.vue'
import Sheet from '~/components/ui/Sheet.vue'
import Skeleton from '~/components/ui/Skeleton.vue'
import Tabs from '~/components/ui/Tabs.vue'
import ChartLine from '~/components/medcredix/ChartLine.vue'
import Hint from '~/components/medcredix/Hint.vue'
import ScoreRing from '~/components/medcredix/ScoreRing.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'
import StatusIcon from '~/components/medcredix/StatusIcon.vue'

const store = useMedCredix()
const data = ref<any>(null)
const loading = ref(false)
const tab = ref('kalkulasi')
const toneHex: Record<string, string> = {
  success: '#059669', warning: '#D97706', danger: '#DC2626', info: '#0284C7', muted: '#5B6B84',
}

watch(
  () => store.whyDomain,
  async (domain) => {
    data.value = null
    tab.value = 'kalkulasi'
    if (!domain) return
    loading.value = true
    try {
      data.value = await api<any>(`/scores?domain=${domain}`)
    } catch (e: any) {
      data.value = { __error: e?.message ?? 'Gagal memuat' }
    }
    loading.value = false
  },
  { immediate: true },
)

const title = computed(() =>
  data.value?.domain?.label ??
  (store.whyDomain === 'bcds' ? 'Skor Kredensial BPJS (BDCS)' : store.whyDomain === 'health' ? 'Hospital Health Score' : 'Skor'),
)
const tabs = [
  { v: 'kalkulasi', label: 'Kalkulasi', icon: Calculator },
  { v: 'kontributor', label: 'Kontributor', icon: CheckCircle2 },
  { v: 'bukti', label: 'Bukti', icon: Fingerprint },
  { v: 'risiko', label: 'Risiko', icon: Radar },
  { v: 'rekomendasi', label: 'Rekomendasi', icon: Lightbulb },
  { v: 'tren', label: 'Tren', icon: TrendingUp },
]
</script>

<template>
  <Sheet :open="!!store.whyDomain" @update:open="(o) => !o && store.closeWhy()">
    <div class="overflow-y-auto scrollbar-thin">
      <!-- Header gradien -->
      <div class="bg-gradient-brand px-6 py-5 text-white">
        <div class="flex flex-wrap items-center gap-4">
          <ScoreRing :value="data?.domain?.score ?? data?.bcds?.final ?? data?.healthScore ?? 0" :size="92" :stroke="8" tone="info" />
          <div>
            <div class="text-[11px] font-semibold uppercase tracking-wider text-white/60">
              Mengapa Skor Ini? · Formula {{ data?.formulaVersion ?? '—' }}
            </div>
            <div class="font-display text-lg font-bold leading-tight">{{ title }}</div>
            <div class="mt-1 flex flex-wrap items-center gap-2 text-xs text-white/80">
              <span v-if="data?.domain?.group">{{ data.domain.group }}</span>
              <template v-if="data?.bcds && store.whyDomain === 'bcds'">
                <span class="rounded-full bg-white/15 px-2 py-0.5 font-semibold">Stage {{ data.bcds.stage }}</span>
                <span :class="['rounded-full px-2 py-0.5 font-semibold', data.bcds.lulus ? 'bg-success text-white' : 'bg-danger text-white']">
                  {{ data.bcds.lulus ? 'LULUS' : 'TIDAK LULUS' }}
                </span>
                <span class="rounded-full bg-white/15 px-2 py-0.5">Trust {{ data.bcds.trust }} → ×{{ data.bcds.multiplier }}</span>
              </template>
            </div>
          </div>
        </div>
      </div>

      <div class="px-4 py-4 sm:px-6">
        <div v-if="loading || !data" class="space-y-3">
          <Skeleton class="h-8 w-full" /><Skeleton class="h-32 w-full" /><Skeleton class="h-32 w-full" />
        </div>
        <div v-else-if="data.__error" class="rounded-xl border border-danger/40 bg-danger-light p-4 text-xs text-danger">
          {{ data.__error }}
        </div>
        <div v-else>
          <Tabs v-model="tab">
            <div class="grid h-auto w-full grid-cols-3 gap-1 rounded-lg bg-muted p-1 sm:grid-cols-6">
              <button
                v-for="t in tabs" :key="t.v" type="button"
                :class="[
                  'inline-flex items-center justify-center gap-1.5 rounded-md px-2 py-1.5 text-[11px] font-medium transition-all',
                  tab === t.v ? 'bg-card text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground',
                ]"
                @click="tab = t.v"
              >
                <component :is="t.icon" class="h-3.5 w-3.5" />{{ t.label }}
              </button>
            </div>

            <!-- KALKULASI -->
            <div v-if="tab === 'kalkulasi'" class="mt-4 space-y-4">
              <div class="rounded-xl border bg-brand-light/50 p-4">
                <div class="mb-1 flex items-center gap-1.5 text-xs font-semibold text-brand"><Calculator class="h-3.5 w-3.5" /> Rumus</div>
                <p class="text-xs leading-relaxed text-foreground/90">{{ data.domain?.formula ?? data.formula }}</p>
              </div>

              <div v-if="data.kind === 'health' && data.domains" class="overflow-hidden rounded-xl border">
                <table class="w-full text-xs">
                  <thead class="bg-muted/60 text-left text-[11px] uppercase tracking-wide text-muted-foreground">
                    <tr><th class="px-3 py-2">Domain</th><th class="px-3 py-2 text-right">Skor</th><th class="px-3 py-2 text-right">Bobot</th><th class="px-3 py-2 text-right">Kontribusi</th></tr>
                  </thead>
                  <tbody>
                    <tr v-for="d in data.domains" :key="d.key" class="border-t">
                      <td class="px-3 py-2 font-medium">{{ d.label }}</td>
                      <td class="px-3 py-2 text-right font-display font-bold" :style="{ color: toneHex[d.score >= 90 ? 'success' : d.score >= 75 ? 'info' : d.score >= 60 ? 'warning' : 'danger'] }">{{ d.score }}</td>
                      <td class="px-3 py-2 text-right text-muted-foreground">{{ Math.round(d.weight * 100) }}%</td>
                      <td class="px-3 py-2 text-right text-muted-foreground">{{ d.contribution }}</td>
                    </tr>
                    <tr class="border-t bg-brand-light/50 font-semibold">
                      <td class="px-3 py-2">Hospital Health Score</td>
                      <td class="px-3 py-2 text-right font-display text-base font-bold text-brand">{{ data.healthScore }}</td>
                      <td class="px-3 py-2 text-right">100%</td>
                      <td class="px-3 py-2 text-right">{{ data.healthScore }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <div v-if="data.kind === 'bcds'" class="space-y-4">
                <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
                  <div v-for="k in ['A', 'B', 'C', 'D']" :key="k" class="rounded-xl border bg-card p-3 text-center">
                    <div class="text-[10px] font-semibold uppercase tracking-wide text-muted-foreground">Kelompok {{ k }}</div>
                    <div class="mt-1 font-display text-xl font-bold text-brand">{{ data.kelompok?.[k] }}</div>
                  </div>
                </div>
                <div class="rounded-xl border p-4">
                  <div class="mb-2 text-xs font-semibold">Stage & Ambang Kelulusan</div>
                  <div class="grid grid-cols-4 gap-1.5 text-center text-[10px] sm:grid-cols-8">
                    <div
                      v-for="s in data.stageTable ?? []" :key="s.stage"
                      :class="['rounded-lg border px-1.5 py-1.5', s.stage >= 5 ? 'border-success/40 bg-success-light font-bold text-success' : 'border-border text-muted-foreground']"
                    >
                      <div class="font-display text-xs">S{{ s.stage }}</div>
                      <div>≥{{ s.min }}</div>
                    </div>
                  </div>
                  <p class="mt-2.5 text-[11px] leading-relaxed text-muted-foreground">
                    Skor mentah {{ data.bcds?.raw }} × Trust Multiplier {{ data.bcds?.multiplier }} (Trust {{ data.bcds?.trust }}) =
                    <strong class="text-foreground">{{ data.bcds?.final }}</strong>.
                    Kelulusan kredensial digital pada Stage ≥ 5 (skor ≥ 66).
                    <span v-if="!data.prereqOk">Catatan: salah satu prasyarat wajib (P1–P3) belum Valid — Stage efektif 0.</span>
                  </p>
                </div>
                <div class="rounded-xl border p-4">
                  <div class="mb-2 flex items-center gap-1.5 text-xs font-semibold">
                    Transparansi Trust Multiplier
                    <Hint>Sistem produk lama memakai multiplier tersembunyi; MedCredix menampilkannya terbuka.</Hint>
                  </div>
                  <div class="space-y-1">
                    <div
                      v-for="(b, i) in data.trustBands ?? []" :key="i"
                      :class="['flex items-center justify-between rounded-lg px-3 py-1.5 text-xs', data.trust?.score >= b.minTrust ? 'bg-brand-2-light font-semibold text-brand-2-dark' : 'text-muted-foreground']"
                    >
                      <span>Trust ≥ {{ b.minTrust }}</span><span>{{ b.label }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <div v-if="data.domain" class="overflow-hidden rounded-xl border">
                <div class="border-b bg-muted/60 px-3 py-2 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground">
                  Indikator Penyusun — nilai · ambang · status · sumber data
                </div>
                <table class="w-full text-xs">
                  <tbody>
                    <tr v-for="(i, idx) in data.domain.indicators" :key="idx" class="border-b last:border-0">
                      <td class="px-3 py-2.5">
                        <div class="flex items-center gap-1 font-medium">
                          {{ i.label }}
                          <Hint>Bobot {{ Math.round(i.weight * 100) }}% — {{ i.note ?? i.source }}{{ i.fhir ? ` · FHIR: ${i.fhir}` : '' }}</Hint>
                        </div>
                        <div v-if="i.note" class="mt-0.5 text-[10px] text-muted-foreground">{{ i.note }}</div>
                      </td>
                      <td class="whitespace-nowrap px-3 py-2.5 text-right">
                        <span class="font-display font-bold">{{ i.value }}{{ i.unit === '%' ? '%' : ` ${i.unit}` }}</span>
                        <div class="text-[10px] text-muted-foreground">ambang {{ i.threshold }}{{ i.unit === '%' ? '%' : ` ${i.unit}` }}</div>
                      </td>
                      <td class="px-3 py-2.5 text-right"><StatusBadge :status="i.status" /></td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- KONTRIBUTOR -->
            <div v-else-if="tab === 'kontributor'" class="mt-4 space-y-4">
              <div class="rounded-xl border border-success/30 bg-success-light/40 p-4">
                <div class="mb-2 flex items-center gap-1.5 text-xs font-semibold text-success"><CheckCircle2 class="h-4 w-4" /> Kontributor Positif</div>
                <ul class="space-y-1.5">
                  <li v-for="(p, i) in data.domain?.positives ?? []" :key="i" class="flex items-start gap-2 text-xs text-foreground/90">
                    <StatusIcon tone="success" class="mt-0.5 h-3.5 w-3.5" />{{ p }}
                  </li>
                  <li v-if="!(data.domain?.positives?.length)" class="text-xs text-muted-foreground">Belum ada kontributor positif yang terdeteksi.</li>
                </ul>
              </div>
              <div class="rounded-xl border border-warning/30 bg-warning-light/40 p-4">
                <div class="mb-2 flex items-center gap-1.5 text-xs font-semibold text-warning"><AlertTriangle class="h-4 w-4" /> Kontributor Negatif (Penurun Skor)</div>
                <ul class="space-y-1.5">
                  <li v-for="(p, i) in data.domain?.negatives ?? []" :key="i" class="flex items-start gap-2 text-xs text-foreground/90">
                    <StatusIcon tone="warning" class="mt-0.5 h-3.5 w-3.5" />{{ p }}
                  </li>
                  <li v-if="!(data.domain?.negatives?.length)" class="text-xs text-muted-foreground">Tidak ada kontributor negatif — pertahankan kinerja.</li>
                </ul>
              </div>
            </div>

            <!-- BUKTI -->
            <div v-else-if="tab === 'bukti'" class="mt-4 space-y-3">
              <div class="flex items-center justify-between">
                <p class="text-xs text-muted-foreground">Setiap angka dapat ditelusuri ke data sumber. Tidak pernah ada bukti palsu.</p>
                <Button variant="outline" size="sm" class="h-7 gap-1.5 text-xs" @click="store.closeWhy(); store.navigate('evidence')">
                  <FileSearch class="h-3.5 w-3.5" /> Evidence Center
                </Button>
              </div>
              <div v-if="(data.evidence ?? []).length === 0" class="rounded-xl border border-dashed p-6 text-center text-xs text-muted-foreground">
                Belum ada item bukti spesifik terikat pada skor ini. Buka Evidence Center untuk melihat seluruh bukti rumah sakit, atau buka risiko terkait untuk melihat bukti per-risiko.
              </div>
              <div v-for="e in data.evidence ?? []" v-else :key="e.id" class="rounded-xl border bg-card p-3.5">
                <div class="flex items-start justify-between gap-2">
                  <div>
                    <div class="text-xs font-semibold">{{ e.title }}</div>
                    <div class="mt-0.5 text-[11px] text-muted-foreground">
                      {{ e.metric }}: <strong class="text-foreground">{{ e.value }}{{ e.unit === '%' ? '%' : ` ${e.unit ?? ''}` }}</strong> · {{ e.recordCount }} record
                    </div>
                  </div>
                  <StatusBadge :status="e.validationStatus" />
                </div>
                <div class="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-[10px] text-muted-foreground">
                  <span>Sumber: {{ e.sourceSystem }}</span>
                  <span v-if="e.fhirResource">· FHIR: {{ e.fhirResource }}</span>
                  <span>· Sinkron: {{ fmtDateTime(e.lastSyncAt) }}</span>
                </div>
              </div>
            </div>

            <!-- RISIKO -->
            <div v-else-if="tab === 'risiko'" class="mt-4 space-y-3">
              <div v-if="(data.risks ?? []).length === 0" class="rounded-xl border border-dashed p-6 text-center text-xs text-muted-foreground">
                Tidak ada risiko terbuka yang terkait skor ini.
              </div>
              <div v-for="r in data.risks ?? []" v-else :key="r.id" class="rounded-xl border bg-card p-3.5">
                <div class="flex items-start justify-between gap-2">
                  <div class="text-xs font-semibold leading-snug">{{ r.title }}</div>
                  <StatusBadge :status="r.level" :tone="r.level === 'Tinggi' ? 'danger' : r.level === 'Sedang' ? 'warning' : 'info'" />
                </div>
                <p class="mt-1 line-clamp-2 text-[11px] leading-relaxed text-muted-foreground">{{ r.description }}</p>
                <div class="mt-2 flex flex-wrap items-center gap-2 text-[10px] text-muted-foreground">
                  <StatusBadge :status="r.status" />
                  <span>· {{ r.category }}</span><span>· Tenggat: {{ fmtDate(r.deadline) }}</span>
                </div>
              </div>
              <Button variant="outline" size="sm" class="w-full gap-1.5" @click="store.closeWhy(); store.navigate('risk-radar')">
                Buka Risk Radar <ArrowRight class="h-3.5 w-3.5" />
              </Button>
            </div>

            <!-- REKOMENDASI -->
            <div v-else-if="tab === 'rekomendasi'" class="mt-4 space-y-3">
              <div v-if="(data.recommendations ?? []).length === 0" class="rounded-xl border border-dashed p-6 text-center text-xs text-muted-foreground">
                Belum ada rekomendasi spesifik — pertahankan kinerja domain ini.
              </div>
              <div v-for="(r, i) in data.recommendations ?? []" v-else :key="i" class="rounded-xl border bg-brand-2-light/30 p-4">
                <div class="flex items-start justify-between gap-3">
                  <div class="min-w-0">
                    <div class="flex items-center gap-2">
                      <StatusBadge :status="r.priority" :tone="r.priority === 'Tinggi' ? 'danger' : r.priority === 'Sedang' ? 'warning' : 'info'" />
                      <span class="text-[10px] text-muted-foreground">Saran owner: {{ r.suggestedOwner }}</span>
                    </div>
                    <p class="mt-1.5 text-xs font-medium leading-snug">{{ r.title }}</p>
                  </div>
                  <Button
                    size="sm" class="h-7 shrink-0 gap-1 bg-brand-2 text-xs hover:bg-brand-2-hover"
                    @click="store.openAction({ title: r.title, sourceType: 'Skor', sourceRef: data.domain ? `DOMAIN:${data.domain.key}` : 'DOMAIN:health', priority: r.priority, suggestedOwner: r.suggestedOwner, description: `Rekomendasi dari analisis skor ${title}.` })"
                  >
                    Buat Tindakan <ArrowRight class="h-3 w-3" />
                  </Button>
                </div>
              </div>
              <p class="rounded-lg bg-info-light px-3 py-2 text-[10px] leading-relaxed text-info">{{ AI_DISCLAIMER }}</p>
            </div>

            <!-- TREN -->
            <div v-else-if="tab === 'tren'" class="mt-4 space-y-4">
              <div class="rounded-xl border bg-card p-4">
                <div class="mb-3 text-xs font-semibold">Tren 6 Bulan — Health Score & BDCS</div>
                <ChartLine
                  v-if="data.trend?.length"
                  :labels="data.trend.map((t: any) => t.period)"
                  :datasets="[
                    { label: 'Health Score', data: data.trend.map((t: any) => t.health), color: '#14285C', width: 2.5 },
                    { label: 'BDCS', data: data.trend.map((t: any) => t.bcds), color: '#2198A5', width: 2.5 },
                  ]"
                  :y-min="40" :y-max="100" :height="176" :legend="true"
                />
              </div>
              <div v-if="(data.actions ?? []).length > 0" class="rounded-xl border p-4">
                <div class="mb-2 text-xs font-semibold">Aksi Terkait</div>
                <div class="space-y-1.5">
                  <div v-for="a in data.actions.slice(0, 6)" :key="a.id" class="flex items-center justify-between gap-2 rounded-lg bg-muted/40 px-3 py-2 text-xs">
                    <span class="min-w-0 truncate font-medium">{{ a.code }} — {{ a.title }}</span>
                    <StatusBadge :status="a.status" />
                  </div>
                </div>
                <Button variant="outline" size="sm" class="mt-3 w-full gap-1.5" @click="store.closeWhy(); store.navigate('actions')">
                  Buka Action Center <ArrowRight class="h-3.5 w-3.5" />
                </Button>
              </div>
            </div>
          </Tabs>
        </div>
      </div>
    </div>
  </Sheet>
</template>
