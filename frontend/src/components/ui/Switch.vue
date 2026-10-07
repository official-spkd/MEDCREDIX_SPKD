<script setup lang="ts">
// Switch — toggle kecil (padanan shadcn/ui Switch).
const model = defineModel<boolean>({ default: false })
</script>
<template>
  <button
    type="button"
    role="switch"
    :aria-checked="model"
    :class="['inline-flex h-5 w-9 shrink-0 items-center rounded-full border-2 border-transparent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring', model ? 'bg-brand-2' : 'bg-input']"
    @click="model = !model"
  >
    <span :class="['pointer-events-none block h-4 w-4 rounded-full bg-white shadow-lg ring-0 transition-transform', model ? 'translate-x-4' : 'translate-x-0']" />
  </button>
</template>
