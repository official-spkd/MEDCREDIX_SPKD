<script setup lang="ts">
// FacilitiesView — port dari views-readiness.tsx (FacilitiesView).
import { Building2 } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import { useIntel } from '~/composables/useIntel'
import { useMedCredix } from '~/stores/medcredix'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'
import StatusBadge from '~/components/medcredix/StatusBadge.vue'

const store = useMedCredix()
const { data: boot, loading } = useIntel()
const srcData = ref<any>(null)

async function loadSrc() {
  try { srcData.value = await api(`/data?type=sources${store.hospitalId ? `&hospitalId=${store.hospitalId}` : ''}`) } catch { srcData.value = {} }
}
onMounted(loadSrc)
watch(() => store.hospitalId, loadSrc)

const facilities = computed(() => srcData.value?.facilities ?? [])
const beds = computed(() => srcData.value?.beds ?? [])
const hospital = computed(() => srcData.value?.hospital)
const inList = (v: string, good: string[]) => good.includes(v)
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Sarana & Prasarana"
      sub="Lifecycle fasilitas: Ada → Berfungsi → Laik Pakai → Kalibrasi Valid → Sesuai Kebutuhan (perbaikan temuan QA: status diperinci, bukan sekadar Ada/Tidak)." demo
    />
    <div class="grid grid-cols-2 gap-3 md:grid-cols-4">
      <div class="rounded-xl border bg-card p-4"><div class="text-[10px] font-bold uppercase text-muted-foreground">Skor Sarana</div><div class="mt-1 font-display text-2xl font-bold text-brand">{{ boot?.intel?.sarana?.score ?? '—' }}</div></div>
      <div class="rounded-xl border bg-card p-4"><div class="text-[10px] font-bold uppercase text-muted-foreground">Fasilitas Wajib</div><div class="mt-1 font-display text-2xl font-bold text-foreground">{{ boot?.intel?.sarana?.wajib ?? '—' }}</div></div>
      <div class="rounded-xl border bg-card p-4"><div class="text-[10px] font-bold uppercase text-muted-foreground">TT Terlisensi</div><div class="mt-1 font-display text-2xl font-bold text-foreground">{{ hospital?.bedLicensed ?? 240 }}</div></div>
      <div class="rounded-xl border bg-card p-4"><div class="text-[10px] font-bold uppercase text-muted-foreground">TT Aktif</div><div class="mt-1 font-display text-2xl font-bold text-success">{{ hospital?.bedActive ?? 196 }}</div></div>
    </div>

    <SectionCard title="Checklist Lifecycle Fasilitas" sub="Standar pelayanan minimal RS mengacu Permenkes 6/2026">
      <template #icon><Building2 class="h-4 w-4" /></template>
      <div class="overflow-x-auto rounded-xl border">
        <table class="w-full text-xs">
          <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
            <tr>
              <th class="px-3 py-2.5">Fasilitas</th><th class="px-3 py-2.5">Ada</th><th class="px-3 py-2.5">Berfungsi</th>
              <th class="px-3 py-2.5">Laik Pakai</th><th class="px-3 py-2.5">Kalibrasi</th><th class="px-3 py-2.5">Kesesuaian</th><th class="px-3 py-2.5">Catatan</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="f in facilities" :key="f.id" class="border-t">
              <td class="px-3 py-2.5">
                <div class="font-medium">{{ f.name }}</div>
                <div class="text-[9px] text-muted-foreground">{{ f.category }}{{ f.standardRef ? ` · ${f.standardRef}` : '' }}</div>
              </td>
              <td class="px-3 py-2.5"><StatusBadge :status="inList(f.availability, ['Ada']) ? 'Ada' : 'Tidak Ada'" /></td>
              <td class="px-3 py-2.5"><StatusBadge :status="inList(f.functional, ['Berfungsi']) ? 'Baik' : 'Kritis'" /></td>
              <td class="px-3 py-2.5">
                <StatusBadge :status="inList(f.feasible, ['Laik Pakai']) ? 'Baik' : inList(f.feasible, ['Perlu Perbaikan']) ? 'Perlu Perhatian' : 'Kritis'" />
              </td>
              <td class="px-3 py-2.5"><StatusBadge :status="inList(f.calibration, ['Valid']) ? 'Baik' : 'Perlu Perhatian'" /></td>
              <td class="px-3 py-2.5 text-[10px]">{{ f.suitable }}</td>
              <td class="max-w-52 px-3 py-2.5 text-[10px] leading-snug text-muted-foreground">{{ f.notes ?? '—' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </SectionCard>

    <SectionCard title="Rincian Ranjang per Kelas" sub="Alasan TT tidak beroperasi dicatat eksplisit (perbaikan temuan QA produk lama)">
      <template #icon><Building2 class="h-4 w-4" /></template>
      <div class="overflow-x-auto rounded-xl border">
        <table class="w-full text-xs">
          <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
            <tr>
              <th class="px-3 py-2.5">Ruangan</th><th class="px-3 py-2.5">Kelas</th><th class="px-3 py-2.5 text-center">Terpasang</th>
              <th class="px-3 py-2.5 text-center">Aktif</th><th class="px-3 py-2.5 text-center">Layak Pakai</th><th class="px-3 py-2.5">Alasan tidak beroperasi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="b in beds" :key="b.id" class="border-t">
              <td class="px-3 py-2.5 font-medium">{{ b.ward }}</td>
              <td class="px-3 py-2.5">{{ b.bedClass }}</td>
              <td class="px-3 py-2.5 text-center font-display font-bold">{{ b.total }}</td>
              <td class="px-3 py-2.5 text-center font-display font-bold text-success">{{ b.active }}</td>
              <td class="px-3 py-2.5 text-center">{{ b.usable }}</td>
              <td class="px-3 py-2.5 text-[10px] text-muted-foreground">{{ b.nonOperatingReason ?? '—' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p v-if="hospital?.notOperatingNote" class="mt-3 rounded-lg bg-info-light px-3 py-2 text-[11px] leading-relaxed text-info">
        {{ hospital.notOperatingNote }}
      </p>
    </SectionCard>
  </div>
</template>
