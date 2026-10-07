<script setup lang="ts">
// SettingsAlertsView — port dari views-settings.tsx.
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import { fmtDate } from '~/utils/format'
import Switch from '~/components/ui/Switch.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const data = ref<any>(null)
async function load() {
  try { data.value = await api('/settings') } catch { data.value = null }
}
onMounted(load)
const canEdit = computed(() => store.user?.permissions.includes('CONFIGURE'))

async function toggle(rule: any) {
  try {
    await api('/settings', { method: 'PATCH', body: { target: 'alert-rule', id: rule.id, active: !rule.active } })
    await load()
    toast({ title: `Aturan ${rule.active ? 'dinonaktifkan' : 'diaktifkan'}` })
  } catch (e: any) {
    toast({ title: 'Gagal', description: e?.message, variant: 'destructive' })
  }
}
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Aturan Notifikasi"
      sub="Kanal: in-app aktif penuh; email & WhatsApp-ready terkonfigurasi pada level platform oleh SPKD." demo
    />
    <div class="grid gap-3 md:grid-cols-2">
      <div v-for="r in data?.rules ?? []" :key="r.id" class="rounded-xl border bg-card p-4">
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="text-xs font-semibold leading-snug">{{ r.name }}</div>
            <div class="mt-1.5 font-mono text-[10px] text-muted-foreground">IF {{ r.condition }}</div>
            <div class="mt-1.5 flex flex-wrap items-center gap-1.5">
              <StatusBadge :status="r.severity" :tone="r.severity === 'Tinggi' ? 'danger' : r.severity === 'Sedang' ? 'warning' : 'info'" />
              <span class="text-[9px] text-muted-foreground">{{ r.channels }}</span>
            </div>
          </div>
          <Switch :model-value="r.active" :disabled="!canEdit" :aria-label="`Aktifkan ${r.name}`" @update:model-value="() => toggle(r)" />
        </div>
        <div class="mt-2.5 border-t pt-2.5 text-[10px] text-muted-foreground">
          Terpicu {{ r.triggeredCount }}× · terakhir {{ fmtDate(r.lastTriggeredAt) }} · Domain {{ r.domain }}
        </div>
      </div>
    </div>
  </div>
</template>
