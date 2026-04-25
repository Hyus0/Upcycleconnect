<template>
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />
    <h1>Projets</h1>

</template>

<script setup>
import {ref, computed} from "vue";
import SiteNavbar from "../components/SiteNavbar.vue";

const isLoggedIn = computed(() => {
  return !!localStorage.getItem("userToken");
});

const userName = computed(() => {
  const prenom = localStorage.getItem("userPrenom") || "";
  const nom = localStorage.getItem("userNom") || "";
  
  return (prenom || nom) ? `${prenom} ${nom}`.trim() : "Utilisateur";
});
</script>