import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import StockDetailView from '../Detail.vue'

const routes: RouteRecordRaw[] = [
  { path: '/', 
    name: 'home', 
    component: HomeView 
  },
  {
    path: '/stocks/:ticker/detail',
    name: 'stock-detail',
    component: StockDetailView,
    // `ticker` queda disponible en route.params.ticker
  },
]

export const router = createRouter({
  history: createWebHistory(), // usa el API History de HTML5
  routes,
})