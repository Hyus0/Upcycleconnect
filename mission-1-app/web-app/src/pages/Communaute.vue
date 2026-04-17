<template>
  <div class="page-shell">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <header class="hero">
      <div class="hero-copy">
        <div class="eyebrow">Communaute</div>
        <h1>Forums, evenements et conseils</h1>
        <p>
          Une page centrale pour suivre les evenements, les echanges et les ressources de la communaute.
        </p>
      </div>

      <aside class="hero-panel">
        <span class="hero-chip">Vie communautaire</span>
        <strong>{{ events.length }} evenement{{ events.length > 1 ? "s" : "" }}</strong>
        <small>Donnees issues de l'API Go</small>
      </aside>
    </header>

    <section class="content-grid">
      <aside class="summary-card">
        <div class="summary-head">
          <h2>Acces rapide</h2>
          <span class="source-badge">BDD</span>
        </div>

        <div class="summary-block">
          <h3>Forums</h3>
          <p>Les discussions communautaires seront branchees ici lorsque les donnees forum seront exploitees dans l'interface.</p>
        </div>

        <div class="summary-block">
          <h3>Conseils</h3>
          <p>Les conseils et guides d'upcycling pourront etre publies depuis le back-office.</p>
        </div>
      </aside>

      <main class="list-card">
        <div class="list-head">
          <div>
            <h2>Evenements communautaires</h2>
            <p>{{ loading ? "Chargement..." : `${events.length} resultat${events.length > 1 ? "s" : ""}` }}</p>
          </div>
        </div>

        <div v-if="loading" class="state-card">Chargement des evenements...</div>
        <div v-else-if="error" class="state-card state-error">Impossible de charger les evenements.</div>
        <div v-else-if="events.length === 0" class="state-card">
          Aucun evenement n'est encore renseigne dans la base.
        </div>

        <section v-else class="annonces-grid">
          <article v-for="event in events" :key="event.id" class="annonce-card">
            <div class="card-top">
              <div>
                <h3>{{ displayValue(event.titre) }}</h3>
              </div>
              <span class="type-badge">{{ displayValue(event.type) }}</span>
            </div>

            <div class="card-meta">
              <span>{{ displayValue(event.ville) }}</span>
              <span>{{ displayValue(event.code_postal) }}</span>
              <span>{{ formatDate(event.date_evenement) }}</span>
            </div>

            <p class="card-description">{{ displayValue(event.description) }}</p>

            <div class="card-foot">
              <div>
                <strong class="price">{{ displayValue(event.adresse) }}</strong>
                <span>{{ displayValue(event.ville) }}</span>
              </div>
            </div>
          </article>
        </section>
      </main>
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import SiteNavbar from "../components/SiteNavbar.vue";

const loading = ref(true);
const error = ref(false);
const events = ref([]);
const isLoggedIn = ref(Boolean(sessionStorage.getItem("userToken") || localStorage.getItem("userToken")));
const userName = ref("Marie Lambert");

function displayValue(value) {
  return value === null || value === undefined || value === "" ? "NULL" : value;
}

function formatDate(value) {
  if (!value) return "NULL";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "NULL";
  return new Intl.DateTimeFormat("fr-FR", {
    day: "2-digit",
    month: "short",
    year: "numeric"
  }).format(date);
}

onMounted(async () => {
  loading.value = true;
  error.value = false;
  try {
    const response = await fetch("/go/evenements");
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const payload = await response.json();
    events.value = Array.isArray(payload) ? payload : [];
  } catch {
    events.value = [];
    error.value = false;
  } finally {
    loading.value = false;
  }
});
</script>
