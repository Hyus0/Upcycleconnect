<template>
  <div class="admin-layout">
    <div v-if="sidebarOpen" class="sidebar-overlay" @click="sidebarOpen = false"></div>

    <aside class="sidebar surface-card" :class="{ open: sidebarOpen }">
      <div class="sidebar-brand">
        <img class="sidebar-logo" :src="logoSrc" alt="UpcycleConnect" />
        <div class="sidebar-badge">Admin</div>
      </div>

      <nav class="nav-group">
        <div class="nav-heading">Pilotage front</div>
        <RouterLink
          v-for="item in adminItems"
          :key="item.name"
          :to="item.to"
          class="nav-link"
          :class="{ 'nav-link-active': currentRouteName === item.name }"
          @click="sidebarOpen = false"
        >
          <span class="nav-icon">{{ item.icon }}</span>
          <span><strong>{{ item.label }}</strong></span>
        </RouterLink>
      </nav>
    </aside>

    <div class="main-column">
      <TopBar @toggle-menu="sidebarOpen = !sidebarOpen" />
      <main class="page-shell">
        <RouterView />
      </main>
    </div>

    <ToastHost />
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { RouterLink, RouterView, useRoute } from "vue-router";
import TopBar from "./TopBar.vue";
import ToastHost from "../ToastHost.vue";
import logoSrc from "../logo.png";

const sidebarOpen = ref(false);
const route = useRoute();

const adminItems = [
  { name: "dashboard", label: "Vue globale", to: { name: "dashboard" }, icon: "01" },
  { name: "users", label: "Comptes", to: { name: "users" }, icon: "02" },
  { name: "prestations", label: "Catalogue", to: { name: "prestations" }, icon: "03" },
  { name: "categories", label: "Categories", to: { name: "categories" }, icon: "04" },
  { name: "events", label: "Planning", to: { name: "events" }, icon: "05" },
  { name: "moderation", label: "Moderation", to: { name: "moderation" }, icon: "06" },
  { name: "finance", label: "Paiements", to: { name: "finance" }, icon: "07" },
  { name: "notifications", label: "Notifications", to: { name: "notifications" }, icon: "08" },
  { name: "system", label: "Messages & ventes", to: { name: "system" }, icon: "09" }
];

const currentRouteName = computed(() => route.name);
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 292px 1fr;
  gap: 20px;
  padding: 20px;
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: rgba(26, 31, 28, 0.45);
  z-index: 40;
}

.sidebar {
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  gap: 18px;
  background:
    radial-gradient(circle at 50% 0%, rgba(98, 196, 136, 0.18), transparent 34%),
    linear-gradient(180deg, rgba(16, 27, 23, 0.98) 0%, rgba(23, 33, 29, 0.98) 100%);
  color: rgba(255, 255, 255, 0.92);
  position: sticky;
  top: 20px;
  height: calc(100vh - 40px);
  overflow: auto;
  z-index: 50;
}

.sidebar-brand {
  margin-bottom: 6px;
  display: grid;
  gap: 18px;
  padding: 8px 8px 18px;
  border-bottom: 1px solid rgba(152, 179, 163, 0.14);
}

.sidebar-logo {
  display: block;
  width: min(100%, 92px);
  object-fit: contain;
}

.sidebar-badge {
  display: inline-flex;
  width: fit-content;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(45, 122, 79, 0.2);
  color: var(--brand-green-light);
  font-family: "Space Mono", monospace;
  font-size: 11px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.nav-group {
  display: grid;
  gap: 12px;
}

.nav-heading {
  padding: 0 4px;
  color: rgba(255, 255, 255, 0.45);
  font-family: "Space Mono", monospace;
  font-size: 11px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.nav-link {
  display: grid;
  grid-template-columns: 46px 1fr;
  gap: 12px;
  align-items: center;
  padding: 13px 14px 13px 12px;
  border-radius: 18px;
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(152, 179, 163, 0.08);
  background: rgba(255, 255, 255, 0.02);
  transition: transform 0.18s ease, border-color 0.18s ease, background 0.18s ease;
}

.nav-link-active {
  background:
    radial-gradient(circle at 0% 50%, rgba(118, 230, 160, 0.28), transparent 38%),
    linear-gradient(180deg, rgba(54, 141, 92, 0.42), rgba(30, 79, 52, 0.44));
  border-color: rgba(98, 196, 136, 0.45);
  color: white;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.08),
    0 12px 28px rgba(15, 57, 34, 0.28);
}

.nav-link-active .nav-icon {
  background: rgba(255, 255, 255, 0.12);
  color: #f8fffb;
}

.nav-link:hover {
  transform: translateX(2px);
  border-color: rgba(75, 171, 115, 0.12);
}

.nav-icon {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  border-radius: 14px;
  font-family: "Space Mono", monospace;
  background: rgba(255, 255, 255, 0.07);
}

.main-column {
  min-width: 0;
}

@media (max-width: 1080px) {
  .admin-layout {
    grid-template-columns: 1fr;
  }

  .sidebar {
    position: fixed;
    left: 20px;
    top: 20px;
    bottom: 20px;
    height: auto;
    width: min(320px, calc(100vw - 40px));
    transform: translateX(calc(-100% - 20px));
    transition: transform 0.22s ease;
  }

  .sidebar.open {
    transform: translateX(0);
  }
}
</style>
