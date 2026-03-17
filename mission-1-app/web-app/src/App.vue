<template>
  <div class="page-shell">
    <header class="hero">
      <div class="hero-copy">
        <div class="eyebrow">Mission 1</div>
        <h1>Annonces disponibles</h1>
        <p>
          Consulte les objets proposes sur UpcycleConnect, filtre par type et trouve rapidement ce
          qui est reutilisable.
        </p>
      </div>

      <aside class="hero-panel">
        <span class="hero-chip">Catalogue public</span>
        <strong>{{ filteredAnnonces.length }} annonce{{ filteredAnnonces.length > 1 ? "s" : "" }}</strong>
        <small>{{ sourceLabel }}</small>
      </aside>
    </header>

    <section class="toolbar-card">
      <label class="search-field">
        <span>Recherche</span>
        <input v-model="filters.search" type="search" placeholder="Titre, ville, description..." />
      </label>

      <label class="select-field">
        <span>Type</span>
        <select v-model="filters.type">
          <option value="">Tous</option>
          <option value="don">Don</option>
          <option value="vente">Vente</option>
        </select>
      </label>

      <label class="select-field">
        <span>Statut</span>
        <select v-model="filters.status">
          <option value="">Tous</option>
          <option value="en ligne">En ligne</option>
          <option value="reserve">Reserve</option>
          <option value="vendu">Vendu</option>
        </select>
      </label>
    </section>

    <section class="content-grid">
      <aside class="summary-card">
        <div class="summary-head">
          <h2>Vue rapide</h2>
          <span class="source-badge" :class="{ fallback: source !== 'api' }">{{ sourceBadge }}</span>
        </div>

        <div class="summary-stats">
          <article>
            <span>Total</span>
            <strong>{{ filteredAnnonces.length }}</strong>
          </article>
          <article>
            <span>Dons</span>
            <strong>{{ donationsCount }}</strong>
          </article>
          <article>
            <span>Ventes</span>
            <strong>{{ salesCount }}</strong>
          </article>
        </div>

        <div class="summary-block">
          <h3>Conseil</h3>
          <p>Affiche d'abord un type d'annonce puis affine avec une ville ou un mot-clé.</p>
        </div>
      </aside>

      <main class="list-card">
        <div class="list-head">
          <div>
            <h2>Catalogue</h2>
            <p>
              {{ loading ? "Chargement..." : `${filteredAnnonces.length} resultat${filteredAnnonces.length > 1 ? "s" : ""}` }}
            </p>
          </div>
        </div>

        <div v-if="loading" class="state-card">Chargement des annonces...</div>
        <div v-else-if="error" class="state-card state-error">Impossible de charger les annonces.</div>
        <div v-else-if="filteredAnnonces.length === 0" class="state-card">
          Aucune annonce ne correspond a ces filtres.
        </div>

        <section v-else class="annonces-grid">
          <article v-for="annonce in filteredAnnonces" :key="annonce.id" class="annonce-card">
            <div class="card-top">
              <div>
                <h3>{{ annonce.titre }}</h3>
              </div>
              <span class="type-badge" :class="annonce.type">{{ annonce.type || "annonce" }}</span>
            </div>

            <div class="card-meta">
              <span>{{ annonce.ville || "Ville non renseignee" }}</span>
              <span>{{ annonce.etat_objet || "Etat non renseigne" }}</span>
              <span>{{ formatDate(annonce.date_creation) }}</span>
            </div>

            <p class="card-description">{{ annonce.description || "Aucune description disponible." }}</p>

            <div class="card-foot">
              <div>
                <strong class="price">{{ formatPrice(annonce.prix, annonce.type) }}</strong>
                <span>{{ annonce.statut || "en ligne" }}</span>
              </div>
              <span>{{ annonce.code_postal || "" }}</span>
            </div>
          </article>
        </section>
      </main>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { fetchAnnonces, fallbackAnnonces } from "./services/annoncesApi";

const loading = ref(true);
const error = ref(false);
const source = ref("api");
const annonces = ref([]);

const filters = reactive({
  search: "",
  type: "",
  status: ""
});

const filteredAnnonces = computed(() => {
  const search = filters.search.trim().toLowerCase();
  return annonces.value.filter((annonce) => {
    const matchesSearch =
      !search ||
      [
        annonce.titre,
        annonce.description,
        annonce.ville,
        annonce.adresse,
        annonce.etat_objet
      ]
        .join(" ")
        .toLowerCase()
        .includes(search);

    const matchesType = !filters.type || annonce.type === filters.type;
    const matchesStatus =
      !filters.status || (annonce.statut || "").toLowerCase() === filters.status.toLowerCase();

    return matchesSearch && matchesType && matchesStatus;
  });
});

const donationsCount = computed(() => filteredAnnonces.value.filter((item) => item.type === "don").length);
const salesCount = computed(() => filteredAnnonces.value.filter((item) => item.type === "vente").length);
const sourceBadge = computed(() => (source.value === "api" ? "API" : "Mode local"));
const sourceLabel = computed(() =>
  source.value === "api" ? "Mise a jour depuis l'API annonces" : "Affichage avec donnees locales"
);

function formatDate(value) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return "-";
  }
  return new Intl.DateTimeFormat("fr-FR", {
    day: "2-digit",
    month: "short",
    year: "numeric"
  }).format(date);
}

function formatPrice(value, type) {
  if (type === "don" || Number(value) === 0) {
    return "Gratuit";
  }
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR"
  }).format(Number(value));
}

onMounted(async () => {
  loading.value = true;
  error.value = false;
  try {
    annonces.value = await fetchAnnonces();
    source.value = "api";
  } catch {
    annonces.value = fallbackAnnonces;
    source.value = "fallback";
    error.value = false;
  } finally {
    loading.value = false;
  }
});
</script>
