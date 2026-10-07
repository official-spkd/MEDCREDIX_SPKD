<script setup lang="ts">
// Input — padanan shadcn/ui Input (v-model).
const model = defineModel<string | number>({ default: '' })
defineProps<{ type?: string; id?: string; placeholder?: string; required?: boolean; disabled?: boolean }>()
</script>
<template>
  <input
    :id="id"
    v-model="model"
    :type="type ?? 'text'"
    :placeholder="placeholder"
    :required="required"
    :disabled="disabled"
    :class="[
      'flex h-9 w-full rounded-lg border border-input bg-card px-3 py-1 text-sm shadow-xs transition-colors',
      'placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
      'disabled:cursor-not-allowed disabled:opacity-50',
    ]"
  />
</template>
