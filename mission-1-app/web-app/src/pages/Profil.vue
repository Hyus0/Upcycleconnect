<template>
  <main class="profile-shell">
    <SiteNavbar
      :is-authenticated="isLoggedIn"
      :user-name="userName"
      user-role="Particulier"
      :user-score="847"
    />
    <router-view />
  </main>
</template>

<script setup>
import { ref } from "vue";
import SiteNavbar from "../components/SiteNavbar.vue";

const isLoggedIn = ref(true);
const userName = ref("Marie Lambert");
</script>

<style scoped>
.profile-shell {
  min-height: 100vh;
  padding: 20px;
  background: var(--bg-light);
}
</style>
