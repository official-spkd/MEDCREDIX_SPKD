<script setup lang="ts">
// WorkforceView — port dari views-readiness.tsx (WorkforceView).
import { Stethoscope, Activity, Info } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { computed } from 'vue'
import { useIntel } from '~/composables/useIntel'
import { fmtID } from '~/utils/format'
import Button from '~/components/ui/Button.vue'
import Skeleton from '~/components/ui/Skeleton.vue'
import Hint from '~/components/medcredix/Hint.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import ScoreRing from '~/components/medcredix/ScoreRing.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const { data, loading } = useIntel()
const wf = computed(() => data.value?.intel?.workforce)
const domain = computed(() => data.value?.intel?.domains?.find((d: any) => d.key === 'workforce'))
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Clinical Workforce"
      sub="Cakupan spesialis wajib per standar Permenkes 56/2014 kelas RS — surplus & defisit tampil eksplisit (perbaikan temuan QA produk lama)." demo
    />
    <template v-if="loading || !data">
      <Skeleton class="h-8 w-64" /><Skeleton class="mt-4 h-72" />
    </template>
    <template v-else>
      <div class="grid gap-4 lg:grid-cols-[280px_1fr]">
        <div class="rounded-2xl border bg-card p-5">
          <ScoreRing :value="domain?.score ?? 0" :size="120" label="Skor Clinical Workforce" @click="store.openWhy('workforce')" />
          <Button size="sm" variant="outline" class="mt-4 w-full gap-1.5 text-xs" @click="store.openWhy('workforce')">
            <Activity class="h-3.5 w-3.5" /> Mengapa Skor Ini?
          </Button>
          <div class="mt-4 space-y-2 border-t pt-4 text-xs">
            <div class="flex justify-between"><span class="text-muted-foreground">Total praktisi</span><strong>{{ fmtID(wf.practitioners) }}</strong></div>
            <div class="flex justify-between">
              <span class="text-muted-foreground">Spesialis terpenuhi</span>
              <strong>{{ wf.specialistCov.filter((s: any) => s.coverage >= 100).length }}/{{ wf.specialistCov.length }}</strong>
            </div>
            <div class="flex justify-between"><span class="text-muted-foreground">Rumus</span><Hint>Skor = cakupan 50% + SIP aktif 30% + tenaga aktif 20%</Hint></div>
          </div>
        </div>

        <SectionCard
          title="Cakupan Spesialis Wajib"
          sub="Tersedia vs kebutuhan — cakupan per spesialis dapat melampaui 100% bila tersedia melebihi kebutuhan"
        >
          <template #icon><Stethoscope class="h-4 w-4" /></template>
          <div class="overflow-x-auto rounded-xl border">
            <table class="w-full text-xs">
              <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
                <tr>
                  <th class="px-3 py-2.5">Spesialisasi (dasar kelas B)</th><th class="px-3 py-2.5 text-center">Wajib</th>
                  <th class="px-3 py-2.5 text-center">Tersedia</th><th class="px-3 py-2.5">Cakupan</th><th class="px-3 py-2.5">Status</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="s in wf.specialistCov" :key="s.specialty" class="border-t">
                  <td class="px-3 py-2.5 font-medium">
                    {{ s.specialty }}<span v-if="s.note" class="ml-1.5 text-[9px] text-muted-foreground">({{ s.note }})</span>
                  </td>
                  <td class="px-3 py-2.5 text-center">{{ s.required }}</td>
                  <td class="px-3 py-2.5 text-center font-display font-bold">{{ s.available }}</td>
                  <td class="px-3 py-2.5">
                    <div class="flex items-center gap-2">
                      <div class="h-1.5 w-24 overflow-hidden rounded-full bg-muted">
                        <div
                          :class="['h-full rounded-full', s.coverage >= 100 ? 'bg-success' : s.coverage > 0 ? 'bg-warning' : 'bg-danger']"
                          :style="{ width: `${Math.min(100, s.coverage)}%` }"
                        />
                      </div>
                      <span class="font-display text-[11px] font-bold">{{ Math.round(s.coverage) }}%</span>
                    </div>
                  </td>
                  <td class="px-3 py-2.5">
                    <StatusBadge v-if="s.available - s.required > 0" status="Surplus" tone="success" />
                    <StatusBadge v-else-if="s.available - s.required === 0" status="Terpenuhi" tone="success" />
                    <StatusBadge v-else :status="`Defisit ${Math.abs(s.available - s.required)}`" tone="danger" />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <p class="mt-3 flex items-start gap-1.5 text-[11px] leading-relaxed text-muted-foreground">
            <Info class="mt-0.5 h-3.5 w-3.5 shrink-0" />
            Total kebutuhan {{ wf.specialistCov.reduce((a: number, s: any) => a + s.required, 0) }} SDM spesialis; tersedia
            {{ wf.specialistCov.reduce((a: number, s: any) => a + s.available, 0) }}. Spesialis defisit memicu aksi rekrutmen otomatis di Action Center.
          </p>
        </SectionCard>
      </div>

      <SectionCard title="Kontributor Positif & Negatif" sub="Dari engine skor — dapat ditelusuri">
        <div class="grid gap-3 md:grid-cols-2">
          <div class="rounded-xl border border-success/30 bg-success-light/40 p-4">
            <div class="mb-2 text-xs font-semibold text-success">Positif</div>
            <ul class="space-y-1 text-xs"><li v-for="(p, i) in domain?.positives ?? []" :key="i">✓ {{ p }}</li></ul>
          </div>
          <div class="rounded-xl border border-warning/30 bg-warning-light/40 p-4">
            <div class="mb-2 text-xs font-semibold text-warning">Perlu Perhatian</div>
            <ul class="space-y-1 text-xs"><li v-for="(p, i) in domain?.negatives ?? []" :key="i">! {{ p }}</li></ul>
            <Button size="sm" variant="outline" class="mt-3 gap-1.5 text-xs" @click="store.navigate('actions')">Lihat aksi rekrutmen</Button>
          </div>
        </div>
      </SectionCard>
    </template>
  </div>
</template>
