<script setup lang="ts">
// Tooltip — info hover (padanan shadcn/ui Tooltip).
import { ref } from 'vue'
defineProps<{ delay?: number }>()
const show = ref(false)
let timer: ReturnType<typeof setTimeout> | null = null
function enter() {
  timer = setTimeout(() => (show.value = true), 200)
}
function leave() {
  if (timer) clearTimeout(timer)
  show.value = false
}
</script>
<template>
  <span class="relative inline-flex" @mouseenter="enter" @mouseleave="leave" @focus="show = true" @blur="show = false">
    <slot name="trigger" />
    <transition name="fade">
      <span
        v-if="show"
        role="tooltip"
        class="absolute bottom-full left-1/2 z-50 mb-1.5 max-w-64 -translate-x-1/2 rounded-md bg-brand px-2.5 py-1.5 text-xs leading-relaxed text-white shadow-md"
      >
        <slot />
      </span>
    </transition>
  </span>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
