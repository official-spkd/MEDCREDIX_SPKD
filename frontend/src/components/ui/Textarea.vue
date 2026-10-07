<script setup lang="ts">
const model = defineModel<string>({ default: '' })
defineProps<{ id?: string; rows?: number; placeholder?: string }>()
</script>
<template>
  <textarea
    :id="id"
    v-model="model"
    :rows="rows ?? 3"
    :placeholder="placeholder"
    class="flex w-full rounded-lg border border-input bg-card px-3 py-2 text-sm shadow-xs placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
  />
</template>
