<script setup lang="ts">
// PractitionersView — port dari views-credentialing.tsx (PractitionersView).
import { Search, FileText, Eye, Clock } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { expiryLabel, daysFromNow, fmtDate } from '~/utils/format'
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import Dialog from '~/components/ui/Dialog.vue'
import Input from '~/components/ui/Input.vue'
import Select from '~/components/ui/Select.vue'
import SelectOption from '~/components/ui/SelectOption.vue'
import DemoBadge from '~/components/medcredix/DemoBadge.vue'
import Hint from '~/components/medcredix/Hint.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const q = ref('')
const prof = ref('all')
const status = ref('all')
const selected = ref<string | null>(null)
const data = ref<any>(null)
const detail = ref<any>(null)

async function load() {
  try {
    const params = new URLSearchParams()
    if (q.value) params.set('q', q.value)
    if (prof.value !== 'all') params.set('prof', prof.value)
    if (status.value !== 'all') params.set('status', status.value)
    if (store.hospitalId) params.set('hospitalId', store.hospitalId)
    data.value = await api(`/practitioners?${params}`)
  } catch { data.value = { practitioners: [], depts: [], total: 0 } }
}
onMounted(load)
watch([prof, status, () => store.hospitalId], load)
let t: ReturnType<typeof setTimeout> | null = null
watch(q, () => { if (t) clearTimeout(t); t = setTimeout(load, 300) })

async function loadDetail() {
  if (!selected.value) { detail.value = null; return }
  try { detail.value = await api(`/practitioners?id=${selected.value}`) } catch { detail.value = null }
}
watch(selected, loadDetail)

const list = computed(() => data.value?.practitioners ?? [])
const professions = computed(() => [...new Set(list.value.map((p: any) => p.profession))])
function toneCls(tone: string) {
  return tone === 'danger' ? 'bg-danger-light text-danger' : tone === 'warning' ? 'bg-warning-light text-warning' : tone === 'success' ? 'bg-success-light text-success' : 'bg-muted text-muted-foreground'
}
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Praktisi"
      sub="Master tenaga medis — status STR & SIP terlihat langsung (perbaikan dari temuan QA produk lama)." demo
    />
    <div class="flex flex-wrap items-center gap-2 rounded-xl border bg-card p-3">
      <div class="relative">
        <Search class="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
        <Input v-model="q" placeholder="Cari nama, spesialisasi, nomor STR/SIP…" class="h-8 w-64 pl-8 text-xs" />
      </div>
      <div class="w-48">
        <Select v-model="prof">
          <SelectOption value="all">Semua Profesi</SelectOption>
          <SelectOption v-for="p in professions" :key="p" :value="p">{{ p }}</SelectOption>
        </Select>
      </div>
      <div class="w-36">
        <Select v-model="status">
          <SelectOption value="all">Semua Status</SelectOption>
          <SelectOption v-for="s in ['Aktif', 'Cuti', 'Non-Aktif']" :key="s" :value="s">{{ s }}</SelectOption>
        </Select>
      </div>
      <span class="ml-auto text-xs text-muted-foreground">{{ data?.total ?? 0 }} praktisi</span>
    </div>

    <div class="overflow-x-auto rounded-xl border bg-card">
      <table class="w-full text-xs">
        <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
          <tr>
            <th class="px-3 py-2.5">Nama</th><th class="px-3 py-2.5">Profesi / Spesialisasi</th>
            <th class="px-3 py-2.5">Unit</th><th class="px-3 py-2.5">STR</th>
            <th class="px-3 py-2.5">SIP</th><th class="px-3 py-2.5">Status</th><th />
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="p in list" :key="p.id"
            class="cursor-pointer border-t transition-colors hover:bg-muted/30" @click="selected = p.id"
          >
            <td class="px-3 py-2.5 font-semibold">{{ p.fullName }}</td>
            <td class="px-3 py-2.5 text-muted-foreground">
              {{ p.profession }}{{ p.specialty ? ` — ${p.specialty}` : '' }}{{ p.subSpecialty ? ` (${p.subSpecialty})` : '' }}
            </td>
            <td class="px-3 py-2.5 text-muted-foreground">{{ p.department }}</td>
            <td class="px-3 py-2.5">
              <span :class="['inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-semibold', toneCls(expiryLabel(p.strExpiry).tone)]">
                <Clock class="h-3 w-3" /> {{ expiryLabel(p.strExpiry).text }}
              </span>
            </td>
            <td class="px-3 py-2.5">
              <StatusBadge v-if="daysFromNow(p.sipExpiry) != null && daysFromNow(p.sipExpiry)! < 0" status="Kedaluwarsa" />
              <span v-else :class="['inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-semibold', toneCls(expiryLabel(p.sipExpiry).tone)]">
                <Clock class="h-3 w-3" /> {{ expiryLabel(p.sipExpiry).text }}
              </span>
            </td>
            <td class="px-3 py-2.5"><StatusBadge :status="p.status" /></td>
            <td class="px-3 py-2.5 text-right"><Button size="sm" variant="ghost" class="h-6 text-[10px] text-brand-2">Detail</Button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Detail drawer -->
    <Dialog :open="!!selected" @update:open="(o) => !o && (selected = null)">
      <template v-if="detail">
        <div class="flex items-center gap-2 text-[11px] uppercase tracking-wider text-muted-foreground">
          {{ detail.practitioner.profession }} · {{ detail.practitioner.department }} <DemoBadge />
        </div>
        <h2 class="mt-1 text-left text-lg font-bold">{{ detail.practitioner.fullName }}</h2>
        <div class="mt-4 space-y-4">
          <div class="grid gap-2 text-xs sm:grid-cols-2">
            <div class="rounded-lg bg-muted/40 p-3">
              <span class="text-muted-foreground">Spesialisasi:</span> <strong>{{ detail.practitioner.specialty ?? '—' }}</strong>
              <template v-if="detail.practitioner.subSpecialty"> · Sub: {{ detail.practitioner.subSpecialty }}</template>
            </div>
            <div class="rounded-lg bg-muted/40 p-3">
              <span class="text-muted-foreground">NIK (disamarkan):</span> <strong>{{ detail.practitioner.nikMasked }}</strong>
              <Hint>Perlindungan Data Pribadi (UU 27/2022) — NIK penuh tidak pernah disimpan di platform analitik.</Hint>
            </div>
            <div class="rounded-lg bg-muted/40 p-3">
              <span class="text-muted-foreground">STR:</span> <strong>{{ detail.practitioner.strNumber }}</strong> · berlaku s.d. {{ fmtDate(detail.practitioner.strExpiry) }}
            </div>
            <div class="rounded-lg bg-muted/40 p-3">
              <span class="text-muted-foreground">SIP:</span> <strong>{{ detail.practitioner.sipNumber }}</strong> ·
              <span :class="['inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-semibold', toneCls(expiryLabel(detail.practitioner.sipExpiry).tone)]">
                <Clock class="h-3 w-3" /> {{ expiryLabel(detail.practitioner.sipExpiry).text }}
              </span>
            </div>
            <div class="rounded-lg bg-muted/40 p-3"><span class="text-muted-foreground">Status kepegawaian:</span> <strong>{{ detail.practitioner.employmentType }}</strong></div>
            <div class="rounded-lg bg-muted/40 p-3">
              <span class="text-muted-foreground">Standar kebutuhan:</span> <strong>{{ detail.standard?.required ?? 1 }} SDM</strong>
              {{ detail.standard?.note ? `· ${detail.standard.note}` : '' }}
            </div>
          </div>

          <SectionCard title="Kredensial & Dokumen" :sub="`${detail.credentials.length} kredensial tercatat — versi & jejak audit tersimpan`">
            <div class="space-y-2">
              <div v-for="c in detail.credentials" :key="c.id" class="rounded-xl border p-3">
                <div class="flex flex-wrap items-start justify-between gap-2">
                  <div class="min-w-0">
                    <div class="flex items-center gap-1.5">
                      <Badge variant="outline" class="text-[9px] font-bold">{{ c.type }}</Badge>
                      <StatusBadge :status="c.status" />
                      <span class="text-[9px] text-muted-foreground">v{{ c.version }}</span>
                    </div>
                    <div class="mt-1 truncate text-xs font-medium">{{ c.name }}</div>
                    <div class="mt-0.5 text-[10px] text-muted-foreground">{{ c.number }} · Terbit {{ fmtDate(c.issueDate) }} · Berakhir {{ fmtDate(c.expiryDate) }}</div>
                  </div>
                  <span :class="['inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-semibold', toneCls(expiryLabel(c.expiryDate).tone)]">
                    <Clock class="h-3 w-3" /> {{ expiryLabel(c.expiryDate).text }}
                  </span>
                </div>
                <div v-for="d in detail.documents.filter((x: any) => x.credentialId === c.id)" :key="d.id" class="mt-2 flex items-center gap-2 rounded-lg bg-muted/50 px-2.5 py-1.5 text-[10px]">
                  <FileText class="h-3.5 w-3.5 text-brand-2" />
                  <span class="truncate font-medium">{{ d.fileName }}</span>
                  <span class="text-muted-foreground">{{ d.sizeKb }} KB</span>
                  <StatusBadge :status="d.status" class="ml-auto" />
                  <Button
                    size="sm" variant="ghost" class="h-5 w-5 p-0" aria-label="Lihat dokumen"
                    @click="toast({ title: 'Pratinjau dokumen', description: `${d.fileName} — pada produksi, file dibuka dari object storage terenkripsi.` })"
                  >
                    <Eye class="h-3 w-3" />
                  </Button>
                </div>
              </div>
            </div>
          </SectionCard>
        </div>
      </template>
    </Dialog>
  </div>
</template>
