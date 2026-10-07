<script setup lang="ts">
// AuditView — port dari views-assurance.tsx (AuditView).
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import { fmtDateTime } from '~/utils/format'
import Badge from '~/components/ui/Badge.vue'
import Input from '~/components/ui/Input.vue'
import Select from '~/components/ui/Select.vue'
import SelectOption from '~/components/ui/SelectOption.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'

const action = ref('all')
const q = ref('')
const data = ref<any>(null)

async function load() {
  try {
    const params = new URLSearchParams()
    if (action.value !== 'all') params.set('action', action.value)
    if (q.value) params.set('q', q.value)
    data.value = await api(`/audit?${params}`)
  } catch { data.value = null }
}
onMounted(load)
watch(action, load)
let t: ReturnType<typeof setTimeout> | null = null
watch(q, () => { if (t) clearTimeout(t); t = setTimeout(load, 300) })
const logs = computed(() => data.value?.logs ?? [])
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Audit Trail"
      sub="Log immutable: siapa, apa, kapan, sebelum/sesudah, IP, correlation ID. Tidak dapat diubah atau dihapus siapa pun." demo
    />
    <div class="flex flex-wrap items-center gap-2 rounded-xl border bg-card p-3">
      <div class="w-40">
        <Select v-model="action">
          <SelectOption value="all">Semua Aksi</SelectOption>
          <SelectOption v-for="a in ['login', 'create', 'update', 'approve', 'verify', 'export', 'config', 'scoring', 'sync', 'detect']" :key="a" :value="a">{{ a }}</SelectOption>
        </Select>
      </div>
      <Input v-model="q" placeholder="Cari aktor, deskripsi, atau ID entitas…" class="h-8 w-64 text-xs" />
      <span class="ml-auto text-xs text-muted-foreground">{{ logs.length }} log (maks 200 terakhir)</span>
    </div>
    <div class="overflow-x-auto rounded-xl border bg-card">
      <table class="w-full text-xs">
        <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
          <tr>
            <th class="px-3 py-2.5">Waktu</th><th class="px-3 py-2.5">Aktor</th><th class="px-3 py-2.5">Aksi</th>
            <th class="px-3 py-2.5">Entitas</th><th class="px-3 py-2.5">Deskripsi</th><th class="px-3 py-2.5">Sebelum → Sesudah</th><th class="px-3 py-2.5">IP</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="l in logs" :key="l.id" class="border-t">
            <td class="whitespace-nowrap px-3 py-2 font-mono text-[10px] text-muted-foreground">{{ fmtDateTime(l.createdAt) }}</td>
            <td class="px-3 py-2"><div class="font-semibold">{{ l.actorName }}</div><div class="text-[9px] text-muted-foreground">{{ l.actorRole }}</div></td>
            <td class="px-3 py-2"><Badge variant="outline" class="font-mono text-[9px]">{{ l.action }}</Badge></td>
            <td class="px-3 py-2 text-[10px] text-muted-foreground">{{ l.entityType }}{{ l.entityId ? ` · ${l.entityId.slice(0, 14)}` : '' }}</td>
            <td class="max-w-56 px-3 py-2"><div class="truncate" :title="l.description">{{ l.description }}</div></td>
            <td class="px-3 py-2 font-mono text-[9px] text-muted-foreground">{{ l.before || l.after ? `${l.before ?? '—'} → ${l.after ?? '—'}` : '—' }}</td>
            <td class="px-3 py-2 font-mono text-[9px] text-muted-foreground">{{ l.ip }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
