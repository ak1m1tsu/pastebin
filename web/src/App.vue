<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <header class="flex items-center justify-center bg-gray-800">
    <div class="container px-[30px] py-2">
      <nav class="flex text-gray-200 items-center gap-4 justify-end">
        <RouterLink class="hover:text-emerald-500" to="/">Home</RouterLink>
        <RouterLink class="hover:text-emerald-500" to="/profile">Profile</RouterLink>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.router-link-active {
  @apply text-emerald-500;
}
</style>
