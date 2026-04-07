<template>
  <div class="dashboard-shell">
    <aside class="sidebar">
      <div class="sidebar-brand">
        <img :src="logoSrc" alt="UpcycleConnect" class="sidebar-brand__logo" />
      </div>

      <article class="sidebar-user">
        <div class="sidebar-user__avatar">{{ displayInitials(currentUser.fullName) }}</div>
        <div>
          <strong>{{ displayValue(currentUser.fullName) }}</strong>
          <span>{{ displayValue(currentUser.role) }} - Score {{ displayValue(stats.score) }} pts</span>
        </div>
      </article>

      <div class="sidebar-group">
        <span class="sidebar-label">Principal</span>
        <a class="sidebar-link active" href="#">Tableau de bord</a>
        <a class="sidebar-link" href="#">Mes annonces</a>
        <a class="sidebar-link" href="#">Mes depots conteneurs</a>
        <a class="sidebar-link" href="#">Upcycling Score</a>
      </div>

      <div class="sidebar-group">
        <span class="sidebar-label">Services</span>
        <a class="sidebar-link" href="#">Formations & Ateliers</a>
        <a class="sidebar-link" href="#">Mon planning</a>
        <a class="sidebar-link" href="#">Espace Conseils</a>
        <a class="sidebar-link" href="#">Catalogue offres</a>
      </div>

      <div class="sidebar-group">
        <span class="sidebar-label">Communaute</span>
        <a class="sidebar-link" href="#">Forums</a>
        <a class="sidebar-link" href="#">Evenements</a>
      </div>

      <button class="sidebar-logout">Se deconnecter</button>
    </aside>

    <main class="dashboard-main">
      <SiteNavbar variant="app" :user-name="currentUser.fullName ?? 'Marie Lambert'" :is-authenticated="true" />

      <header class="dashboard-header">
        <div>
          <div class="breadcrumbs">Accueil > Tableau de bord</div>
          <h1>Bonjour {{ displayValue(currentUser.firstName) }} 👋</h1>
          <p>Voici un resume de votre activite sur UpcycleConnect.</p>
        </div>
        <button class="primary-action">+ Deposer une annonce</button>
      </header>

      <section class="headline-grid">
        <article class="score-card">
          <div class="card-kicker">Upcycling Score</div>
          <strong>{{ displayValue(stats.score, " pts") }}</strong>
          <p>Niveau : {{ displayValue(stats.level) }}</p>
          <div class="score-details">
            <div>
              <strong>{{ displayValue(stats.co2, " kg") }}</strong>
              <span>CO2 evite</span>
            </div>
            <div>
              <strong>{{ displayValue(stats.recycled) }}</strong>
              <span>Objets recycles</span>
            </div>
            <div>
              <strong>{{ displayValue(stats.saved, " EUR") }}</strong>
              <span>Economise</span>
            </div>
          </div>
        </article>

        <article class="metric-card">
          <strong>{{ displayValue(metrics.activeAnnonces) }}</strong>
          <span>Annonces actives</span>
          <small>{{ metrics.activeAnnonces === null ? "NULL" : "+1 ce mois" }}</small>
        </article>

        <article class="metric-card warning">
          <strong>{{ displayValue(metrics.pendingDeposits) }}</strong>
          <span>Depots en attente</span>
          <small>{{ metrics.pendingDeposits === null ? "NULL" : "EN COURS" }}</small>
        </article>
      </section>

      <section class="panel">
        <div class="panel-head">
          <h2>Mes dernieres annonces</h2>
          <div class="panel-actions">
            <input type="search" placeholder="Rechercher..." />
            <select>
              <option>Tous statuts</option>
            </select>
            <button class="secondary-action">+ Nouvelle annonce</button>
          </div>
        </div>

        <table class="data-table">
          <thead>
            <tr>
              <th>Objet</th>
              <th>Categorie</th>
              <th>Type</th>
              <th>Statut</th>
              <th>Date</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in latestAnnonces" :key="item.id">
              <td>{{ displayValue(item.title) }}</td>
              <td>{{ displayValue(item.category) }}</td>
              <td><span class="pill">{{ displayValue(item.type) }}</span></td>
              <td><span class="pill soft">{{ displayValue(item.status) }}</span></td>
              <td>{{ displayValue(item.date) }}</td>
              <td class="table-actions">
                <button>Voir</button>
                <button class="ghost">Modifier</button>
              </td>
            </tr>
          </tbody>
        </table>
      </section>

      <section class="panel">
        <div class="panel-head">
          <h2>Mon planning - semaine</h2>
          <button class="outline-action">Vue mensuelle</button>
        </div>

        <div class="calendar-grid">
          <article v-for="day in planningDays" :key="day.label" class="calendar-day">
            <span class="calendar-day__label">{{ day.label }}</span>
            <strong>{{ day.date }}</strong>
            <div class="calendar-event">
              {{ displayValue(day.event) }}
            </div>
          </article>
        </div>
      </section>

      <section class="bottom-grid">
        <article class="info-card">
          <div class="card-kicker">Conseil du jour</div>
          <h3>{{ displayValue(advice.title) }}</h3>
          <p>{{ displayValue(advice.content) }}</p>
          <a href="#">Lire la suite ></a>
        </article>

        <article class="info-card warning-card">
          <div class="card-kicker">Notification</div>
          <h3>{{ displayValue(notification.title) }}</h3>
          <p>{{ displayValue(notification.content) }}</p>
          <a href="#">Voir le projet ></a>
        </article>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import SiteNavbar from "./components/SiteNavbar.vue";
import logoSrc from "./components/logo.png";
import { fetchUserDashboard } from "./services/userDashboardApi";

const state = ref({
  user: null,
  annonces: null,
  planning: null,
  advice: null,
  notification: null
});

const currentUser = computed(() => ({
  firstName: state.value.user?.firstName ?? null,
  fullName: state.value.user?.fullName ?? null,
  role: state.value.user?.role ?? null
}));

const stats = computed(() => ({
  score: null,
  level: null,
  co2: null,
  recycled: state.value.annonces ? state.value.annonces.length : null,
  saved: null
}));

const metrics = computed(() => ({
  activeAnnonces: state.value.annonces
    ? state.value.annonces.filter((item) => (item.statut ?? "").toLowerCase() !== "vendu").length
    : null,
  pendingDeposits: state.value.planning
    ? state.value.planning.filter((item) => (item.task ?? "").trim() !== "").length
    : null
}));

const latestAnnonces = computed(() => {
  if (!state.value.annonces || state.value.annonces.length === 0) {
    return [
      {
        id: "null-annonce",
        title: null,
        category: null,
        type: null,
        status: null,
        date: null
      }
    ];
  }

  return state.value.annonces.slice(0, 3).map((item) => ({
    id: item.id,
    title: item.titre ?? null,
    category: item.categorie ?? null,
    type: item.type ?? null,
    status: item.statut ?? null,
    date: formatDate(item.date_creation)
  }));
});

const planningDays = computed(() => {
  const labels = ["Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"];
  const items = state.value.planning ?? [];

  return labels.map((label, index) => ({
    label,
    date: String(23 + index),
    event: items[index]?.task ?? null
  }));
});

const advice = computed(() => ({
  title: state.value.advice?.title ?? null,
  content: state.value.advice?.content ?? null
}));

const notification = computed(() => ({
  title: state.value.notification?.title ?? null,
  content: state.value.notification?.content ?? null
}));

function displayValue(value, suffix = "") {
  return value === null || value === undefined || value === "" ? "NULL" : `${value}${suffix}`;
}

function displayInitials(value) {
  if (!value) return "NU";
  return value
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("");
}

function formatDate(value) {
  if (!value) return null;
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return null;
  return new Intl.DateTimeFormat("fr-FR", {
    day: "2-digit",
    month: "short",
    year: "numeric"
  }).format(date);
}

onMounted(async () => {
  state.value = await fetchUserDashboard();
});
</script>
