<template>
    <div class="layout-wrapper">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            user-role="Particulier"
            :user-score="userScore"
        />

        <main class="page-container">
            <header class="content-header">
                <div class="header-left">
                    <p class="sidebar__category2">ACCUEIL &gt; COMMUNAUTE &gt; {{ projet?.titre || "NULL" }}</p>
                    <h1 class="hero-title1">{{ projet?.titre || "Chargement..." }}</h1>
                    <div v-if="projet" class="project-sub-header">
                        <p class="creation-date">Cree le {{ formatDateLong(projet.date_creation) }}</p>
                        <div class="stats-bar-meta">
                            <span>{{ projet.nb_vues || 0 }} vues</span>
                            <span>{{ projet.nb_likes || 0 }} likes</span>
                        </div>
                    </div>
                </div>
                <div class="header-actions">
                    <button class="btn-secondary" @click="$router.back()">Retour</button>
                </div>
            </header>

            <div v-if="loading" class="loading-state">Recuperation des donnees...</div>

            <div v-else-if="projet?.id" class="split-layout">
                <div class="info-card">
                    <div class="card-header-flex">
                        <h2 class="card-title">Histoire de la transformation</h2>
                        <span class="type-badge">PROJET UPCYCLING</span>
                    </div>

                    <div class="description-box">{{ projet.description_courte || "Description NULL" }}</div>

                    <div class="steps-timeline">
                        <div v-for="(etape, index) in projet.etapes || []" :key="etape.id || index" class="step-item-card">
                            <div class="step-num-circle">{{ index + 1 }}</div>
                            <div class="step-info">
                                <h4 class="step-title">{{ etape.titre || `Etape ${index + 1}` }}</h4>
                                <p class="step-desc">{{ etape.description || "Description NULL" }}</p>
                                <img v-if="etape.image_url" :src="etape.image_url" class="step-img" alt="Illustration etape" />
                            </div>
                        </div>
                    </div>
                </div>

                <aside class="right-column">
                    <div class="info-card side-card">
                        <h3>Createur</h3>
                        <p>{{ userCreator.prenom || "NULL" }} {{ userCreator.nom || "" }}</p>
                    </div>

                    <div class="info-card side-card">
                        <h3>Impact ecologique</h3>
                        <div class="data-row"><span>CO2 evite :</span><strong>{{ projet.co2_evite_kg || 0 }} kg</strong></div>
                        <div class="data-row"><span>Score :</span><strong>{{ projet.score_impact || 0 }} / 100</strong></div>
                    </div>

                    <div class="action-card">
                        <button class="btn-main-action" :class="{ 'btn-liked-active': isLiked }" @click="toggleLike">
                            <span v-if="isLiked">Projet like</span>
                            <span v-else>Liker la creation</span>
                        </button>
                    </div>
                </aside>
            </div>
        </main>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import {
    fetchProjet,
    fetchProjetLikeStatus,
    fetchUser,
    toggleProjetLike
} from "../services/publicApi";

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const projet = ref(null);
const userCreator = ref({});
const isLiked = ref(false);
const userScore = ref(0);
const currentUserId = Number(localStorage.getItem("userId") || 0);

const isLoggedIn = computed(() => !!localStorage.getItem("userToken"));

const userName = computed(() => {
    const prenom = localStorage.getItem("userPrenom") || "";
    const nom = localStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const formatDateLong = (value) => {
    if (!value) return "Date inconnue";
    return new Date(value).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
        year: "numeric"
    });
};

const loadLikeStatus = async () => {
    if (!isLoggedIn.value || !currentUserId) return;
    try {
        const payload = await fetchProjetLikeStatus(route.params.id, currentUserId);
        isLiked.value = Boolean(payload?.liked);
    } catch (error) {
        console.error("Erreur statut like :", error);
    }
};

const loadProjet = async () => {
    loading.value = true;
    try {
        const payload = await fetchProjet(route.params.id);
        projet.value = payload;

        if (payload?.id_createur) {
            userCreator.value = (await fetchUser(payload.id_createur)) || {};
        }
        await loadLikeStatus();
    } catch (error) {
        console.error("Erreur detail projet :", error);
        projet.value = null;
    } finally {
        loading.value = false;
    }
};

const toggleLike = async () => {
    if (!isLoggedIn.value || !currentUserId) {
        alert("Connectez-vous pour liker.");
        router.push("/connexion");
        return;
    }

    try {
        await toggleProjetLike(route.params.id, currentUserId);
        isLiked.value = !isLiked.value;
        if (projet.value) {
            projet.value.nb_likes = (projet.value.nb_likes || 0) + (isLiked.value ? 1 : -1);
        }
    } catch (error) {
        console.error("Erreur toggle like :", error);
        alert(error.message || "Action impossible.");
    }
};

onMounted(loadProjet);
</script>

<style scoped>
.layout-wrapper {
    min-height: 100vh;
    padding: 20px;
    background: #f7f9f7;
    max-width: 1600px;
    margin: 0 auto;
}

.page-container {
    padding: 0 20px;
    margin-top: 2rem;
}

.content-header {
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 2rem;
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
}

.split-layout {
    display: grid;
    grid-template-columns: 1.6fr 1fr;
    gap: 1.5rem;
}

.info-card,
.action-card {
    background: white;
    padding: 1.2rem;
    border-radius: 12px;
    border: 1px solid #eee;
    margin-bottom: 1rem;
}

.description-box {
    background: #f9f9f9;
    padding: 1rem;
    border-radius: 8px;
    margin: 1rem 0;
}

.step-item-card {
    display: grid;
    grid-template-columns: 48px 1fr;
    gap: 1rem;
    margin-bottom: 1rem;
}

.step-num-circle {
    width: 48px;
    height: 48px;
    border-radius: 999px;
    display: grid;
    place-items: center;
    background: #eaf4ed;
    color: #2d6a4f;
    font-weight: 800;
}

.step-img {
    width: 100%;
    max-height: 220px;
    object-fit: cover;
    border-radius: 10px;
    margin-top: 0.75rem;
}
</style>
