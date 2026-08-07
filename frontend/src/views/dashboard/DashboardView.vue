<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Line, Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js'
import * as dashboardApi from '@/api/dashboard'
import type { DashboardSummary } from '@/types'
import StatCard from '@/components/ui/StatCard.vue'
import AppCard from '@/components/ui/AppCard.vue'
import { formatMoney } from '@/lib/format'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend)

const summary = ref<DashboardSummary | null>(null)
const salesChartData = ref<{ labels: string[]; datasets: any[] }>({ labels: [], datasets: [] })
const profitChartData = ref<{ labels: string[]; datasets: any[] }>({ labels: [], datasets: [] })
const topProductsData = ref<{ labels: string[]; datasets: any[] }>({ labels: [], datasets: [] })
const loading = ref(true)

const chartOptions = { responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } } }

async function load() {
  loading.value = true
  const [summaryRes, salesRes, profitRes, topRes] = await Promise.all([
    dashboardApi.getSummary(),
    dashboardApi.getSalesOverTime(30),
    dashboardApi.getProfitTrend(30),
    dashboardApi.getTopProducts(30, 5),
  ])

  summary.value = summaryRes.data.data

  salesChartData.value = {
    labels: salesRes.data.data.map((p) => p.date),
    datasets: [{ label: 'Revenue', data: salesRes.data.data.map((p) => p.total), borderColor: '#2563eb', backgroundColor: '#93c5fd', tension: 0.3 }],
  }

  profitChartData.value = {
    labels: profitRes.data.data.map((p) => p.date),
    datasets: [{ label: 'Profit', data: profitRes.data.data.map((p) => p.profit), borderColor: '#059669', backgroundColor: '#6ee7b7', tension: 0.3 }],
  }

  topProductsData.value = {
    labels: topRes.data.data.map((p) => p.product_name),
    datasets: [{ label: 'Units Sold', data: topRes.data.data.map((p) => p.quantity_sold), backgroundColor: '#2563eb' }],
  }

  loading.value = false
}

onMounted(load)
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-xl font-semibold text-slate-900">Dashboard</h1>
      <p class="text-sm text-slate-500">Today's performance and inventory health</p>
    </div>

    <div v-if="summary" class="grid grid-cols-2 gap-4 lg:grid-cols-4">
      <StatCard label="Today's Sales" :value="formatMoney(summary.today.sales_amount)" tone="brand" />
      <StatCard label="Today's Profit" :value="formatMoney(summary.today.profit)" tone="green" />
      <StatCard label="Transactions" :value="String(summary.today.transaction_count)" />
      <StatCard label="Items Sold" :value="String(summary.today.items_sold)" />
    </div>

    <div v-if="summary" class="grid grid-cols-1 gap-4 sm:grid-cols-3">
      <StatCard label="Total Products" :value="String(summary.inventory.total_products)" />
      <StatCard label="Low Stock" :value="String(summary.inventory.low_stock_products)" tone="amber" />
      <StatCard label="Out of Stock" :value="String(summary.inventory.out_of_stock_products)" tone="red" />
    </div>

    <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
      <AppCard title="Sales Over Time (30 days)">
        <div class="h-64">
          <Line v-if="!loading" :data="salesChartData" :options="chartOptions" />
        </div>
      </AppCard>
      <AppCard title="Profit Trend (30 days)">
        <div class="h-64">
          <Line v-if="!loading" :data="profitChartData" :options="chartOptions" />
        </div>
      </AppCard>
    </div>

    <AppCard title="Top Selling Products (30 days)">
      <div class="h-64">
        <Bar v-if="!loading" :data="topProductsData" :options="chartOptions" />
      </div>
    </AppCard>
  </div>
</template>
