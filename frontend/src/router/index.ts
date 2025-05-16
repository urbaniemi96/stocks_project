import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import StockDetailView from '../views/Detail.vue'
import RecommendationsView from '../views/RecommendationsView.vue'

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
  { 
    path: '/recommendations', 
    name: 'Recommendations', 
    component: RecommendationsView },
]

export const router = createRouter({
  history: createWebHistory(), // usa el API History de HTML5
  routes,
})