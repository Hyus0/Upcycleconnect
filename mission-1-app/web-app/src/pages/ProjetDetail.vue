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
                    <p class="sidebar__category2">
                        ACCUEIL > COMMUNAUTÉ > {{ projet?.titre }}
                    </p>
                    <h1 class="hero-title1">
                        {{ projet?.titre || "Chargement..." }}
                    </h1>
                    <div v-if="projet" class="project-sub-header">
                        <p class="creation-date">
                            Créé le {{ formatDateLong(projet.date_creation) }}
                        </p>
                        <div class="stats-bar-meta">
                            <span>{{ projet.nb_vues }} vues</span>
                            <span class="dot">•</span>
                            <span>{{ projet.nb_likes }} likes</span>
                        </div>
                    </div>
                </div>
                <div class="header-actions">
                    <button class="btn-secondary" @click="$router.back()">
                        🠔 Retour
                    </button>
                </div>
            </header>

            <div v-if="loading" class="loading-state">
                Récupération des données...
            </div>

            <div v-else-if="projet?.id" class="split-layout">
                <div class="left-column">
                    <div class="info-card">
                        <div class="card-header-flex">
                            <h2 class="card-title">
                                Histoire de la transformation
                            </h2>
                            <span class="type-badge">PROJET UPCYCLING</span>
                        </div>

                        <div class="description-section">
                            <label class="info-label"
                                >Description & Programme</label
                            >
                            <div class="description-box">
                                {{ projet.description_courte }}
                            </div>
                        </div>

                        <div class="steps-section">
                            <label class="info-label"
                                >Étapes de fabrication</label
                            >
                            <div class="steps-timeline">
                                <div
                                    v-for="(etape, index) in projet.etapes"
                                    :key="etape.id"
                                    class="step-item-card"
                                >
                                    <div class="step-num-circle">
                                        {{ index + 1 }}
                                    </div>
                                    <div class="step-info">
                                        <h4 class="step-title">
                                            {{ etape.titre }}
                                        </h4>
                                        <p class="step-desc">
                                            {{ etape.description }}
                                        </p>
                                        <img
                                            v-if="etape.image_url"
                                            :src="etape.image_url"
                                            class="step-img"
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="right-column">
                    <div class="info-card side-card">
                        <h3>Organisateur de la session</h3>
                        <div class="trainer-preview">
                            <div class="mini-avatar">
                                {{ user_creator.prenom?.charAt(0)
                                }}{{ user_creator.nom?.charAt(0) }}
                            </div>
                            <div>
                                <p class="trainer-name">
                                    {{ user_creator.prenom }}
                                    {{ user_creator.nom }}
                                </p>
                                <button
                                    class="link-btn"
                                    @click="viewProfile(user_creator.id)"
                                >
                                    Voir le profil expert
                                </button>
                            </div>
                        </div>
                    </div>

                    <div class="info-card side-card">
                        <h3>Impact Écologique</h3>
                        <div class="data-row">
                            <span class="data-label">CO2 Évité :</span>
                            <span class="text-success"
                                >{{ projet.co2_evite_kg }} kg</span
                            >
                        </div>
                        <div class="data-row">
                            <span class="data-label">Score d'impact :</span>
                            <span class="status-badge"
                                >{{ projet.score_impact }} / 100</span
                            >
                        </div>
                    </div>

                    <div class="action-card">
                        <button
                            @click="toggleLike"
                            class="btn-main-action"
                            :class="{ 'btn-liked-active': isLiked }"
                        >
                            <span v-if="isLiked">Projet liké</span>
                            <span v-else>Liker la création</span>
                        </button>
                    </div>
                </div>
            </div>
        </main>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const projet = ref(null);
const user_creator = ref([]);
const isLiked = ref(false);
const userScore = ref(0);

const currentUserId = sessionStorage.getItem("userId") || 0;
const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const formatDateLong = (d) => {
    if (!d) return "Date inconnue";
    return new Date(d).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
        year: "numeric",
    });
};

const fetchCreateur = async (id) => {
    try {
        const res = await fetch(`http://localhost:8081/users/${id}`, {
            method: "GET",
        });
        if (res.ok) user_creator.value = await res.json();
    } catch (error) {
        console.log("Erreur de récupération du créateur: ", error);
    } finally {
        loading.value = false;
    }
};

const fetchDetail = async () => {
    const id = route.params.id;
    try {
        const res = await fetch(`http://localhost:8081/projet/${id}`);
        if (res.ok) {
            projet.value = await res.json();
            if (projet.value.id_createur) {
                await fetchCreateur(projet.value.id_createur);
            }
            if (isLoggedIn.value) checkLikeStatus();
        }
    } catch (error) {
        console.error("Erreur projet detail :", error);
    } finally {
        loading.value = false;
    }
};

const checkLikeStatus = async () => {
    try {
        const res = await fetch(
            `http://localhost:8081/projets/${route.params.id}/like-status/${currentUserId}`,
        );
        const data = await res.json();
        isLiked.value = data.liked;
    } catch (e) {
        console.error(e);
    }
};

const toggleLike = async () => {
    if (!isLoggedIn.value) {
        alert("Connectez-vous pour liker !");
        router.push("/connexion");
        return;
    }
    try {
        const res = await fetch(
            `http://localhost:8081/projets/${route.params.id}/like/${currentUserId}`,
            { method: "POST" },
        );
        if (res.ok) {
            isLiked.value = !isLiked.value;
            projet.value.nb_likes += isLiked.value ? 1 : -1;
        }
    } catch (e) {
        console.error(e);
    }
};

const viewProfile = (id) => {
    router.push({ name: "public-profile", params: { id } });
};

onMounted(fetchDetail);
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

.project-sub-header {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.creation-date {
    font-size: 0.85rem;
    color: #888;
    margin: 0;
}

.stats-bar-meta {
    font-size: 0.85rem;
    font-weight: 700;
    color: #2d6a4f;
    display: flex;
    gap: 6px;
}

.split-layout {
    display: grid;
    grid-template-columns: 1.6fr 1fr;
    gap: 1.5rem;
}

.info-card {
    background: white;
    padding: 1.2rem;
    border-radius: 12px;
    border: 1px solid #eee;
    margin-bottom: 1rem;
}

.card-title {
    font-size: 1.1rem;
    font-weight: 800;
}

.description-box {
    background: #f9f9f9;
    padding: 1rem;
    border-radius: 8px;
    font-size: 0.95rem;
    line-height: 1.5;
    border: 1px solid #f0f0f0;
}

.steps-timeline {
    margin-top: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.step-item-card {
    display: flex;
    gap: 12px;
    padding: 1rem;
    background: #fff;
    border-radius: 10px;
    border: 1px solid #f0f0f0;
}

.step-num-circle {
    width: 24px;
    height: 24px;
    background: #2d6a4f;
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 800;
    font-size: 0.75rem;
    flex-shrink: 0;
}

.step-img {
    width: 100%;
    border-radius: 6px;
    max-height: 250px;
    object-fit: cover;
    margin-top: 10px;
}

.trainer-preview {
    display: flex;
    gap: 12px;
    align-items: center;
}

.mini-avatar {
    width: 44px;
    height: 44px;
    background: #2d6a4f;
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 1rem;
}

.trainer-name {
    font-weight: 800;
    margin: 0;
    font-size: 1rem;
}

.link-btn {
    background: none;
    border: none;
    color: #2d6a4f;
    text-decoration: underline;
    cursor: pointer;
    font-size: 0.8rem;
    padding: 0;
}

.data-row {
    display: flex;
    justify-content: space-between;
    font-size: 0.85rem;
    margin-bottom: 8px;
}

.btn-main-action {
    width: 100%;
    background: #2d6a4f;
    color: white;
    padding: 1rem;
    border-radius: 10px;
    font-weight: bold;
    cursor: pointer;
    border: none;
}

.btn-liked-active {
    background: #ff4d4d;
}

.sidebar__category2 {
    font-size: 0.7rem;
    color: #aaa;
    font-weight: 700;
}

.type-badge {
    background: #e8f5e9;
    color: #2d6a4f;
    padding: 4px 10px;
    border-radius: 15px;
    font-size: 0.7rem;
    font-weight: 800;
}

.info-label {
    font-size: 0.7rem;
    color: #999;
    text-transform: uppercase;
    font-weight: 700;
    margin-bottom: 5px;
    display: block;
}
</style>
