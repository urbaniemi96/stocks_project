import { defineStore } from 'pinia'
import { api } from '../services/api' // Mi api axios para peticiones

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: {
      id: 'demo-user',
      role: null as 'admin' | 'user' | null,
    },
  }),
  getters: {
    isAdmin: (state) => state.user.role === 'admin',
    isUser: (state) => state.user.role === 'user',
  },
  actions: {
    setUser(userData: { id: string, role: 'admin' | 'user' }) {
      this.user = userData
    },
    async fetchUser() {
      // Leo el usuario logueado
      const res = await api.get('/read/user')
      this.setUser(res.data)
    }
  }
})