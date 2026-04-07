<template>
  <header class="topbar surface-card">
    <div class="left">
      <button class="menu-trigger" @click="$emit('toggle-menu')">
        <span></span>
        <span></span>
        <span></span>
      </button>
      <div>
        <div class="eyebrow">UpcycleConnect</div>
        <strong class="title">{{ pageTitle }}</strong>
        <span class="subtitle">{{ pageSubtitle }}</span>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed } from "vue";
import { useRoute } from "vue-router";

defineEmits(["toggle-menu"]);

const route = useRoute();
const pageTitle = computed(() => route.meta.title ?? "Admin");
const pageSubtitle = computed(() => route.meta.subtitle ?? "Back office");
</script>

<style scoped>
.topbar {
  margin: 20px 20px 0 0;
  padding: 14px 18px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  border-radius: 20px;
}

.left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.menu-trigger {
  display: none;
  width: 42px;
  height: 42px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.86);
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 4px;
}

.menu-trigger span {
  width: 16px;
  height: 2px;
  background: var(--brand-dark);
  border-radius: 999px;
}

.title {
  display: block;
  margin-top: 4px;
  font-size: 1.1rem;
}

.subtitle {
  display: block;
  margin-top: 4px;
  color: var(--text-secondary);
  font-size: 0.92rem;
}

@media (max-width: 1080px) {
  .topbar {
    margin-right: 0;
  }

  .menu-trigger {
    display: inline-flex;
  }
}

@media (max-width: 700px) {
  .topbar {
    flex-direction: row;
    align-items: center;
  }
}
</style>
