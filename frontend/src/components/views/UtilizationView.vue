<script setup lang="ts">
// UtilizationView — port dari views-clinical.tsx (UtilizationView).
import { TrendingUp, BedDouble } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import { fmtID } from '~/utils/format'
import ChartBar from '~/components/medcredix/ChartBar.vue'
import ChartLine from '~/components/medcredix/ChartLine.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const data = ref<any>(null)
onMounted(async () => {
  try { data.value = await api('/clinical?type=utilization') } catch { data.value = null }
})
const utils = computed(() => data.value?.utilization ?? [])
const latest = computed(() => utils.value[utils.value.length - 1])
const complaints = computed(() => data.value?.complaints ?? [])
const lastC = computed(() => complaints.value[complaints.value.length - 1])
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Utilization Review"
      sub="ALOS ≤ 5 hari · Readmission ≤ 5% · BOR 60–85% (target rujukan BPJS). Benchmark eksternal selalu mencantumkan sumber & konteksnya." demo
    />
    <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
      <div
        v-for="m in [
          { label: 'ALOS', value: latest?.alos, unit: 'hari', target: '≤ 5,0', ok: (latest?.alos ?? 0) <= 5 },
          { label: 'Readmission Rate', value: latest?.readmission, unit: '%', target: '≤ 5,0%', ok: (latest?.readmission ?? 0) <= 5 },
          { label: 'BOR', value: latest?.bor, unit: '%', target: '60–85%', ok: (latest?.bor ?? 0) >= 60 && (latest?.bor ?? 0) <= 85 },
        ]" :key="m.label" class="rounded-xl border bg-card p-4"
      >
        <div class="flex items-center justify-between">
          <div class="text-[10px] font-bold uppercase tracking-wide text-muted-foreground">{{ m.label }}</div>
          <StatusBadge :status="m.ok ? 'Baik' : 'Perlu Perhatian'" />
        </div>
        <div class="mt-1 font-display text-2xl font-bold text-foreground">{{ m.value }}<span class="ml-1 text-xs text-muted-foreground">{{ m.unit }}</span></div>
        <div class="mt-0.5 text-[10px] text-muted-foreground">Target {{ m.target }} · 6 bulan terakhir</div>
      </div>
    </div>

    <SectionCard title="Tren 6 Bulan" sub="ALOS menurun, readmission membaik, BOR mendekati batas atas optimal">
      <template #icon><TrendingUp class="h-4 w-4" /></template>
      <ChartLine
        :labels="utils.map((u: any) => u.monthLabel)"
        :datasets="[
          { label: 'ALOS (hari)', data: utils.map((u: any) => u.alos), color: '#14285C', width: 2.5 },
          { label: 'Readmission (%)', data: utils.map((u: any) => u.readmission), color: '#D97706', width: 2.5 },
          { label: 'BOR (%)', data: utils.map((u: any) => u.bor), color: '#2198A5', width: 2.5 },
        ]"
        :y-min="0" :y-max="100" :height="224" :legend="true"
      />
    </SectionCard>

    <div class="grid gap-4 lg:grid-cols-2">
      <SectionCard
        title="Keluhan & Kepuasan JKN"
        :sub="`Sumber: JKN Mobile (P-Care) & portal pengaduan BPJS · ${lastC ? `${fmtID(lastC.total)} keluhan, kepuasan ${lastC.satisfaction}%` : ''}`"
      >
        <ChartBar
          :labels="complaints.map((c: any) => c.month.slice(5))"
          :datasets="[
            { label: 'Antrian', data: complaints.map((c: any) => c.antrian), color: '#14285C', stack: 'a' },
            { label: 'Pelayanan', data: complaints.map((c: any) => c.pelayanan), color: '#2198A5', stack: 'a' },
            { label: 'Obat', data: complaints.map((c: any) => c.obat), color: '#7C93C4', stack: 'a' },
            { label: 'Administrasi', data: complaints.map((c: any) => c.administrasi), color: '#66C3CE', stack: 'a' },
            { label: 'Fasilitas', data: complaints.map((c: any) => c.fasilitas), color: '#D97706', stack: 'a' },
            { label: 'Lainnya', data: complaints.map((c: any) => c.lainnya), color: '#94A3B8', stack: 'a' },
          ]"
          :height="208" stacked :legend="true"
        />
        <p class="mt-2 text-[11px] leading-relaxed text-muted-foreground">
          Total keluhan per 1.000 kunjungan: <strong class="text-foreground">{{ lastC?.per1000 }}</strong> vs rata-rata RS sejenis
          <strong class="text-foreground">{{ lastC?.peerAvg }}</strong> (peer benchmark, periode sama).
        </p>
      </SectionCard>

      <SectionCard title="Riwayat Metrik Utilization" sub="Definisi: ALOS = hari perawatan ÷ pasien keluar; BOR = hari rawat ÷ kapasitas">
        <template #icon><BedDouble class="h-4 w-4" /></template>
        <div class="overflow-x-auto rounded-xl border">
          <table class="w-full text-xs">
            <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
              <tr><th class="px-3 py-2.5">Periode</th><th class="px-3 py-2.5 text-right">ALOS</th><th class="px-3 py-2.5 text-right">Readmission</th><th class="px-3 py-2.5 text-right">BOR</th></tr>
            </thead>
            <tbody>
              <tr v-for="u in [...utils].reverse()" :key="u.id" class="border-t">
                <td class="px-3 py-2 font-medium">{{ u.monthLabel }} {{ u.month.slice(0, 4) }}</td>
                <td :class="['px-3 py-2 text-right font-display font-bold', u.alos <= 5 ? 'text-success' : 'text-warning']">{{ u.alos }}</td>
                <td :class="['px-3 py-2 text-right font-display font-bold', u.readmission <= 5 ? 'text-success' : 'text-warning']">{{ u.readmission }}%</td>
                <td :class="['px-3 py-2 text-right font-display font-bold', u.bor >= 60 && u.bor <= 85 ? 'text-success' : 'text-warning']">{{ u.bor }}%</td>
              </tr>
            </tbody>
          </table>
        </div>
      </SectionCard>
    </div>
  </div>
</template>
