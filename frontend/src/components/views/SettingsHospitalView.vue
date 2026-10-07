<script setup lang="ts">
// SettingsHospitalView — port dari views-settings.tsx.
import { Building, Cable, Save } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import Input from '~/components/ui/Input.vue'
import Label from '~/components/ui/Label.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const data = ref<any>(null)
const src = ref<any>(null)
const form = ref<any>(null)
const saving = ref(false)

onMounted(async () => {
  try {
    data.value = await api('/settings')
    src.value = await api('/data?type=sources')
  } catch { /* ignore */ }
})
const h = computed(() => form.value ?? data.value?.hospital)
const canEdit = computed(() => store.user?.permissions.includes('CONFIGURE'))
const fields = [
  ['name', 'Nama RS'], ['class', 'Kelas (A/B/C/D)'], ['city', 'Kota'], ['province', 'Provinsi'],
  ['bpjsFacilityCode', 'Kode Faskes BPJS'], ['satusehatOrgId', 'SATUSEHAT Organization ID'],
]

function setField(k: string, v: any) {
  form.value = { ...h.value, [k]: v }
}

async function save() {
  saving.value = true
  try {
    await api('/settings', {
      method: 'PATCH',
      body: {
        target: 'hospital', name: h.value.name, class: h.value.class, city: h.value.city,
        province: h.value.province, bedLicensed: Number(h.value.bedLicensed),
        bpjsFacilityCode: h.value.bpjsFacilityCode, satusehatOrgId: h.value.satusehatOrgId,
      },
    })
    toast({ title: 'Profil RS tersimpan', description: 'Perubahan tercatat di Audit Trail.' })
    data.value = await api('/settings')
    form.value = null
  } catch (e: any) {
    toast({ title: 'Gagal menyimpan', description: e?.message, variant: 'destructive' })
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Profil RS & Integrasi"
      sub="Konfigurasi koneksi API RS (GET-only), SATUSEHAT Organization ID, dan kredensial BPJS. Secret tidak pernah dikirim balik ke browser." demo
    />
    <div class="grid gap-4 lg:grid-cols-2">
      <SectionCard title="Profil Rumah Sakit" sub="Hanya Admin RS & Super Admin dapat menyimpan">
        <template #icon><Building class="h-4 w-4" /></template>
        <template #action>
          <Button v-if="canEdit" size="sm" class="h-7 gap-1 bg-brand text-xs hover:bg-brand-hover" :disabled="saving" @click="save">
            <Save class="h-3 w-3" /> Simpan
          </Button>
        </template>
        <div class="grid gap-3">
          <div v-for="[k, label] in fields" :key="k" class="grid gap-1">
            <Label class="text-xs font-semibold">{{ label }}</Label>
            <Input :model-value="(h?.[k] as any) ?? ''" :disabled="!canEdit" class="text-xs" @update:model-value="(v) => setField(k, v)" />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="grid gap-1">
              <Label class="text-xs font-semibold">Jumlah TT Terlisensi</Label>
              <Input :model-value="h?.bedLicensed ?? 0" type="number" :disabled="!canEdit" class="text-xs" @update:model-value="(v) => setField('bedLicensed', v)" />
            </div>
            <div class="grid gap-1">
              <Label class="text-xs font-semibold">Tier Langganan SPKD</Label>
              <Input :model-value="h?.tier ?? ''" disabled class="text-xs" />
            </div>
          </div>
        </div>
      </SectionCard>

      <SectionCard title="Konfigurasi API RS" sub="Base URL FHIR internal + autentikasi OAuth 2.0 (direkomendasikan) — diuji lewat tombol Test Koneksi">
        <template #icon><Cable class="h-4 w-4" /></template>
        <div class="space-y-2.5">
          <div v-for="s in (src?.sources ?? []).slice(0, 6)" :key="s.id" class="flex items-center justify-between gap-3 rounded-xl border p-3">
            <div class="min-w-0">
              <div class="truncate text-xs font-semibold">{{ s.name }}</div>
              <div class="text-[10px] text-muted-foreground">{{ s.authMethod }}</div>
            </div>
            <div class="flex items-center gap-2">
              <Badge variant="outline" class="text-[9px]">Secret tersimpan ●</Badge>
              <StatusBadge :status="s.status" />
            </div>
          </div>
          <div class="rounded-lg bg-info-light px-3 py-2.5 text-[11px] leading-relaxed text-info">
            Platform hanya melakukan panggilan <strong>GET</strong> ke sistem RS — tidak pernah POST/PUT/DELETE. Data pasien mentah dibuang segera setelah dihitung; hanya hasil komputasi &amp; metadata yang disimpan.
          </div>
        </div>
      </SectionCard>
    </div>
  </div>
</template>
