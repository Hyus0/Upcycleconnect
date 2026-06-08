<template>
  <div class="page-container">
    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > MES FAVORIS</p>
        <h1 class="hero-title1">ANNONCES LIKÉES</h1>
        <p class="classic-text">
          Retrouvez ici tous les matériaux et objets que vous avez sauvegardés pour vos futurs projets.
        </p>
      </div>
      <button class="btn-secondary-back" @click="$router.push('/catalogue')">
        🠔 Retour au catalogue
      </button>
    </header>

    <div class="section-container">
      <div v-if="loading" class="state-card">
        Chargement de vos favoris...
      </div>
      
      <div v-else-if="favoris.length === 0" class="empty-state">
        <span class="empty-icon">🤍</span>
        <p>Aucun favori pour le moment. Parcourez le catalogue et cliquez sur le cœur pour sauvegarder des annonces ici.</p>
        <router-link to="/catalogue" class="btn-main-action" style="text-decoration: none">
          Découvrir le catalogue
        </router-link>
      </div>

      <div v-else class="annonces-grid">
        <article v-for="annonce in favoris" :key="annonce.id" class="annonce-card">
          <div class="annonce-card__image-wrapper" @click="goToAnnonce(annonce.id)">
            <img 
                :src="(annonce.image && annonce.image.trim() !== '') ? annonce.image : imageParDefaut" 
                alt="Image de l'annonce" 
                class="annonce-card__image" 
            />
            <div class="annonce-card__badges">
              <span :class="annonce.type === 'Vente' ? 'badge badge--orange' : 'badge badge--green'">
                {{ (annonce.type || "N/A").toUpperCase() }}
              </span>
            </div>
            <div class="favorite-indicator" title="Retirer des favoris" @click.stop="retirerFavori(annonce.id)">❤️</div>
          </div>
          
          <div class="annonce-card__content" @click="goToAnnonce(annonce.id)">
            <div class="annonce-card__header">
              <h3 class="annonce-card__title">{{ annonce.titre || "Sans titre" }}</h3>
              <p class="annonce-card__price">{{ formatPrice(annonce.prix, annonce.type) }}</p>
            </div>
                      
            <div class="annonce-card__meta">
              <span>📍 {{ annonce.ville || "N/A" }}</span>
            </div>
          </div>

          <div class="annonce-card__footer">
            <button class="btn-view" type="button" @click="goToAnnonce(annonce.id)">
              Voir
            </button>
            <button class="btn-remove" type="button" @click.stop="retirerFavori(annonce.id)">
              Retirer
            </button>
          </div>
        </article>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import imageParDefaut from "../../components/upcycling-concept.jpg";

const router = useRouter();
const API_URL = "http://localhost:8081";

const loading = ref(true);
const favoris = ref([]);

const currentUserId = computed(() => {
  const storedId = sessionStorage.getItem("id") || sessionStorage.getItem("userId");
  return Number(storedId) || 0;
});

const fetchFavoris = async () => {
  if (currentUserId.value === 0) return;
  
  loading.value = true;
  try {
    const res = await fetch(`${API_URL}/users/${currentUserId.value}/favoris`);
    if (res.ok) {
      favoris.value = await res.json();
    } else {
      console.error("Erreur HTTP:", res.status);
    }
  } catch (error) {
    console.error("Erreur lors de la récupération des favoris :", error);
  } finally {
    loading.value = false;
  }
};

const retirerFavori = async (annonceId) => {
  try {
    const res = await fetch(`${API_URL}/annonces/${annonceId}/favori/${currentUserId.value}`, {
      method: "POST"
    });
    
    if (res.ok) {
      favoris.value = favoris.value.filter(a => a.id !== annonceId);
    }
  } catch (error) {
    console.error("Erreur suppression favori", error);
  }
};

function formatPrice(value, type) {
  if (value === null || value === undefined || value === "") return "N/A";
  if ((type || "").toLowerCase() === "don" || Number(value) === 0) return "Gratuit";
  
  return new Intl.NumberFormat("fr-FR", {
    style: "currency", currency: "EUR"
  }).format(Number(value));
}

function goToAnnonce(id) {
  router.push(`/annonce/${id}`); 
}

onMounted(fetchFavoris);
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
}

.header-left {
  display: flex;
  flex-direction: column;
}

.sidebar__category2 {
  font-size: 0.65rem;
  color: #8fa396;
  letter-spacing: 1px;
  margin: 0 0 0.5rem 0;
  text-transform: uppercase;
}

.hero-title1 {
  font-size: 2rem;
  font-weight: 800;
  color: #1a1a1a;
  margin: 1.5rem 0 0.5rem;
}

.classic-text {
  color: #6d7b72;
  font-size: 0.95rem;
  margin: 0;
}

.btn-secondary-back {
  padding: 8px 16px;
  border-radius: 10px;
  border: 1px solid #ddd;
  background: white;
  cursor: pointer;
  font-weight: 500;
  transition: 0.2s;
}

.btn-secondary-back:hover {
  background: #f0f4f1;
}

.btn-main-action {
  display: inline-block;
  background-color: #2d7a4f;
  color: #ffffff;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  transition: 0.2s;
}

.btn-main-action:hover {
  background-color: #23653e;
}

.section-container {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e5ede7;
  padding: 20px;
}

.state-card {
  padding: 40px;
  text-align: center;
  color: #8fa396;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: #888;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}

.empty-icon {
  font-size: 3rem;
}

.annonces-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
}

.annonce-card {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e5ede7;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  cursor: pointer;
}

.annonce-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(44, 126, 79, 0.08);
}

.annonce-card__image-wrapper {
  position: relative;
  width: 100%;
  aspect-ratio: 4/3; 
  background-color: #f0f4f1;
  overflow: hidden;
}

.annonce-card__image {
  position: absolute; 
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover; 
  object-position: center; 
}

.annonce-card__badges {
  position: absolute;
  top: 10px;
  left: 10px;
}

.favorite-indicator {
  position: absolute;
  top: 10px;
  right: 10px;
  background: rgba(255, 255, 255, 0.9);
  padding: 6px;
  border-radius: 50%;
  font-size: 1rem;
  line-height: 1;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  transition: transform 0.2s;
}
.favorite-indicator:hover {
  transform: scale(1.1);
}

.badge {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 700;
}
.badge--orange { background: #fff4e6; color: #cc6600; }
.badge--green { background: #e9f5ed; color: #1e5636; }

.annonce-card__content {
  padding: 16px;
  flex: 1;
}

.annonce-card__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 6px;
}

.annonce-card__title {
  font-size: 1rem;
  font-weight: 700;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  color: #1a1a1a;
}

.annonce-card__price {
  font-size: 1rem;
  font-weight: 800;
  color: #2c7e4f;
  margin: 0;
}

.annonce-card__desc {
  font-size: 0.85rem;
  color: #6d7b72;
  margin: 0 0 16px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.annonce-card__meta {
  font-size: 0.8rem;
  color: #8fa396;
}

.annonce-card__footer {
  padding: 12px 16px;
  border-top: 1px solid #f0f4f1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.btn-view, .btn-remove {
  width: 100%;
  height: 32px;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: 0.2s;
  border: none;
}

.btn-view {
  background-color: #f0f4f1;
  color: #2c7e4f;
}
.btn-view:hover { background-color: #e1ede5; }

.btn-remove {
  background-color: #fdf2f2;
  color: #e74c3c;
}
.btn-remove:hover { background-color: #fceaea; }
</style>