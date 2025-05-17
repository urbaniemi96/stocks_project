import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: {
      id: 'demo-user',
      role: 'admin' as 'admin' | 'user',
    },
  }),
  getters: {
    isAdmin: (state) => state.user.role === 'admin'
  }
})