<script setup lang="ts">
// StatusBadge — badge status dengan ikon + label (no color-only — WCAG 2.2 AA).
import { computed } from 'vue'
import StatusIcon from '~/components/medcredix/StatusIcon.vue'
const props = defineProps<{ status: string; tone?: string }>()
const tone = computed(() => props.tone ?? statusToneOf(props.status))

function statusToneOf(status: string): string {
  const success = ['Aktif', 'Valid', 'Disetujui', 'Selesai', 'Ditutup', 'Tervalidasi', 'Berhasil', 'Terhubung', 'Patuh', 'Baik', 'Teratasi', 'Diperbaiki', 'Diklarifikasi', 'Terbuka', 'Diverifikasi', 'LULUS']
  const warning = ['Segera Berakhir', 'Menunggu Verifikasi', 'Ditinjau', 'Diketahui', 'Dalam Pengerjaan', 'Perlu Perhatian', 'Ditugaskan', 'Sebagian', 'Terhubung Sebagian', 'Diajukan', 'Ditunda', 'Sedang Ditinjau', 'Verifikasi', 'Tinjauan Komite', 'Tinjauan Dokumen', 'GOOD']
  const danger = ['Kedaluwarsa', 'Ditolak', 'Gagal', 'Kritis', 'Tidak Tervalidasi', 'Ditangguhkan', 'Dicabut', 'Tidak Ada', 'Terputus', 'TIDAK LULUS', 'POOR', 'Under Review']
  if (success.includes(status)) return 'success'
  if (warning.includes(status)) return 'warning'
  if (danger.includes(status)) return 'danger'
  return 'info'
}
const toneCls: Record<string, string> = {
  success: 'bg-success-light text-success',
  warning: 'bg-warning-light text-warning',
  danger: 'bg-danger-light text-danger',
  info: 'bg-info-light text-info',
  muted: 'bg-muted text-muted-foreground',
}
</script>

<template>
  <span :class="['inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-semibold', toneCls[tone] ?? toneCls.info]">
    <StatusIcon :tone="tone" class="h-3.5 w-3.5" />
    {{ status }}
  </span>
</template>
