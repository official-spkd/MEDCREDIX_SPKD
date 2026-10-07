<script setup lang="ts">
// VerificationView — port dari views-credentialing.tsx (VerificationView).
import { BadgeCheck, CheckCircle2, XCircle, Clock } from 'lucide-vue-next'
import { expiryLabel, fmtDate } from '~/utils/format'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import EmptyState from '~/components/medcredix/EmptyState.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const data = ref<any>(null)
async function load() {
  try { data.value = await api('/credentials') } catch { data.value = { credentials: [] } }
}
onMounted(load)
const queue = computed(() =>
  (data.value?.credentials ?? []).filter((c: any) => ['Diajukan', 'Tinjauan Dokumen', 'Verifikasi', 'Tinjauan Komite'].includes(c.status)),
)
async function transition(c: any, action: string) {
  try {
    await api('/credentials', { method: 'PATCH', body: { id: c.id, action } })
    toast({ title: `Berhasil: ${action}` })
    load()
  } catch (e: any) {
    toast({ title: 'Ditolak', description: e?.message, variant: 'destructive' })
  }
}
function toneCls(tone: string) {
  return tone === 'danger' ? 'bg-danger-light text-danger' : tone === 'warning' ? 'bg-warning-light text-warning' : tone === 'success' ? 'bg-success-light text-success' : 'bg-muted text-muted-foreground'
}
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Verifikasi & Persetujuan"
      sub="Alur kerja petugas kredensial → komite medis. Setiap keputusan tercatat immutable di Audit Trail." demo
    />
    <EmptyState
      v-if="queue.length === 0" title="Antrean verifikasi kosong"
      desc="Tidak ada kredensial yang menunggu verifikasi atau persetujuan komite saat ini."
    >
      <template #icon><BadgeCheck class="h-5 w-5" /></template>
    </EmptyState>
    <div v-else class="grid gap-3 lg:grid-cols-2">
      <div v-for="c in queue" :key="c.id" class="rounded-xl border bg-card p-4">
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-1.5">
              <Badge variant="outline" class="text-[9px] font-bold">{{ c.type }}</Badge>
              <StatusBadge :status="c.status" />
            </div>
            <div class="mt-1.5 text-sm font-semibold leading-snug">{{ c.practitionerName }}</div>
            <div class="mt-0.5 text-xs text-muted-foreground">{{ c.name.split('—')[0] }} · {{ c.number }} · berakhir {{ fmtDate(c.expiryDate) }}</div>
          </div>
          <span :class="['inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-semibold', toneCls(expiryLabel(c.expiryDate).tone)]">
            <Clock class="h-3 w-3" /> {{ expiryLabel(c.expiryDate).text }}
          </span>
        </div>
        <div class="mt-3 flex flex-wrap gap-1.5 border-t pt-3">
          <Button v-if="['Diajukan', 'Tinjauan Dokumen'].includes(c.status)" size="sm" variant="outline" class="h-7 gap-1 text-xs" @click="transition(c, 'verify')">
            <CheckCircle2 class="h-3.5 w-3.5" /> Verifikasi Dokumen
          </Button>
          <Button v-if="c.status === 'Verifikasi'" size="sm" variant="outline" class="h-7 gap-1 text-xs" @click="transition(c, 'review')">
            <BadgeCheck class="h-3.5 w-3.5" /> Teruskan ke Komite
          </Button>
          <template v-if="c.status === 'Tinjauan Komite'">
            <Button size="sm" class="h-7 gap-1 bg-success text-xs hover:bg-success" @click="transition(c, 'approve')">
              <CheckCircle2 class="h-3.5 w-3.5" /> Setujui
            </Button>
            <Button size="sm" variant="outline" class="h-7 gap-1 text-xs text-danger" @click="transition(c, 'suspend')">
              <XCircle class="h-3.5 w-3.5" /> Tangguhkan
            </Button>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>
