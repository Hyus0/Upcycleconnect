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
                <p class="sidebar__category2">ACCUEIL &gt; COMMUNAUTE</p>
                <h1 class="hero-title1">GALERIE DES PROJETS</h1>
                <p class="classic-text">Les projets affiches ici viennent du back et de l'API publique.</p>
            </div>

            <div class="header-actions">
                <div class="search-box">
                    <input v-model="searchQuery" type="text" placeholder="Rechercher un projet..." class="search-input" />
                </div>
            </div>
        </header>

        <div class="section-container">
            <div v-if="loading" class="loading-state">Chargement de la galerie...</div>
            <div v-else-if="filteredProjets.length === 0" class="empty-msg">Aucun projet ne correspond a votre recherche.</div>

            <div v-else class="annonces-grid">
                <article v-for="projet in filteredProjets" :key="projet.id" class="annonce-card formation-card">
                    <div class="card-body">
                        <div class="card-meta">
                            <span class="tag-type">PROJET</span>
                            <span class="date-text">{{ formatDate(projet.date_creation) }}</span>
                        </div>

                        <h3 class="item-title">{{ projet.titre }}</h3>
                        <p class="project-description-short">{{ projet.description_courte || "Description NULL" }}</p>

                        <div class="stats-row">
                            <span class="stat-item">{{ projet.nb_likes || 0 }} likes</span>
                            <span class="stat-item">{{ projet.nb_vues || 0 }} vues</span>
                            <span class="stat-item">{{ projet.co2_evite_kg || 0 }} kg CO2</span>
                        </div>
                    </div>

                    <div class="card-footer">
                        <button class="btn-main-action-full" @click="goToProjet(projet.id)">Decouvrir les etapes</button>
                    </div>
                </article>
            </div>
        </div>
    </main>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import { fetchProjets } from "../services/publicApi";

const projets = ref([]);
const loading = ref(true);
const searchQuery = ref("");
const userScore = ref(0);
const router = useRouter();

const isLoggedIn = computed(() => !!localStorage.getItem("userToken"));

const userName = computed(() => {
    const prenom = localStorage.getItem("userPrenom") || "";
    const nom = localStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const filteredProjets = computed(() =>
    projets.value.filter((item) => {
        const term = searchQuery.value.toLowerCase();
        return (
            !term ||
            (item.titre || "").toLowerCase().includes(term) ||
            (item.description_courte || "").toLowerCase().includes(term)
        );
    })
);

const formatDate = (dateStr) => {
    if (!dateStr) return "Recemment";
    return new Date(dateStr).toLocaleDateString("fr-FR", { day: "numeric", month: "short" });
};

const loadProjets = async () => {
    loading.value = true;
    try {
        const payload = await fetchProjets();
        projets.value = Array.isArray(payload) ? payload : [];
    } catch (error) {
        console.error("Erreur chargement projets :", error);
        projets.value = [];
    } finally {
        loading.value = false;
    }
};

const goToProjet = (id) => {
    router.push({ name: "projet-detail", params: { id } });
};

onMounted(loadProjets);
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
}

.card-body,
.card-footer {
    padding: 1.2rem;
}

.card-meta {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.8rem;
}

.tag-type {
    background: #eaf4ed;
    color: #2d6a4f;
    font-weight: 800;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.7rem;
}

.item-title {
    font-size: 1.15rem;
    font-weight: 700;
    margin-bottom: 0.8rem;
}

.project-description-short,
.stats-row {
    color: #666;
}
</style>
