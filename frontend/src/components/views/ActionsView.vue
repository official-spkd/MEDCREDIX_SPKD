<script setup lang="ts">
// ActionsView — Action Center closed-loop (port dari views-actions.tsx).
import { CheckCircle2, PlusCircle, MessageSquare, History, UserPlus, RotateCcw, ClipboardCheck, Send, Clock } from 'lucide-vue-next'
import { ACTION_STATUSES, fmtDate, fmtDateTime } from '~/utils/format'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import Button from '~/components/ui/Button.vue'
import Dialog from '~/components/ui/Dialog.vue'
import Input from '~/components/ui/Input.vue'
import Select from '~/components/ui/Select.vue'
import SelectOption from '~/components/ui/SelectOption.vue'
import Textarea from '~/components/ui/Textarea.vue'
import DemoBadge from '~/components/medcredix/DemoBadge.vue'
import EmptyState from '~/components/medcredix/EmptyState.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const status = ref('all')
const priority = ref('all')
const search = ref('')
const selected = ref<string | null>(null)
const comment = ref('')
const data = ref<any>(null)
const detail = ref<any>(null)

async function load() {
  try {
    const params = new URLSearchParams()
    if (status.value !== 'all') params.set('status', status.value)
    if (priority.value !== 'all') params.set('priority', priority.value)
    if (search.value) params.set('q', search.value)
    if (store.hospitalId) params.set('hospitalId', store.hospitalId)
    data.value = await api(`/actions?${params}`)
  } catch { data.value = { actions: [], stats: {} } }
}
onMounted(load)
watch([status, priority, () => store.hospitalId], load)
let searchTimer: ReturnType<typeof setTimeout> | null = null
watch(search, () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(load, 300)
})

async function loadDetail() {
  if (!selected.value) { detail.value = null; return }
  try { detail.value = await api(`/actions?id=${selected.value}`) } catch { detail.value = null }
}
watch(selected, loadDetail)

async function patch(body: any, msg: string) {
  try {
    await api('/actions', { method: 'PATCH', body: { id: selected.value, ...body } })
    toast({ title: msg })
    comment.value = ''
    loadDetail()
    load()
  } catch (e: any) {
    toast({ title: 'Gagal', description: e?.message, variant: 'destructive' })
  }
}

const actions = computed(() => data.value?.actions ?? [])
const stats = computed(() => data.value?.stats)
function isOverdue(a: any) {
  return a.dueDate && new Date(a.dueDate) < new Date() && !['Selesai', 'Ditutup'].includes(a.status)
}
const assignOwner = ref('')
</script>

<template>
  <div class="space-y-5">
    <PageHeader title="Action Center" sub="Platform yang mengubah analitik menjadi pekerjaan — hingga skor membaik." demo>
      <template #actions>
        <Button size="sm" class="gap-1.5 bg-brand hover:bg-brand-hover" @click="store.openAction({ title: '', sourceType: 'Manual', description: '' })">
          <PlusCircle class="h-3.5 w-3.5" /> Tindakan Baru
        </Button>
      </template>
    </PageHeader>

    <div class="grid grid-cols-2 gap-3 md:grid-cols-4">
      <div v-for="s in [
        { label: 'Total Aksi', value: stats?.total ?? 0, tone: 'info' },
        { label: 'Aktif', value: stats?.open ?? 0, tone: 'warning' },
        { label: 'Prioritas Tinggi', value: stats?.high ?? 0, tone: 'danger' },
        { label: 'Terselesaikan', value: stats?.resolved ?? 0, tone: 'success' },
      ]" :key="s.label" class="rounded-xl border bg-card p-4">
        <div class="text-[10px] font-bold uppercase tracking-wide text-muted-foreground">{{ s.label }}</div>
        <div :class="['mt-1 font-display text-2xl font-bold', s.tone === 'danger' ? 'text-danger' : s.tone === 'warning' ? 'text-warning' : s.tone === 'success' ? 'text-success' : 'text-brand']">{{ s.value }}</div>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-2 rounded-xl border bg-card p-3">
      <Input v-model="search" placeholder="Cari kode, judul, atau owner…" class="h-8 w-56 text-xs" />
      <div class="w-44">
        <Select v-model="status">
          <SelectOption value="all">Semua Status</SelectOption>
          <SelectOption v-for="s in ACTION_STATUSES" :key="s" :value="s">{{ s }}</SelectOption>
        </Select>
      </div>
      <div class="w-36">
        <Select v-model="priority">
          <SelectOption value="all">Semua Prioritas</SelectOption>
          <SelectOption v-for="p in ['Tinggi', 'Sedang', 'Rendah']" :key="p" :value="p">{{ p }}</SelectOption>
        </Select>
      </div>
      <span class="ml-auto text-xs text-muted-foreground">{{ actions.length }} aksi</span>
    </div>

    <EmptyState
      v-if="actions.length === 0" title="Tidak ada aksi pada filter ini"
      desc="Buat aksi baru dari rekomendasi skor, risiko, atau peringatan dini."
    >
      <template #icon><CheckCircle2 class="h-5 w-5" /></template>
    </EmptyState>
    <div v-else class="grid gap-3 lg:grid-cols-2">
      <button
        v-for="a in actions" :key="a.id"
        class="rounded-xl border bg-card p-4 text-left transition-all hover:border-brand-2/50 hover:shadow-sm"
        @click="selected = a.id"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-1.5">
              <StatusBadge :status="a.priority" :tone="a.priority === 'Tinggi' ? 'danger' : a.priority === 'Sedang' ? 'warning' : 'info'" />
              <StatusBadge :status="a.status" />
              <span v-if="isOverdue(a)" class="rounded-full bg-danger-light px-2 py-0.5 text-[10px] font-bold text-danger">Terlambat</span>
            </div>
            <h3 class="mt-2 font-display text-sm font-semibold leading-snug">{{ a.title }}</h3>
            <div class="mt-1.5 flex flex-wrap items-center gap-x-2 gap-y-1 text-[10px] text-muted-foreground">
              <span class="flex items-center gap-1"><UserPlus class="h-3 w-3" />{{ a.owner.split(',')[0] }}</span>
              <span :class="['flex items-center gap-1', isOverdue(a) ? 'font-bold text-danger' : '']"><Clock class="h-3 w-3" />{{ fmtDate(a.dueDate) }}</span>
              <span>· {{ a.sourceType }}{{ a.sourceRef ? ` (${a.sourceRef})` : '' }}</span>
            </div>
          </div>
          <span class="font-display text-xs font-bold text-muted-foreground">{{ a.code }}</span>
        </div>
      </button>
    </div>

    <!-- Detail drawer -->
    <Dialog :open="!!selected" @update:open="(o) => !o && (selected = null)">
      <template v-if="detail">
        <div class="flex flex-wrap items-center gap-2 text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">
          {{ detail.action.code }} · dari {{ detail.action.sourceType }}{{ detail.action.sourceRef ? ` (${detail.action.sourceRef})` : '' }} <DemoBadge />
        </div>
        <h2 class="mt-1 text-left text-lg font-bold leading-snug">{{ detail.action.title }}</h2>

        <div class="mt-4 space-y-4">
          <div class="flex flex-wrap gap-1.5">
            <StatusBadge :status="detail.action.priority" :tone="detail.action.priority === 'Tinggi' ? 'danger' : detail.action.priority === 'Sedang' ? 'warning' : 'info'" />
            <StatusBadge :status="detail.action.status" />
          </div>
          <p class="text-sm leading-relaxed text-foreground/90">{{ detail.action.description }}</p>
          <div class="grid gap-2 text-xs sm:grid-cols-2">
            <div class="rounded-lg bg-muted/40 p-3"><span class="font-semibold">Owner:</span> {{ detail.action.owner }}</div>
            <div class="rounded-lg bg-muted/40 p-3"><span class="font-semibold">Tenggat:</span> {{ fmtDate(detail.action.dueDate) }}</div>
            <div class="rounded-lg bg-muted/40 p-3"><span class="font-semibold">Hasil diharapkan:</span> {{ detail.action.expectedOutcome ?? '—' }}</div>
            <div class="rounded-lg bg-muted/40 p-3"><span class="font-semibold">Dibuat:</span> {{ fmtDateTime(detail.action.createdAt) }}</div>
          </div>
          <div v-if="detail.action.progressNote" class="rounded-lg border border-info/30 bg-info-light/40 p-3 text-xs">
            <span class="font-semibold">Progres terakhir:</span> {{ detail.action.progressNote }}
          </div>
          <div v-if="detail.action.resolutionNote" class="rounded-lg border border-success/30 bg-success-light/40 p-4">
            <div class="flex items-center gap-1.5 text-xs font-semibold text-success"><ClipboardCheck class="h-4 w-4" /> Resolusi (closed-loop)</div>
            <p class="mt-1 text-xs leading-relaxed">{{ detail.action.resolutionNote }}</p>
            <p class="mt-2 text-[10px] text-muted-foreground">Selesai {{ fmtDateTime(detail.action.resolvedAt) }} — skor terkait diusulkan untuk dihitung ulang; lihat Tren pada drawer skor.</p>
          </div>

          <!-- Alur kerja -->
          <div class="rounded-xl border p-4">
            <div class="mb-2.5 flex items-center gap-1.5 text-xs font-semibold"><History class="h-4 w-4 text-brand-2" /> Alur Kerja</div>
            <div class="flex flex-wrap gap-1.5">
              <button
                v-for="s in ACTION_STATUSES" :key="s"
                :class="[
                  'rounded-lg border px-2.5 py-1.5 text-[10px] font-semibold transition-colors',
                  detail.action.status === s ? 'border-brand bg-brand text-white'
                    : ACTION_STATUSES.indexOf(s) < ACTION_STATUSES.indexOf(detail.action.status) ? 'border-success/40 bg-success-light/50 text-success'
                    : 'text-muted-foreground hover:border-brand-2/50',
                ]"
                @click="patch({ action: 'status', status: s }, `Status → ${s}`)"
              >
                {{ s }}
              </button>
            </div>
            <div class="mt-3 flex flex-wrap items-center gap-1.5">
              <div class="w-52">
                <Select v-model="assignOwner" placeholder="Tugaskan ulang…">
                  <SelectOption v-for="u in detail.users" :key="u.id" :value="u.name">{{ u.name }}</SelectOption>
                </Select>
              </div>
              <Button v-if="assignOwner" size="sm" variant="outline" class="h-8 text-xs" @click="patch({ action: 'assign', owner: assignOwner }, 'Penanggung jawab diperbarui'); assignOwner = ''">
                <UserPlus class="h-3.5 w-3.5" /> Tugaskan
              </Button>
              <Input type="date" class="h-8 w-40 text-xs" aria-label="Ubah tenggat" @change="(e: any) => e.target.value && patch({ action: 'assign', dueDate: e.target.value }, 'Tenggat diperbarui')" />
              <Button size="sm" variant="outline" class="h-8 gap-1 text-xs" @click="patch({ action: 'status', status: 'Dibuka Kembali' }, 'Aksi dibuka kembali')">
                <RotateCcw class="h-3.5 w-3.5" /> Buka Kembali
              </Button>
            </div>
          </div>

          <!-- Komentar & riwayat -->
          <div class="rounded-xl border p-4">
            <div class="mb-2.5 flex items-center gap-1.5 text-xs font-semibold"><MessageSquare class="h-4 w-4 text-brand-2" /> Komentar & Riwayat</div>
            <div class="mb-3 flex gap-2">
              <Textarea v-model="comment" :rows="2" placeholder="Tambahkan catatan progres…" class="text-xs" />
              <Button size="sm" class="gap-1 self-end bg-brand-2 hover:bg-brand-2-hover" @click="comment.trim() && patch({ action: 'comment', note: comment }, 'Komentar ditambahkan')">
                <Send class="h-3.5 w-3.5" /> Kirim
              </Button>
            </div>
            <div class="max-h-56 space-y-2 overflow-y-auto pr-1 scrollbar-thin">
              <div v-for="e in [...detail.events].reverse()" :key="e.id" class="flex gap-2.5 rounded-lg bg-muted/40 p-2.5">
                <div class="mt-0.5 rounded-md bg-card p-1.5 text-brand-2">
                  <PlusCircle v-if="e.type === 'created'" class="h-3 w-3" />
                  <MessageSquare v-else-if="e.type === 'comment'" class="h-3 w-3" />
                  <UserPlus v-else-if="e.type === 'assigned'" class="h-3 w-3" />
                  <ClipboardCheck v-else-if="e.type === 'resolved'" class="h-3 w-3" />
                  <History v-else class="h-3 w-3" />
                </div>
                <div class="min-w-0 text-[11px] leading-relaxed">
                  <div>
                    <span class="font-semibold">{{ e.byName }}</span>
                    <span class="text-muted-foreground">
                      {{ e.type === 'status' ? `mengubah status ${e.fromValue} → ${e.toValue}`
                        : e.type === 'created' ? 'membuat aksi'
                        : e.type === 'assigned' ? `menugaskan ke ${e.toValue}`
                        : e.type === 'comment' ? 'menulis'
                        : e.type === 'resolved' ? `menyelesaikan (${e.toValue})` : e.type }}
                    </span>
                  </div>
                  <div v-if="e.note" class="mt-0.5 text-muted-foreground">{{ e.note }}</div>
                  <div class="mt-0.5 text-[9px] text-muted-foreground/70">{{ fmtDateTime(e.createdAt) }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </Dialog>
  </div>
</template>
