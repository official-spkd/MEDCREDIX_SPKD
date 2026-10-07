<script setup lang="ts">
// ChartTsLine — wrapper LineChart/AreaChart chart.js (pengganti recharts).
// Mendukung: multi dataset garis, area fill (gradient), dashed, opsi sumbu.
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend,
} from 'chart.js'
import { Line } from 'vue-chartjs'
import { computed } from 'vue'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Filler, Tooltip, Legend)

export type LineDataset = {
  label: string
  data: number[]
  color: string
  fill?: boolean
  dashed?: boolean
  width?: number
}

const props = defineProps<{
  labels: (string | number)[]
  datasets: LineDataset[]
  yMin?: number
  yMax?: number
  height?: number
  legend?: boolean
}>()

const chartData = computed(() => ({
  labels: props.labels,
  datasets: props.datasets.map((d) => ({
    label: d.label,
    data: d.data,
    borderColor: d.color,
    backgroundColor: (ctx: any) => {
      if (!d.fill) return 'transparent'
      const { chart } = ctx
      const { ctx: c, chartArea } = chart
      if (!chartArea) return d.color + '20'
      const gradient = c.createLinearGradient(0, chartArea.top, 0, chartArea.bottom)
      gradient.addColorStop(0, d.color + '2e')
      gradient.addColorStop(1, d.color + '05')
      return gradient
    },
    borderWidth: d.width ?? 2.5,
    borderDash: d.dashed ? [5, 4] : undefined,
    pointRadius: 3,
    pointHoverRadius: 5,
    pointBackgroundColor: d.color,
    pointBorderWidth: 0,
    tension: 0.35,
    fill: !!d.fill,
  })),
}))

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: props.legend ?? false },
    tooltip: {
      backgroundColor: 'rgba(20, 40, 92, 0.92)',
      titleFont: { size: 11, family: 'Poppins' },
      bodyFont: { size: 11, family: 'Open Sans' },
      padding: 10,
      cornerRadius: 8,
      displayColors: true,
      boxWidth: 8,
      boxHeight: 8,
      usePointStyle: true,
    },
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: { color: '#5B6B84', font: { size: 10, family: 'Open Sans' } },
      border: { color: '#E2E8F0' },
    },
    y: {
      min: props.yMin,
      max: props.yMax,
      grid: { color: '#E2E8F0', drawTicks: false },
      ticks: { color: '#5B6B84', font: { size: 10, family: 'Open Sans' }, padding: 6 },
      border: { display: false },
    },
  },
}))
</script>

<template>
  <div :style="{ height: (height ?? 208) + 'px' }">
    <Line :data="chartData" :options="chartOptions as any" />
  </div>
</template>
