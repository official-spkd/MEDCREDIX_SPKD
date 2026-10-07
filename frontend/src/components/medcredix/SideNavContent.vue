<script setup lang="ts">
// SideNavContent — isi navigasi sidebar (dipakai desktop & mobile).
import {
  LayoutDashboard, Radar, AlertTriangle, CheckCircle2, Users, IdCard, FileText, BadgeCheck, ShieldCheck,
  Stethoscope, Building2, MonitorSmartphone, Network, GitBranch, BedDouble, HeartPulse, FileCheck2,
  ScanSearch, Database, ShieldHalf, Cable, Fingerprint, Scale, History, FileBarChart, Settings, UserCog,
  Bell, Globe2,
} from 'lucide-vue-next'
import type { NavGroup } from '~/utils/nav'
import type { SessionUser } from '~/stores/medcredix'

const props = defineProps<{ nav: NavGroup[]; view: string; user: SessionUser | null; badges?: any }>()
const emit = defineEmits<{ navigate: [view: string] }>()

const ICONS: Record<string, any> = {
  LayoutDashboard, Radar, AlertTriangle, CheckCircle2, Users, IdCard, FileText, BadgeCheck, ShieldCheck,
  Stethoscope, Building2, MonitorSmartphone, Network, GitBranch, BedDouble, HeartPulse, FileCheck2,
  ScanSearch, Database, ShieldHalf, Cable, Fingerprint, Scale, History, FileBarChart, Settings, UserCog,
  Bell, Globe2,
}
function badge(key?: string): number {
  if (!key || !props.badges) return 0
  return props.badges[key] ?? 0
}
</script>

<template>
  <nav class="flex h-full flex-col" aria-label="Navigasi utama">
    <div class="flex items-center gap-2.5 px-4 pb-4 pt-5">
      <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-white/95">
        <svg viewBox="0 0 32 32" class="h-6 w-6"><rect width="32" height="32" rx="7" fill="#14285C" /><path d="M13 7h6v6h6v6h-6v6h-6v-6H7v-6h6z" fill="#2198A5" /></svg>
      </div>
      <div class="min-w-0">
        <div class="font-display text-[15px] font-bold leading-none text-white">MedCredix</div>
        <div class="mt-0.5 truncate text-[9px] font-semibold uppercase tracking-[0.14em] text-white/50">by SPKD</div>
      </div>
    </div>
    <div class="flex-1 space-y-4 overflow-y-auto px-2.5 pb-4 scrollbar-thin">
      <div v-for="g in nav" :key="g.id">
        <div class="px-2.5 pb-1.5 text-[10px] font-semibold uppercase tracking-wider text-white/35">{{ g.label }}</div>
        <div class="space-y-0.5">
          <button
            v-for="item in g.items" :key="item.id"
            :aria-current="view === item.id ? 'page' : undefined"
            :class="[
              'group flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-left text-[13px] font-medium transition-all',
              view === item.id ? 'bg-sidebar-accent text-white shadow-sm' : 'text-sidebar-foreground/85 hover:bg-white/5 hover:text-white',
            ]"
            @click="emit('navigate', item.id)"
          >
            <component :is="ICONS[item.icon] ?? LayoutDashboard" :class="['h-4 w-4 shrink-0', view === item.id ? 'text-brand-2' : 'text-sidebar-foreground/60 group-hover:text-white/80']" />
            <span class="flex-1 truncate">{{ item.label }}</span>
            <span v-if="badge(item.badgeKey) > 0" class="rounded-full bg-brand-2 px-1.5 py-0.5 text-[10px] font-bold text-white">
              {{ badge(item.badgeKey) }}
            </span>
          </button>
        </div>
      </div>
    </div>
    <div class="border-t border-sidebar-border px-4 py-3">
      <div class="truncate text-[11px] font-medium text-white/70">{{ user?.name }}</div>
      <div class="truncate text-[10px] text-white/40">{{ user?.roleLabel }}</div>
    </div>
  </nav>
</template>
