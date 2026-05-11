<template>
  <main class="public-dashboard">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <header class="content-header annonces-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > CATALOGUE</p>
        <h1 class="hero-title1">Annonces disponibles</h1>
        <p class="classic-text">
          Consulte les objets proposés sur UpcycleConnect, filtre par type et trouve rapidement ce qui est réutilisable.
        </p>
      </div>
      <RouterLink class="btn-main-action" to="/profil/annonces">+ Déposer une annonce</RouterLink>
    </header>

    <div class="stats-grid annonces-stats">
      <div class="card card--score">
        <p class="tag-score">CATALOGUE PUBLIC</p>
        <div class="score-value">{{ filteredAnnonces.length }} <span>annonces</span></div>
        <p class="score-level">{{ sourceLabel }}</p>
        <div class="score-footer">
          <div class="mini-stat">
            <strong>{{ donationsCount }}</strong><br />Dons
          </div>
          <div class="mini-stat">
            <strong>{{ salesCount }}</strong><br />Ventes
          </div>
          <div class="mini-stat">
            <strong>{{ filteredAnnonces.length }}</strong><br />Disponibles
          </div>
        </div>
      </div>
      <div class="card card--white">
        <div class="card-num">{{ donationsCount }}</div>
        <p class="text-dm">Objets gratuits</p>
        <span class="badge badge--green">DON</span>
      </div>
      <div class="card card--white">
        <div class="card-num2">{{ salesCount }}</div>
        <p class="text-dm">Objets en vente</p>
        <span class="badge badge--orange">VENTE</span>
      </div>
    </div>

    <section class="section-container">
      <div class="section-header">
        <div>
          <h2>Catalogue des annonces</h2>
          <p class="classic-text">{{ loading ? "Chargement..." : `${filteredAnnonces.length} résultat${filteredAnnonces.length > 1 ? "s" : ""}` }}</p>
        </div>
        <div class="header-actions">
          <input
            v-model="filters.search"
            type="search"
            placeholder="Rechercher..."
            class="search-input"
          />
          <select v-model="filters.type" class="btn-secondary">
            <option value="">Tous types</option>
            <option value="Don">Don</option>
            <option value="Vente">Vente</option>
          </select>
        </div>
      </div>

      <div v-if="loading" class="state-card">Chargement des annonces...</div>
      <div v-else-if="filteredAnnonces.length === 0" class="state-card">
        Aucune annonce ne correspond à ces filtres.
      </div>

      <div v-else class="annonces-grid">
        <article v-for="annonce in filteredAnnonces" :key="annonce.id" class="annonce-card">
          <div class="annonce-card__image-wrapper">
            <img 
              :src="annonce.imageUrl || imageParDefaut" 
              alt="Image de l'annonce" 
              class="annonce-card__image" 
            />
            <div class="annonce-card__badges">
              <span :class="annonce.type === 'Vente' ? 'badge badge--orange' : 'badge badge--green'">
                {{ displayValue(annonce.type).toUpperCase() }}
              </span>
            </div>
          </div>
          
          <div class="annonce-card__content">
            <div class="annonce-card__header">
              <h3 class="annonce-card__title">{{ displayValue(annonce.titre) }}</h3>
              <p class="annonce-card__price">{{ formatPrice(annonce.prix, annonce.type) }}</p>
            </div>
            
            <p class="annonce-card__desc">{{ displayValue(annonce.description) }}</p>
            
            <div class="annonce-card__meta">
              <span>📍 {{ displayValue(annonce.ville) }} ({{ displayValue(annonce.code_postal) }})</span>
              <span>📅 {{ formatDate(annonce.date_creation) }}</span>
            </div>
          </div>

          <div class="annonce-card__footer">
              <button class="btn-view btn-small" type="button" @click="goToAnnonce(annonce.id)">Voir</button>
              <RouterLink class="btn-main-action btn-small" to="/profil/annonces">Contacter</RouterLink>
          </div>
        </article>
      </div>
    </section>
  </main>
  <SiteFooter />

</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";
import imageParDefaut from "../components/upcycling-concept.jpg";

const router = useRouter();
const loading = ref(true);
const source = ref("api");
const annonces = ref([]);
const API_URL = "http://localhost:8081";

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));

const userName = computed(() => {
  const prenom = sessionStorage.getItem("userPrenom") || "";
  const nom = sessionStorage.getItem("userNom") || "";
  return (prenom || nom) ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const filters = reactive({
  search: "",
  type: ""
});

const filteredAnnonces = computed(() => {
  const search = filters.search.trim().toLowerCase();
  
  return annonces.value.filter((annonce) => {
    const matchesSearch =
      !search ||
      [annonce.titre, annonce.description, annonce.ville].join(" ").toLowerCase().includes(search);

    const matchesType = !filters.type || (annonce.type || "").toLowerCase() === filters.type.toLowerCase();

    return matchesSearch && matchesType;
  });
});

const donationsCount = computed(() => filteredAnnonces.value.filter((item) => (item.type || "").toLowerCase() === "don").length);
const salesCount = computed(() => filteredAnnonces.value.filter((item) => (item.type || "").toLowerCase() === "vente").length);

const sourceLabel = computed(() =>
  source.value === "api" ? "Données issues de l'API annonces" : "Aucune donnée disponible"
);

function displayValue(value) {
  return value === null || value === undefined || value === "" ? "N/A" : value;
}

function formatDate(value) {
  if (!value) return "N/A";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "N/A";
  
  return new Intl.DateTimeFormat("fr-FR", {
    day: "2-digit", month: "short", year: "numeric"
  }).format(date);
}

function goToAnnonce(id) {
  router.push(`/annonce/${id}`); 
}

function formatPrice(value, type) {
  if (value === null || value === undefined || value === "") return "N/A";
  if ((type || "").toLowerCase() === "don" || Number(value) === 0) return "Gratuit";
  
  return new Intl.NumberFormat("fr-FR", {
    style: "currency", currency: "EUR"
  }).format(Number(value));
}

onMounted(async () => {
  loading.value = true;
  try {
    const res = await fetch(`${API_URL}/annonces`);
    if (res.ok) {
      annonces.value = await res.json() || [];
      source.value = "api";
    } else {
      throw new Error("Erreur HTTP");
    }
  } catch (error) {
    console.error("Erreur lors de la récupération des annonces :", error);
    annonces.value = [];
    source.value = "empty";
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.public-dashboard {
  min-height: 100vh;
  padding: 20px;
  background: var(--bg-light, #f7f9f7);
}

.annonces-header .btn-main-action {
  display: inline-flex;
  align-items: center;
  text-decoration: none;
}

.annonces-stats {
  grid-template-columns: 1.4fr 0.8fr 0.8fr;
}

.state-card {
  border: 1px dashed #cfe0d4;
  border-radius: 14px;
  padding: 26px;
  color: var(--text-grey);
  background: #fbfdfb;
}

.annonces-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 20px;
  margin-top: 24px;
}

.annonce-card {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #e5ede7;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  text-align: left;
}

.annonce-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 24px rgba(44, 126, 79, 0.08);
  border-color: #9bcbae;
}

.annonce-card__image-wrapper {
  position: relative;
  width: 100%;
  aspect-ratio: 4/3;
  background-color: #f0f4f1;
}

.annonce-card__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.annonce-card__badges {
  position: absolute;
  top: 12px;
  left: 12px;
  display: flex;
  gap: 8px;
}

.annonce-card__content {
  padding: 16px;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.annonce-card__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
  margin-bottom: 8px;
}

.annonce-card__title {
  font-family: "Syne", sans-serif;
  font-size: 1.1rem;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.annonce-card__price {
  font-size: 1.1rem;
  font-weight: 800;
  color: #2c7e4f;
  margin: 0;
  white-space: nowrap;
}

.annonce-card__desc {
  font-size: 0.85rem;
  color: #6d7b72;
  margin: 0 0 16px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.annonce-card__meta {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 0.75rem;
  color: #8fa396;
}

/* =========================================
   FOOTER DE LA CARTE : BOUTONS STRICTEMENT ÉGAUX
   ========================================= */
.annonce-card__footer {
  padding: 14px 16px;
  border-top: 1px solid #f0f4f1;
  display: grid; /* Force 2 colonnes identiques */
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

/* On force TOUS les boutons dans le footer à avoir la même structure */
.annonce-card__footer > button,
.annonce-card__footer > a {
  box-sizing: border-box !important;
  width: 100% !important;
  height: 32px !important;
  min-height: 32px !important;
  max-height: 32px !important;
  margin: 0 !important;
  padding: 0 8px !important;
  
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  
  border-radius: 8px !important;
  font-size: 0.85rem !important;
  font-family: "Syne", sans-serif !important;
  font-weight: 600 !important;
  cursor: pointer !important;
  text-decoration: none !important;
  transition: all 0.2s !important;
  line-height: 1 !important;
  
  /* Anti-saut de ligne */
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  
  border: 1px solid transparent !important;
}

/* Bouton VOIR */
.btn-view {
  background-color: #f0f4f1 !important;
  color: #2c7e4f !important;
  border-color: #f0f4f1 !important;
}

.btn-view:hover {
  background-color: #e1ede5 !important;
  border-color: #e1ede5 !important;
}

/* Bouton CONTACTER (btn-main-action) */
.annonce-card__footer .btn-main-action {
  background-color: #2c7e4f !important;
  color: #ffffff !important;
  border-color: #2c7e4f !important;
}

.annonce-card__footer .btn-main-action:hover {
  background-color: #23653e !important;
  border-color: #23653e !important;
}

/* =========================================
   MEDIA QUERIES
   ========================================= */
@media (max-width: 1500px) {
  .annonces-grid { grid-template-columns: repeat(4, 1fr); }
}

@media (max-width: 1200px) {
  .annonces-grid { grid-template-columns: repeat(3, 1fr); }
}

@media (max-width: 920px) {
  .annonces-stats { grid-template-columns: 1fr; }
  .annonces-grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 600px) {
  .annonces-grid { grid-template-columns: 1fr; }
}
</style>