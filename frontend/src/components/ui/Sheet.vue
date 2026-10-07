<script setup lang="ts">
// Sheet — panel samping (padanan shadcn/ui Sheet). Slot: default (konten).
import { onBeforeUnmount, watch } from 'vue'
const props = defineProps<{ open: boolean; side?: 'right' | 'left' }>()
const emit = defineEmits<{ 'update:open': [value: boolean] }>()
const side = props.side ?? 'right'

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
      <transition :name="side === 'left' ? 'slide-l' : 'slide-r'">
        <div
          role="dialog"
          aria-modal="true"
          :class="[
            'absolute top-0 z-10 flex h-full w-full flex-col bg-card shadow-xl outline-none',
            side === 'right' ? 'right-0 sm:max-w-2xl' : 'left-0 w-64 border-r border-sidebar-border bg-sidebar',
          ]"
        >
          <slot />
        </div>
      </transition>
    </div>
  </Teleport>
</template>

<style scoped>
.slide-r-enter-active,
.slide-r-leave-active,
.slide-l-enter-active,
.slide-l-leave-active {
  transition: transform 0.35s cubic-bezier(0.32, 0.72, 0, 1);
}
.slide-r-enter-from,
.slide-r-leave-to {
  transform: translateX(100%);
}
.slide-l-enter-from,
.slide-l-leave-to {
  transform: translateX(-100%);
}
</style>
