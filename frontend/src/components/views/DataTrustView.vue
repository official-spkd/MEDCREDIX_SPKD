<script setup lang="ts">
// DataTrustView — port dari views-data.tsx (DataTrustView).
import { ShieldAlert } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { api } from '~/composables/useApi'
import Skeleton from '~/components/ui/Skeleton.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import ScoreRing from '~/components/medcredix/ScoreRing.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'

const data = ref<any>(null)
onMounted(async () => {
  try { data.value = await api('/data?type=trust') } catch { data.value = null }
})
const t = computed(() => data.value?.trust)
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Data Trust Score"
      sub="Kepercayaan data dihitung transparan dari 9 faktor berbobot. Tidak ada multiplier tersembunyi — jika skor berubah, Anda tahu penyebabnya." demo
    />
    <template v-if="!data">
      <Skeleton class="h-8 w-64" /><Skeleton class="mt-4 h-72" />
    </template>
    <template v-else>
      <div class="grid gap-4 lg:grid-cols-[300px_1fr]">
        <div class="space-y-4">
          <div class="rounded-2xl border bg-card p-5 text-center">
            <ScoreRing :value="t.score" :size="124" :tone="t.score >= 90 ? 'success' : t.score >= 70 ? 'warning' : 'danger'" label="Data Trust Score" />
            <div class="mt-4 rounded-xl bg-muted/50 p-3 text-left font-mono text-[10px] leading-relaxed text-muted-foreground">
              Trust {{ t.score }} → multiplier <strong class="text-foreground">×{{ t.multiplier }}</strong><br />
              BDCS mentah {{ data.bcds.raw }} × {{ t.multiplier }} = <strong class="text-foreground">{{ data.bcds.final }}</strong>
            </div>
          </div>
          <div :class="['rounded-xl border p-4', t.underReview ? 'border-danger/40 bg-danger-light/40' : 'border-success/30 bg-success-light/40']">
            <div class="flex items-center gap-2 text-xs font-semibold">
              <ShieldAlert :class="['h-4 w-4', t.underReview ? 'text-danger' : 'text-success']" />
              {{ t.underReview ? 'Status: UNDER REVIEW' : `Status: Normal (${t.anomalyFlags} anomaly flag aktif)` }}
            </div>
            <p class="mt-1.5 text-[11px] leading-relaxed text-muted-foreground">
              {{ t.underReview
                ? '3 atau lebih sinyal anomali aktif — rumah sakit berstatus Under Review sampai investigasi selesai.'
                : 'Di bawah 3 sinyal anomali aktif — status normal. Sinyal anomali memantau potensi manipulasi data (mis. lonjakan volume tidak wajar) dan selalu memerlukan konfirmasi.' }}
            </p>
          </div>
        </div>

        <div class="space-y-4">
          <SectionCard title="9 Faktor Penyusun Trust" sub="Bobot & skor per faktor — dapat ditelusuri">
            <div class="space-y-2">
              <div v-for="f in t.factors ?? []" :key="f.id" class="flex items-center gap-3">
                <span class="w-44 shrink-0 text-xs font-medium">{{ f.factor }}</span>
                <div class="h-2 flex-1 overflow-hidden rounded-full bg-muted">
                  <div :class="['h-full rounded-full', f.score >= 90 ? 'bg-success' : f.score >= 75 ? 'bg-warning' : 'bg-danger']" :style="{ width: `${f.score}%` }" />
                </div>
                <span class="w-10 shrink-0 text-right font-display text-xs font-bold">{{ f.score }}</span>
                <span class="w-10 shrink-0 text-right text-[10px] text-muted-foreground">{{ Math.round(f.weight * 100) }}%</span>
              </div>
            </div>
          </SectionCard>
          <SectionCard
            title="Pengali Trust (Transparan)"
            sub="Tabel pengali yang memengaruhi skor kredensial akhir — ditampilkan terbuka, berbeda dari produk lama yang menyembunyikannya"
          >
            <div class="space-y-1.5">
              <div
                v-for="b in data.bands" :key="b.minTrust"
                :class="['flex items-center justify-between rounded-lg px-3.5 py-2.5 text-xs', t.score >= b.minTrust ? 'bg-brand-2-light font-semibold text-brand-2-dark' : 'bg-muted/40 text-muted-foreground']"
              >
                <span>Trust ≥ {{ b.minTrust }}</span>
                <span>{{ b.label }}</span>
              </div>
            </div>
          </SectionCard>
        </div>
      </div>
    </template>
  </div>
</template>
