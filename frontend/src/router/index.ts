import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import type { RouterScrollBehavior } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AdminView from '../views/AdminView.vue'
import StockView from '../views/StockDashboard.vue'
import StockDetailView from '../views/Detail.vue'
import RecommendationsView from '../views/RecommendationsView.vue'
import ThankYouView from '../views/ThankYouView.vue'
import PayView from '../views/PayView.vue'
import { useAuthStore } from '../stores/auth'
import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import TetrisGrid from '../components/TetrisGrid.vue'


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
    component: RecommendationsView 
  },
  {
    path: '/suscription',
    name: 'suscription',
    component: PayView,
  },
  {
    path: '/thanks',
    name: 'thanks',
    component: ThankYouView,
  },
  {
    path: '/tetris',
    name: 'tetris',
    component: TetrisGrid,
  },
]

// Para reiniciar el scroll al cambiar de página
const scrollBehavior: RouterScrollBehavior = (to, from, savedPosition) => {
  // Si venimos con "historial" (back/forward) volvemos a la posición guardada
  if (savedPosition) {
    return savedPosition
  }
  // si no volvemos al top
  return { top: 0 }
}

export const router = createRouter({
  history: createWebHistory(), // usa el API History de HTML5
  routes,
  scrollBehavior
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
