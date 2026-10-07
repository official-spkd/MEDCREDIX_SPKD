<script setup lang="ts">
// Tabs — kontainer tab (padanan shadcn/ui Tabs) via provide/inject.
import { computed, provide } from 'vue'
const props = defineProps<{ modelValue: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
provide('tabs', {
  current: computed(() => props.modelValue),
  select: (v: string) => emit('update:modelValue', v),
})
</script>
<template>
  <div><slot /></div>
</template>
