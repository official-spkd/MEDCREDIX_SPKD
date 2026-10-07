<script setup lang="ts">
// DocumentsView — port dari views-credentialing.tsx (DocumentsView).
import { Clock, Eye } from 'lucide-vue-next'
import { expiryLabel } from '~/utils/format'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const data = ref<any>(null)
onMounted(async () => {
  try { data.value = await api('/credentials') } catch { data.value = { credentials: [] } }
})
const creds = computed(() => (data.value?.credentials ?? []).filter((c: any) => c.number).slice(0, 30))
function toneCls(tone: string) {
  return tone === 'danger' ? 'bg-danger-light text-danger' : tone === 'warning' ? 'bg-warning-light text-warning' : tone === 'success' ? 'bg-success-light text-success' : 'bg-muted text-muted-foreground'
}
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Dokumen Kredensial"
      sub="Semua dokumen terunggah dapat dibuka kembali (perbaikan dari temuan QA produk lama: tombol mati & dokumen tak terlihat)." demo
    />
    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="c in creds" :key="c.id" class="rounded-xl border bg-card p-4">
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <Badge variant="outline" class="text-[9px] font-bold">{{ c.type }}</Badge>
            <div class="mt-1.5 truncate text-xs font-semibold" :title="c.name">{{ c.name.split('—')[0] }}</div>
            <div class="mt-0.5 text-[10px] text-muted-foreground">{{ c.practitionerName }}</div>
          </div>
          <StatusBadge :status="c.status" />
        </div>
        <div class="mt-3 flex items-center justify-between border-t pt-2.5">
          <span :class="['inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-semibold', toneCls(expiryLabel(c.expiryDate).tone)]">
            <Clock class="h-3 w-3" /> {{ expiryLabel(c.expiryDate).text }}
          </span>
          <Button
            size="sm" variant="outline" class="h-6 gap-1 px-2 text-[10px]"
            @click="toast({ title: 'Pratinjau dokumen', description: `${c.type} — ${c.practitionerName}. Pada produksi, dokumen dibuka dari object storage terenkripsi (S3) dengan tautan bertanda tangan.` })"
          >
            <Eye class="h-3 w-3" /> Lihat
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>
