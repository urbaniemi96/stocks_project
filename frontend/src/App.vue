<template>
  <router-view />
</template>

<script setup lang="ts">
import { onBeforeMount } from 'vue'
import { useAuthStore } from './stores/auth'

const auth = useAuthStore()

onBeforeMount(async () => {
  // Traigo las credenciales del usuario
  try {
    await auth.fetchUser()
  } catch {
    // Usuario de demostración por defecto
    auth.setUser({ id: "demo-user", role: 'user' })
  }
})
</script>