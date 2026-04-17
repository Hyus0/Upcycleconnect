<template>
  <header class="site-navbar" :class="`site-navbar--${variant}`">
    <div class="site-navbar__inner">
      <RouterLink class="site-navbar__brand" to="/" aria-label="Accueil UpcycleConnect">
        <img :src="logoSrc" alt="UpcycleConnect" class="site-navbar__logo" />
      </RouterLink>

      <nav class="site-navbar__links" aria-label="Navigation principale">
        <RouterLink
          v-for="item in items"
          :key="item.label"
          :to="item.to"
          class="site-navbar__link"
          :class="{ 'is-active': isActive(item) }"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="site-navbar__actions" v-if="isAuthenticated">
        <RouterLink class="site-navbar__button site-navbar__button--primary" to="/profil/annonces">
          + Deposer
        </RouterLink>
        <RouterLink class="site-navbar__avatar" to="/profil" aria-label="Profil utilisateur">
          {{ userInitials }}
        </RouterLink>
      </div>

      <div class="site-navbar__actions" v-else>
        <RouterLink class="site-navbar__button site-navbar__button--ghost" to="/connexion">
          Connexion
        </RouterLink>
        <RouterLink class="site-navbar__button site-navbar__button--primary" to="/inscription">
          S'inscrire
        </RouterLink>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed } from "vue";
import { RouterLink, useRoute } from "vue-router";
import logoSrc from "./logo_texte.png";

const route = useRoute();

const props = defineProps({
  variant: {
    type: String,
    default: "public"
  },
  userName: {
    type: String,
    default: "Marie Lambert"
  },
  isAuthenticated: {
    type: Boolean,
    default: false
  }
});

const publicItems = [
  { label: "Comment ca marche", to: { path: "/", hash: "#processus" } },
  { label: "Annonces", to: "/annonces" },
  { label: "Formations", to: "/formations" },
  { label: "Communaute", to: "/communaute" }
];

const appItems = [
  { label: "Tableau de bord", to: "/profil" },
  { label: "Annonces", to: "/profil/annonces" },
  { label: "Formations", to: "/formations" },
  { label: "Communaute", to: "/communaute" }
];

const items = computed(() => (props.variant === "public" ? publicItems : appItems));

const userInitials = computed(() =>
  props.userName
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("")
);

function itemPath(item) {
  return typeof item.to === "string" ? item.to : item.to.path;
}

function isActive(item) {
  const path = itemPath(item);
  if (path === "/" && item.to.hash) {
    return route.path === "/" && route.hash === item.to.hash;
  }
  if (path === "/profil") {
    return route.path === "/profil";
  }
  return route.path === path || route.path.startsWith(`${path}/`);
}
</script>
