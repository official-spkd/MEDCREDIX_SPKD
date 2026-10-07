<script setup lang="ts">
// ClaimAnomaliesView — port dari views-claims.tsx (ClaimAnomaliesView).
import { ShieldAlert, PlusCircle, CheckCircle2, ScanSearch, Eye } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const status = ref('all')
const data = ref<any>(null)

async function load() {
  try { data.value = await api('/claims') } catch { data.value = null }
}
onMounted(load)

async function review(a: any, action: string) {
  try {
    await api('/claims', { method: 'PATCH', body: { anomalyId: a.id, action } })
    toast({ title: `Anomali ${a.claimNumber}: ${action === 'confirm' ? 'terkonfirmasi bermasalah' : action === 'clarify' ? 'diklarifikasi' : 'ditinjau'}` })
    load()
  } catch (e: any) {
    toast({ title: 'Gagal', description: e?.message, variant: 'destructive' })
  }
}
const anomalies = computed(() => (data.value?.anomalies ?? []).filter((a: any) => status.value === 'all' || a.status === status.value))
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Anomali Klaim"
      sub="Upcoding/downcoding, medication mismatch, stock-claim triangle, timeline — semua sinyal netral yang memerlukan peninjauan manusia." demo
    />
    <div class="flex flex-wrap items-center gap-2 rounded-xl border bg-card p-3">
      <ShieldAlert class="h-4 w-4 text-warning" />
      <button
        v-for="s in ['all', 'Baru', 'Ditinjau', 'Terkonfirmasi Bermasalah', 'Diklarifikasi']" :key="s"
        :class="['inline-flex h-7 items-center rounded-lg px-3 text-xs font-semibold transition-colors', status === s ? 'bg-brand text-white' : 'border border-input bg-card hover:bg-muted/60']"
        @click="status = s"
      >{{ s === 'all' ? 'Semua' : s }}</button>
      <span class="ml-auto text-xs text-muted-foreground">{{ anomalies.length }} sinyal</span>
    </div>

    <div class="grid gap-3 lg:grid-cols-2">
      <div
        v-for="a in anomalies" :key="a.id"
        :class="['rounded-xl border bg-card p-4', a.severity === 'Tinggi' && 'border-danger/30', a.severity === 'Sedang' && 'border-warning/30']"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-1.5">
              <StatusBadge :status="a.severity" :tone="a.severity === 'Tinggi' ? 'danger' : a.severity === 'Sedang' ? 'warning' : 'info'" />
              <Badge variant="outline" class="text-[10px]">{{ a.category }}</Badge>
              <StatusBadge :status="a.status" />
            </div>
            <h3 class="mt-2 font-display text-sm font-semibold">{{ a.claimNumber }} — {{ a.title }}</h3>
            <p class="mt-1 line-clamp-2 text-xs leading-relaxed text-muted-foreground">{{ a.description }}</p>
          </div>
        </div>
        <div class="mt-3 rounded-lg bg-brand-2-light/40 px-3 py-2 text-[11px] leading-relaxed text-brand-2-dark">
          <strong>Cara penyelesaian:</strong> {{ a.recommendation }}
        </div>
        <div class="mt-2.5 flex flex-wrap gap-1.5 border-t pt-3">
          <Button v-if="a.status === 'Baru'" size="sm" variant="outline" class="h-7 gap-1 text-[11px]" @click="review(a, 'review')"><Eye class="h-3 w-3" /> Tandai Ditinjau</Button>
          <Button v-if="a.status !== 'Diklarifikasi'" size="sm" variant="outline" class="h-7 gap-1 text-[11px]" @click="review(a, 'clarify')"><CheckCircle2 class="h-3 w-3" /> Klarifikasi (terjelaskan)</Button>
          <Button v-if="a.status !== 'Terkonfirmasi Bermasalah'" size="sm" variant="outline" class="h-7 gap-1 text-[11px] text-danger" @click="review(a, 'confirm')"><ScanSearch class="h-3 w-3" /> Konfirmasi Bermasalah</Button>
          <Button
            size="sm" variant="ghost" class="ml-auto h-7 gap-1 text-[11px]"
            @click="store.openAction({ title: `Selesaikan sinyal ${a.category} pada ${a.claimNumber}`, description: a.description, sourceType: 'Klaim', sourceRef: a.claimNumber ?? '', priority: a.severity })"
          ><PlusCircle class="h-3 w-3" /> Buat Tindakan</Button>
        </div>
      </div>
    </div>
  </div>
</template>
