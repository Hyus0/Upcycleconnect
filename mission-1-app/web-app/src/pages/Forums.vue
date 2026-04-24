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

const isLoggedIn = computed(() => Boolean(localStorage.getItem("userToken")));
const userScore = computed(() => localStorage.getItem("userScore") || 0);
const userName = computed(() => {
  const prenom = localStorage.getItem("userPrenom") || "";
  const nom = localStorage.getItem("userNom") || "";
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
