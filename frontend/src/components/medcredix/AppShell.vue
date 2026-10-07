<script setup lang="ts">
// AppShell — port dari components/medcredix/app-shell.tsx:
// sidebar SPKD navy, topbar, breadcrumb, pencarian global (Ctrl+K), notifikasi.
import {
  Building, Sparkles, Menu, Search, LogOut, ChevronRight, Moon, SunMedium, Bell,
} from 'lucide-vue-next'
import { NAV, VIEW_TITLES } from '~/utils/nav'
import { useMedCredix, useHashRouter } from '~/stores/medcredix'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import Badge from '~/components/ui/Badge.vue'
import Button from '~/components/ui/Button.vue'
import Dialog from '~/components/ui/Dialog.vue'
import Dropdown from '~/components/ui/Dropdown.vue'
import DropdownItem from '~/components/ui/DropdownItem.vue'
import DropdownLabel from '~/components/ui/DropdownLabel.vue'
import DropdownSeparator from '~/components/ui/DropdownSeparator.vue'
import Sheet from '~/components/ui/Sheet.vue'
import SideNavContent from '~/components/medcredix/SideNavContent.vue'

useHashRouter()
const store = useMedCredix()

const mobileOpen = ref(false)
const searchOpen = ref(false)
const term = ref('')
const boot = ref<any>(null)

async function loadBoot() {
  try {
    boot.value = await api<any>(`/bootstrap${store.hospitalId ? `?hospitalId=${store.hospitalId}` : ''}`)
  } catch {
    /* biarkan data lama */
  }
}
onMounted(() => {
  loadBoot()
  setInterval(loadBoot, 60_000)
})
watch(() => store.hospitalId, loadBoot)

// pencarian global (debounce)
const searchResults = ref<any[]>([])
const searching = ref(false)
let searchTimer: ReturnType<typeof setTimeout> | null = null
watch(term, (v) => {
  if (searchTimer) clearTimeout(searchTimer)
  searchResults.value = []
  if (v.trim().length < 2) return
  searchTimer = setTimeout(async () => {
    searching.value = true
    try {
      const r = await api<{ results: any[] }>(`/search?q=${encodeURIComponent(v)}`)
      searchResults.value = r.results
    } catch {
      searchResults.value = []
    }
    searching.value = false
  }, 250)
})

const visibleNav = computed(() =>
  NAV.map((g) => ({ ...g, items: g.items.filter((i) => !i.roles || (store.user && i.roles.includes(store.user.role))) }))
    .filter((g) => g.items.length),
)
const crumbs = computed(() => {
  const view = store.view
  const group = NAV.find((g) => g.items.some((i) => i.id === view))
  return { group: group?.label ?? '', item: VIEW_TITLES[view]?.title ?? view }
})
function badge(key?: string): number {
  if (!key || !boot.value) return 0
  return boot.value.badges?.[key] ?? 0
}
function onKey(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    searchOpen.value = true
  }
  if (e.key === 'Escape') searchOpen.value = false
}
onMounted(() => window.addEventListener('keydown', onKey))
onBeforeUnmount(() => window.removeEventListener('keydown', onKey))
const userInitials = computed(() => {
  const name = store.user?.name ?? 'U'
  return (
    name.split(' ').map((w) => w[0]).slice(0, 2).join('').replace(/[^A-Za-z]/g, '').slice(0, 2) || 'U'
  )
})
function openNotificationLink(link: string) {
  const view = link.replace('#/', '').split('/')[0]
  store.navigate(view)
}
</script>

<template>
  <div class="flex min-h-screen bg-background">
    <!-- Sidebar desktop -->
    <aside class="fixed inset-y-0 left-0 z-40 hidden w-60 bg-sidebar lg:block">
      <SideNavContent :nav="visibleNav" :view="store.view" :user="store.user" :badges="boot?.badges" @navigate="store.navigate" />
    </aside>

    <!-- Sidebar mobile -->
    <Sheet :open="mobileOpen" side="left" @update:open="mobileOpen = $event">
      <SideNavContent :nav="visibleNav" :view="store.view" :user="store.user" :badges="boot?.badges" @navigate="(v) => { store.navigate(v); mobileOpen = false }" />
    </Sheet>

    <div class="flex min-w-0 flex-1 flex-col lg:pl-60">
      <!-- Topbar -->
      <header class="sticky top-0 z-30 flex h-14 items-center gap-2 border-b border-border bg-background/85 px-3 backdrop-blur-md sm:px-5">
        <Button variant="ghost" size="icon" class="lg:hidden" aria-label="Buka menu" @click="mobileOpen = true">
          <Menu class="h-5 w-5" />
        </Button>
        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-1.5 text-[11px] text-muted-foreground">
            <Building class="h-3 w-3" />
            <span class="max-w-40 truncate font-medium sm:max-w-none">{{ boot?.hospital?.name ?? '…' }}</span>
            <span class="hidden rounded bg-brand-light px-1.5 py-px text-[10px] font-bold text-brand sm:inline">Kelas {{ boot?.hospital?.class }}</span>
            <span class="hidden items-center gap-1 rounded bg-brand-2-light px-1.5 py-px text-[10px] font-bold text-brand-2-dark md:inline-flex">
              <Sparkles class="h-3 w-3" /> Stage {{ boot?.intel?.bcds?.stage ?? '—' }} {{ boot?.intel?.bcds?.lulus ? '· LULUS' : '· BELUM LULUS' }}
            </span>
          </div>
        </div>

        <!-- Hospital switcher (super admin) -->
        <Dropdown v-if="boot?.hospitals?.length > 1">
          <template #trigger>
            <Button variant="outline" size="sm" class="hidden h-8 gap-1.5 text-xs sm:inline-flex">
              <Building class="h-3.5 w-3.5" /> Ganti RS
            </Button>
          </template>
          <DropdownItem v-for="h in boot.hospitals" :key="h.id" :class="store.hospitalId === h.id ? 'bg-brand-light font-semibold' : ''" @click="store.setHospital(h.id)">
            {{ h.name }} <span class="text-[10px] text-muted-foreground">{{ h.city }}</span>
          </DropdownItem>
        </Dropdown>

        <Button variant="outline" size="sm" class="h-8 gap-2 text-xs" @click="searchOpen = true">
          <Search class="h-3.5 w-3.5" />
          <span class="hidden sm:inline">Cari…</span>
          <kbd class="hidden rounded bg-muted px-1.5 py-0.5 text-[10px] font-semibold text-muted-foreground md:inline">Ctrl K</kbd>
        </Button>

        <Button variant="ghost" size="icon" class="h-8 w-8" aria-label="Ganti tema" @click="store.toggleTheme()">
          <Moon v-if="store.theme === 'light'" class="h-4 w-4" /><SunMedium v-else class="h-4 w-4" />
        </Button>

        <!-- Notifikasi -->
        <Dropdown>
          <template #trigger>
            <Button variant="ghost" size="icon" class="relative h-8 w-8" aria-label="Notifikasi">
              <Bell class="h-4 w-4" />
              <span v-if="(boot?.notifications?.length ?? 0) > 0" class="absolute right-1 top-1 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-danger text-[9px] font-bold text-white">
                {{ boot.notifications.length }}
              </span>
            </Button>
          </template>
          <DropdownLabel class="text-xs">Notifikasi ({{ boot?.notifications?.length ?? 0 }})</DropdownLabel>
          <DropdownSeparator />
          <div class="max-h-80 w-80 overflow-y-auto scrollbar-thin">
            <button
              v-for="n in boot?.notifications ?? []" :key="n.id"
              class="flex w-full flex-col gap-1 border-b border-border/60 px-3 py-2.5 text-left transition-colors last:border-0 hover:bg-muted/50"
              @click="n.link && openNotificationLink(n.link)"
            >
              <span class="flex items-center gap-2">
                <span :class="['h-1.5 w-1.5 rounded-full', n.severity === 'Tinggi' ? 'bg-danger' : n.severity === 'Sedang' ? 'bg-warning' : 'bg-info']" />
                <span class="text-xs font-semibold leading-snug">{{ n.title }}</span>
              </span>
              <span class="line-clamp-2 pl-3.5 text-[11px] leading-snug text-muted-foreground">{{ n.body }}</span>
            </button>
            <div v-if="!(boot?.notifications?.length)" class="px-3 py-4 text-center text-xs text-muted-foreground">Tidak ada notifikasi.</div>
          </div>
        </Dropdown>

        <!-- User menu -->
        <Dropdown>
          <template #trigger>
            <Button variant="ghost" class="h-8 gap-2 px-2">
              <span class="flex h-7 w-7 items-center justify-center rounded-full bg-gradient-brand text-[11px] font-bold text-white">{{ userInitials }}</span>
              <span class="hidden max-w-28 truncate text-xs font-semibold sm:block">{{ store.user?.name?.split(',')[0] }}</span>
            </Button>
          </template>
          <DropdownLabel>
            <div class="text-xs font-semibold">{{ store.user?.name }}</div>
            <div class="text-[11px] font-normal text-muted-foreground">{{ store.user?.email }}</div>
            <Badge class="mt-1.5 text-[10px] font-semibold">{{ store.user?.roleLabel }}</Badge>
          </DropdownLabel>
          <DropdownSeparator />
          <DropdownItem danger @click="store.logout()"><LogOut class="h-4 w-4" /> Keluar</DropdownItem>
        </Dropdown>
      </header>

      <!-- Breadcrumb -->
      <div class="flex items-center gap-1.5 border-b border-border/60 bg-muted/30 px-4 py-2 text-[11px] text-muted-foreground sm:px-6">
        <button class="hover:text-brand-2" @click="store.navigate('dashboard')">MedCredix</button>
        <ChevronRight class="h-3 w-3" />
        <span>{{ crumbs.group }}</span>
        <ChevronRight class="h-3 w-3" />
        <span class="font-semibold text-foreground">{{ crumbs.item }}</span>
        <span class="ml-auto hidden text-[10px] md:block">
          Skor Kredensial BPJS:
          <button class="font-semibold text-brand-2 hover:underline" @click="store.openWhy('bcds')">{{ boot?.intel?.bcds?.final ?? '—' }}</button>
          · Trust {{ boot?.intel?.trust?.score ?? '—' }}
        </span>
      </div>

      <main id="konten-utama" class="min-w-0 flex-1 px-4 py-6 sm:px-6"><slot /></main>

      <footer class="mt-auto border-t border-border/60 bg-muted/30 px-6 py-3">
        <div class="flex flex-wrap items-center justify-between gap-2 text-[10px] text-muted-foreground">
          <span>© 2026 PT Sistem Pelayanan Kesehatan dan Data (SPKD) — MedCredix · From Hospital Data to Action.</span>
          <span>Dirancang untuk mendukung regulasi kesehatan Indonesia yang berlaku, tunduk pada verifikasi regulasi terkini dan implementasi spesifik rumah sakit.</span>
        </div>
      </footer>
    </div>

    <!-- Global Search Palette -->
    <Dialog :open="searchOpen" @update:open="searchOpen = $event">
      <div class="flex items-center gap-2 border-b px-3 pb-3">
        <Search class="h-4 w-4 shrink-0 text-muted-foreground" />
        <input
          v-model="term" type="text" placeholder="Cari praktisi, kredensial, klaim, risiko, aksi, regulasi, bukti…"
          class="h-8 w-full bg-transparent text-sm outline-none placeholder:text-muted-foreground"
        >
      </div>
      <div class="mt-3 max-h-80 overflow-y-auto scrollbar-thin">
        <div v-if="term.trim().length < 2" class="px-4 py-8 text-center text-xs text-muted-foreground">Ketik minimal 2 karakter untuk mencari lintas modul.</div>
        <div v-else-if="searching" class="px-4 py-8 text-center text-xs text-muted-foreground">Mencari…</div>
        <div v-else-if="searchResults.length === 0" class="px-4 py-8 text-center text-xs text-muted-foreground">Tidak ada hasil untuk "{{ term }}".</div>
        <div v-else class="space-y-1">
          <button
            v-for="r in searchResults" :key="`${r.type}-${r.id}`"
            class="flex w-full items-center gap-3 rounded-lg px-2.5 py-2 text-left transition-colors hover:bg-brand-light"
            @click="store.navigate(r.view); searchOpen = false; term = ''"
          >
            <span class="rounded-md bg-brand-light p-1.5 text-brand font-display text-[11px] font-bold">{{ r.type.slice(0, 2).toUpperCase() }}</span>
            <span class="min-w-0 flex-1">
              <span class="block truncate text-xs font-semibold">{{ r.label }}</span>
              <span class="block truncate text-[11px] text-muted-foreground">{{ r.sub }}</span>
            </span>
            <Badge variant="outline" class="shrink-0 text-[9px]">{{ r.type }}</Badge>
          </button>
        </div>
      </div>
    </Dialog>
  </div>
</template>
