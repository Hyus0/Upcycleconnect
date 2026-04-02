<template>
  <div class="admin-layout">
    <div v-if="sidebarOpen" class="sidebar-overlay" @click="sidebarOpen = false"></div>

    <aside class="sidebar surface-card" :class="{ open: sidebarOpen }">
      <div class="sidebar-brand">
        <div class="sidebar-badge">Admin</div>
        <img class="sidebar-logo" :src="logoSrc" alt="UpcycleConnect" />
      </div>

      <nav class="nav-group">
        <div class="nav-heading">Admin</div>
        <RouterLink
          v-for="item in adminItems"
          :key="item.name"
          :to="item.to"
          class="nav-link"
          active-class="nav-link-active"
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
import { ref } from "vue";
import { RouterLink, RouterView } from "vue-router";
import TopBar from "./TopBar.vue";
import ToastHost from "../ToastHost.vue";
import logoSrc from "../logo.png";

const sidebarOpen = ref(false);

const adminItems = [
  { name: "dashboard", label: "Dashboard", to: { name: "dashboard" }, icon: "01" },
  { name: "users", label: "Utilisateurs", to: { name: "users" }, icon: "02" },
  { name: "prestations", label: "Prestations", to: { name: "prestations" }, icon: "03" },
  { name: "categories", label: "Categories", to: { name: "categories" }, icon: "04" },
  { name: "events", label: "Evenements", to: { name: "events" }, icon: "05" }
];
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 272px 1fr;
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
  padding: 28px 22px;
  display: flex;
  flex-direction: column;
  gap: 18px;
  background:
    radial-gradient(circle at top, rgba(75, 171, 115, 0.12), transparent 28%),
    linear-gradient(180deg, rgba(24, 29, 27, 0.99) 0%, rgba(39, 45, 42, 0.98) 100%);
  color: rgba(255, 255, 255, 0.92);
  position: sticky;
  top: 20px;
  height: calc(100vh - 40px);
  overflow: auto;
  z-index: 50;
}

.sidebar-brand {
  margin-bottom: 6px;
}

.sidebar-logo {
  display: block;
  width: min(100%, 88px);
  margin-top: 14px;
  object-fit: contain;
}

.sidebar-badge {
  display: inline-flex;
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
  grid-template-columns: 42px 1fr;
  gap: 12px;
  align-items: center;
  padding: 14px 14px 14px 12px;
  border-radius: 20px;
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.03);
  transition: transform 0.18s ease, border-color 0.18s ease, background 0.18s ease;
}

.nav-link-active {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.06));
  border-color: rgba(75, 171, 115, 0.2);
  color: white;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.05);
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
