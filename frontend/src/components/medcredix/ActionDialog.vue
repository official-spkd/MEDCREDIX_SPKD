<script setup lang="ts">
// ActionDialog — buat tindakan dari mana pun (skor, risiko, peringatan, klaim).
// Port dari components/medcredix/action-dialog.tsx.
import { Loader2, PlusCircle } from 'lucide-vue-next'
import { useMedCredix } from '~/stores/medcredix'
import { toast } from '~/composables/useToast'
import { ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import Button from '~/components/ui/Button.vue'
import Dialog from '~/components/ui/Dialog.vue'
import Input from '~/components/ui/Input.vue'
import Label from '~/components/ui/Label.vue'
import Select from '~/components/ui/Select.vue'
import SelectOption from '~/components/ui/SelectOption.vue'
import Textarea from '~/components/ui/Textarea.vue'

const store = useMedCredix()
const priority = ref('Sedang')
const owner = ref('')
const dueDate = ref('')
const title = ref('')
const desc = ref('')
const expected = ref('')
const saving = ref(false)
const users = ref<{ id: string; name: string; role: string }[]>([])

watch(
  () => store.actionDraft,
  async (draft) => {
    if (!draft) return
    title.value = draft.title
    desc.value = draft.description ?? ''
    expected.value = draft.expectedOutcome ?? ''
    priority.value = draft.priority ?? 'Sedang'
    owner.value = draft.suggestedOwner ?? ''
    dueDate.value = ''
    try {
      const r = await api<{ users: { id: string; name: string; role: string }[] }>('/settings')
      users.value = r.users ?? []
    } catch { /* abaikan */ }
  },
  { immediate: true },
)

async function submit() {
  if (!title.value.trim()) {
    toast({ title: 'Judul tindakan wajib diisi', variant: 'destructive' })
    return
  }
  saving.value = true
  try {
    await api('/actions', {
      method: 'POST',
      body: {
        title: title.value,
        description: desc.value,
        priority: priority.value,
        sourceType: store.actionDraft?.sourceType ?? 'Manual',
        sourceRef: store.actionDraft?.sourceRef ?? null,
        expectedOutcome: expected.value || null,
        owner: owner.value || 'Belum ditugaskan',
        dueDate: dueDate.value || null,
      },
    })
    toast({ title: 'Tindakan dibuat', description: `${title.value} — lacak progresnya di Action Center.` })
    store.closeAction()
  } catch (e: any) {
    toast({ title: 'Gagal membuat tindakan', description: e?.message, variant: 'destructive' })
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <Dialog :open="!!store.actionDraft" @update:open="(o) => !o && store.closeAction()">
    <div class="flex items-center gap-2">
      <PlusCircle class="h-4 w-4 text-brand-2" />
      <h2 class="font-display text-base font-semibold">Buat Tindakan Perbaikan</h2>
    </div>
    <p class="mt-1 text-xs leading-relaxed text-muted-foreground">
      Dari {{ store.actionDraft?.sourceType ?? 'Manual' }}{{ store.actionDraft?.sourceRef ? ` (${store.actionDraft.sourceRef})` : '' }} —
      memasuki alur closed-loop: aksi selesai → verifikasi → skor dihitung ulang.
    </p>
    <div class="mt-4 grid gap-3.5">
      <div class="grid gap-1.5">
        <Label for="act-title" class="text-xs font-semibold">Judul</Label>
        <Input id="act-title" v-model="title" placeholder="Mis. Tinjau 26 klaim dengan sinyal inkonsistensi" />
      </div>
      <div class="grid gap-1.5">
        <Label for="act-desc" class="text-xs font-semibold">Deskripsi</Label>
        <Textarea id="act-desc" v-model="desc" :rows="3" placeholder="Konteks, langkah yang diperlukan, kriteria selesai…" />
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div class="grid gap-1.5">
          <Label class="text-xs font-semibold">Prioritas</Label>
          <Select v-model="priority">
            <SelectOption value="Tinggi">Tinggi</SelectOption>
            <SelectOption value="Sedang">Sedang</SelectOption>
            <SelectOption value="Rendah">Rendah</SelectOption>
          </Select>
        </div>
        <div class="grid gap-1.5">
          <Label for="act-due" class="text-xs font-semibold">Tenggat</Label>
          <Input id="act-due" v-model="dueDate" type="date" />
        </div>
      </div>
      <div class="grid gap-1.5">
        <Label class="text-xs font-semibold">Penanggung Jawab</Label>
        <Select v-model="owner" placeholder="Pilih penanggung jawab">
          <SelectOption v-for="u in users" :key="u.id" :value="u.name">{{ u.name }}</SelectOption>
        </Select>
      </div>
      <div class="grid gap-1.5">
        <Label for="act-expected" class="text-xs font-semibold">Hasil yang Diharapkan</Label>
        <Input id="act-expected" v-model="expected" placeholder="Mis. CCS naik dari 79 ke ≥85" />
      </div>
    </div>
    <div class="mt-5 flex justify-end gap-2">
      <Button variant="outline" @click="store.closeAction()">Batal</Button>
      <Button :disabled="saving" class="bg-brand hover:bg-brand-hover" @click="submit">
        <Loader2 v-if="saving" class="h-4 w-4 animate-spin" /> Buat Tindakan
      </Button>
    </div>
  </Dialog>
</template>
