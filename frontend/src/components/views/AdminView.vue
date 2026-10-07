<script setup lang="ts">
// AdminView — Panel Admin SPKD (port dari views-settings.tsx).
import { Building, Award, TrendingUp, ShieldCheck, Globe2 } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import Badge from '~/components/ui/Badge.vue'
import ChartBar from '~/components/medcredix/ChartBar.vue'
import KpiCard from '~/components/medcredix/KpiCard.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const data = ref<any>(null)
onMounted(async () => {
  try { data.value = await api('/admin') } catch { data.value = null }
})
const hospitals = computed(() => data.value?.hospitals ?? [])
const stageColors = (stage: number) => (stage >= 5 ? '#059669' : stage >= 3 ? '#D97706' : '#DC2626')
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Panel Admin SPKD"
      sub="Gambaran nasional lintas rumah sakit — untuk Super Admin & Auditor BPJS. Isolasi tenant tetap ditegakkan: auditor melihat agregat, bukan data mentah RS lain." demo
    />
    <div class="grid grid-cols-2 gap-3 md:grid-cols-5">
      <KpiCard label="Total RS" :value="data?.summary?.totalRS ?? 0" tone="info">
        <template #icon><Building class="h-4 w-4" /></template>
      </KpiCard>
      <KpiCard label="Rata-rata BDCS" :value="data?.summary?.rataBdc ?? 0" tone="info">
        <template #icon><Award class="h-4 w-4" /></template>
      </KpiCard>
      <KpiCard label="Rata-rata Health" :value="data?.summary?.rataHealth ?? 0" tone="success">
        <template #icon><TrendingUp class="h-4 w-4" /></template>
      </KpiCard>
      <KpiCard label="Lulus Kredensial" :value="data?.summary?.lulus ?? 0" unit="RS" tone="success">
        <template #icon><ShieldCheck class="h-4 w-4" /></template>
      </KpiCard>
      <KpiCard label="Under Review" :value="data?.summary?.underReview ?? 0" unit="RS" :tone="(data?.summary?.underReview ?? 0) > 0 ? 'danger' : 'muted'">
        <template #icon><Globe2 class="h-4 w-4" /></template>
      </KpiCard>
    </div>

    <div class="grid gap-4 lg:grid-cols-[1fr_360px]">
      <SectionCard title="Manajemen Rumah Sakit" sub="Klik untuk membuka konteks RS (super admin)">
        <div class="overflow-x-auto rounded-xl border">
          <table class="w-full text-xs">
            <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
              <tr>
                <th class="px-3 py-2.5">Rumah Sakit</th><th class="px-3 py-2.5">Kelas</th><th class="px-3 py-2.5 text-right">Health</th>
                <th class="px-3 py-2.5 text-right">BDCS</th><th class="px-3 py-2.5 text-center">Stage</th>
                <th class="px-3 py-2.5 text-right">Trust</th><th class="px-3 py-2.5 text-right">CCS</th><th class="px-3 py-2.5">Status</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="h in hospitals" :key="h.id"
                class="cursor-pointer border-t transition-colors hover:bg-muted/30"
                @click="store.setHospital(h.id); store.navigate('dashboard')"
              >
                <td class="px-3 py-2.5">
                  <div class="font-semibold">{{ h.name }}</div>
                  <div class="text-[9px] text-muted-foreground">{{ h.city }}, {{ h.province }} · Tier {{ h.tier }}</div>
                </td>
                <td class="px-3 py-2.5">{{ h.class }}</td>
                <td class="px-3 py-2.5 text-right font-display font-bold">{{ h.healthScore }}</td>
                <td class="px-3 py-2.5 text-right font-display font-bold text-brand">{{ h.bcds }}</td>
                <td class="px-3 py-2.5 text-center">
                  <Badge :class="h.lulus ? 'bg-success text-[10px] hover:bg-success' : 'bg-danger text-[10px] hover:bg-danger'">S{{ h.stage }} {{ h.lulus ? 'LULUS' : 'TIDAK' }}</Badge>
                </td>
                <td class="px-3 py-2.5 text-right">{{ h.trust }}</td>
                <td class="px-3 py-2.5 text-right">{{ h.ccs }}</td>
                <td class="px-3 py-2.5"><StatusBadge :status="h.underReview ? 'Under Review' : 'Normal'" :tone="h.underReview ? 'danger' : 'success'" /></td>
              </tr>
            </tbody>
          </table>
        </div>
      </SectionCard>

      <div class="space-y-4">
        <SectionCard title="Distribusi Stage Nasional" sub="Stage 0–7 (kelulusan ≥ Stage 5)">
          <ChartBar
            :labels="(data?.stageDist ?? []).map((s: any) => `S${s.stage}`)"
            :datasets="[{ label: 'Jumlah RS', data: (data?.stageDist ?? []).map((s: any) => s.count), color: (data?.stageDist ?? []).map((s: any) => stageColors(s.stage)) }]"
            :height="176" :suggested-max="4"
          />
        </SectionCard>
        <SectionCard title="Tren BDCS per RS" sub="6 bulan terakhir">
          <div class="space-y-2.5">
            <div v-for="h in hospitals" :key="h.id">
              <div class="mb-1 flex items-center justify-between text-[10px]">
                <span class="font-semibold">{{ h.name }}</span>
                <span class="font-display font-bold text-brand">{{ h.trend[h.trend.length - 1]?.bcds ?? '—' }}</span>
              </div>
              <div class="flex h-6 items-end gap-1">
                <div
                  v-for="t in h.trend" :key="t.period"
                  class="flex-1 rounded-sm bg-brand-2/80" :style="{ height: `${Math.max(12, ((t.bcds - 40) / 60) * 100)}%` }"
                  :title="`${t.period}: ${t.bcds}`"
                />
              </div>
            </div>
          </div>
        </SectionCard>
      </div>
    </div>
  </div>
</template>
