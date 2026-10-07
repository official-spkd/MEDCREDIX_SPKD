<script setup lang="ts">
// SafetyView — port dari views-clinical.tsx (SafetyView).
import { Activity } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import Button from '~/components/ui/Button.vue'
import Hint from '~/components/medcredix/Hint.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import ScoreRing from '~/components/medcredix/ScoreRing.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const data = ref<any>(null)
onMounted(async () => {
  try { data.value = await api('/clinical?type=safety') } catch { data.value = null }
})
const s = computed(() => data.value?.safety)
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Patient Safety Intelligence"
      sub="Skor Keselamatan adalah AKUMULASI berbobot 3 domain (bukan rata-rata sederhana) — divisualisasikan agar tidak salah paham (perbaikan temuan QA produk lama)." demo
    />
    <div class="grid gap-4 lg:grid-cols-[300px_1fr]">
      <div class="rounded-2xl border bg-card p-5">
        <ScoreRing :value="s?.score ?? 0" :size="124" label="Skor Keselamatan Pasien" @click="store.openWhy('safety')" />
        <div class="mt-4 rounded-xl bg-muted/50 p-3 font-mono text-[10px] leading-relaxed text-muted-foreground">
          Skor = Alergi × 40% + High-Alert &amp; Closed-Loop × 30% + Dok. Pemulangan × 30% =
          <strong class="text-foreground">{{ s?.score }}</strong>
        </div>
        <Button size="sm" variant="outline" class="mt-4 w-full gap-1.5 text-xs" @click="store.openWhy('safety')">
          <Activity class="h-3.5 w-3.5" /> Mengapa Skor Ini?
        </Button>
      </div>
      <div class="space-y-3">
        <div
          v-for="d in [
            { label: 'Keamanan Alergi Obat', value: s?.allergy, weight: 40, threshold: 95, desc: 'Pasien beralergi tidak menerima obat kontraindikasi', fhir: 'AllergyIntolerance' },
            { label: 'Obat High-Alert & Closed-Loop', value: s?.closedLoop, weight: 30, threshold: 80, desc: 'Siklus resep → dispense → administrasi terdokumentasi penuh', fhir: 'MedicationAdministration' },
            { label: 'Dokumentasi Pemulangan Pasien', value: s?.discharge, weight: 30, threshold: 80, desc: 'Status pemulangan terdokumentasi di RME', fhir: 'Composition' },
          ]" :key="d.label" class="rounded-xl border bg-card p-4"
        >
          <div class="flex items-start justify-between gap-3">
            <div>
              <div class="flex items-center gap-2">
                <span class="text-xs font-semibold">{{ d.label }}</span>
                <Hint>Bobot {{ d.weight }}% · Ambang {{ d.threshold }}% · FHIR: {{ d.fhir }}</Hint>
              </div>
              <p class="mt-0.5 text-[10px] text-muted-foreground">{{ d.desc }}</p>
            </div>
            <div class="text-right">
              <div :class="['font-display text-xl font-bold', (d.value ?? 0) >= (d.threshold ?? 0) ? 'text-success' : 'text-warning']">{{ d.value }}%</div>
              <StatusBadge :status="(d.value ?? 0) >= (d.threshold ?? 0) ? 'Baik' : 'Perlu Perhatian'" />
            </div>
          </div>
          <div class="mt-2.5 h-2 overflow-hidden rounded-full bg-muted">
            <div :class="['h-full rounded-full transition-all', (d.value ?? 0) >= (d.threshold ?? 0) ? 'bg-success' : 'bg-warning']" :style="{ width: `${d.value}%` }" />
          </div>
          <div class="mt-1 text-right text-[9px] text-muted-foreground">ambang {{ d.threshold }}%</div>
        </div>
      </div>
    </div>
  </div>
</template>
