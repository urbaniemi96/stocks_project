import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AdminView from '../views/AdminView.vue'
import StockView from '../views/StockDashboard.vue'
import StockDetailView from '../views/Detail.vue'
import RecommendationsView from '../views/RecommendationsView.vue'
import { useAuthStore } from '../stores/auth'
import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/', 
    name: 'home', 
    component: HomeView 
  },
  { path: '/stocks', 
    name: 'dashboard', 
    component: StockView 
  },
  { path: '/configuration', 
    name: 'configuration', 
    component: AdminView,
    meta: { requiresAdmin: true }
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

// Analizo rol de las rutas
router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
  const auth = useAuthStore()

  // Si la ruta requiere admin, y no es admin, redirijo al home
  if (to.meta.requiresAdmin && !auth.isAdmin) {
    return next({ name: 'home' })
  }

  // En cualquier otro caso, dejo pasar
  return next()
})