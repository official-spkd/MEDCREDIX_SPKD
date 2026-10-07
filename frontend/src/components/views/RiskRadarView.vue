<script setup lang="ts">
// RiskRadarView — port dari views-risk.tsx.
import { Radar, CheckCircle2, UserPlus, PlusCircle, Fingerprint, Filter } from 'lucide-vue-next'
import { RISK_CATEGORIES, fmtDate } from '~/utils/format'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import Dialog from '~/components/ui/Dialog.vue'
import Select from '~/components/ui/Select.vue'
import SelectOption from '~/components/ui/SelectOption.vue'
import DemoBadge from '~/components/medcredix/DemoBadge.vue'
import EmptyState from '~/components/medcredix/EmptyState.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const level = ref('all')
const category = ref('all')
const status = ref('all')
const selected = ref<any>(null)
const data = ref<{ risks: any[]; evidence: any[] } | null>(null)
const loading = ref(true)

const LEVELS = ['Tinggi', 'Sedang', 'Rendah']

async function load() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (level.value !== 'all') params.set('level', level.value)
    if (category.value !== 'all') params.set('category', category.value)
    if (status.value !== 'all') params.set('status', status.value)
    if (store.hospitalId) params.set('hospitalId', store.hospitalId)
    data.value = await api(`/risks?${params}`)
  } catch { data.value = { risks: [], evidence: [] } }
  loading.value = false
}
onMounted(load)
watch([level, category, status, () => store.hospitalId], load)

async function patchRisk(r: any, body: any, msg: string) {
  try {
    await api('/risks', { method: 'PATCH', body: { id: r.id, ...body } })
    toast({ title: msg })
    load()
  } catch (e: any) {
    toast({ title: 'Gagal', description: e?.message, variant: 'destructive' })
  }
}

const summary = computed(() => {
  const s: Record<string, number> = { Tinggi: 0, Sedang: 0, Rendah: 0 }
  for (const r of data.value?.risks ?? []) s[r.level] = (s[r.level] ?? 0) + 1
  return s
})
</script>

<template>
  <div class="space-y-5">
    <PageHeader title="Hospital Risk Radar" sub="Pemantauan risiko lintas 9 kategori — setiap risiko berbasis bukti, bukan opini." demo>
      <template #actions><DemoBadge /></template>
    </PageHeader>

    <div class="grid grid-cols-3 gap-3">
      <div
        v-for="lv in LEVELS" :key="lv"
        :class="['rounded-xl border p-4', lv === 'Tinggi' ? 'border-danger/30 bg-danger-light/40' : lv === 'Sedang' ? 'border-warning/30 bg-warning-light/40' : 'border-info/30 bg-info-light/40']"
      >
        <div class="text-[10px] font-bold uppercase tracking-wide text-muted-foreground">Risiko {{ lv }}</div>
        <div :class="['mt-1 font-display text-2xl font-bold', lv === 'Tinggi' ? 'text-danger' : lv === 'Sedang' ? 'text-warning' : 'text-info']">{{ summary[lv] ?? 0 }}</div>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-2 rounded-xl border bg-card p-3">
      <Filter class="h-4 w-4 text-muted-foreground" />
      <div class="w-36"><Select v-model="level"><SelectOption value="all">Semua Level</SelectOption><SelectOption v-for="l in LEVELS" :key="l" :value="l">{{ l }}</SelectOption></Select></div>
      <div class="w-44"><Select v-model="category"><SelectOption value="all">Semua Kategori</SelectOption><SelectOption v-for="c in RISK_CATEGORIES" :key="c" :value="c">{{ c }}</SelectOption></Select></div>
      <div class="w-36"><Select v-model="status"><SelectOption value="all">Semua Status</SelectOption><SelectOption v-for="s in ['Terbuka', 'Ditinjau', 'Dalam Pengerjaan', 'Terkelola', 'Ditutup']" :key="s" :value="s">{{ s }}</SelectOption></Select></div>
      <span class="ml-auto text-xs text-muted-foreground">{{ data?.risks?.length ?? 0 }} risiko</span>
    </div>

    <div v-if="loading" class="py-16 text-center text-sm text-muted-foreground">Memuat risiko…</div>
    <EmptyState
      v-else-if="!data?.risks?.length" title="Tidak ada risiko pada filter ini"
      desc="Coba ubah filter — kondisi ideal ketika radar kosong."
    >
      <template #icon><Radar class="h-5 w-5" /></template>
    </EmptyState>
    <div v-else class="grid gap-3 lg:grid-cols-2">
      <div v-for="r in data.risks" :key="r.id" class="rounded-xl border bg-card p-4 transition-shadow hover:shadow-sm">
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-1.5">
              <StatusBadge :status="r.level" :tone="r.level === 'Tinggi' ? 'danger' : r.level === 'Sedang' ? 'warning' : 'info'" />
              <Badge variant="outline" class="text-[10px]">{{ r.category }}</Badge>
              <StatusBadge :status="r.status" />
            </div>
            <h3 class="mt-2 font-display text-sm font-semibold leading-snug">{{ r.title }}</h3>
          </div>
          <Button size="sm" variant="outline" class="h-7 shrink-0 text-[11px]" @click="selected = r">Detail</Button>
        </div>
        <p class="mt-1.5 line-clamp-2 text-xs leading-relaxed text-muted-foreground">{{ r.description }}</p>
        <div class="mt-3 flex flex-wrap items-center justify-between gap-2 border-t pt-3 text-[10px] text-muted-foreground">
          <span>{{ r.code }} · Deteksi {{ fmtDate(r.detectedAt) }} · Tenggat {{ fmtDate(r.deadline) }}</span>
          <span class="flex items-center gap-1.5 font-medium text-brand-2-dark">Saran owner: {{ r.suggestedOwner.split(' (')[0] }}</span>
        </div>
        <div class="mt-2.5 flex flex-wrap gap-1.5">
          <Button
            size="sm" variant="outline" class="h-7 gap-1 text-[11px]"
            @click="store.openAction({ title: r.recommendedAction.split(';')[0] || `Tindak lanjuti ${r.title}`, description: `Dari risiko ${r.code}: ${r.description}`, sourceType: 'Risiko', sourceRef: r.code, priority: r.level, suggestedOwner: r.suggestedOwner.split(' (')[0] })"
          >
            <PlusCircle class="h-3 w-3" /> Buat Tindakan
          </Button>
          <Button size="sm" variant="ghost" class="h-7 gap-1 text-[11px]" @click="patchRisk(r, { action: 'ack', status: 'Ditinjau' }, `${r.code} ditandai ditinjau`)">
            <CheckCircle2 class="h-3 w-3" /> Tandai Ditinjau
          </Button>
        </div>
      </div>
    </div>

    <!-- Detail drawer -->
    <Dialog :open="!!selected" @update:open="(o) => !o && (selected = null)">
      <template v-if="selected">
        <div class="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">
          {{ selected.code }} · Domain {{ selected.domain }} <DemoBadge />
        </div>
        <h2 class="mt-1 flex flex-wrap items-center gap-2 text-left text-lg font-bold leading-snug">
          {{ selected.title }}
          <StatusBadge :status="selected.level" :tone="selected.level === 'Tinggi' ? 'danger' : selected.level === 'Sedang' ? 'warning' : 'info'" />
        </h2>
        <div class="mt-4 space-y-4 text-sm">
          <p class="leading-relaxed text-foreground/90">{{ selected.description }}</p>
          <div class="grid gap-3 sm:grid-cols-2">
            <div class="rounded-xl border p-3.5">
              <div class="text-[10px] font-bold uppercase tracking-wide text-muted-foreground">Dampak Bisnis</div>
              <p class="mt-1 text-xs leading-relaxed">{{ selected.businessImpact }}</p>
            </div>
            <div class="rounded-xl border p-3.5">
              <div class="text-[10px] font-bold uppercase tracking-wide text-muted-foreground">Aksi Direkomendasikan</div>
              <p class="mt-1 text-xs leading-relaxed">{{ selected.recommendedAction }}</p>
            </div>
          </div>
          <div class="grid gap-3 text-xs sm:grid-cols-2">
            <div class="rounded-xl bg-muted/40 p-3.5"><span class="font-semibold">Kategori:</span> {{ selected.category }}</div>
            <div class="rounded-xl bg-muted/40 p-3.5"><span class="font-semibold">Record terdampak:</span> {{ selected.affectedRecords }}</div>
            <div class="rounded-xl bg-muted/40 p-3.5"><span class="font-semibold">Owner disarankan:</span> {{ selected.suggestedOwner }}</div>
            <div class="rounded-xl bg-muted/40 p-3.5"><span class="font-semibold">Tenggat:</span> {{ fmtDate(selected.deadline) }}</div>
          </div>
          <div class="flex flex-wrap gap-2">
            <Button
              class="gap-1.5 bg-brand hover:bg-brand-hover"
              @click="store.openAction({ title: selected.recommendedAction.split(';')[0], description: `Dari risiko ${selected.code}: ${selected.description}`, sourceType: 'Risiko', sourceRef: selected.code, priority: selected.level, suggestedOwner: selected.suggestedOwner.split(' (')[0] }); selected = null"
            >
              <PlusCircle class="h-4 w-4" /> Buat Tindakan
            </Button>
            <Button variant="outline" class="gap-1.5" @click="store.openEvidence({ refType: 'Risiko', refId: selected.code, title: `Bukti risiko ${selected.code}` })">
              <Fingerprint class="h-4 w-4" /> Lihat Bukti
            </Button>
          </div>
        </div>
      </template>
    </Dialog>
  </div>
</template>
