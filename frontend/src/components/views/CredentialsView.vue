<script setup lang="ts">
// CredentialsView — port dari views-credentialing.tsx (CredentialsView).
import { Clock } from 'lucide-vue-next'
import { CRED_LIFECYCLE, expiryLabel, fmtDate } from '~/utils/format'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import Input from '~/components/ui/Input.vue'
import Select from '~/components/ui/Select.vue'
import SelectOption from '~/components/ui/SelectOption.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const status = ref('all')
const type = ref('all')
const q = ref('')
const data = ref<any>(null)

async function load() {
  try {
    const params = new URLSearchParams()
    if (status.value !== 'all') params.set('status', status.value)
    if (type.value !== 'all') params.set('type', type.value)
    if (q.value) params.set('q', q.value)
    if (store.hospitalId) params.set('hospitalId', store.hospitalId)
    data.value = await api(`/credentials?${params}`)
  } catch { data.value = { credentials: [], total: 0 } }
}
onMounted(load)
watch([status, type, () => store.hospitalId], load)
let t: ReturnType<typeof setTimeout> | null = null
watch(q, () => { if (t) clearTimeout(t); t = setTimeout(load, 300) })

async function transition(c: any, action: string) {
  try {
    await api('/credentials', { method: 'PATCH', body: { id: c.id, action } })
    toast({ title: `${c.type} — ${action}`, description: 'Transisi status berhasil & tercatat di audit trail.' })
    load()
  } catch (e: any) {
    toast({ title: 'Transisi ditolak', description: e?.message, variant: 'destructive' })
  }
}

const canVerify = computed(() => store.user?.permissions.includes('VERIFY'))
const canApprove = computed(() => store.user?.permissions.includes('APPROVE'))
const list = computed(() => data.value?.credentials ?? [])
function toneCls(tone: string) {
  return tone === 'danger' ? 'bg-danger-light text-danger' : tone === 'warning' ? 'bg-warning-light text-warning' : tone === 'success' ? 'bg-success-light text-success' : 'bg-muted text-muted-foreground'
}
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Kredensial"
      :sub="`Siklus hidup: ${CRED_LIFECYCLE.slice(0, 6).join(' → ')} → … — record yang telah disetujui tidak pernah diubah diam-diam (versioning otomatis).`" demo
    />
    <div class="flex flex-wrap gap-1.5 rounded-xl border bg-card p-3">
      <button
        :class="['inline-flex h-7 items-center rounded-lg px-3 text-xs font-semibold', status === 'all' ? 'bg-brand text-white' : 'border border-input bg-card hover:bg-muted/60']"
        @click="status = 'all'"
      >Semua ({{ data?.total ?? 0 }})</button>
      <button
        v-for="s in CRED_LIFECYCLE" :key="s"
        :class="['inline-flex h-7 items-center rounded-lg px-3 text-xs font-semibold', status === s ? 'bg-brand text-white' : 'border border-input bg-card hover:bg-muted/60']"
        @click="status = s"
      >{{ s }}</button>
    </div>
    <div class="flex flex-wrap items-center gap-2">
      <Input v-model="q" placeholder="Cari nama dokumen / nomor…" class="h-8 w-60 text-xs" />
      <div class="w-44">
        <Select v-model="type">
          <SelectOption value="all">Semua Jenis</SelectOption>
          <SelectOption v-for="t in ['STR', 'SIP', 'Ijazah', 'Sertifikat Kompetensi', 'Pelatihan']" :key="t" :value="t">{{ t }}</SelectOption>
        </Select>
      </div>
    </div>

    <div class="overflow-x-auto rounded-xl border bg-card">
      <table class="w-full text-xs">
        <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
          <tr><th class="px-3 py-2.5">Praktisi</th><th class="px-3 py-2.5">Jenis / Dokumen</th><th class="px-3 py-2.5">Berakhir</th><th class="px-3 py-2.5">Status</th><th class="px-3 py-2.5">Aksi</th></tr>
        </thead>
        <tbody>
          <tr v-for="c in list.slice(0, 100)" :key="c.id" class="border-t">
            <td class="px-3 py-2.5"><div class="font-semibold">{{ c.practitionerName }}</div><div class="text-[10px] text-muted-foreground">{{ c.profession }}</div></td>
            <td class="px-3 py-2.5">
              <div class="flex items-center gap-1.5"><Badge variant="outline" class="text-[9px]">{{ c.type }}</Badge><span class="text-muted-foreground">{{ c.name.split('—')[0] }}</span></div>
              <div class="text-[10px] text-muted-foreground">{{ c.number }} · v{{ c.version }}</div>
            </td>
            <td class="px-3 py-2.5">
              <span :class="['inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-semibold', toneCls(expiryLabel(c.expiryDate).tone)]">
                <Clock class="h-3 w-3" /> {{ expiryLabel(c.expiryDate).text }}
              </span>
              <div class="mt-0.5 text-[10px] text-muted-foreground">{{ fmtDate(c.expiryDate) }}</div>
            </td>
            <td class="px-3 py-2.5"><StatusBadge :status="c.status" /></td>
            <td class="px-3 py-2.5">
              <div class="flex flex-wrap gap-1">
                <Button v-if="c.status === 'Diajukan' && canVerify" size="sm" variant="outline" class="h-6 px-2 text-[10px]" @click="transition(c, 'verify')">Verifikasi</Button>
                <Button v-if="c.status === 'Verifikasi' && canVerify" size="sm" variant="outline" class="h-6 px-2 text-[10px]" @click="transition(c, 'review')">Ke Komite</Button>
                <Button v-if="c.status === 'Tinjauan Komite' && canApprove" size="sm" class="h-6 bg-brand-2 px-2 text-[10px] hover:bg-brand-2-hover" @click="transition(c, 'approve')">Setujui</Button>
                <Button v-if="c.status === 'Disetujui' && canApprove" size="sm" class="h-6 bg-success px-2 text-[10px] hover:bg-success" @click="transition(c, 'activate')">Aktifkan</Button>
                <Button v-if="['Aktif', 'Segera Berakhir'].includes(c.status) && canApprove" size="sm" variant="outline" class="h-6 px-2 text-[10px] text-danger" @click="transition(c, 'suspend')">Tangguhkan</Button>
                <Button v-if="c.status === 'Ditangguhkan' && canApprove" size="sm" variant="outline" class="h-6 px-2 text-[10px]" @click="transition(c, 'reactivate')">Aktifkan</Button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
