<script setup lang="ts">
// SettingsUsersView — port dari views-settings.tsx.
import { Users, UserCog } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import { fmtDateTime } from '~/utils/format'
import Badge from '~/components/ui/Badge.vue'
import Select from '~/components/ui/Select.vue'
import SelectOption from '~/components/ui/SelectOption.vue'
import PageHeader from '~/components/medcredix/PageHeader.vue'
import SectionCard from '~/components/medcredix/SectionCard.vue'

const data = ref<any>(null)
const roleSel = ref('')
onMounted(async () => {
  try { data.value = await api('/settings') } catch { data.value = null }
})
const roles = computed(() => data.value?.roles ?? [])
const selectedRole = computed(() => roles.value.find((r: any) => r.key === roleSel.value) ?? roles.value[0])
watch(selectedRole, (r) => { if (r && !roleSel.value) roleSel.value = r.key })
</script>

<template>
  <div class="space-y-5">
    <PageHeader
      title="Pengguna & Peran"
      sub="RBAC — 13 izin dasar, ditegakkan di sisi server untuk setiap endpoint. Akun dibuatkan oleh Admin RS / Super Admin (tanpa registrasi mandiri)." demo
    />
    <div class="grid gap-4 lg:grid-cols-[1.3fr_1fr]">
      <SectionCard title="Pengguna" :sub="`${data?.users?.length ?? 0} akun aktif`">
        <template #icon><Users class="h-4 w-4" /></template>
        <div class="overflow-x-auto rounded-xl border">
          <table class="w-full text-xs">
            <thead class="bg-muted/60 text-left text-[10px] uppercase tracking-wide text-muted-foreground">
              <tr><th class="px-3 py-2.5">Nama</th><th class="px-3 py-2.5">Email</th><th class="px-3 py-2.5">Peran</th><th class="px-3 py-2.5">Login Terakhir</th></tr>
            </thead>
            <tbody>
              <tr v-for="u in data?.users ?? []" :key="u.id" class="border-t">
                <td class="px-3 py-2.5"><div class="font-semibold">{{ u.name }}</div><div class="text-[9px] text-muted-foreground">{{ u.title }}</div></td>
                <td class="px-3 py-2 text-[10px] text-muted-foreground">{{ u.email }}</td>
                <td class="px-3 py-2.5"><Badge class="bg-brand-light text-[10px] font-semibold text-brand hover:bg-brand-light">{{ u.roleLabel }}</Badge></td>
                <td class="px-3 py-2.5 text-[10px] text-muted-foreground">{{ fmtDateTime(u.lastLogin) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </SectionCard>

      <SectionCard title="Matriks Peran & Izin" sub="Pilih peran untuk melihat izinnya">
        <template #icon><UserCog class="h-4 w-4" /></template>
        <div v-if="roles.length">
          <Select v-model="roleSel">
            <SelectOption v-for="r in roles" :key="r.key" :value="r.key">{{ r.label }}</SelectOption>
          </Select>
          <div class="mt-3 flex flex-wrap gap-1.5">
            <Badge
              v-for="p in selectedRole?.permissions ?? []" :key="p" variant="outline"
              :class="['text-[10px] font-semibold', p === 'VIEW' ? 'text-muted-foreground' : 'border-brand-2/40 bg-brand-2-light text-brand-2-dark']"
            >{{ p }}</Badge>
          </div>
          <p class="mt-3 text-[11px] leading-relaxed text-muted-foreground">
            Izin diperiksa server pada setiap request. Peran tanpa izin CONFIGURE tetap dapat membuka halaman pengaturan, tetapi setiap upaya penyimpanan ditolak sistem dan tercatat di audit.
          </p>
        </div>
      </SectionCard>
    </div>
  </div>
</template>
