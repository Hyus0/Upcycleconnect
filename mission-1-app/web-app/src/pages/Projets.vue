<template>
    <main class="page-main-content">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            user-role="Particulier"
            :user-score="userScore"
        />
        
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > COMMUNAUTÉ</p>
                <h1 class="hero-title1">GALERIE DES PROJETS</h1>
                <p class="classic-text">
                    Laissez-vous inspirer par les créations de la communauté et découvrez les secrets de leur fabrication.
                </p>
            </div>

            <div class="header-actions">
                <div class="search-section">
                    <div class="search-box">
                        <input 
                            v-model="searchQuery" 
                            type="text" 
                            placeholder="Rechercher un projet, un matériau..."
                            class="search-input"
                        />
                    </div>
                </div>
            </div>
        </header>

        <div class="section-container">
            <div v-if="loading" class="loading-state">
                Chargement de la galerie...
            </div>

            <div v-else-if="filteredProjets.length === 0" class="empty-msg">
                Aucune création ne correspond à votre recherche.
            </div>

            <div v-else class="annonces-grid">
                <div
                    v-for="projet in filteredProjets"
                    :key="projet.id"
                    class="annonce-card formation-card"
                >
                    <div class="card-image">
                        <img 
                            src='../components/upcycling-concept.jpg' alt='Image de recyclage'
                        >
                        <div class="impact-overlay">
                            🍃 {{ projet.co2_evite_kg }}kg CO2
                        </div>
                    </div>

                    <div class="card-body">
                        <div class="card-meta">
                            <span class="tag-type">PROJET</span>
                            <span class="date-text">{{ formatDate(projet.date_creation) }}</span>
                        </div>

                        <h3 class="item-title">{{ projet.titre }}</h3>
                        
                        <p class="project-description-short">
                            {{ projet.description_courte }}
                        </p>

                        <div class="stats-row">
                            <span class="stat-item">❤️ {{ projet.nb_likes }}</span>
                            <span class="stat-item">👁️ {{ projet.nb_vues }}</span>
                        </div>
                    </div>

                    <div class="card-footer">
                        <button
                            class="btn-main-action-full"
                            @click="goToProjet(projet.id)"
                        >
                            Découvrir les étapes
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </main>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";

const projets = ref([]);
const loading = ref(true);
const searchQuery = ref(""); 
const userScore = ref(0); 
const router = useRouter();

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));

const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const filteredProjets = computed(() => {
    return projets.value.filter(p => {
        const term = searchQuery.value.toLowerCase();
        return !term || 
               p.titre?.toLowerCase().includes(term) || 
               p.description_courte?.toLowerCase().includes(term);
    });
});

const formatDate = (dateStr) => {
    if (!dateStr) return "Récemment";
    const date = new Date(dateStr);
    return date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' });
};

const fetchProjets = async () => {
    loading.value = true;
    try {
        const res = await fetch(`http://localhost:8081/projets`);
        if (res.ok) projets.value = await res.json();
    } catch (error) {
        console.error("Erreur:", error);
    } finally {
        loading.value = false;
    }
}

const goToProjet = (id) => {
    router.push({ name: 'projet-detail', params: { id: id } });
};

onMounted(fetchProjets);
</script>

<style scoped>
.page-main-content {
    min-height: 100vh;
    padding: 20px;
    background: #f7f9f7;
    max-width: 1600px; 
    margin: 0 auto;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #eee;
}

.card-image {
    height: 180px;
    width: 100%;
    position: relative;
    overflow: hidden;
}

.card-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.5s ease;
}

.formation-card:hover .card-image img {
    transform: scale(1.05);
}

.impact-overlay {
    position: absolute;
    top: 12px;
    right: 12px;
    background: rgba(255, 255, 255, 0.9);
    padding: 4px 10px;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 800;
    color: #2d6a4f;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.annonces-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
}

.formation-card {
    border: 1px solid #f0f0f0;
    border-radius: 16px;
    background: white;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.formation-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
}

.card-body {
    padding: 1.2rem;
    flex-grow: 1;
}

.card-meta {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.8rem;
}

.item-title {
    font-size: 1.15rem;
    font-weight: 700;
    margin-bottom: 0.8rem;
    color: #222;
}

.project-description-short {
    font-size: 0.85rem;
    color: #666;
    margin-bottom: 1rem;
    line-height: 1.5;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.stats-row {
    display: flex;
    gap: 15px;
    font-size: 0.8rem;
    color: #888;
    margin-bottom: 1.2rem;
}

.card-footer {
    padding: 1.2rem;
    padding-top: 0;
}

.btn-main-action-full {
    width: 100%;
    background: #2d6a4f;
    color: white;
    border: none;
    padding: 12px;
    border-radius: 10px;
    font-weight: 700;
    cursor: pointer;
}

.search-input {
    padding: 12px 20px;
    border-radius: 25px;
    border: 1px solid #ddd;
    width: 320px;
    outline: none;
}
</style>