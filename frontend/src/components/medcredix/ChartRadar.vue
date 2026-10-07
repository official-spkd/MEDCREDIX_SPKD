<script setup lang="ts">
// ChartRadar — wrapper RadarChart chart.js (Risk Radar).
import {
  Chart as ChartJS, RadialLinearScale, PointElement, LineElement, Filler, Tooltip, Legend,
} from 'chart.js'
import { Radar } from 'vue-chartjs'
import { computed } from 'vue'

ChartJS.register(RadialLinearScale, PointElement, LineElement, Filler, Tooltip, Legend)

const props = defineProps<{
  labels: string[]
  datasets: { label: string; data: number[]; color: string }[]
  height?: number
  max?: number
}>()

const chartData = computed(() => ({
  labels: props.labels,
  datasets: props.datasets.map((d) => ({
    label: d.label,
    data: d.data,
    borderColor: d.color,
    backgroundColor: d.color + '2b',
    pointBackgroundColor: d.color,
    pointBorderColor: '#fff',
    pointHoverRadius: 5,
    borderWidth: 2,
  })),
}))

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: true, labels: { color: '#5B6B84', font: { size: 10, family: 'Open Sans' }, boxWidth: 8, usePointStyle: true } },
    tooltip: {
      backgroundColor: 'rgba(20, 40, 92, 0.92)',
      titleFont: { size: 11, family: 'Poppins' },
      bodyFont: { size: 11, family: 'Open Sans' },
      padding: 10, cornerRadius: 8,
    },
  },
  scales: {
    r: {
      min: 0,
      max: props.max ?? 100,
      ticks: { display: false, stepSize: 25 },
      grid: { color: '#E2E8F0' },
      angleLines: { color: '#E2E8F0' },
      pointLabels: { color: '#5B6B84', font: { size: 10, family: 'Open Sans' } },
    },
  },
}))
</script>

<template>
  <div :style="{ height: (height ?? 260) + 'px' }">
    <Radar :data="chartData" :options="chartOptions as any" />
  </div>
</template>
