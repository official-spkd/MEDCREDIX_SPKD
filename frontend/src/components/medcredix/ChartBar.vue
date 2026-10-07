<script setup lang="ts">
// ChartBar — wrapper BarChart chart.js (pengganti recharts BarChart).
import {
  Chart as ChartJS, CategoryScale, LinearScale, BarElement, Tooltip, Legend,
} from 'chart.js'
import { Bar } from 'vue-chartjs'
import { computed } from 'vue'

ChartJS.register(CategoryScale, LinearScale, BarElement, Tooltip, Legend)

export type BarDataset = { label: string; data: number[]; color: string | string[]; stack?: string }

const props = defineProps<{
  labels: (string | number)[]
  datasets: BarDataset[]
  height?: number
  stacked?: boolean
  horizontal?: boolean
  legend?: boolean
  suggestedMax?: number
}>()

const chartData = computed(() => ({
  labels: props.labels,
  datasets: props.datasets.map((d) => ({
    label: d.label,
    data: d.data,
    backgroundColor: d.color,
    borderRadius: 6,
    maxBarThickness: 34,
    stack: d.stack,
  })),
}))

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  indexAxis: props.horizontal ? ('y' as const) : ('x' as const),
  plugins: {
    legend: { display: props.legend ?? false },
    tooltip: {
      backgroundColor: 'rgba(20, 40, 92, 0.92)',
      titleFont: { size: 11, family: 'Poppins' },
      bodyFont: { size: 11, family: 'Open Sans' },
      padding: 10, cornerRadius: 8, usePointStyle: true,
    },
  },
  scales: {
    x: {
      stacked: !!props.stacked,
      grid: { display: props.horizontal, color: '#E2E8F0' },
      ticks: { color: '#5B6B84', font: { size: 10, family: 'Open Sans' }, autoSkip: false, maxRotation: 0 },
      border: { color: '#E2E8F0' },
      suggestedMax: props.suggestedMax,
    },
    y: {
      stacked: !!props.stacked,
      grid: { display: !props.horizontal, color: '#E2E8F0', drawTicks: false },
      ticks: { color: '#5B6B84', font: { size: 10, family: 'Open Sans' }, padding: 6 },
      border: { display: false },
      suggestedMax: props.suggestedMax,
    },
  },
}))
</script>

<template>
  <div :style="{ height: (height ?? 208) + 'px' }">
    <Bar :data="chartData" :options="chartOptions as any" />
  </div>
</template>
