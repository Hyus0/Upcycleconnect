<template>
  <section class="dashboard-shell">
    <header class="dashboard-hero surface-card">
      <div>
        <div class="eyebrow">Back office connecte</div>
        <h2>Piloter UpcycleConnect</h2>
        <p>
          Administration centrale du front : comptes, catalogue, planning, moderation,
          messagerie, abonnements, paiements et notifications.
        </p>
      </div>
      <div class="hero-actions">
        <span class="api-pill" :class="{ degraded: dashboard.source === 'partial' }">
          {{ dashboard.source === "partial" ? "API partielle" : "API connectee" }}
        </span>
        <button class="button button-primary" @click="loadDashboard">Actualiser</button>
      </div>
    </header>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadDashboard" />

    <template v-else>
      <div class="stats-grid">
        <StatCard
          v-for="stat in dashboard.stats"
          :key="stat.label"
          :label="stat.label"
          :value="stat.value"
          :tone="stat.tone"
        />
      </div>

      <section class="domain-grid">
        <RouterLink
          v-for="domain in dashboard.domains"
          :key="domain.title"
          class="domain-card surface-card"
          :to="{ name: domain.route }"
        >
          <span>{{ domain.title }}</span>
          <strong>{{ domain.value }}</strong>
          <small>{{ domain.label }}</small>
          <p>{{ domain.description }}</p>
        </RouterLink>
      </section>

      <section class="control-grid">
        <article class="surface-card section-card">
          <div class="panel-head">
            <div>
              <span class="panel-kicker">Sante API</span>
              <h3>Regles et alertes</h3>
            </div>
            <span class="mini-note">{{ dashboard.quickNotes.length }} controles</span>
          </div>

          <ul class="ops-list">
            <li v-for="note in dashboard.quickNotes" :key="note.text" :class="`tone-${note.tone}`">
              <span></span>
              <p>{{ note.text }}</p>
            </li>
          </ul>

          <div v-if="dashboard.alerts.length" class="alert-box">
            <strong>Endpoints a verifier</strong>
            <p v-for="alert in dashboard.alerts" :key="alert.index">{{ alert.message }}</p>
          </div>
        </article>

        <article class="surface-card section-card">
          <div class="panel-head">
            <div>
              <span class="panel-kicker">Activite front</span>
              <h3>Derniers signaux</h3>
            </div>
          </div>

          <ul class="activity-list">
            <li v-for="item in dashboard.recentActivity" :key="item.id">
              <span>{{ item.kind }}</span>
              <strong>{{ item.title || "Sans titre" }}</strong>
              <p>{{ item.subtitle || "Aucune description" }}</p>
            </li>
            <li v-if="dashboard.recentActivity.length === 0" class="muted">Aucune activite disponible.</li>
          </ul>
        </article>
      </section>

      <div class="two-up">
        <SimpleBarChart title="Comptes par role" subtitle="Utilisateurs front" :items="dashboard.charts.usersByRole" />
        <SimpleBarChart title="Ressources gerees" subtitle="Domaines back" :items="dashboard.charts.resources" />
      </div>
    </template>
  </section>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import ErrorState from "../components/ErrorState.vue";
import LoadingState from "../components/LoadingState.vue";
import SimpleBarChart from "../components/SimpleBarChart.vue";
import StatCard from "../components/StatCard.vue";
import { adminApi } from "../services/api";

const loading = ref(true);
const error = ref("");
const dashboard = ref({
  source: "partial",
  stats: [],
  charts: { usersByRole: [], resources: [] },
  domains: [],
  quickNotes: [],
  alerts: [],
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
.dashboard-shell {
  display: grid;
  gap: 20px;
}

.dashboard-hero {
  position: relative;
  display: flex;
  justify-content: space-between;
  gap: 24px;
  align-items: center;
  padding: clamp(24px, 4vw, 38px);
  overflow: hidden;
  background:
    radial-gradient(circle at 78% 20%, rgba(98, 196, 136, 0.24), transparent 24%),
    linear-gradient(135deg, rgba(12, 36, 25, 0.98), rgba(8, 16, 13, 0.96));
}

.dashboard-hero::after {
  content: "";
  position: absolute;
  right: -80px;
  bottom: -120px;
  width: 320px;
  height: 320px;
  border: 1px solid rgba(136, 220, 168, 0.24);
  border-radius: 999px;
}

.dashboard-hero h2 {
  position: relative;
  z-index: 1;
  margin: 12px 0 0;
  max-width: 780px;
  font-family: "Syne", sans-serif;
  font-size: clamp(2.8rem, 7vw, 6.4rem);
  line-height: 0.86;
  letter-spacing: -0.08em;
}

.dashboard-hero p {
  position: relative;
  z-index: 1;
  max-width: 68ch;
  margin: 18px 0 0;
  color: var(--text-secondary);
  font-size: 1.05rem;
}

.hero-actions {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.api-pill {
  display: inline-flex;
  border: 1px solid rgba(118, 230, 160, 0.34);
  border-radius: 999px;
  padding: 12px 16px;
  color: var(--brand-green-light);
  background: rgba(98, 196, 136, 0.12);
  font-weight: 800;
}

.api-pill.degraded {
  color: #ffd68a;
  border-color: rgba(240, 178, 75, 0.4);
  background: rgba(240, 178, 75, 0.12);
}

.domain-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(230px, 1fr));
  gap: 16px;
}

.domain-card {
  display: grid;
  min-height: 210px;
  padding: 22px;
  color: inherit;
  background:
    linear-gradient(160deg, rgba(36, 78, 54, 0.34), transparent 46%),
    rgba(12, 24, 19, 0.86);
  transition: transform 0.18s ease, border-color 0.18s ease;
}

.domain-card:hover {
  transform: translateY(-3px);
  border-color: rgba(118, 230, 160, 0.42);
}

.domain-card span {
  color: var(--brand-green-light);
  font-family: "Space Mono", monospace;
  font-size: 0.72rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.domain-card strong {
  align-self: end;
  font-family: "Syne", sans-serif;
  font-size: 3.4rem;
  letter-spacing: -0.08em;
  line-height: 0.88;
}

.domain-card small {
  color: var(--text-secondary);
  font-weight: 800;
}

.domain-card p {
  margin: 12px 0 0;
  color: var(--text-secondary);
}

.control-grid {
  display: grid;
  grid-template-columns: minmax(300px, 0.9fr) minmax(420px, 1.1fr);
  gap: 16px;
}

.ops-list,
.activity-list {
  display: grid;
  gap: 12px;
  margin: 18px 0 0;
  padding: 0;
  list-style: none;
}

.ops-list li {
  display: grid;
  grid-template-columns: 12px 1fr;
  gap: 12px;
  align-items: start;
  padding: 14px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.035);
  color: var(--text-secondary);
}

.ops-list li span {
  width: 10px;
  height: 10px;
  margin-top: 5px;
  border-radius: 999px;
  background: var(--brand-green);
}

.ops-list li.tone-amber span {
  background: var(--accent-amber);
}

.ops-list p,
.activity-list p {
  margin: 0;
}

.alert-box {
  margin-top: 16px;
  padding: 16px;
  border: 1px solid rgba(229, 111, 92, 0.36);
  border-radius: 18px;
  background: rgba(229, 111, 92, 0.1);
}

.alert-box p {
  margin: 8px 0 0;
  color: #ffc5bd;
}

.activity-list li {
  display: grid;
  gap: 6px;
  padding: 16px;
  border: 1px solid rgba(139, 210, 166, 0.12);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.035);
}

.activity-list span {
  color: var(--brand-green-light);
  font-family: "Space Mono", monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.activity-list strong {
  font-family: "Syne", sans-serif;
  font-size: 1.12rem;
}

.activity-list p {
  color: var(--text-secondary);
}

@media (max-width: 900px) {
  .dashboard-hero,
  .control-grid {
    grid-template-columns: 1fr;
  }

  .dashboard-hero {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
