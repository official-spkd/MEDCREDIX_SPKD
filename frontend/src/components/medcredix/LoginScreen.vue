<script setup lang="ts">
// LoginScreen — port dari components/medcredix/login.tsx (SPKD brand).
import { Activity, Eye, EyeOff, Loader2, Lock, LogIn, Mail, ShieldCheck, TrendingUp, Radar } from 'lucide-vue-next'
import { useMedCredix, type SessionUser } from '~/stores/medcredix'
import { ref } from 'vue'
import { api } from '~/composables/useApi'
import Button from '~/components/ui/Button.vue'
import Input from '~/components/ui/Input.vue'
import Label from '~/components/ui/Label.vue'

const store = useMedCredix()
const email = ref('admin@nusantarasehat.co.id')
const password = ref('medcredix2026')
const showPw = ref(false)
const loading = ref(false)
const error = ref<string | null>(null)

const DEMO_ACCOUNTS = [
  { email: 'admin@nusantarasehat.co.id', role: 'Admin RS — dr. Aditya Nugraha' },
  { email: 'manajemen@nusantarasehat.co.id', role: 'Manajemen RS (eksekutif)' },
  { email: 'kredensial@nusantarasehat.co.id', role: 'Petugas Kredensial' },
  { email: 'komite@nusantarasehat.co.id', role: 'Komite Medis' },
  { email: 'klaim@nusantarasehat.co.id', role: 'Tim Klaim' },
  { email: 'it@nusantarasehat.co.id', role: 'IT / SIMRS' },
  { email: 'superadmin@spkd.co.id', role: 'SPKD Super Admin' },
  { email: 'auditorbpjs@spkd.co.id', role: 'Auditor BPJS' },
]
const PW = 'medcredix2026'

async function submit() {
  loading.value = true
  error.value = null
  try {
    const user = await api<SessionUser>('/auth', { method: 'POST', body: { email: email.value, password: password.value } })
    store.login(user)
  } catch (e: any) {
    error.value = e?.message ?? 'Gagal masuk.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="grid min-h-screen lg:grid-cols-[1.1fr_1fr]">
    <!-- Panel kiri: brand -->
    <div class="relative hidden flex-col justify-between overflow-hidden bg-gradient-brand p-10 text-white lg:flex">
      <div class="absolute -right-24 -top-24 h-96 w-96 rounded-full bg-white/5" />
      <div class="absolute -bottom-32 -left-16 h-96 w-96 rounded-full bg-white/5" />
      <div class="relative">
        <div class="flex items-center gap-3">
          <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-white/95">
            <svg viewBox="0 0 32 32" class="h-7 w-7"><rect width="32" height="32" rx="7" fill="#14285C" /><path d="M13 7h6v6h6v6h-6v6h-6v-6H7v-6h6z" fill="#2198A5" /></svg>
          </div>
          <div>
            <div class="font-display text-lg font-bold leading-none tracking-tight">MedCredix</div>
            <div class="text-[11px] font-medium uppercase tracking-[0.18em] text-white/70">Medical Credentialing Intelligence</div>
          </div>
        </div>
        <h1 class="mt-14 max-w-md font-display text-4xl font-bold leading-tight">
          From Hospital Data<br />to Action.
        </h1>
        <p class="mt-4 max-w-md text-sm leading-relaxed text-white/80">
          Platform kredensial & inteligensi rumah sakit dari <strong>PT Sistem Pelayanan Kesehatan dan Data (SPKD)</strong>.
          Pahami kondisi rumah sakit Anda, lihat risikonya, ketahui alasannya — lalu ambil tindakan.
        </p>
      </div>
      <div class="relative grid max-w-lg grid-cols-3 gap-3">
        <div v-for="f in [
          { icon: Activity, t: 'Skor', d: '9 domain kesehatan RS, dapat ditelusuri' },
          { icon: Radar, t: 'Risiko', d: 'Deteksi dini & bukti, bahasa netral' },
          { icon: TrendingUp, t: 'Aksi', d: 'Closed-loop hingga skor membaik' },
        ]" :key="f.t" class="rounded-xl bg-white/10 p-3.5 backdrop-blur-sm">
          <component :is="f.icon" class="mb-2 h-5 w-5 text-white" aria-hidden />
          <div class="font-display text-sm font-semibold">{{ f.t }}</div>
          <div class="mt-0.5 text-[11px] leading-snug text-white/70">{{ f.d }}</div>
        </div>
      </div>
      <div class="relative flex items-center gap-2 text-[11px] text-white/60">
        <ShieldCheck class="h-3.5 w-3.5" /> RBAC · Audit Trail · FHIR-ready · SATUSEHAT · Isolasi multi-tenant
      </div>
    </div>

    <!-- Panel kanan: form login -->
    <div class="flex items-center justify-center bg-background p-6">
      <div class="w-full max-w-md rounded-xl border border-border/80 bg-card shadow-sm">
        <div class="p-8">
          <div class="mb-6 flex items-center gap-3 lg:hidden">
            <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-brand">
              <svg viewBox="0 0 32 32" class="h-6 w-6"><path d="M13 7h6v6h6v6h-6v6h-6v-6H7v-6h6z" fill="#fff" /></svg>
            </div>
            <div>
              <div class="font-display text-base font-bold leading-none">MedCredix</div>
              <div class="text-[10px] font-medium uppercase tracking-widest text-muted-foreground">by SPKD</div>
            </div>
          </div>
          <h2 class="font-display text-xl font-bold">Masuk ke MedCredix</h2>
          <p class="mt-1 text-sm text-muted-foreground">Gunakan akun yang dibuatkan Super Admin platform Anda.</p>

          <form class="mt-6 space-y-4" @submit.prevent="submit">
            <div class="space-y-1.5">
              <Label for="email" class="text-xs font-semibold">Email</Label>
              <div class="relative">
                <Mail class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input id="email" v-model="email" type="email" class="pl-9" placeholder="nama@rumahsakit.co.id" required />
              </div>
            </div>
            <div class="space-y-1.5">
              <Label for="password" class="text-xs font-semibold">Kata Sandi</Label>
              <div class="relative">
                <Lock class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input id="password" v-model="password" :type="showPw ? 'text' : 'password'" class="pl-9 pr-10" required />
                <button type="button" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground" :aria-label="showPw ? 'Sembunyikan kata sandi' : 'Tampilkan kata sandi'" @click="showPw = !showPw">
                  <EyeOff v-if="showPw" class="h-4 w-4" /><Eye v-else class="h-4 w-4" />
                </button>
              </div>
            </div>
            <div v-if="error" class="rounded-lg bg-danger-light px-3 py-2 text-xs font-medium text-danger">{{ error }}</div>
            <Button type="submit" :disabled="loading" class="w-full bg-brand hover:bg-brand-hover">
              <Loader2 v-if="loading" class="h-4 w-4 animate-spin" /><LogIn v-else class="h-4 w-4" /> Masuk
            </Button>
          </form>

          <div class="mt-6 rounded-xl border bg-muted/40 p-4">
            <div class="mb-2 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground">Akun demo (kata sandi: {{ PW }})</div>
            <div class="max-h-40 space-y-1 overflow-y-auto scrollbar-thin">
              <button
                v-for="a in DEMO_ACCOUNTS" :key="a.email" type="button"
                class="flex w-full items-center justify-between gap-2 rounded-lg px-2 py-1.5 text-left text-xs transition-colors hover:bg-brand-light"
                @click="email = a.email; password = PW"
              >
                <span class="truncate font-medium text-foreground">{{ a.email }}</span>
                <span class="shrink-0 text-[10px] text-muted-foreground">{{ a.role }}</span>
              </button>
            </div>
          </div>
          <p class="mt-4 text-center text-[10px] leading-relaxed text-muted-foreground">
            Seluruh data pada lingkungan demo ini fiktif dan berlabel DEMO. Tidak ada data pasien nyata.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
