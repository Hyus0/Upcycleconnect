<template>
  <main class="profile-shell">
    <SiteNavbar
      :is-authenticated="isLoggedIn"
      :user-name="userName"
      user-role="Particulier"
      :user-score="userScore"
    />
    <router-view />
  </main>
  <SiteFooter />

</template>

<script setup>
import { ref, computed } from "vue";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

const isLoggedIn = computed(() => {
  return !!sessionStorage.getItem("userToken");
});

const userName = computed(() => {
  const prenom = sessionStorage.getItem("userPrenom") || "";
  const nom = sessionStorage.getItem("userNom") || "";
  
  return (prenom || nom) ? `${prenom} ${nom}`.trim() : "Utilisateur";
});
</script>

<style scoped>
.profile-shell {
  min-height: 100vh;
  padding: 0 20px 20px 20px; 
  background: var(--bg-light);
}
</style>