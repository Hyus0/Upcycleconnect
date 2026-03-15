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
        <strong class="title">Admin</strong>
      </div>
    </div>
    <div class="topbar-meta">
      <span class="pill pill-green">Local ready</span>
    </div>
  </header>
</template>

<script setup>
defineEmits(["toggle-menu"]);
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

.topbar-meta {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.pill {
  border-radius: 999px;
  padding: 8px 11px;
  background: var(--brand-sand);
  color: var(--brand-earth);
  font-size: 0.78rem;
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
