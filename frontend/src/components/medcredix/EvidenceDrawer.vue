<script setup lang="ts">
// EvidenceDrawer — telusuri bukti dari item mana pun.
// Port dari components/medcredix/action-dialog.tsx (EvidenceDrawer).
import { Database, ExternalLink } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import { fmtDateTime } from '~/utils/format'
import Button from '~/components/ui/Button.vue'
import Sheet from '~/components/ui/Sheet.vue'
import DemoBadge from '~/components/medcredix/DemoBadge.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const items = ref<any[]>([])
const loading = ref(false)

watch(
  () => store.evidenceRef,
  async (ref) => {
    items.value = []
    if (!ref) return
    loading.value = true
    try {
      const params = new URLSearchParams()
      if (ref.refType) params.set('refType', ref.refType)
      if (ref.refId) params.set('refId', ref.refId)
      const r = await api<{ items: any[] }>(`/evidence?${params.toString()}`)
      items.value = r.items ?? []
    } catch { items.value = [] }
    loading.value = false
  },
  { immediate: true },
)
</script>

<template>
  <Sheet :open="!!store.evidenceRef" @update:open="(o) => !o && store.closeEvidence()">
    <div class="overflow-y-auto scrollbar-thin" style="max-width: 28rem">
      <div class="border-b p-4 pb-2">
        <div class="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">
          <Database class="h-3.5 w-3.5" /> Bukti Data <DemoBadge />
        </div>
        <h2 class="mt-1 font-display text-base font-semibold leading-snug">{{ store.evidenceRef?.title ?? 'Telusuri Bukti' }}</h2>
      </div>
      <div class="space-y-3 p-4 pb-6">
        <div v-if="loading" class="py-8 text-center text-xs text-muted-foreground">Memuat bukti…</div>
        <div v-else-if="items.length === 0" class="rounded-xl border border-dashed p-6 text-center text-xs leading-relaxed text-muted-foreground">
          Tidak ada item bukti spesifik yang terikat pada entitas ini.
        </div>
        <div v-for="e in items" v-else :key="e.id" class="rounded-xl border bg-card p-3.5">
          <div class="flex items-start justify-between gap-2">
            <div class="text-xs font-semibold leading-snug">{{ e.title }}</div>
            <StatusBadge :status="e.validationStatus" />
          </div>
          <div class="mt-1.5 font-display text-lg font-bold text-brand">
            {{ e.value }}{{ e.unit === '%' ? '%' : ` ${e.unit ?? ''}` }}
          </div>
          <div class="mt-0.5 text-[11px] text-muted-foreground">{{ e.metric }} · {{ e.recordCount }} record</div>
          <div class="mt-2 flex flex-wrap gap-x-3 gap-y-1 text-[10px] text-muted-foreground">
            <span>Sumber: {{ e.sourceSystem }}</span>
            <span v-if="e.fhirResource">· FHIR: {{ e.fhirResource }}</span>
            <span>· Sinkron: {{ fmtDateTime(e.lastSyncAt) }}</span>
          </div>
        </div>
        <Button variant="outline" size="sm" class="w-full gap-1.5" @click="store.closeEvidence(); store.navigate('evidence')">
          Buka Evidence Center <ExternalLink class="h-3.5 w-3.5" />
        </Button>
      </div>
    </div>
  </Sheet>
</template>
