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

    <div class="topbar-actions">
      <span class="sync-dot"></span>
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
</script>

<style scoped>
.topbar {
  margin: 20px 20px 0 0;
  padding: 16px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  border-radius: 28px;
  background:
    linear-gradient(90deg, rgba(14, 26, 21, 0.98), rgba(21, 33, 28, 0.94)),
    rgba(23, 30, 27, 0.95);
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
  width: 164px;
  height: auto;
  object-fit: contain;
}

.topbar-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

.page-chip {
  padding: 0 14px;
  min-height: 40px;
  border-radius: 14px;
  display: inline-flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.06);
  color: var(--text-secondary);
  font-weight: 700;
}

.sync-dot {
  width: 12px;
  height: 12px;
  border-radius: 999px;
  background: var(--brand-green);
  box-shadow: 0 0 0 8px rgba(98, 196, 136, 0.12);
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
