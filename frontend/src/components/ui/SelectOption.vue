<script setup lang="ts">
// SelectOption — item opsi untuk Select (registrasi via provide/inject).
import { inject, onMounted, useSlots } from 'vue'
const props = defineProps<{ value: string }>()
const register = inject<((value: string, label: string) => void) | null>('select-register', null)
if (register) register(props.value, props.value)
const slots = useSlots()
onMounted(() => {
  // label mengikuti slot default bila tersedia
  const label = slots.default?.()?.[0]?.children
  if (register && typeof label === 'string') register(props.value, label)
})
</script>
<template>
  <span class="hidden"><slot /></span>
</template>
