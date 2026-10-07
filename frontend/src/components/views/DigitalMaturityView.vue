<script setup lang="ts">
// DigitalMaturityView — port dari views-readiness.tsx (DigitalMaturityView).
import { MonitorSmartphone } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { computed } from 'vue'
import { useIntel } from '~/composables/useIntel'
import Button from '~/components/ui/Button.vue'
import Skeleton from '~/components/ui/Skeleton.vue'
import Hint from '~/components/medcredix/Hint.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import ScoreRing from '~/components/medcredix/ScoreRing.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const { data, loading } = useIntel()
const dims = computed(() => (data.value?.intel?.domains ?? []).filter((d: any) => ['digital', 'medication', 'diagnostic'].includes(d.key)))
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Digital Maturity"
      sub="EMR/RME, infrastruktur, medikasi, diagnostik — setiap indikator punya ambang, status, dan sumber FHIR-nya." demo
    />
    <template v-if="loading || !data">
      <Skeleton class="h-8 w-64" /><Skeleton class="mt-4 h-72" />
    </template>
    <template v-else>
      <div class="grid gap-4 md:grid-cols-3">
        <div v-for="d in dims" :key="d.key" class="rounded-2xl border bg-card p-5 text-center">
          <ScoreRing :value="d.score" :size="104" :label="d.label" @click="store.openWhy(d.key)" />
          <Button size="sm" variant="outline" class="mt-3 w-full text-xs" @click="store.openWhy(d.key)">Mengapa Skor Ini?</Button>
        </div>
      </div>
      <SectionCard v-for="d in dims" :key="d.key" :title="d.label" :sub="d.formula">
        <template #icon><MonitorSmartphone class="h-4 w-4" /></template>
        <div class="overflow-x-auto rounded-xl border">
          <table class="w-full text-xs">
            <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
              <tr>
                <th class="px-3 py-2.5">Indikator</th><th class="px-3 py-2.5 text-right">Nilai</th>
                <th class="px-3 py-2.5 text-right">Ambang</th><th class="px-3 py-2.5">Status</th><th class="px-3 py-2.5">Sumber Data</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(i, idx) in d.indicators" :key="idx" class="border-t">
                <td class="px-3 py-2.5 font-medium">
                  {{ i.label }}
                  <Hint>{{ i.note ?? i.source }}{{ i.fhir ? ` · FHIR: ${i.fhir}` : '' }}</Hint>
                </td>
                <td class="px-3 py-2.5 text-right font-display font-bold">{{ i.value }}{{ i.unit === '%' ? '%' : ` ${i.unit}` }}</td>
                <td class="px-3 py-2.5 text-right text-muted-foreground">{{ i.threshold }}{{ i.unit === '%' ? '%' : ` ${i.unit}` }}</td>
                <td class="px-3 py-2.5"><StatusBadge :status="i.status" /></td>
                <td class="px-3 py-2.5 text-[10px] text-muted-foreground">{{ i.source }}{{ i.fhir ? ` · ${i.fhir}` : '' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </SectionCard>
    </template>
  </div>
</template>
