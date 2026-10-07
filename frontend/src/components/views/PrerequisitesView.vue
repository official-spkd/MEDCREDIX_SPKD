<script setup lang="ts">
// PrerequisitesView — port dari views-credentialing.tsx (PrerequisitesView).
import { Info, FileText, Eye, Upload, CheckCircle2, Clock } from 'lucide-vue-next'
import { expiryLabel, fmtDate } from '~/utils/format'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import Dialog from '~/components/ui/Dialog.vue'
import Input from '~/components/ui/Input.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const data = ref<any>(null)
const verifyTarget = ref<any>(null)
const docNumber = ref('')
const docExpiry = ref('')

async function load() {
  try { data.value = await api('/requirements') } catch { data.value = { requirements: [] } }
}
onMounted(load)
const reqs = computed(() => data.value?.requirements ?? [])
const p3 = computed(() => reqs.value.find((r: any) => r.code === 'P3'))

async function submitReq(r: any, action: string, extra?: any) {
  try {
    await api('/requirements', { method: 'PATCH', body: { id: r.id, action, ...extra } })
    toast({ title: `${r.code}: ${action === 'verify' ? 'Terverifikasi' : action === 'submit' ? 'Diajukan' : 'Diperbarui'}` })
    load()
    verifyTarget.value = null
  } catch (e: any) {
    toast({ title: 'Gagal', description: e?.message, variant: 'destructive' })
  }
}
function toneCls(tone: string) {
  return tone === 'danger' ? 'bg-danger-light text-danger' : tone === 'warning' ? 'bg-warning-light text-warning' : tone === 'success' ? 'bg-success-light text-success' : 'bg-muted text-muted-foreground'
}
const canEdit = computed(() => store.user?.permissions.includes('EDIT'))
const canVerify = computed(() => store.user?.permissions.includes('VERIFY'))
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Prasyarat Wajib"
      sub="Requirement Engine — dokumen legal dapat dikonfigurasi (kode, sumber regulasi, versi, aturan validasi). Salah satu P1–P3 tidak Valid → Stage 0." demo
    />
    <div class="rounded-xl border border-warning/40 bg-warning-light/40 p-4">
      <div class="flex items-start gap-2.5">
        <Info class="mt-0.5 h-4 w-4 shrink-0 text-warning" />
        <p class="text-xs leading-relaxed text-foreground/90">
          <strong>Pakta Integritas Data aktif</strong> — data yang diakses MedCredix dinyatakan sebagai data operasional nyata, bukan data yang disiapkan khusus untuk penilaian.
          {{ !p3 || p3.status !== 'Valid' ? ' Perhatian: prasyarat belum lengkap.' : ` Akreditasi ${p3.docNumber} berlaku hingga ${fmtDate(p3.docExpiry)}.` }}
        </p>
      </div>
    </div>

    <div class="grid gap-3 lg:grid-cols-2">
      <div v-for="r in reqs" :key="r.id" class="rounded-xl border bg-card p-4">
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-1.5">
              <Badge class="bg-brand text-[10px] font-bold hover:bg-brand">{{ r.code }}</Badge>
              <Badge v-if="r.required" variant="outline" class="text-[9px] font-bold text-danger">WAJIB</Badge>
              <Badge v-else variant="outline" class="text-[9px]">Opsional</Badge>
              <StatusBadge :status="r.status" />
            </div>
            <h3 class="mt-2 font-display text-sm font-semibold leading-snug">{{ r.name }}</h3>
            <div class="mt-1 space-y-0.5 text-[10px] text-muted-foreground">
              <div>Sumber regulasi: {{ r.regulatorySource }} · {{ r.version }}</div>
              <div>Aturan validasi: {{ r.validationRule }}</div>
              <div>
                No. dokumen: <strong class="text-foreground">{{ r.docNumber ?? '—' }}</strong>
                {{ r.docExpiry ? ` · Berakhir ${fmtDate(r.docExpiry)}` : ' · Tanpa masa berlaku' }}
              </div>
              <div v-if="r.fileName" class="flex items-center gap-1">
                <FileText class="h-3 w-3" />{{ r.fileName }}
                <Button
                  size="sm" variant="ghost" class="h-5 w-5 p-0" aria-label="Lihat"
                  @click="toast({ title: 'Pratinjau dokumen', description: `${r.fileName} — dibuka dari object storage terenkripsi.` })"
                ><Eye class="h-3 w-3" /></Button>
              </div>
            </div>
          </div>
          <span :class="['inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-semibold shrink-0', toneCls(expiryLabel(r.docExpiry).tone)]">
            <Clock class="h-3 w-3" /> {{ expiryLabel(r.docExpiry).text }}
          </span>
        </div>
        <div class="mt-3 flex flex-wrap gap-1.5 border-t pt-3">
          <Button v-if="r.status === 'Belum Dilengkapi' && canEdit" size="sm" variant="outline" class="h-7 gap-1 text-xs" @click="verifyTarget = r">
            <Upload class="h-3.5 w-3.5" /> Unggah & Ajukan
          </Button>
          <Button v-if="r.status === 'Menunggu Verifikasi' && canVerify" size="sm" class="h-7 gap-1 bg-success text-xs hover:bg-success" @click="submitReq(r, 'verify')">
            <CheckCircle2 class="h-3.5 w-3.5" /> Verifikasi
          </Button>
          <Button
            v-if="r.fileName" size="sm" variant="ghost" class="h-7 gap-1 text-xs"
            @click="toast({ title: 'Pratinjau', description: `${r.fileName} dibuka dari penyimpanan terenkripsi.` })"
          ><Eye class="h-3.5 w-3.5" /> Lihat Dokumen</Button>
          <span class="ml-auto self-center text-[9px] text-muted-foreground">Owner: {{ r.responsibleRole }}</span>
        </div>
      </div>
    </div>

    <!-- Dialog unggah -->
    <Dialog :open="!!verifyTarget" @update:open="(o) => !o && (verifyTarget = null)">
      <template v-if="verifyTarget">
        <h2 class="font-display text-base font-semibold">Unggah & Ajukan — {{ verifyTarget.code }}</h2>
        <p class="mt-1 text-xs text-muted-foreground">{{ verifyTarget.name }}. Setelah diajukan, status berubah menjadi Menunggu Verifikasi.</p>
        <div class="mt-4 grid gap-3">
          <div>
            <label class="text-xs font-semibold">Nomor dokumen</label>
            <Input v-model="docNumber" class="mt-1 text-xs" placeholder="Mis. PKS/0310/1026/2026" />
          </div>
          <div>
            <label class="text-xs font-semibold">Tanggal kedaluwarsa</label>
            <Input v-model="docExpiry" type="date" class="mt-1 text-xs" />
          </div>
          <div class="rounded-lg border border-dashed p-4 text-center text-[11px] text-muted-foreground">
            <Upload class="mx-auto mb-1.5 h-4 w-4" /> Tarik berkas PDF/JPG/PNG ke sini (maks 10 MB) — tersimpan terenkripsi.
          </div>
        </div>
        <div class="mt-5 flex justify-end gap-2">
          <Button variant="outline" @click="verifyTarget = null">Batal</Button>
          <Button class="bg-brand hover:bg-brand-hover" @click="verifyTarget && submitReq(verifyTarget, 'submit', { docNumber, docExpiry })">
            Simpan & Ajukan Verifikasi
          </Button>
        </div>
      </template>
    </Dialog>
  </div>
</template>
