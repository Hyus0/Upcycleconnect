<template>
  <header class="topbar surface-card">
    <div>
      <div class="eyebrow">Administration separee</div>
      <strong class="title">Vue.js admin connected to the existing APIs</strong>
    </div>
    <div class="topbar-meta">
      <span class="pill">Vue 3 Composition API</span>
      <span class="pill pill-green">REST orchestration</span>
    </div>
  </header>
</template>

<style scoped>
.topbar {
  margin: 18px 18px 0 0;
  padding: 18px 22px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.title {
  display: block;
  margin-top: 6px;
  font-size: 1rem;
}

.topbar-meta {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.pill {
  border-radius: 999px;
  padding: 8px 12px;
  background: var(--brand-sand);
  color: var(--brand-earth);
  font-size: 0.85rem;
  font-weight: 700;
}

.pill-green {
  background: var(--brand-green-pale);
  color: var(--brand-green);
}

@media (max-width: 1080px) {
  .topbar {
    margin-right: 0;
  }
}

@media (max-width: 700px) {
  .topbar {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
