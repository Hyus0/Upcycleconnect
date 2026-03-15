<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Espace particulier</div>
        <h2 class="page-title">Mon espace</h2>
        <p class="page-description">Premiere base front pour les parcours utilisateurs.</p>
      </div>
    </header>

    <div class="two-up">
      <article class="surface-card section-card">
        <h3>Actions rapides</h3>
        <div class="quick-grid">
          <button v-for="item in snapshot.quickActions" :key="item" class="quick-card">{{ item }}</button>
        </div>
      </article>

      <article class="surface-card section-card">
        <h3>Prestations recommandees</h3>
        <ul class="portal-list">
          <li v-for="item in snapshot.prestations" :key="item.id">
            <strong>{{ item.title }}</strong>
            <span>{{ item.description }}</span>
          </li>
        </ul>
      </article>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { adminApi } from "../../services/api";

const snapshot = ref({ quickActions: [], prestations: [] });

onMounted(async () => {
  const data = await adminApi.getPortalSnapshot();
  snapshot.value = data.particulier;
});
</script>

<style scoped>
.quick-grid {
  display: grid;
  gap: 12px;
}

.quick-card {
  border: 1px solid var(--border);
  background: linear-gradient(180deg, rgba(45, 122, 79, 0.08), rgba(255, 255, 255, 0.8));
  border-radius: 18px;
  padding: 16px;
  text-align: left;
  font-weight: 700;
}

.portal-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 12px;
}

.portal-list li {
  display: grid;
  gap: 4px;
  padding: 14px 16px;
  border-radius: 16px;
  background: rgba(45, 122, 79, 0.05);
}

.portal-list span {
  color: var(--text-secondary);
}
</style>
