<script setup lang="ts">
import { computed, inject } from 'vue'
const props = defineProps<{ value: string }>()
const tabs = inject<{ current: { value: string } }>('tabs')!
const show = computed(() => tabs.current.value === props.value)
</script>
<template>
  <div v-if="show"><slot /></div>
</template>
