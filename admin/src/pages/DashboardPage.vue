<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Dashboard admin</div>
        <h2 class="page-title">Synthese des ressources critiques</h2>
        <p class="page-description">
          Les chiffres consolident uniquement les endpoints disponibles aujourd'hui, et signalent
          visiblement les zones backend encore absentes.
        </p>
      </div>
      <button class="button button-primary" @click="loadDashboard">Actualiser</button>
    </header>

    <LoadingState v-if="loading" />
    <ErrorState
      v-else-if="error"
      :message="error"
      retry-label="Recharger"
      @retry="loadDashboard"
    />

    <template v-else>
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
        <SimpleBarChart
          title="Utilisateurs par type"
          subtitle="Repartition selon les donnees admin disponibles"
          :items="dashboard.charts.usersByRole"
        />
        <SimpleBarChart
          title="Ressources cataloguees"
          subtitle="Prestations connectees, categories et evenements a completer"
          :items="dashboard.charts.resources"
        />
      </div>

      <div class="split-grid">
        <ResourceNotice
          title="Notes de statut rapide"
          message="Points d'attention remontes lors de la lecture des APIs."
          :items="dashboard.quickNotes.map((note) => note.text)"
          tone="green"
          badge="surveille"
        />

        <article class="surface-card section-card">
          <h3>Activites recentes</h3>
          <p class="muted">Flux derive des prestations disponibles.</p>
          <ul class="activity-list">
            <li v-for="item in dashboard.recentActivity" :key="item.id">
              <strong>{{ item.title }}</strong>
              <span>{{ item.subtitle }}</span>
            </li>
            <li v-if="dashboard.recentActivity.length === 0" class="muted">
              Aucune activite recente disponible.
            </li>
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
  stats: [],
  charts: {
    usersByRole: [],
    resources: []
  },
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
.activity-list {
  display: grid;
  gap: 12px;
  margin: 16px 0 0;
  padding: 0;
  list-style: none;
}

.activity-list li {
  display: grid;
  gap: 4px;
  padding: 12px 14px;
  border-radius: 16px;
  background: rgba(45, 122, 79, 0.05);
}

.activity-list span {
  color: var(--text-secondary);
}
</style>
