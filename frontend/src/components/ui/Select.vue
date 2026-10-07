<script setup lang="ts">
// Select — listbox ringan (padanan shadcn/ui Select) dengan v-model.
// Opsi: <SelectOption value="...">Label</SelectOption>
import { onClickOutside } from '@vueuse/core'
import { ChevronDown } from 'lucide-vue-next'
import { computed, inject, provide, ref } from 'vue'

const model = defineModel<string>({ default: '' })
defineProps<{ placeholder?: string }>()

const open = ref(false)
const root = ref<HTMLElement | null>(null)
onClickOutside(root, () => (open.value = false))

const options = inject<{ items: { value: string; label: string }[] }>('select-options', { items: [] })
provide('select-register', (value: string, label: string) => {
  if (!options.items.find((i) => i.value === value)) options.items.push({ value, label })
})
const selectedLabel = computed(() => options.items.find((i) => i.value === model.value)?.label ?? (model.value || ''))
function pick(value: string) {
  model.value = value
  open.value = false
}
</script>

<template>
  <div ref="root" class="relative">
    <button
      type="button"
      :class="[
        'flex h-9 w-full items-center justify-between gap-2 rounded-lg border border-input bg-card px-3 py-2 text-sm shadow-xs',
        'focus:outline-none focus:ring-2 focus:ring-ring',
      ]"
      @click="open = !open"
    >
      <span :class="selectedLabel ? '' : 'text-muted-foreground'">{{ selectedLabel || placeholder || 'Pilih…' }}</span>
      <ChevronDown class="h-4 w-4 shrink-0 text-muted-foreground" :class="open ? 'rotate-180' : ''" />
    </button>
    <div
      v-if="open"
      class="absolute z-50 mt-1 max-h-64 w-full overflow-y-auto rounded-lg border border-border bg-popover p-1 text-popover-foreground shadow-md scrollbar-thin"
    >
      <button
        v-for="opt in options.items"
        :key="opt.value"
        type="button"
        :class="[
          'relative flex w-full cursor-pointer select-none items-center rounded-md px-2 py-1.5 text-left text-xs outline-none transition-colors',
          opt.value === model ? 'bg-brand-light font-semibold text-brand' : 'hover:bg-muted/70',
        ]"
        @click="pick(opt.value)"
      >
        {{ opt.label }}
      </button>
    </div>
  </div>
</template>
