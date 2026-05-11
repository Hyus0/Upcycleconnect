<template>
  <main class="public-dashboard">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">
          ACCUEIL > CATALOGUE > {{ annonce?.titre ? annonce.titre.toUpperCase() : "..." }}
        </p>
      </div>
      <div class="header-actions">
        <button class="btn-secondary-back" @click="$router.push('/catalogue')">
          🠔 Retour au catalogue
        </button>
      </div>
    </header>

    <div v-if="loading" class="loading-state">Récupération de l'annonce...</div>

    <div v-else-if="annonce?.id">
      
      <div class="split-layout">
        
        <div class="left-column">
          <div class="image-container-stretch">
            <img 
              :src="annonce.imageUrl || imageParDefaut" 
              alt="Image de l'objet" 
              class="main-image-fit"
            />
            <span :class="annonce.type === 'Don' ? 'tag-don' : 'tag-vente'" class="type-badge-image">
              {{ annonce.type === "Don" ? "DON" : "VENTE" }}
            </span>
            
            <button 
              class="image-like-btn" 
              :class="{ 'is-liked': isFavorited }" 
              @click.prevent="toggleFavori"
            >
              <span class="heart-icon">{{ isFavorited ? '❤️' : '🤍' }}</span>
              <span v-if="likesCount > 0" class="likes-count">{{ likesCount }}</span>
            </button>
          </div>
        </div>

        <div class="right-column">
          <div class="info-card-compact sticky-card">
            
            <div class="product-header">
              <h1 class="product-title">{{ annonce.titre }}</h1>
            </div>

            <hr class="divider" />
            
            <ul class="product-details-list">
              <li>
                <span class="detail-label">État</span>
                <span class="detail-value"><strong>{{ annonce.etat_objet || 'Non renseigné' }}</strong></span>
              </li>
              <li>
                <span class="detail-label">Date de poste</span>
                <span class="detail-value">{{ formatDate(annonce.date_creation) }}</span>
              </li>
              <li>
                <span class="detail-label">Catégorie</span>
                <span class="detail-value">{{ categoryName || "Non renseignée" }}</span>
              </li>
              <li>
                <span class="detail-label">Localisation</span>
                <span class="detail-value">{{ annonce.ville }}</span>
              </li>
            </ul>

            <hr class="divider" />
                        
            <div class="product-vendor clickable-vendor" @click="goToProfilVendeur">
                <div class="vendor-avatar">👤</div>
                <div class="vendor-info">
                <span class="vendor-name">{{ vendeur.prenom || 'Utilisateur #' + annonce.id_vendeur }}</span>
                <span class="vendor-role">Membre certifié</span>
                </div>
            </div>

            <div class="purchase-area">
              <div class="product-price">{{ formatPrice(annonce.prix, annonce.type) }}</div>
              
              <div class="action-buttons">
                <button class="btn-action-primary" @click="acheter">
                  {{ annonce.type === "Don" ? "Réserver" : "Ajouter au Panier" }}
                </button>
                <button class="btn-contact" @click="contacterVendeur">
                  Contacter
                </button>
              </div>
            </div>

          </div>
        </div>
      </div>

      <div class="description-section-full description-box">
        <h3 class="section-title">Description du produit</h3>
        {{ annonce.description || "Aucune description fournie par le vendeur." }}
      </div>

      <section v-if="autresAnnonces.length > 0" class="other-annonces-section">
        <h2 class="section-title-large">Autres annonces de l'utilisateur</h2>
        
        <div class="annonces-grid">
          <article v-for="ann in autresAnnonces" :key="ann.id" class="annonce-card">
            <div class="annonce-card__image-wrapper" @click="goToAnnonce(ann.id)">
              <img 
                :src="ann.imageUrl || imageParDefaut" 
                alt="Image de l'annonce" 
                class="annonce-card__image" 
              />
              <div class="annonce-card__badges">
                <span :class="ann.type === 'Vente' ? 'badge badge--orange' : 'badge badge--green'">
                  {{ ann.type.toUpperCase() }}
                </span>
              </div>
            </div>
            
            <div class="annonce-card__content" @click="goToAnnonce(ann.id)">
              <div class="annonce-card__header">
                <h3 class="annonce-card__title">{{ ann.titre }}</h3>
                <p class="annonce-card__price">{{ formatPrice(ann.prix, ann.type) }}</p>
              </div>
              <p class="annonce-card__desc">{{ ann.description }}</p>
            </div>
          </article>
        </div>

        <div class="center-btn">
          <button class="btn-secondary btn-more" @click="voirToutVendeur">
            Voir toutes ses annonces
          </button>
        </div>
      </section>
    </div>
  </main>
  <SiteFooter />
</template>

<script setup>
import { ref, onMounted, computed, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

import imageParDefaut from "../components/upcycling-concept.jpg";

const route = useRoute();
const router = useRouter();
const API_URL = "http://localhost:8081";

const loading = ref(true);
const annonce = ref(null);
const autresAnnonces = ref([]);
const categoryName = ref("");
const vendeur = ref({}); 

const isFavorited = ref(false);
const likesCount = ref(0);

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const currentUserId = computed(() => sessionStorage.getItem("userId") || 0);

const userName = computed(() => {
  const prenom = sessionStorage.getItem("userPrenom") || "";
  const nom = sessionStorage.getItem("userNom") || "";
  return (prenom || nom) ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const formatDate = (d) => {
  if (!d) return "N/A";
  return new Intl.DateTimeFormat("fr-FR", { day: "2-digit", month: "long", year: "numeric" }).format(new Date(d));
};

const formatPrice = (value, type) => {
  if (value === null || value === undefined) return "N/A";
  if ((type || "").toLowerCase() === "don" || Number(value) === 0) return "Gratuit";
  return new Intl.NumberFormat("fr-FR", { style: "currency", currency: "EUR" }).format(Number(value));
};

const fetchCategory = async (catId) => {
  try {
    const res = await fetch(`${API_URL}/category/${catId}`);
    if (res.ok) {
      const catData = await res.json();
      categoryName.value = catData.nom;
    }
  } catch (e) { console.error(e); }
};

const fetchVendeur = async (vendeurId) => {
  try {
    const res = await fetch(`${API_URL}/users/${vendeurId}`);
    if (res.ok) { vendeur.value = await res.json(); }
  } catch (e) { console.error(e); }
};

const fetchAutresAnnonces = async (vendeurId) => {
  try {
    const res = await fetch(`${API_URL}/users/${vendeurId}/annonces`);
    if (res.ok) {
      const allAnnonces = await res.json() || [];
      const filtered = allAnnonces.filter(a => a.id !== annonce.value.id && a.statut === "Disponible" && a.est_valide === "Valide");
      autresAnnonces.value = filtered.slice(0, 4);
    }
  } catch (e) { console.error(e); }
};

const fetchFavoriStatus = async () => {
  try {
    const res = await fetch(`${API_URL}/annonces/${route.params.id}/favori/${currentUserId.value}`);
    if (res.ok) {
      const data = await res.json();
      likesCount.value = data.total || 0;
      isFavorited.value = data.is_favorited || false;
    }
  } catch (e) { console.error("Erreur récupération favori", e); }
};

const goToProfilVendeur = () => {
  if (annonce.value && annonce.value.id_vendeur) {
    router.push(`/user/${annonce.value.id_vendeur}`); 
  }
};
const toggleFavori = async () => {
  if (!isLoggedIn.value) {
    router.push("/connexion");
    return;
  }

  isFavorited.value = !isFavorited.value;
  likesCount.value += isFavorited.value ? 1 : -1;

  try {
    const res = await fetch(`${API_URL}/annonces/${route.params.id}/favori/${currentUserId.value}`, {
      method: "POST"
    });
    
    if (!res.ok) {
      throw new Error("Erreur serveur");
    }
  } catch (e) {
    isFavorited.value = !isFavorited.value;
    likesCount.value += isFavorited.value ? 1 : -1;
    console.error("Erreur sauvegarde favori", e);
  }
};

const fetchAnnoncePrincipale = async () => {
  loading.value = true;
  try {
    const res = await fetch(`${API_URL}/annonces/${route.params.id}`);
    if (res.ok) {
      annonce.value = await res.json();
      if (annonce.value.id_categorie) fetchCategory(annonce.value.id_categorie);
      if (annonce.value.id_vendeur) {
        fetchVendeur(annonce.value.id_vendeur);
        fetchAutresAnnonces(annonce.value.id_vendeur);
      }
      fetchFavoriStatus(); 
    }
  } catch (e) { console.error(e); } 
  finally { loading.value = false; }
};

const goToAnnonce = (id) => { router.push(`/annonce/${id}`); };

watch(() => route.params.id, (newId) => {
  if (newId) fetchAnnoncePrincipale();
});

const voirToutVendeur = () => { alert(`Toutes les annonces de ${vendeur.value.prenom || 'ce membre'}`); };
const acheter = () => alert("Redirection achat/réservation");
const contacterVendeur = () => alert("Ouverture messagerie");

onMounted(fetchAnnoncePrincipale);
</script>

<style scoped>
.public-dashboard {
  min-height: 100vh;
  padding: 20px 40px 60px 40px;
  background: var(--bg-light, #f7f9f7);
  font-family: "Syne", sans-serif;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
  margin-bottom: 24px;
}

.sidebar__category2 { font-size: 0.8rem; color: #8fa396; letter-spacing: 1px; }

.btn-secondary-back {
    padding: 8px 16px;
    border-radius: 10px;
    border: 1px solid #ddd;
    background: white;
    cursor: pointer;
    font-weight: 400;
}

.btn-secondary-back:hover { background: #f0f4f1; }

.split-layout {
  display: grid;
  grid-template-columns: 1.85fr 1fr; 
  gap: 30px;
  align-items: stretch; 
  margin-bottom: 30px;
}

.left-column {
  height: 100%;
}

.image-container-stretch {
  position: relative;
  width: 100%;
  height: 100%; 
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  background: #fff;
  border: 1px solid #eaeaea;
}

.main-image-fit {
  width: 100%;
  height: 100%;
  object-fit: cover; 
}

.type-badge-image {
  position: absolute;
  top: 12px; left: 12px;
  padding: 6px 12px;
  border-radius: 20px;
  font-weight: 800; font-size: 0.8rem;
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
  z-index: 1;
}

.tag-don { background: #e9f5ed; color: #1e5636; }
.tag-vente { background: #fff4e6; color: #cc6600; }

.image-like-btn {
  position: absolute;
  bottom: 16px;
  right: 16px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(4px);
  border: 1px solid rgba(0,0,0,0.05);
  border-radius: 30px;
  padding: 8px 14px;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  z-index: 10;
}

.image-like-btn:hover {
  transform: scale(1.05);
  background: #ffffff;
}

.heart-icon {
  font-size: 1.2rem;
  line-height: 1;
}

.likes-count {
  font-weight: 700;
  font-size: 0.95rem;
  color: #1a1a1a;
}

.description-section-full {
  width: 100%; 
  margin-bottom: 60px; 
}

.section-title { font-size: 1.05rem; color: #1a1a1a; margin-bottom: 10px; font-weight: 700; }
.description-box {
  background: #ffffff; padding: 20px; border-radius: 12px;
  color: #555; line-height: 1.6; font-size: 0.95rem; border: 1px solid #e5ede7;
}

.right-column {
  height: 100%;
}

.sticky-card {
  position: sticky;
  top: 30px;
  background: #ffffff;
  padding: 24px; 
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
  border: 1px solid #eaeaea;
  height: 100%;
  display: flex;
  flex-direction: column; 
}

.product-header { margin-bottom: 10px; }

.product-title {
  font-size: 1.8rem;
  font-weight: 800;
  color: #1a1a1a;
  margin: 0;
  line-height: 1.2;
}

.divider {
  border: 0;
  border-top: 1px solid #f0f0f0;
  margin: 16px 0;
}

.product-details-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.product-details-list li {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  font-size: 0.95rem;
  border-bottom: 1px dashed #f0f0f0;
}

.product-details-list li:last-child { border-bottom: none; }

.detail-label { color: #8fa396; }
.detail-value { color: #2b302d; text-align: right; }

.product-vendor {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
  margin-bottom: auto; 
  padding-bottom: 15px;
  cursor: pointer;
}

.vendor-avatar {
  width: 40px; height: 40px;
  background: #f0f4f1; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 1.1rem;
}

.vendor-name {
  display: block; font-weight: 700; color: #1a1a1a; font-size: 0.95rem;
}

.vendor-role {
  display: block; font-size: 0.8rem; color: #8fa396;
}

.purchase-area {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 20px;
  background: #fcfcfc;
  padding: 16px;
  border-radius: 12px;
  border: 1px solid #f0f4f1;
}

.product-price {
  font-size: 1.5rem;
  font-weight: 900;
  margin-bottom: 20px;
  color: #2c7e4f;
}

.action-buttons {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
    width: 100%;
}

.btn-action-primary, 
.btn-contact {
    box-sizing: border-box;
    width: 100%;
    height: 32px !important;
    min-height: 32px !important; 
    max-height: 32px !important;
    margin: 0;
    padding: 0 14px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 8px;
    font-size: 13px;
    font-family: inherit;
    font-weight: 500;
    cursor: pointer;
    text-decoration: none;
    transition: all 0.2s;
    line-height: 1;
}

.btn-action-primary {
    background-color: #2d7a4f;
    color: #ffffff;
    border: 1px solid #2d7a4f;
}

.btn-action-primary:hover {
    background-color: #23653e;
    border-color: #23653e;
}

.btn-contact {
    background-color: #ffffff;
    color: #2d7a4f;
    border: 1px solid #2d7a4f;
}

.btn-contact:hover {
    background-color: #f0f4f1;
}

.other-annonces-section {
  padding-top: 30px;
  border-top: 1px solid #e5ede7;
}

.section-title-large {
  font-size: 1.4rem; font-weight: 700; margin-bottom: 20px; color: #1a1a1a;
}

.annonces-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.annonce-card {
  background: #ffffff; border-radius: 12px;
  border: 1px solid #eaeaea; overflow: hidden;
  display: flex; flex-direction: column;
  transition: transform 0.2s; cursor: pointer;
}
.annonce-card:hover { transform: translateY(-4px); box-shadow: 0 8px 16px rgba(0,0,0,0.06); }

.annonce-card__image-wrapper { position: relative; width: 100%; aspect-ratio: 4/3; background-color: #f0f4f1; }
.annonce-card__image { width: 100%; height: 100%; object-fit: cover; }
.annonce-card__badges { position: absolute; top: 8px; left: 8px; }
.badge { padding: 4px 8px; border-radius: 12px; font-size: 0.7rem; font-weight: 700; }
.badge--green { background: #e9f5ed; color: #1e5636; }
.badge--orange { background: #fff4e6; color: #cc6600; }

.annonce-card__content { padding: 12px; flex: 1; }
.annonce-card__header { display: flex; justify-content: space-between; align-items: flex-start; gap: 8px; margin-bottom: 6px; }
.annonce-card__title { font-size: 0.95rem; font-weight: 600; margin: 0; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.annonce-card__price { font-size: 1rem; font-weight: 700; color: #2c7e4f; margin: 0; }
.annonce-card__desc { font-size: 0.8rem; color: #8fa396; margin: 0; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }

.center-btn { display: flex; justify-content: center; }

.btn-secondary {
    padding: 8px 16px;
    border-radius: 20px;
    border: 1px solid #ddd;
    background: white;
    cursor: pointer;
    font-weight: 600;
}

.btn-more { padding: 10px 20px; }

@media (max-width: 1024px) {
  .split-layout { grid-template-columns: 1fr; }
  .sticky-card { position: relative; top: 0; height: auto; }
  .image-container-stretch { height: auto; aspect-ratio: 16/9; }
  .annonces-grid { grid-template-columns: repeat(2, 1fr); }
  .other-annonces-section { padding-top: 20px; }
}
@media (max-width: 600px) {
  .annonces-grid { grid-template-columns: 1fr; }
  .action-buttons { grid-template-columns: 1fr; } 
}
</style>