<template>
  <div class="admin-layout">
    <aside class="sidebar surface-card">
      <div>
        <div class="sidebar-badge">Admin console</div>
        <h1>Upcycle<span>Connect</span></h1>
        <p>
          Pilotage des ressources, supervision des flux et suivi des points API encore
          manquants.
        </p>
      </div>

      <nav class="nav-list">
        <RouterLink
          v-for="item in navItems"
          :key="item.name"
          :to="item.to"
          class="nav-link"
          active-class="nav-link-active"
        >
          <span class="nav-icon">{{ item.icon }}</span>
          <span>
            <strong>{{ item.label }}</strong>
            <small>{{ item.caption }}</small>
          </span>
        </RouterLink>
      </nav>

      <div class="sidebar-foot">
        <p>Theme aligne sur la charte graphique UpcycleConnect.</p>
      </div>
    </aside>

    <div class="main-column">
      <TopBar />
      <main class="page-shell">
        <RouterView />
      </main>
    </div>

    <ToastHost />
  </div>
</template>

<script setup>
import { RouterLink, RouterView } from "vue-router";
import TopBar from "./TopBar.vue";
import ToastHost from "../ToastHost.vue";

const navItems = [
  {
    name: "dashboard",
    label: "Dashboard",
    caption: "Vue d'ensemble",
    to: { name: "dashboard" },
    icon: "01"
  },
  {
    name: "users",
    label: "Utilisateurs",
    caption: "Particuliers, prestataires, admins",
    to: { name: "users" },
    icon: "02"
  },
  {
    name: "prestations",
    label: "Prestations",
    caption: "Catalogue et suivi",
    to: { name: "prestations" },
    icon: "03"
  },
  {
    name: "categories",
    label: "Categories",
    caption: "Taxonomie et hierarchie",
    to: { name: "categories" },
    icon: "04"
  },
  {
    name: "events",
    label: "Evenements",
    caption: "Programmation et statuts",
    to: { name: "events" },
    icon: "05"
  }
];
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 18px;
  padding: 18px;
}

.sidebar {
  padding: 24px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  background:
    linear-gradient(180deg, rgba(26, 31, 28, 0.98) 0%, rgba(46, 53, 48, 0.96) 100%),
    linear-gradient(120deg, rgba(75, 171, 115, 0.18), transparent);
  color: rgba(255, 255, 255, 0.88);
  position: sticky;
  top: 18px;
  height: calc(100vh - 36px);
  overflow: auto;
}

.sidebar h1 {
  margin: 12px 0 8px;
  font-family: "Syne", sans-serif;
  font-size: 2.4rem;
  line-height: 0.9;
}

.sidebar h1 span {
  color: var(--brand-green-light);
}

.sidebar p {
  color: rgba(255, 255, 255, 0.58);
}

.sidebar-badge {
  display: inline-flex;
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(75, 171, 115, 0.16);
  color: var(--brand-green-light);
  font-family: "Space Mono", monospace;
  font-size: 11px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.nav-list {
  display: grid;
  gap: 10px;
  margin-top: 28px;
}

.nav-link {
  display: grid;
  grid-template-columns: 42px 1fr;
  gap: 12px;
  align-items: center;
  padding: 14px;
  border-radius: 18px;
  color: rgba(255, 255, 255, 0.74);
  border: 1px solid transparent;
}

.nav-link small {
  display: block;
  margin-top: 2px;
  color: rgba(255, 255, 255, 0.44);
}

.nav-link-active {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(75, 171, 115, 0.25);
  color: white;
}

.nav-icon {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  border-radius: 14px;
  font-family: "Space Mono", monospace;
  background: rgba(255, 255, 255, 0.08);
}

.sidebar-foot {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.main-column {
  min-width: 0;
}

@media (max-width: 1080px) {
  .admin-layout {
    grid-template-columns: 1fr;
  }

  .sidebar {
    position: static;
    height: auto;
  }
}
</style>
