<template>
  <header class="topbar surface-card">
    <div class="left">
      <button class="menu-trigger" @click="$emit('toggle-menu')">
        <span></span>
        <span></span>
        <span></span>
      </button>
      <RouterLink class="brand-link" :to="{ name: 'dashboard' }">
        <img :src="logoTextSrc" alt="UpcycleConnect" class="brand-logo" />
      </RouterLink>
    </div>

    <nav class="topbar-nav">
      <RouterLink
        v-for="item in navItems"
        :key="item.name"
        :to="item.to"
        class="topbar-link"
        :class="{ 'topbar-link-active': route.name === item.name }"
      >
        {{ item.label }}
      </RouterLink>
    </nav>

    <div class="topbar-actions">
      <span class="page-chip">{{ pageTitle }}</span>
      <div class="admin-avatar">AD</div>
    </div>
  </header>
</template>

<script setup>
import { computed } from "vue";
import { RouterLink, useRoute } from "vue-router";
import logoTextSrc from "../logo_texte.png";

defineEmits(["toggle-menu"]);

const route = useRoute();
const pageTitle = computed(() => route.meta.title ?? "Admin");
const navItems = [
  { name: "dashboard", label: "Tableau de bord", to: { name: "dashboard" } },
  { name: "prestations", label: "Annonces", to: { name: "prestations" } },
  { name: "events", label: "Evenements", to: { name: "events" } },
  { name: "notifications", label: "Notifications", to: { name: "notifications" } }
];
</script>

<style scoped>
.topbar {
  margin: 20px 20px 0 0;
  padding: 14px 18px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  border-radius: 24px;
  background: rgba(23, 30, 27, 0.95);
}

.left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.menu-trigger {
  display: none;
  width: 42px;
  height: 42px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.08);
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 4px;
}

.menu-trigger span {
  width: 16px;
  height: 2px;
  background: #f4faf6;
  border-radius: 999px;
}

.brand-link {
  display: inline-flex;
  align-items: center;
}

.brand-logo {
  width: 180px;
  height: auto;
  object-fit: contain;
}

.topbar-nav {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

.topbar-link {
  display: inline-flex;
  align-items: center;
  min-height: 42px;
  padding: 0 16px;
  border-radius: 14px;
  color: rgba(235, 243, 238, 0.72);
  text-decoration: none;
  transition: background 0.2s ease, color 0.2s ease;
}

.topbar-link:hover,
.topbar-link-active {
  background: rgba(56, 124, 81, 0.26);
  color: #eef9f2;
}

.topbar-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-chip {
  padding: 0 14px;
  min-height: 40px;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.06);
  color: var(--text-secondary);
  font-weight: 700;
}

.admin-avatar {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  background: linear-gradient(180deg, #3e985f, #30794b);
  color: white;
  font-size: 0.82rem;
  font-weight: 700;
}

@media (max-width: 1080px) {
  .topbar {
    margin-right: 0;
    flex-wrap: wrap;
  }

  .menu-trigger {
    display: inline-flex;
  }

  .topbar-nav {
    order: 3;
    width: 100%;
    overflow-x: auto;
    margin-left: 0;
  }
}

@media (max-width: 700px) {
  .topbar {
    padding: 14px;
  }

  .brand-logo {
    width: 132px;
  }

  .page-chip {
    display: none;
  }
}
</style>
