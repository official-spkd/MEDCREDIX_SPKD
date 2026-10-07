<script setup lang="ts">
// DashboardView — port dari views-dashboard.tsx (LEVEL 1-6).
import {
  Activity, Radar, AlertTriangle, CheckCircle2, TrendingUp, ArrowRight, ShieldCheck, FileCheck2,
  GitBranch, HeartPulse, Award, Layers,
} from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import { fmtDate, fmtID } from '~/utils/format'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import Skeleton from '~/components/ui/Skeleton.vue'
import ChartLine from '~/components/medcredix/ChartLine.vue'
import KpiCard from '~/components/medcredix/KpiCard.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import ScoreChip from '~/components/medcredix/ScoreChip.vue'
import ScoreRing from '~/components/medcredix/ScoreRing.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const data = ref<any>(null)
const loading = ref(true)

async function load() {
  loading.value = true
  try {
    data.value = await api<any>(`/dashboard${store.hospitalId ? `?hospitalId=${store.hospitalId}` : ''}`)
  } catch { data.value = null }
  loading.value = false
}
onMounted(load)
watch(() => store.hospitalId, load)

const intel = computed(() => data.value?.intel)
const kLabel = computed(() => {
  const hs = intel.value?.healthScore ?? 0
  if (hs >= 90) return 'Sangat Sehat'
  if (hs >= 75) return 'Sehat dengan Area Perbaikan'
  if (hs >= 60) return 'Perlu Perhatian'
  return 'Berisiko Tinggi'
})
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Dashboard Eksekutif"
      sub="Kondisi rumah sakit dalam 30 detik — setiap angka dapat diklik untuk penjelasan, bukti, dan tindakan."
      demo
    >
      <template #actions>
        <Button variant="outline" size="sm" class="gap-1.5" @click="store.navigate('reports')">
          <FileCheck2 class="h-3.5 w-3.5" /> Laporan Kredensial
        </Button>
      </template>
    </PageHeader>

    <template v-if="loading || !data">
      <div class="grid gap-4 lg:grid-cols-3">
        <Skeleton class="h-64" /><Skeleton class="h-64 lg:col-span-2" />
      </div>
      <Skeleton class="h-40" />
    </template>

    <template v-else>
      <!-- LEVEL 1 — Hospital Health Score -->
      <div class="grid gap-4 lg:grid-cols-[340px_1fr]">
        <div class="relative overflow-hidden rounded-2xl bg-gradient-brand p-6 text-white shadow-sm">
          <div class="absolute -right-16 -top-16 h-48 w-48 rounded-full bg-white/5" />
          <div class="relative">
            <div class="flex items-center justify-between">
              <div class="text-[11px] font-semibold uppercase tracking-wider text-white/70">Hospital Health Score</div>
              <Badge class="border-white/20 bg-white/10 text-[10px] text-white">Formula {{ intel.formulaVersion ?? '1.0' }}</Badge>
            </div>
            <div class="mt-4 flex items-center gap-5">
              <ScoreRing :value="intel.healthScore" :size="124" :tone="intel.healthScore >= 75 ? 'info' : intel.healthScore >= 60 ? 'warning' : 'danger'" />
              <div>
                <div class="font-display text-lg font-bold leading-tight">{{ kLabel }}</div>
                <div class="mt-2 space-y-1 text-[11px] text-white/75">
                  <div>
                    Skor Kredensial BPJS:
                    <button class="font-bold text-white underline decoration-white/40 hover:decoration-white" @click="store.openWhy('bcds')">{{ intel.bcds.final }}</button>
                    · Stage {{ intel.bcds.stage }} {{ intel.bcds.lulus ? '· LULUS' : '' }}
                  </div>
                  <div>
                    Data Trust: <button class="font-bold text-white underline decoration-white/40 hover:decoration-white" @click="store.navigate('data-trust')">{{ intel.trust.score }}</button>
                    (multiplier ×{{ intel.trust.multiplier }})
                  </div>
                  <div>CCS {{ intel.ccs.score }} ({{ intel.ccs.grade }}) · CPCS {{ intel.cpcs.score }}</div>
                </div>
              </div>
            </div>
            <div class="mt-5 flex flex-wrap gap-1.5">
              <button
                v-for="k in ['A', 'B', 'C', 'D']" :key="k"
                class="rounded-lg bg-white/10 px-2.5 py-1.5 text-[10px] font-semibold backdrop-blur-sm transition-colors hover:bg-white/20"
                @click="store.openWhy('bcds')"
              >
                Kelompok {{ k }}: <span class="font-display text-xs">{{ intel.kelompok[k] }}</span>
              </button>
            </div>
            <Button size="sm" class="mt-5 w-full gap-1.5 bg-white text-brand hover:bg-white/90" @click="store.openWhy('health')">
              <Activity class="h-3.5 w-3.5" /> Mengapa Skor Ini?
            </Button>
          </div>
        </div>

        <!-- LEVEL 2 — Risiko & Peringatan Dini -->
        <div class="grid gap-4 sm:grid-cols-2">
          <SectionCard
            :title="'Risk Radar'"
            :sub="`${data.riskSummary.tinggi} tinggi · ${data.riskSummary.sedang} sedang · ${data.riskSummary.rendah} rendah`"
          >
            <template #icon><Radar class="h-4 w-4" /></template>
            <template #action>
              <Button variant="ghost" size="sm" class="h-7 gap-1 text-xs text-brand-2" @click="store.navigate('risk-radar')">Semua</Button>
            </template>
            <div class="max-h-64 space-y-2 overflow-y-auto pr-1 scrollbar-thin">
              <button
                v-for="r in data.risks.slice(0, 5)" :key="r.id"
                class="w-full rounded-xl border p-3 text-left transition-colors hover:border-brand-2/50"
                @click="store.navigate('risk-radar')"
              >
                <div class="flex items-start justify-between gap-2">
                  <span class="text-xs font-semibold leading-snug">{{ r.title }}</span>
                  <StatusBadge :status="r.level" :tone="r.level === 'Tinggi' ? 'danger' : r.level === 'Sedang' ? 'warning' : 'info'" />
                </div>
                <div class="mt-1.5 flex items-center gap-2 text-[10px] text-muted-foreground">
                  <span>{{ r.category }}</span><span>· {{ fmtDate(r.detectedAt) }}</span>
                </div>
              </button>
            </div>
          </SectionCard>

          <SectionCard title="Early Warnings" sub="Deteksi dini berbasis ambang batas">
            <template #icon><AlertTriangle class="h-4 w-4" /></template>
            <template #action>
              <Button variant="ghost" size="sm" class="h-7 gap-1 text-xs text-brand-2" @click="store.navigate('early-warnings')">Semua</Button>
            </template>
            <div class="max-h-64 space-y-2 overflow-y-auto pr-1 scrollbar-thin">
              <button
                v-for="w in data.warnings.slice(0, 5)" :key="w.id"
                class="w-full rounded-xl border p-3 text-left transition-colors hover:border-brand-2/50"
                @click="store.navigate('early-warnings')"
              >
                <div class="flex items-start justify-between gap-2">
                  <span class="text-xs font-semibold leading-snug">{{ w.title }}</span>
                  <StatusBadge :status="w.severity" :tone="w.severity === 'Tinggi' ? 'danger' : w.severity === 'Sedang' ? 'warning' : 'info'" />
                </div>
                <div class="mt-1.5 text-[10px] text-muted-foreground">Pemicu: {{ w.thresholdRule }} → {{ w.triggerValue }}</div>
              </button>
            </div>
          </SectionCard>
        </div>
      </div>

      <!-- KPI cards -->
      <div class="grid grid-cols-2 gap-3 md:grid-cols-4">
        <KpiCard label="Skor Kredensial BPJS" :value="intel.bcds.final" unit="/100" tone="info" @click="store.openWhy('bcds')">
          <template #icon><Award class="h-4 w-4" /></template>
          <template #sub>
            <StatusBadge :status="intel.bcds.lulus ? 'LULUS' : 'TIDAK LULUS'" /><span>Stage {{ intel.bcds.stage }}</span>
          </template>
        </KpiCard>
        <KpiCard label="Claim Consistency" :value="intel.ccs.score" :unit="`· ${intel.ccs.grade}`" :tone="intel.ccs.score >= 75 ? 'success' : 'warning'" @click="store.navigate('claim-integrity')">
          <template #icon><FileCheck2 class="h-4 w-4" /></template>
          <template #sub><span>{{ intel.ccs.flaggedClaims }}/{{ intel.ccs.totalClaims }} klaim perlu tinjauan</span></template>
        </KpiCard>
        <KpiCard label="Kepatuhan Pathway" :value="intel.cpcs.score" unit="/100" :tone="intel.cpcs.score >= 75 ? 'success' : 'warning'" @click="store.openWhy('pathway')">
          <template #icon><GitBranch class="h-4 w-4" /></template>
          <template #sub><span>Adherence {{ intel.cpcs.adherence }}%</span></template>
        </KpiCard>
        <KpiCard label="Keselamatan Pasien" :value="intel.safety.score" unit="/100" :tone="intel.safety.score >= 85 ? 'success' : 'warning'" @click="store.openWhy('safety')">
          <template #icon><HeartPulse class="h-4 w-4" /></template>
          <template #sub><span>Alergi {{ intel.safety.allergy }}% · CL {{ intel.safety.closedLoop }}%</span></template>
        </KpiCard>
      </div>

      <!-- LEVEL 3 — Aksi Rekomendasi -->
      <div class="grid gap-4 lg:grid-cols-[1.2fr_1fr]">
        <SectionCard title="Tindakan Prioritas" :sub="`${data.topActions.length} aksi aktif — closed-loop hingga skor membaik`">
          <template #icon><CheckCircle2 class="h-4 w-4" /></template>
          <template #action>
            <Button variant="ghost" size="sm" class="h-7 gap-1 text-xs text-brand-2" @click="store.navigate('actions')">Action Center</Button>
          </template>
          <div class="space-y-2">
            <div
              v-for="a in data.topActions" :key="a.id"
              :class="[
                'flex items-start justify-between gap-3 rounded-xl border border-l-4 p-3',
                a.priority === 'Tinggi' ? 'border-l-danger bg-danger-light/25' : a.priority === 'Sedang' ? 'border-l-warning bg-warning-light/25' : 'border-l-info bg-info-light/20',
              ]"
            >
              <div class="min-w-0">
                <div class="flex flex-wrap items-center gap-1.5">
                  <StatusBadge :status="a.priority" :tone="a.priority === 'Tinggi' ? 'danger' : a.priority === 'Sedang' ? 'warning' : 'info'" />
                  <span class="text-[10px] text-muted-foreground">{{ a.code }} · {{ a.owner.split(',')[0] }}</span>
                </div>
                <div class="mt-1 truncate text-xs font-semibold">{{ a.title }}</div>
                <div class="mt-0.5 text-[10px] text-muted-foreground">
                  Tenggat {{ fmtDate(a.dueDate) }} · dari {{ a.sourceType }}{{ a.sourceRef ? ` (${a.sourceRef})` : '' }}
                </div>
              </div>
              <div class="flex shrink-0 flex-col items-end gap-1.5">
                <StatusBadge :status="a.status" />
                <Button size="sm" variant="outline" class="h-6 px-2 text-[10px]" @click="store.navigate('actions')">Detail</Button>
              </div>
            </div>
          </div>
        </SectionCard>

        <!-- LEVEL 4 — Rincian skor domain -->
        <SectionCard title="Rincian Skor — 9 Domain" sub="Klik domain untuk melihat mengapa">
          <template #icon><Layers class="h-4 w-4" /></template>
          <div class="grid gap-2 sm:grid-cols-2">
            <ScoreChip
              v-for="d in intel.domains" :key="d.key" :label="d.label" :score="d.score"
              :group="d.group.replace(/^[A-D]\. /, '')" @click="store.openWhy(d.key)"
            />
          </div>
        </SectionCard>
      </div>

      <!-- LEVEL 5-6 — Bukti & Tren -->
      <div class="grid gap-4 lg:grid-cols-[1.4fr_1fr]">
        <SectionCard title="Tren Perbaikan 6 Bulan" sub="Hospital Health Score & Skor Kredensial BPJS">
          <template #icon><TrendingUp class="h-4 w-4" /></template>
          <ChartLine
            :labels="data.intel.trend.map((t: any) => t.period)"
            :datasets="[
              { label: 'Health Score', data: data.intel.trend.map((t: any) => t.health), color: '#14285C', fill: true, width: 3 },
              { label: 'BDCS', data: data.intel.trend.map((t: any) => t.bcds), color: '#2198A5', fill: true, width: 3 },
              { label: 'Trust', data: data.intel.trend.map((t: any) => t.trust), color: '#D97706', dashed: true, width: 2 },
            ]"
            :y-min="50" :y-max="100" :height="208" :legend="true"
          />
          <div class="mt-3 rounded-xl border bg-brand-2-light/30 p-3">
            <div class="flex items-center gap-1.5 text-[11px] font-semibold text-brand-2-dark">
              <TrendingUp class="h-3.5 w-3.5" /> Timeline Perbaikan
            </div>
            <div class="mt-2 max-h-32 space-y-1.5 overflow-y-auto pr-1 scrollbar-thin">
              <div v-for="(t, i) in data.timeline.slice(-8)" :key="i" class="flex items-start gap-2 text-[11px]">
                <span :class="['mt-1 h-1.5 w-1.5 shrink-0 rounded-full', t.type === 'Skor' ? 'bg-brand' : 'bg-success']" />
                <span class="w-14 shrink-0 font-semibold text-muted-foreground">{{ t.period }}</span>
                <span class="min-w-0 flex-1 truncate text-foreground/90" :title="t.label">{{ t.label }}</span>
              </div>
            </div>
          </div>
        </SectionCard>

        <div class="space-y-4">
          <SectionCard title="Bukti Terbaru" sub="Setiap metrik punya jejak data">
            <template #icon><ShieldCheck class="h-4 w-4" /></template>
            <template #action>
              <Button variant="ghost" size="sm" class="h-7 gap-1 text-xs text-brand-2" @click="store.navigate('evidence')">Semua</Button>
            </template>
            <div class="space-y-2">
              <div v-for="e in data.evidence.slice(0, 4)" :key="e.id" class="rounded-xl border p-3">
                <div class="flex items-start justify-between gap-2">
                  <span class="text-[11px] font-semibold leading-snug">{{ e.title }}</span>
                  <StatusBadge :status="e.validationStatus" />
                </div>
                <div class="mt-1 text-[10px] text-muted-foreground">
                  {{ e.metric }}: <strong class="text-foreground">{{ e.value }}{{ e.unit === '%' ? '%' : '' }}</strong> · {{ e.sourceSystem }} · {{ fmtID(e.recordCount ?? 0) }} record
                </div>
              </div>
            </div>
          </SectionCard>

          <SectionCard title="Clinical Pathway Terlemah" sub="Prioritas audit komite medis">
            <template #icon><GitBranch class="h-4 w-4" /></template>
            <div class="space-y-2">
              <button
                v-for="p in data.pathwayWorst" :key="p.id"
                class="flex w-full items-center justify-between gap-2 rounded-xl border p-3 text-left transition-colors hover:border-brand-2/50"
                @click="store.navigate('clin-pathway')"
              >
                <div class="min-w-0">
                  <div class="truncate text-xs font-semibold">{{ p.diagnosis }}</div>
                  <div class="text-[10px] text-muted-foreground">LOS standar {{ p.standardLos }} hr · cost ratio {{ p.costRatio }}</div>
                </div>
                <span :class="['font-display text-lg font-bold', p.avgAdherence < 70 ? 'text-danger' : p.avgAdherence < 80 ? 'text-warning' : 'text-success']">{{ p.avgAdherence }}%</span>
              </button>
            </div>
          </SectionCard>
        </div>
      </div>

      <!-- Prasyarat -->
      <SectionCard
        title="Prasyarat Wajib Kredensial BPJS"
        sub="P1 Izin Operasional · P2 PKS BPJS · P3 Akreditasi — wajib Valid untuk skor penuh"
      >
        <template #icon><ShieldCheck class="h-4 w-4" /></template>
        <template #action>
          <Button variant="outline" size="sm" class="h-7 gap-1 text-xs" @click="store.navigate('cred-prerequisites')">Kelola <ArrowRight class="h-3.5 w-3.5" /></Button>
        </template>
        <div class="grid gap-2 sm:grid-cols-3">
          <div v-for="p in data.intel.prereq" :key="p.code" class="rounded-xl border p-3.5">
            <div class="flex items-center justify-between gap-2">
              <span class="text-[10px] font-bold uppercase tracking-wide text-muted-foreground">{{ p.code }}</span>
              <StatusBadge :status="p.status" />
            </div>
            <div class="mt-1.5 text-xs font-semibold leading-snug">{{ p.name }}</div>
            <div class="mt-1 text-[10px] text-muted-foreground">
              {{ p.daysLeft != null ? (p.daysLeft > 0 ? `Berlaku ${p.daysLeft} hari lagi` : `Kedaluwarsa ${Math.abs(p.daysLeft)} hari lalu`) : 'Tanpa masa berlaku' }}
            </div>
          </div>
        </div>
      </SectionCard>
    </template>
  </div>
</template>
