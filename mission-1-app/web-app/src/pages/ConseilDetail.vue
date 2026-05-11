<template>
    <main class="page-container">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            :user-role="userRole"
            :user-score="userScore"
        />

        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ACCUEIL > ESPACE CONSEILS > {{ tip?.titre || "DÉTAIL" }}
                </p>
                <h1 class="hero-title1">
                    {{ tip?.titre || "Chargement..." }}
                </h1>
                <p v-if="tip" class="classic-text">
                    Astuce publiée le {{ formatDate(tip.date_creation) }}
                </p>
            </div>
            <div class="header-actions">
                <button class="btn-secondary" @click="$router.back()">
                    🠔 Retour
                </button>
            </div>
        </header>

        <div v-if="loading" class="loading-state">
            Récupération de l'astuce...
        </div>

        <div v-else-if="!tip" class="state-card error-card">
            <h2>Oups ! Conseil introuvable.</h2>
            <p>Ce conseil n'existe plus ou a été désactivé.</p>
        </div>

        <div v-else class="split-layout">
            <div class="left-column">
                <div
                    class="info-card"
                    style="
                        padding: 0;
                        overflow: hidden;
                        border: none;
                        box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
                    "
                >
                    <div v-if="tip.video_url" class="video-wrapper">
                        <iframe
                            :src="getEmbedUrl(tip.video_url)"
                            title="YouTube video player"
                            frameborder="0"
                            allow="
                                accelerometer;
                                autoplay;
                                clipboard-write;
                                encrypted-media;
                                gyroscope;
                                picture-in-picture;
                            "
                            allowfullscreen
                        >
                        </iframe>
                    </div>
                    <div v-else class="no-video-placeholder">
                        <span class="icon">💡</span>
                        <p>Aucune vidéo disponible pour cette astuce</p>
                    </div>
                </div>
            </div>

            <div class="right-column">
                <div class="info-card status-card">
                    <div
                        class="card-header-flex"
                        style="
                            margin-bottom: 0.5rem;
                            padding-bottom: 0.5rem;
                            border: none;
                        "
                    >
                        <h2
                            class="card-title"
                            style="margin: 0; font-size: 1.2rem"
                        >
                            À propos de ce conseil
                        </h2>
                        <span
                            :class="[
                                'type-badge',
                                getRoleClass(tip.role_cible),
                            ]"
                        >
                            POUR {{ tip.role_cible?.toUpperCase() }}
                        </span>
                    </div>
                </div>

                <div class="info-card description-card">
                    <div class="description-section">
                        <label class="info-label">Description & Étapes</label>
                        <div class="description-box article-text">
                            {{ tip.description }}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </main>
    <SiteFooter />
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

const route = useRoute();
const router = useRouter();
const API_URL = "http://localhost:8081";

const tip = ref(null);
const loading = ref(true);

const userScore = ref(parseInt(sessionStorage.getItem("userScore")) || 0);
const userRole = ref(sessionStorage.getItem("userRole") || "Particulier");

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const getRoleClass = (role) => {
    if (role === "Particulier") return "type-cours";
    if (role === "Prestataire") return "type-webinaire";
    if (role === "Salarie") return "type-atelier";
    return "";
};

const formatDate = (val) => {
    if (!val) return "";
    const date = new Date(val);
    return isNaN(date.getTime())
        ? ""
        : new Intl.DateTimeFormat("fr-FR", {
              day: "numeric",
              month: "long",
              year: "numeric",
          }).format(date);
};

const getEmbedUrl = (url) => {
    if (!url) return "";
    let videoId = "";
    if (url.includes("youtube.com/watch?v=")) {
        videoId = url.split("v=")[1].split("&")[0];
    } else if (url.includes("youtu.be/")) {
        videoId = url.split("youtu.be/")[1].split("?")[0];
    }
    return videoId ? `https://www.youtube.com/embed/${videoId}` : url;
};

const loadTip = async (id) => {
    loading.value = true;
    try {
        const res = await fetch(`${API_URL}/tips/${id}`);
        if (res.ok) {
            tip.value = await res.json();
        } else {
            tip.value = null;
        }
    } catch (e) {
        console.error("Erreur chargement du tip :", e);
        tip.value = null;
    } finally {
        loading.value = false;
    }
};

onMounted(() => {
    const tipId = route.params.id;
    if (tipId) {
        loadTip(tipId);
    } else {
        loading.value = false;
    }
});
</script>

<style scoped>
.page-container {
    min-height: 100vh;
    padding: 20px;
    background: #f7f9f7;
    max-width: 1600px;
    margin: 0 auto;
}

.content-header {
    padding-top: 1rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 2rem;
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
}

.sidebar__category2 {
    font-size: 0.75rem;
    color: #999;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.hero-title1 {
    font-size: 2.2rem;
    font-weight: 900;
    margin: 0.5rem 0 0;
    color: #1a1a1a;
}

.classic-text {
    color: #666;
    margin-top: 0.5rem;
}

.btn-secondary {
    background: white;
    color: #666;
    border: 1px solid #ddd;
    padding: 8px 16px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: bold;
    transition: background 0.2s;
}
.btn-secondary:hover {
    background: #f0f0f0;
}

.loading-state,
.state-card {
    text-align: center;
    padding: 3rem;
    color: #666;
    font-size: 1.1rem;
    background: white;
    border-radius: 12px;
    border: 1px dashed #ddd;
}
.error-card {
    border-color: #fca5a5;
    background: #fef2f2;
    color: #b91c1c;
}

.split-layout {
    display: grid;
    grid-template-columns: 1.5fr 1fr;
    gap: 2rem;
    margin-bottom: 3rem;
    align-items: start;
}

@media (max-width: 900px) {
    .split-layout {
        grid-template-columns: 1fr;
    }
}

.info-card {
    background: white;
    padding: 1.5rem;
    border-radius: 16px;
    border: 1px solid #eee;
    margin-bottom: 1.5rem;
}

.video-wrapper {
    position: relative;
    padding-bottom: 56.25%;
    height: 0;
    overflow: hidden;
    background: #000;
}
.video-wrapper iframe {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
}
.no-video-placeholder {
    padding: 6rem 2rem;
    text-align: center;
    background: #f0f0f0;
    color: #888;
    font-weight: bold;
}
.no-video-placeholder .icon {
    font-size: 3rem;
    display: block;
    margin-bottom: 1rem;
}

.card-header-flex {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.type-badge {
    padding: 6px 14px;
    border-radius: 20px;
    font-weight: 800;
    font-size: 0.75rem;
}
.type-cours {
    background: #eaf4ed;
    color: #2d7a4f;
}
.type-webinaire {
    background: #e0f2fe;
    color: #0369a1;
}
.type-atelier {
    background: #fef3c7;
    color: #92400e;
}

.description-card {
    min-height: 300px;
}
.info-label {
    display: block;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
    color: #999;
    margin-bottom: 1rem;
    letter-spacing: 0.5px;
}
.description-box {
    background: #f9f9f9;
    padding: 1.5rem;
    border-radius: 12px;
    color: #444;
    line-height: 1.8;
    font-size: 1.05rem;
    border: 1px solid #eee;
}
.article-text {
    white-space: pre-wrap;
}
</style>
