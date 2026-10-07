<script setup lang="ts">
// InteropView — port dari views-readiness.tsx (InteropView).
import { Activity, Network } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import { useIntel } from '~/composables/useIntel'
import Button from '~/components/ui/Button.vue'
import Skeleton from '~/components/ui/Skeleton.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import ScoreRing from '~/components/medcredix/ScoreRing.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const { data, loading } = useIntel()
const src = ref<any>(null)

async function loadSrc() {
  try { src.value = await api(`/data?type=sources${store.hospitalId ? `&hospitalId=${store.hospitalId}` : ''}`) } catch { src.value = { sources: [] } }
}
onMounted(loadSrc)
watch(() => store.hospitalId, loadSrc)

const d = computed(() => data.value?.intel?.domains?.find((x: any) => x.key === 'interop'))
const interopSources = computed(() => (src.value?.sources ?? []).filter((s: any) => ['SATUSEHAT', 'BPJS VClaim', 'BPJS P-Care', 'SIMRS'].includes(s.type)))
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Interoperabilitas"
      sub="Pertukaran data dengan SATUSEHAT, BPJS, dan validasi FHIR (TDS v2.0.0 — 24 resource wajib)." demo
    />
    <template v-if="loading || !data">
      <Skeleton class="h-8 w-64" /><Skeleton class="mt-4 h-72" />
    </template>
    <template v-else>
      <div class="grid gap-4 lg:grid-cols-[280px_1fr]">
        <div class="rounded-2xl border bg-card p-5">
          <ScoreRing :value="d?.score ?? 0" :size="120" label="Skor Interoperabilitas" @click="store.openWhy('interop')" />
          <Button size="sm" variant="outline" class="mt-4 w-full gap-1.5 text-xs" @click="store.openWhy('interop')">
            <Activity class="h-3.5 w-3.5" /> Mengapa Skor Ini?
          </Button>
        </div>
        <SectionCard title="Indikator Interoperabilitas" :sub="d?.formula">
          <template #icon><Network class="h-4 w-4" /></template>
          <div class="overflow-x-auto rounded-xl border">
            <table class="w-full text-xs">
              <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
                <tr><th class="px-3 py-2.5">Indikator</th><th class="px-3 py-2.5 text-right">Nilai</th><th class="px-3 py-2.5 text-right">Ambang</th><th class="px-3 py-2.5">Status</th></tr>
              </thead>
              <tbody>
                <tr v-for="(i, idx) in d?.indicators ?? []" :key="idx" class="border-t">
                  <td class="px-3 py-2.5 font-medium">{{ i.label }}</td>
                  <td class="px-3 py-2.5 text-right font-display font-bold">{{ i.value }}%</td>
                  <td class="px-3 py-2.5 text-right text-muted-foreground">{{ i.threshold }}%</td>
                  <td class="px-3 py-2.5"><StatusBadge :status="i.status" /></td>
                </tr>
              </tbody>
            </table>
          </div>
        </SectionCard>
      </div>

      <SectionCard title="Kanal Integrasi Eksternal" sub="GET-only ke sistem RS — platform tidak pernah menulis ke SIMRS">
        <div class="grid gap-3 sm:grid-cols-2">
          <div v-for="s in interopSources" :key="s.id" class="rounded-xl border bg-card p-4">
            <div class="flex items-start justify-between gap-2">
              <div>
                <div class="text-sm font-semibold">{{ s.name }}</div>
                <div class="mt-0.5 text-[10px] text-muted-foreground">{{ s.configSummary }}</div>
              </div>
              <StatusBadge :status="s.status" />
            </div>
            <div class="mt-3 grid grid-cols-3 gap-2 border-t pt-3 text-center">
              <div><div class="font-display text-sm font-bold text-foreground">{{ s.successRate }}%</div><div class="text-[9px] text-muted-foreground">Sukses</div></div>
              <div><div class="font-display text-sm font-bold text-foreground">{{ s.avgLatencyMs }} ms</div><div class="text-[9px] text-muted-foreground">Latensi</div></div>
              <div><div class="font-display text-sm font-bold text-foreground">{{ s.uptime30d }}%</div><div class="text-[9px] text-muted-foreground">Uptime 30d</div></div>
            </div>
            <p v-if="s.notes" class="mt-2 rounded-lg bg-warning-light px-2.5 py-1.5 text-[10px] leading-snug text-warning">{{ s.notes }}</p>
          </div>
        </div>
      </SectionCard>
    </template>
  </div>
</template>
