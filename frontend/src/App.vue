<script setup lang="ts">
// SPKD MedCredix — Medical Credentialing Intelligence (Vue 3 + Vite)
// "From Hospital Data to Action."
// Port dari src/app/page.tsx: boot sesi + hash-router + 29 view + drawer global.
import { useMedCredix, useHashRouter, type SessionUser } from '~/stores/medcredix'
import { computed, onMounted, ref } from 'vue'
import type { Component } from 'vue'
import { api } from '~/composables/useApi'
import Toaster from '~/components/ui/Toaster.vue'
import ActionDialog from '~/components/medcredix/ActionDialog.vue'
import AppShell from '~/components/medcredix/AppShell.vue'
import EvidenceDrawer from '~/components/medcredix/EvidenceDrawer.vue'
import LoginScreen from '~/components/medcredix/LoginScreen.vue'
import WhyDrawer from '~/components/medcredix/WhyDrawer.vue'
import ActionsView from '~/components/views/ActionsView.vue'
import AdminView from '~/components/views/AdminView.vue'
import AuditView from '~/components/views/AuditView.vue'
import ClaimAnomaliesView from '~/components/views/ClaimAnomaliesView.vue'
import ClaimIntegrityView from '~/components/views/ClaimIntegrityView.vue'
import ComplianceView from '~/components/views/ComplianceView.vue'
import CredentialsView from '~/components/views/CredentialsView.vue'
import DashboardView from '~/components/views/DashboardView.vue'
import DataQualityView from '~/components/views/DataQualityView.vue'
import DataSourcesView from '~/components/views/DataSourcesView.vue'
import DataTrustView from '~/components/views/DataTrustView.vue'
import DigitalMaturityView from '~/components/views/DigitalMaturityView.vue'
import DocumentsView from '~/components/views/DocumentsView.vue'
import EarlyWarningView from '~/components/views/EarlyWarningView.vue'
import EvidenceView from '~/components/views/EvidenceView.vue'
import FacilitiesView from '~/components/views/FacilitiesView.vue'
import InteropView from '~/components/views/InteropView.vue'
import PathwayView from '~/components/views/PathwayView.vue'
import PractitionersView from '~/components/views/PractitionersView.vue'
import PrerequisitesView from '~/components/views/PrerequisitesView.vue'
import ReportsView from '~/components/views/ReportsView.vue'
import RiskRadarView from '~/components/views/RiskRadarView.vue'
import SafetyView from '~/components/views/SafetyView.vue'
import SettingsAlertsView from '~/components/views/SettingsAlertsView.vue'
import SettingsHospitalView from '~/components/views/SettingsHospitalView.vue'
import SettingsUsersView from '~/components/views/SettingsUsersView.vue'
import UtilizationView from '~/components/views/UtilizationView.vue'
import VerificationView from '~/components/views/VerificationView.vue'
import WorkforceView from '~/components/views/WorkforceView.vue'

const store = useMedCredix()
useHashRouter()
const booted = ref(false)

const VIEWS: Record<string, Component> = {
  dashboard: DashboardView,
  'risk-radar': RiskRadarView,
  'early-warnings': EarlyWarningView,
  actions: ActionsView,
  'cred-practitioners': PractitionersView,
  'cred-credentials': CredentialsView,
  'cred-documents': DocumentsView,
  'cred-verification': VerificationView,
  'cred-prerequisites': PrerequisitesView,
  'ready-workforce': WorkforceView,
  'ready-facilities': FacilitiesView,
  'ready-digital': DigitalMaturityView,
  'ready-interop': InteropView,
  'clin-pathway': PathwayView,
  'clin-utilization': UtilizationView,
  'clin-safety': SafetyView,
  'claim-integrity': ClaimIntegrityView,
  'claim-anomalies': ClaimAnomaliesView,
  'data-quality': DataQualityView,
  'data-trust': DataTrustView,
  'data-sources': DataSourcesView,
  evidence: EvidenceView,
  compliance: ComplianceView,
  audit: AuditView,
  reports: ReportsView,
  'settings-hospital': SettingsHospitalView,
  'settings-users': SettingsUsersView,
  'settings-alerts': SettingsAlertsView,
  admin: AdminView,
}

onMounted(async () => {
  // pulihkan tema (sinkronisasi dengan DOM)
  const theme = window.localStorage.getItem('medcredix-theme')
  if (theme === 'dark') {
    document.documentElement.classList.add('dark')
    store.theme = 'dark'
  }
  // validasi sesi tersimpan terhadap data demo
  const saved = window.localStorage.getItem('medcredix-session')
  if (saved) {
    try {
      const session = JSON.parse(saved)
      // validasi sesi tersimpan ke backend Golang (atau data demo lokal bila VITE_DEMO_MODE=true)
      await api('/bootstrap')
      store.restoreSession(session)
    } catch (e: any) {
      // hanya sesi yang benar-benar tidak dikenal yang dikeluarkan; error lain dibiarkan
      if (String(e?.message ?? '').startsWith('Sesi tidak valid')) {
        window.localStorage.removeItem('medcredix-session')
        store.user = null
      } else {
        try { store.restoreSession(JSON.parse(saved)) } catch { /* abaikan */ }
      }
    }
  }
  booted.value = true
})

const activeView = computed(() => VIEWS[store.view])
</script>

<template>
  <div>
    <!-- Memuat -->
    <div v-if="!booted" class="flex min-h-screen items-center justify-center bg-background">
      <div class="flex flex-col items-center gap-3">
        <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-brand">
          <svg viewBox="0 0 32 32" class="h-7 w-7 animate-pulse"><path d="M13 7h6v6h6v6h-6v6h-6v-6H7v-6h6z" fill="#fff" /></svg>
        </div>
        <div class="text-xs font-medium text-muted-foreground">Memuat MedCredix…</div>
      </div>
    </div>

    <!-- Login -->
    <LoginScreen v-else-if="!store.user" />

    <!-- Aplikasi -->
    <template v-else>
      <AppShell>
        <component :is="activeView" v-if="activeView" />
        <div v-else class="flex flex-col items-center justify-center gap-3 py-24 text-center">
          <div class="font-display text-lg font-bold">Halaman tidak ditemukan</div>
          <p class="max-w-sm text-sm text-muted-foreground">
            Tampilan "{{ store.view }}" tidak tersedia. Gunakan menu di sisi kiri atau tekan Ctrl+K untuk mencari.
          </p>
        </div>
        <WhyDrawer />
        <ActionDialog />
        <EvidenceDrawer />
      </AppShell>
    </template>
    <Toaster />
  </div>
</template>
