<template>
  <div class="page-container">
    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > MES PROJETS LIKÉS</p>
        <h1 class="hero-title1">PROJETS LIKÉS</h1>
        <p class="classic-text">
          Retrouvez ici toutes les créations de la communauté que vous avez aimées et soutenues.
        </p>
      </div>
      <button class="btn-secondary-back" @click="$router.push('/projets')">
        🠔 Retour à la galerie
      </button>
    </header>

    <div class="section-container">
      <div v-if="loading" class="state-card">
        Chargement de vos projets favoris...
      </div>
      
      <div v-else-if="projetsLikes.length === 0" class="empty-state">
        <p>Aucun projet liké pour le moment. Parcourez la galerie et cliquez sur le cœur pour soutenir des créations.</p>
        <router-link to="/projets" class="btn-main-action" style="text-decoration: none">
          Découvrir les projets
        </router-link>
      </div>

      <div v-else class="annonces-grid">
        <article v-for="projet in projetsLikes" :key="projet.id" class="annonce-card">
          <div class="annonce-card__image-wrapper" @click="goToProjet(projet.id)">
            <img 
                :src="resolveImageUrl(projet.image_url || projet.ImageUrl)" 
                alt="Image du projet" 
                class="annonce-card__image" 
            />
            <div class="annonce-card__badges">
              <span class="badge badge--green">
                PROJET UPCYCLING
              </span>
            </div>
            <div class="favorite-indicator" title="Ne plus aimer" @click.stop="toggleLike(projet.id)">❤️</div>
          </div>
          
          <div class="annonce-card__content" @click="goToProjet(projet.id)">
            <div class="annonce-card__header">
              <h3 class="annonce-card__title">{{ projet.titre || "Sans titre" }}</h3>
            </div>
                      
            <div class="annonce-card__meta" style="display: flex; flex-direction: column; gap: 4px;">
              <span>🍃 {{ projet.co2_evite_kg || 0 }} kg CO2 évités</span>
              <span>👁️ {{ projet.nb_vues || 0 }} vues</span>
            </div>
          </div>

          <div class="annonce-card__footer">
            <button class="btn-view" type="button" @click="goToProjet(projet.id)">
              Voir
            </button>
            <button class="btn-remove" type="button" @click.stop="toggleLike(projet.id)">
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
const API_URL = "/go";

const loading = ref(true);
const projetsLikes = ref([]);

const currentUserId = computed(() => {
  const storedId = sessionStorage.getItem("id") || sessionStorage.getItem("userId");
  return Number(storedId) || 0;
});

const resolveImageUrl = (url) => {
    if (!url) return imageParDefaut;
    if (url.startsWith("http") || url.startsWith("data:") || url.startsWith("blob:")) {
        return url;
    }
    if (url.startsWith("uploads/")) {
        return `${API_URL}/img/${url.replace('uploads/', '')}`;
    }
    if (!url.includes("/")) {
        return `${API_URL}/img/projets/${url}`;
    }
    if (url.startsWith("/")) {
        return `${API_URL}${url}`;
    }
    return `${API_URL}/${url}`;
};

const fetchProjetsLikes = async () => {
  if (currentUserId.value === 0) return;

  loading.value = true;
  try {
    const res = await fetch(`${API_URL}/users/${currentUserId.value}/projets-likes`);
    if (res.ok) {
      const data = await res.json();
      projetsLikes.value = Array.isArray(data) ? data : [];
    } else {
      console.error("Erreur HTTP:", res.status);
    }
  } catch (error) {
    console.error("Erreur lors de la récupération des projets likés :", error);
  } finally {
    loading.value = false;
  }
};


const toggleLike = async (projetId) => {
  try {
    const res = await fetch(`${API_URL}/projets/${projetId}/like/${currentUserId.value}`, {
      method: "POST"
    });
    
    if (res.ok) {
      projetsLikes.value = projetsLikes.value.filter(p => p.id !== projetId);
    }
  } catch (error) {
    console.error("Erreur lors de la modification du like", error);
  }
};

function goToProjet(id) {
  router.push({ name: 'projet-detail', params: { id: id } }); 
}

onMounted(fetchProjetsLikes);
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