<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Dashboard</div>
        <h2 class="page-title">Vue d'ensemble</h2>
      </div>
      <button class="button button-primary" @click="loadDashboard">Actualiser</button>
    </header>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadDashboard" />

    <template v-else>
      <article class="surface-card section-card dashboard-hero">
        <div>
          <div class="eyebrow">Mode</div>
          <strong>{{ dashboard.source === "local" ? "Base locale active" : "API connectee" }}</strong>
        </div>
        <div class="hero-actions">
          <span class="hero-chip">{{ dashboard.stats[0]?.value ?? 0 }} utilisateurs</span>
          <span class="hero-chip">{{ dashboard.stats[1]?.value ?? 0 }} prestations</span>
        </div>
      </article>

      <div class="stats-grid">
        <StatCard
          v-for="stat in dashboard.stats"
          :key="stat.label"
          :label="stat.label"
          :value="stat.value"
          :tone="stat.tone"
          :missing="stat.missing"
        />
      </div>

      <div class="two-up">
        <SimpleBarChart title="Utilisateurs par type" subtitle="Repartition" :items="dashboard.charts.usersByRole" />
        <SimpleBarChart title="Ressources" subtitle="Vue globale" :items="dashboard.charts.resources" />
      </div>

      <div class="split-grid">
        <ResourceNotice
          title="Statut"
          message="Points a surveiller."
          :items="dashboard.quickNotes.map((note) => note.text)"
          tone="green"
          badge="surveille"
        />

        <article class="surface-card section-card">
          <h3>Activites</h3>
          <ul class="activity-list">
            <li v-for="item in dashboard.recentActivity" :key="item.id">
              <strong>{{ item.title }}</strong>
              <span>{{ item.subtitle }}</span>
            </li>
            <li v-if="dashboard.recentActivity.length === 0" class="muted">Aucune activite recente disponible.</li>
          </ul>
        </article>
      </div>
    </template>
  </section>
</template>

<script setup>
import { onMounted, ref } from "vue";
import ErrorState from "../components/ErrorState.vue";
import LoadingState from "../components/LoadingState.vue";
import ResourceNotice from "../components/ResourceNotice.vue";
import SimpleBarChart from "../components/SimpleBarChart.vue";
import StatCard from "../components/StatCard.vue";
import { adminApi } from "../services/api";

const loading = ref(true);
const error = ref("");
const dashboard = ref({
  source: "local",
  stats: [],
  charts: { usersByRole: [], resources: [] },
  quickNotes: [],
  recentActivity: []
});

async function loadDashboard() {
  loading.value = true;
  error.value = "";
  try {
    dashboard.value = await adminApi.getDashboard();
  } catch (err) {
    error.value = err.message ?? "Impossible de charger le dashboard.";
  } finally {
    loading.value = false;
  }
}

onMounted(loadDashboard);
</script>

<style scoped>
.dashboard-hero {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.dashboard-hero strong {
  display: block;
  margin-top: 8px;
  font-size: 1.2rem;
}

.hero-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.hero-chip {
  padding: 10px 12px;
  border-radius: 999px;
  background: rgba(45, 122, 79, 0.08);
  color: var(--brand-green);
  font-weight: 700;
}

.activity-list {
  display: grid;
  gap: 12px;
  margin: 14px 0 0;
  padding: 0;
  list-style: none;
}

.activity-list li {
  display: grid;
  gap: 4px;
  padding: 14px 16px;
  border-radius: 16px;
  background: linear-gradient(180deg, rgba(45, 122, 79, 0.06), rgba(45, 122, 79, 0.03));
}

.activity-list span {
  color: var(--text-secondary);
}

@media (max-width: 700px) {
  .dashboard-hero {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
