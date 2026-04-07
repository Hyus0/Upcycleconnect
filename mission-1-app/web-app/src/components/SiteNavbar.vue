<template>
  <header class="site-navbar" :class="`site-navbar--${variant}`">
    <div class="site-navbar__inner">
      <a class="site-navbar__brand" href="#">
        <img :src="logoSrc" alt="UpcycleConnect" class="site-navbar__logo" />
      </a>

      <nav class="site-navbar__links">
        <a
          v-for="item in items"
          :key="item.label"
          :href="item.href ?? '#'"
          class="site-navbar__link"
          :class="{ 'is-active': item.active }"
        >
          {{ item.label }}
        </a>
      </nav>

      <div class="site-navbar__actions" v-if="isAuthenticated">
        <button class="site-navbar__button site-navbar__button--primary">+ Deposer</button>
        <div class="site-navbar__avatar">{{ userInitials }}</div>
      </div>

      <div class="site-navbar__actions" v-else>
        <button class="site-navbar__button site-navbar__button--ghost">Connexion</button>
        <button class="site-navbar__button site-navbar__button--primary">S'inscrire</button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed } from "vue";
import logoSrc from "./logo_texte.png";

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

const items = computed(() =>
  props.variant === "public"
    ? [
        { label: "Comment ca marche", active: false },
        { label: "Annonces", active: true },
        { label: "Formations", active: false },
        { label: "Communaute", active: false }
      ]
    : [
        { label: "Tableau de bord", active: true },
        { label: "Annonces", active: false },
        { label: "Formations", active: false }
      ]
);

const userInitials = computed(() =>
  props.userName
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("")
);
</script>
