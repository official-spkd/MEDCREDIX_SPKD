<script setup lang="ts">
// Dialog — modal overlay (padanan shadcn/ui Dialog). Slot: default (konten).
import { onBeforeUnmount, watch } from 'vue'
const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ 'update:open': [value: boolean] }>()

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('update:open', false)
}
watch(
  () => props.open,
  (v) => {
    if (typeof document !== 'undefined') {
      document.removeEventListener('keydown', onKey)
      if (v) document.addEventListener('keydown', onKey)
    }
  },
  { immediate: true },
)
onBeforeUnmount(() => document?.removeEventListener('keydown', onKey))
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50">
      <div class="absolute inset-0 bg-black/60 backdrop-blur-[2px]" @click="emit('update:open', false)" />
      <div class="absolute inset-0 flex items-center justify-center p-4 overflow-y-auto">
        <div
          role="dialog"
          aria-modal="true"
          class="relative z-10 w-full max-w-lg rounded-xl border border-border bg-card p-6 shadow-lg my-8"
        >
          <slot />
        </div>
      </div>
    </div>
  </Teleport>
</template>
