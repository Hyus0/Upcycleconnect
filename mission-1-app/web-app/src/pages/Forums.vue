<template>
  <main class="forums-shell">
    <SiteNavbar
      :is-authenticated="isLoggedIn"
      :user-name="userName"
      user-role="Particulier"
      :user-score="userScore"
      variant="public"
    />
    <ForumBoard />
  </main>
</template>

<script setup>
import { computed } from "vue";
import SiteNavbar from "../components/SiteNavbar.vue";
import ForumBoard from "../components/ForumBoard.vue";

const isLoggedIn = computed(() => Boolean(sessionStorage.getItem("userToken")));
const userScore = computed(() => sessionStorage.getItem("userScore") || 0);
const userName = computed(() => {
  const prenom = sessionStorage.getItem("userPrenom") || "";
  const nom = sessionStorage.getItem("userNom") || "";
  return (prenom || nom) ? `${prenom} ${nom}`.trim() : "Utilisateur";
});
</script>

<style scoped>
.forums-shell {
  min-height: 100vh;
  padding: 0 20px 20px;
  background: var(--bg-light);
}
</style>
