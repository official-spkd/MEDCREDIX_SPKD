<script setup lang="ts">
// ClaimIntegrityView — port dari views-claims.tsx (ClaimIntegrityView).
import { FileCheck2, ScanSearch, Fingerprint, PlusCircle, CheckCircle2, Search } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import { fmtDate, fmtRupiah } from '~/utils/format'
import Button from '~/components/ui/Button.vue'
import Dialog from '~/components/ui/Dialog.vue'
import Input from '~/components/ui/Input.vue'
import ChartRadar from '~/components/medcredix/ChartRadar.vue'
import DemoBadge from '~/components/medcredix/DemoBadge.vue'
import EmptyState from '~/components/medcredix/EmptyState.vue'
import FlagIcon from '~/components/medcredix/FlagIcon.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import ScoreRing from '~/components/medcredix/ScoreRing.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const selected = ref<string | null>(null)
const q = ref('')
const data = ref<any>(null)
const detail = ref<any>(null)

async function load() {
  try { data.value = await api(`/claims${q.value ? `?q=${encodeURIComponent(q.value)}` : ''}`) } catch { data.value = null }
}
onMounted(load)
let t: ReturnType<typeof setTimeout> | null = null
watch(q, () => { if (t) clearTimeout(t); t = setTimeout(load, 300) })

async function loadDetail() {
  if (!selected.value) { detail.value = null; return }
  try { detail.value = await api(`/claims?id=${selected.value}`) } catch { detail.value = null }
}
watch(selected, loadDetail)

const claims = computed(() => data.value?.claims ?? [])
const ccs = computed(() => data.value?.ccs)
const CAT_LABELS: Record<string, string> = { Diagnosis: 'Diagnosis', Medication: 'Medication', Procedure: 'Procedure', Diagnostic: 'Diagnostic', Stock: 'Stock', Timeline: 'Timeline' }
const CAT_KEYS = Object.keys(CAT_LABELS)
function parseFlags(json: string | null) {
  try { return json ? JSON.parse(json) : {} } catch { return {} }
}
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Claim Integrity"
      sub="Mengukur konsistensi klaim BPJS vs rekam medis elektronik (EHR sebagai sumber kebenaran). Sinyal bersifat netral — verifikasi manual wajib." demo
    >
      <template #actions>
        <Button variant="outline" size="sm" class="gap-1.5" @click="store.navigate('claim-anomalies')">
          <ScanSearch class="h-3.5 w-3.5" /> Daftar Anomali
        </Button>
      </template>
    </PageHeader>

    <div class="grid gap-4 lg:grid-cols-[240px_1fr_1fr]">
      <div class="rounded-2xl border bg-card p-5 text-center">
        <ScoreRing :value="ccs?.score ?? 0" :size="116" @click="store.openWhy('claim')" />
        <div class="mt-3">
          <StatusBadge :status="ccs?.grade ?? '—'" :tone="ccs?.grade === 'EXCELLENT' || ccs?.grade === 'GOOD' ? 'success' : ccs?.grade === 'FAIR' ? 'warning' : 'danger'" />
        </div>
        <Button size="sm" variant="outline" class="mt-3 w-full text-xs" @click="store.openWhy('claim')">Mengapa Skor Ini?</Button>
      </div>
      <SectionCard title="6 Sub-Skor Konsistensi" sub="Bobot: Diagnosis 25% · Medication 20% · Procedure 15% · Diagnostic 15% · Stock 15% · Timeline 10%">
        <ChartRadar
          :labels="(ccs?.subs ?? []).map((s: any) => s.label.split(' ')[0])"
          :datasets="[{ label: 'Sub-skor', data: (ccs?.subs ?? []).map((s: any) => s.score), color: '#2198A5' }]"
          :height="208"
        />
      </SectionCard>
      <SectionCard title="Distribusi Sinyal" :sub="`${ccs?.flaggedClaims ?? 0} dari ${ccs?.totalClaims ?? 0} klaim memiliki minimal satu sinyal`">
        <div class="space-y-2">
          <div v-for="s in ccs?.subs ?? []" :key="s.key" class="flex items-center gap-2.5">
            <span class="w-20 shrink-0 text-[10px] font-semibold text-muted-foreground">{{ CAT_LABELS[s.key] }}</span>
            <div class="flex h-6 flex-1 overflow-hidden rounded-md" :title="`${s.ok} ok · ${s.warn} perlu tinjauan · ${s.bad} tinggi`">
              <div class="bg-success" :style="{ width: `${(s.ok / (s.ok + s.warn + s.bad)) * 100}%` }" />
              <div class="bg-warning" :style="{ width: `${(s.warn / (s.ok + s.warn + s.bad)) * 100}%` }" />
              <div class="bg-danger" :style="{ width: `${(s.bad / (s.ok + s.warn + s.bad)) * 100}%` }" />
            </div>
            <span :class="['w-12 shrink-0 text-right font-display text-xs font-bold', s.score < 75 ? 'text-warning' : 'text-success']">{{ s.score }}%</span>
          </div>
          <div class="flex gap-3 pt-1 text-[9px] text-muted-foreground">
            <span class="flex items-center gap-1"><span class="h-2 w-2 rounded-sm bg-success" /> Konsisten</span>
            <span class="flex items-center gap-1"><span class="h-2 w-2 rounded-sm bg-warning" /> Perlu tinjauan</span>
            <span class="flex items-center gap-1"><span class="h-2 w-2 rounded-sm bg-danger" /> Sinyal tinggi</span>
          </div>
        </div>
      </SectionCard>
    </div>

    <SectionCard title="Detail per Klaim" sub="Klik klaim untuk menelusuri sinyal ke bukti & rekomendasi penyelesaian">
      <template #icon><FileCheck2 class="h-4 w-4" /></template>
      <template #action>
        <div class="relative">
          <Search class="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
          <Input v-model="q" placeholder="Cari nomor klaim/SEP/INA-CBG…" class="h-8 w-60 pl-8 text-xs" />
        </div>
      </template>
      <div class="max-h-[420px] overflow-auto rounded-xl border scrollbar-thin">
        <table class="w-full text-xs">
          <thead class="sticky top-0 bg-muted/95 text-left text-[10px] uppercase tracking-wide text-muted-foreground backdrop-blur">
            <tr>
              <th class="px-3 py-2.5">Klaim</th><th class="px-3 py-2.5">INA-CBGs</th>
              <th class="px-3 py-2.5">Nilai</th><th class="px-3 py-2.5">Konsistensi</th><th class="px-3 py-2.5">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="c in claims" :key="c.id"
              :class="['cursor-pointer border-t transition-colors hover:bg-muted/30', Object.values(parseFlags(c.flagsJson)).includes('bad') && 'bg-danger-light/20']"
              @click="selected = c.id"
            >
              <td class="px-3 py-2.5"><div class="font-mono text-[10px] font-semibold">{{ c.claimNumber }}</div><div class="text-[9px] text-muted-foreground">SEP {{ c.sepNumber }} · {{ c.patientInit }}</div></td>
              <td class="px-3 py-2.5"><div class="font-mono text-[10px]">{{ c.inacbgCode }}</div><div class="max-w-36 truncate text-[9px] text-muted-foreground">{{ c.inacbgDesc }}</div></td>
              <td class="px-3 py-2.5 font-mono text-[10px]">{{ fmtRupiah(c.billedAmount) }}</td>
              <td class="px-3 py-2.5"><div class="flex gap-1"><FlagIcon v-for="k in CAT_KEYS" :key="k" :flag="parseFlags(c.flagsJson)[k] ?? 'ok'" /></div></td>
              <td class="px-3 py-2.5"><StatusBadge :status="c.status" /></td>
            </tr>
          </tbody>
        </table>
      </div>
    </SectionCard>

    <!-- Detail drawer -->
    <Dialog :open="!!selected" @update:open="(o) => !o && (selected = null)">
      <template v-if="detail">
        <div class="flex flex-wrap items-center gap-2 text-[11px] uppercase tracking-wider text-muted-foreground">
          {{ detail.claim.careType }} · SEP {{ detail.claim.sepNumber }} · Pasien {{ detail.claim.patientInit }} (inisial saja) <DemoBadge />
        </div>
        <h2 class="mt-1 flex flex-wrap items-center gap-2.5 text-left text-lg font-bold">
          {{ detail.claim.claimNumber }} <StatusBadge :status="detail.claim.status" />
        </h2>
        <div class="mt-4 space-y-4">
          <div class="grid gap-2 text-xs sm:grid-cols-2">
            <div class="rounded-lg bg-muted/40 p-3"><span class="text-muted-foreground">INA-CBGs:</span> <strong>{{ detail.claim.inacbgCode }}</strong> — {{ detail.claim.inacbgDesc }}</div>
            <div class="rounded-lg bg-muted/40 p-3"><span class="text-muted-foreground">Nilai diklaim:</span> <strong>{{ fmtRupiah(detail.claim.billedAmount) }}</strong></div>
            <div class="rounded-lg bg-muted/40 p-3"><span class="text-muted-foreground">Masuk:</span> {{ fmtDate(detail.claim.admissionDate) }}</div>
            <div class="rounded-lg bg-muted/40 p-3"><span class="text-muted-foreground">Keluar:</span> {{ fmtDate(detail.claim.dischargeDate) }}</div>
          </div>

          <SectionCard title="Matriks Konsistensi" sub="✓ konsisten · ! perlu ditinjau · ✕ sinyal tinggi — EHR sebagai sumber kebenaran">
            <div class="grid grid-cols-3 gap-1.5 sm:grid-cols-6">
              <div
                v-for="(label, k) in CAT_LABELS" :key="k"
                :class="['flex flex-col items-center gap-1 rounded-lg border px-2 py-2',
                  parseFlags(detail.claim.flagsJson)[k] === 'ok' ? 'border-success/30 bg-success-light/40' : parseFlags(detail.claim.flagsJson)[k] === 'warn' ? 'border-warning/30 bg-warning-light/40' : 'border-danger/30 bg-danger-light/40']"
              >
                <FlagIcon :flag="parseFlags(detail.claim.flagsJson)[k] ?? 'ok'" />
                <span class="text-[9px] font-semibold text-muted-foreground">{{ label }}</span>
              </div>
            </div>
          </SectionCard>

          <SectionCard :title="`Sinyal Terdeteksi (${detail.anomalies.length})`" sub="Bahasa netral: 'potensi inkonsistensi' — bukan tuduhan. Verifikasi manual wajib.">
            <EmptyState v-if="detail.anomalies.length === 0" title="Tidak ada sinyal pada klaim ini">
              <template #icon><CheckCircle2 class="h-5 w-5" /></template>
            </EmptyState>
            <div v-else class="space-y-2">
              <div v-for="a in detail.anomalies" :key="a.id" class="rounded-xl border p-3.5">
                <div class="flex flex-wrap items-center justify-between gap-2">
                  <div class="flex items-center gap-2">
                    <FlagIcon :flag="a.severity === 'Tinggi' ? 'bad' : 'warn'" />
                    <span class="text-xs font-semibold">{{ a.category }}</span>
                  </div>
                  <StatusBadge :status="a.status" />
                </div>
                <p class="mt-1.5 text-xs leading-relaxed text-foreground/90">{{ a.description }}</p>
                <div class="mt-2 rounded-lg bg-brand-2-light/40 px-3 py-2 text-[11px] leading-relaxed text-brand-2-dark">
                  <strong>Rekomendasi:</strong> {{ a.recommendation }}
                </div>
              </div>
            </div>
          </SectionCard>

          <div class="flex flex-wrap gap-2">
            <Button
              class="gap-1.5 bg-brand hover:bg-brand-hover"
              @click="store.openAction({
                title: `Tinjau klaim ${detail.claim.claimNumber}`,
                description: `Menyelesaikan ${detail.anomalies.length} sinyal pada klaim ${detail.claim.claimNumber} (${detail.claim.inacbgDesc}).`,
                sourceType: 'Klaim', sourceRef: detail.claim.claimNumber,
                priority: detail.anomalies.some((a: any) => a.severity === 'Tinggi') ? 'Tinggi' : 'Sedang',
              })"
            ><PlusCircle class="h-4 w-4" /> Buat Tindakan</Button>
            <Button variant="outline" class="gap-1.5" @click="store.openEvidence({ refId: detail.claim.claimNumber, title: `Bukti klaim ${detail.claim.claimNumber}` })">
              <Fingerprint class="h-4 w-4" /> Lihat Bukti
            </Button>
          </div>
        </div>
      </template>
    </Dialog>
  </div>
</template>
