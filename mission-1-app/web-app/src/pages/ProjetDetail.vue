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
                        ACCUEIL > COMMUNAUTÉ > PROJET
                    </p>
                    <h1 class="hero-title1">
                        {{ projet?.titre || "Chargement..." }}
                    </h1>
                    <div v-if="projet" class="project-sub-header">
                        <p class="classic-text">
                            Créé le {{ formatDateLong(projet.date_creation) }}
                            <span class="dot-separator">•</span>
                            {{ projet.nb_vues }} vues
                            <span class="dot-separator">•</span>
                            {{ projet.nb_likes }} likes
                        </p>
                    </div>
                </div>
                <button class="btn-secondary" @click="$router.back()">
                    🠔 Retour
                </button>
            </header>

            <div v-if="loading" class="loading-state">
                Récupération des données...
            </div>

            <div v-else-if="projet?.id" class="split-layout">
                <div class="left-column">
                    <div class="form-card">
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
                                    class="step-card-clean"
                                >
                                    <div class="step-header-clean">
                                        <div class="step-number-wrapper">
                                            <span class="step-number-badge">{{
                                                index + 1
                                            }}</span>
                                            <span class="step-title-label">{{
                                                etape.titre
                                            }}</span>
                                        </div>
                                    </div>
                                    <div class="step-body-clean">
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
                    <div class="form-card side-card">
                        <h2 class="card-title-side">
                            Organisateur de la session
                        </h2>
                        <div class="trainer-preview">
                            <div class="mini-avatar">
                                <div class="mini-avatar">
                                    <img 
                                        v-if="user_creator.image_profil" 
                                        :src="user_creator.image_profil" 
                                        class="mini-avatar-img" 
                                    />
                                    <span v-else>
                                        {{ user_creator.prenom?.charAt(0) }}{{ user_creator.nom?.charAt(0) }}
                                    </span>
                                </div>                       </div>
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

                    <div class="form-card side-card">
                        <h2 class="card-title-side">Impact Écologique</h2>
                        <div class="data-row">
                            <span class="data-label">CO2 Évité :</span>
                            <span class="text-success">{{ projet.co2_evite_kg }} kg</span>
                        </div>
                        <div class="data-row">
                            <span class="data-label">Score d'impact :</span>
                            <span class="status-badge">{{ projet.score_impact }} / 100</span>
                        </div>
                    </div>

                    <div class="form-actions-card flex-actions">
                        <button
                            @click="toggleLike"
                            class="btn-like-icon"
                            :class="{ 'is-liked': isLiked }"
                            title="Liker ce projet"
                        >
                            <Heart :fill="isLiked ? '#ef4444' : 'none'" :color="isLiked ? '#ef4444' : '#666'" :size="24" />
                        </button>

                        <button
                            v-if="!isOwner"
                            @click="contactCreateur"
                            class="btn-main-action flex-grow"
                            :class="{ 'btn-disabled': !projet.prix || projet.prix <= 0 || projet.statut !== 'Disponible' }"
                            :disabled="!projet.prix || projet.prix <= 0 || projet.statut !== 'Disponible'"
                        >
                            <span v-if="projet.prix > 0 && projet.statut === 'Disponible'">
                                Contacter pour acheter ({{ projet.prix }} €)
                            </span>
                            <span v-else-if="projet.prix > 0 && projet.statut !== 'Disponible'">
                                Objet déjà vendu
                            </span>
                            <span v-else>
                                Objet pas encore en vente
                            </span>
                        </button>
                        <button v-else class="btn-main-action flex-grow"
                        @click="goToModify(projet.id)">
                            <span>
                                Modifier mon projet
                            </span>
                        </button>
                    </div>

                </div>
            </div>
        </main>
    </div>
    <SiteFooter />
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";
import { Heart } from "lucide-vue-next";

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
        const res = await fetch(`/go/users/${id}`);
        if (res.ok) user_creator.value = await res.json();
    } catch (error) {
        console.error("Erreur créateur: ", error);
    } finally {
        loading.value = false;
    }
};

const fetchDetail = async () => {
    const id = route.params.id;
    try {
        const res = await fetch(`/go/projet/${id}`);
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
            `/go/projets/${route.params.id}/like-status/${currentUserId}`,
        );
        const data = await res.json();
        isLiked.value = data.liked;
    } catch (e) {
        console.error(e);
    }
};

const toggleLike = async () => {
    if (!isLoggedIn.value) {
        router.push("/connexion");
        return;
    }
    try {
        const res = await fetch(
            `/go/projets/${route.params.id}/like/${currentUserId}`,
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

const isOwner = computed(() =>
    projet.value?.id_createur === parseInt(sessionStorage.getItem("userId") || "0")
);

const goToModify = (id) => {
  router.push({ name: "modify-projet", params: { id: id } });
};

const contactCreateur = () => {
    if (!isLoggedIn.value) {
        router.push("/connexion");
        return;
    }
    router.push({
        path: "/messages",
        query: {
            user: projet.value.id_createur,
            projet: projet.value.id,
        },
    });
};


const viewProfile = (id) => {
    if (id) {
        router.push(`/user/${id}`);
    } else {
        console.error("ID manquant");
    }
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
}

.flex-actions {
    display: flex;
    flex-direction: row !important; 
    align-items: center;
    gap: 12px;
    margin-top: 0.5rem;
    width: 100%;
}

.btn-like-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 50px;
    height: 50px;
    border-radius: 12px;
    background-color: white;
    border: 2px solid #eee;
    cursor: pointer;
    transition: all 0.2s ease;
    flex-shrink: 0;
}

.btn-like-icon:hover {
    border-color: #ffcccc;
    background-color: #fff5f5;
}

.btn-like-icon.is-liked {
    border-color: #ef4444;
    background-color: #fef2f2;
}

.flex-grow {
    flex-grow: 1;
}

.btn-main-action {
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #2d7a4f;
    color: white;
    padding: 14px;
    border: none;
    border-radius: 12px;
    font-weight: 700;
    font-size: 1rem;
    cursor: pointer;
    transition: background 0.2s;
    height: 50px; 
}

.btn-main-action:hover:not(:disabled) {
    background-color: #246343;
}

.btn-disabled {
    background-color: #e5e7eb !important;
    color: #9ca3af !important;
    cursor: not-allowed !important;
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
    font-size: 0.8rem;
    color: #8fa396;
    letter-spacing: 1px;
    margin: 0 0 0.5rem 0;
    text-transform: uppercase;
}

.hero-title1 {
    font-size: 2rem;
    font-weight: 800;
    margin: 0 0 0.5rem 0;
    color: #1a1a1a;
}

.classic-text {
    color: #666;
    margin: 0;
    font-size: 0.95rem;
}

.dot-separator {
    margin: 0 4px;
    color: #ccc;
}

.btn-secondary {
    padding: 8px 16px;
    border-radius: 10px;
    border: 1px solid #ddd;
    background: white;
    cursor: pointer;
    font-weight: 500;
}

.split-layout {
    display: grid;
    grid-template-columns: 1.5fr 1fr;
    gap: 2rem;
    width: 100%;
}

.right-column {
    display: flex;
    flex-direction: column;
}

.form-card {
    background: #ffffff;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
    display: flex;
    flex-direction: column;
    gap: 2rem;
    margin-bottom: 2rem;
}

.card-header-flex {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #f0f0f0;
    padding-bottom: 1rem;
}

.card-title {
    font-size: 1.4rem;
    font-weight: 700;
    margin: 0;
    color: #1a1a1a;
}

.type-badge {
    background: #e9f5ed;
    color: #2d7a4f;
    padding: 6px 12px;
    border-radius: 15px;
    font-size: 0.75rem;
    font-weight: 800;
}

.description-section {
    display: flex;
    flex-direction: column;
    gap: 0.8rem;
}

.info-label {
    font-size: 0.85rem;
    color: #2d7a4f;
    text-transform: uppercase;
    font-weight: 800;
    letter-spacing: 0.5px;
}

.description-box {
    background: #fcfcfc;
    padding: 1.2rem;
    border-radius: 10px;
    font-size: 0.95rem;
    line-height: 1.6;
    border: 1px solid #eee;
    color: #333;
}

.steps-section {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.steps-timeline {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
}

.step-card-clean {
    background: #ffffff;
    border: 1px solid #eee;
    border-radius: 12px;
}

.step-header-clean {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 14px 16px;
    border-bottom: 1px dashed #eee;
    background: #fafafa;
    border-radius: 12px 12px 0 0;
}

.step-number-wrapper {
    display: flex;
    align-items: center;
    gap: 12px;
}

.step-number-badge {
    background: #2d7a4f;
    color: white;
    font-weight: 700;
    font-size: 0.85rem;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
}

.step-title-label {
    font-size: 1rem;
    font-weight: 700;
    color: #1a1a1a;
}

.step-body-clean {
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.step-desc {
    font-size: 0.95rem;
    color: #444;
    line-height: 1.5;
    margin: 0;
}

.step-img {
    width: 100%;
    border-radius: 10px;
    max-height: 350px;
    object-fit: cover;
}

.side-card {
    padding: 1.5rem 2rem;
    gap: 1rem;
}

.card-title-side {
    font-size: 1.1rem;
    font-weight: 700;
    margin: 0;
    color: #1a1a1a;
}

.trainer-preview {
    display: flex;
    gap: 12px;
    align-items: center;
}

.mini-avatar {
    width: 44px;
    height: 44px;
    background: #2d7a4f;
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 1rem;
    text-transform: uppercase;
    overflow: hidden;
}

.mini-avatar-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.trainer-name {
    font-weight: 700;
    margin: 0 0 2px 0;
    font-size: 0.95rem;
}

.link-btn {
    background: none;
    border: none;
    color: #2d7a4f;
    text-decoration: underline;
    cursor: pointer;
    font-size: 0.85rem;
    padding: 0;
}

.data-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.95rem;
}

.data-label {
    color: #666;
}

.text-success {
    color: #2d7a4f;
    font-weight: 700;
}

.status-badge {
    background: #f5f5f5;
    padding: 4px 8px;
    border-radius: 6px;
    font-weight: 700;
    font-size: 0.85rem;
}

.form-actions-card {
    margin-top: 0.2rem;
    padding: 0;
    display: flex;
    flex-direction: column;
}

.btn-save {
    background-color: #2d7a4f;
    color: white;
    padding: 1rem;
    border: none;
    border-radius: 12px;
    font-weight: 700;
    font-size: 1rem;
    cursor: pointer;
    transition: background 0.2s;
    width: 100%;
}

.btn-save:hover {
    background-color: #246343;
}

.btn-liked-active {
    background-color: #ef4444;
}
.btn-liked-active:hover {
    background-color: #dc2626;
}

.loading-state {
    text-align: center;
    padding: 3rem;
    color: #8fa396;
    font-style: italic;
}
</style>