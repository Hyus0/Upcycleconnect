<template>
  <div class="page-shell">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <header class="hero">
      <div class="hero-copy">
        <div class="eyebrow">Apprendre</div>
        <h1>Formations et ateliers</h1>
        <p>
          Retrouvez les ateliers, cours et webinaires proposes par la communaute UpcycleConnect.
        </p>
      </div>

      <aside class="hero-panel">
        <span class="hero-chip">Catalogue formations</span>
        <strong>{{ formations.length }} session{{ formations.length > 1 ? "s" : "" }}</strong>
        <small>Donnees issues de l'API Go</small>
      </aside>
    </header>

    <section class="content-grid content-grid--single">
      <main class="list-card">
        <div class="list-head">
          <div>
            <h2>Sessions disponibles</h2>
            <p>{{ loading ? "Chargement..." : `${formations.length} resultat${formations.length > 1 ? "s" : ""}` }}</p>
          </div>
        </div>

        <div v-if="loading" class="state-card">Chargement des formations...</div>
        <div v-else-if="error" class="state-card state-error">Impossible de charger les formations.</div>
        <div v-else-if="formations.length === 0" class="state-card">
          Aucune formation n'est encore disponible. Les valeurs resteront vides tant que la BDD ne contient rien.
        </div>

        <section v-else class="annonces-grid">
          <article v-for="formation in formations" :key="formation.id" class="annonce-card">
            <div class="card-top">
              <div>
                <h3>{{ displayValue(formation.titre) }}</h3>
              </div>
              <span class="type-badge">{{ displayValue(formation.type) }}</span>
            </div>

            <div class="card-meta">
              <span>{{ displayValue(formation.ville) }}</span>
              <span>{{ displayValue(formation.statut) }}</span>
              <span>{{ formatDate(formation.date_debut) }}</span>
            </div>

            <p class="card-description">{{ displayValue(formation.description) }}</p>

            <div class="card-foot">
              <div>
                <strong class="price">{{ formatPrice(formation.prix_unitaire) }}</strong>
                <span>{{ displayValue(formation.capacite_max) }} places</span>
              </div>
              <span>{{ displayValue(formation.code_postal) }}</span>
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
const formations = ref([]);
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

function formatPrice(value) {
  if (value === null || value === undefined || value === "") return "NULL";
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR"
  }).format(Number(value));
}

onMounted(async () => {
  loading.value = true;
  error.value = false;
  try {
    const response = await fetch("/go/formations");
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const payload = await response.json();
    formations.value = Array.isArray(payload) ? payload : [];
  } catch {
    formations.value = [];
    error.value = false;
  } finally {
    loading.value = false;
  }
});
</script>
