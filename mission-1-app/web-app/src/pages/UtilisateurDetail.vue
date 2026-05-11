<template>
    <main class="public-dashboard">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            variant="public"
        />

        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ACCUEIL > MEMBRES >
                    {{ user?.prenom ? user.prenom.toUpperCase() : "PROFIL" }}
                </p>
            </div>
            <div class="header-actions">
                <button class="btn-secondary-back" @click="$router.back()">
                    🠔 Retour
                </button>
            </div>
        </header>

        <div v-if="loading" class="loading-state">Chargement du profil...</div>
        <div v-else-if="!user" class="loading-state">
            Utilisateur introuvable.
        </div>

        <div v-else class="profile-container">
            <section class="profile-header-card">
                <div class="profile-banner"></div>

                <div class="profile-info-wrapper">
                    <div class="profile-top-row">
                        <img
                            :src="user.photo_url || defaultAvatar"
                            alt="Avatar utilisateur"
                            class="profile-thumbnail"
                        />

                        <div class="action-buttons-wrapper">
                            <button
                                v-if="currentUserId !== user.id"
                                class="btn-follow-joint"
                                :class="{ 'is-following': isFollowing }"
                                @click="toggleFollow"
                            >
                                {{ isFollowing ? "✓ Suivi(e)" : "+ Suivre" }}
                            </button>
                        </div>
                    </div>

                    <div class="profile-identity">
                        <h1 class="user-name">
                            {{ user.prenom }} {{ user.nom }}
                        </h1>
                        <span class="badge badge--green"
                            >Membre UpcycleConnect - Particulier</span
                        >
                    </div>

                    <div class="profile-meta-details">
                        <div class="meta-item rating-row">
                            <div class="stars-wrapper">
                                <span
                                    v-for="n in 5"
                                    :key="n"
                                    class="star-container"
                                >
                                    <span class="star-bg">★</span>
                                    <span
                                        class="star-fill"
                                        :style="{
                                            width: getStarFillPercentage(
                                                n,
                                                averageNote,
                                            ),
                                        }"
                                        >★</span
                                    >
                                </span>
                            </div>
                            <span class="reviews-count">
                                {{ displayNote }} / 5 •
                                {{ avisList.length }} avis client(s)
                            </span>
                        </div>

                        <div class="meta-item location-row">
                            📍 {{ user.ville || "Ville non renseignée" }}
                            <template v-if="user.code_postal"
                                >({{ user.code_postal }})</template
                            >
                        </div>
                        <div class="meta-item location-row">
                            👥 {{ followersCount }} Abonné(s),
                            {{ followingCount }} Suivi(e)s
                        </div>
                    </div>
                </div>

                <nav class="profile-tabs">
                    <button
                        class="tab-btn"
                        :class="{ active: activeTab === 'annonces' }"
                        @click="activeTab = 'annonces'"
                    >
                        Annonces ({{ annonces.length }})
                    </button>
                    <button
                        class="tab-btn"
                        :class="{ active: activeTab === 'evaluations' }"
                        @click="activeTab = 'evaluations'"
                    >
                        Évaluations ({{ avisList.length }})
                    </button>
                </nav>
            </section>

            <section class="tab-content-area">
                <div v-if="activeTab === 'annonces'">
                    <div v-if="annonces.length === 0" class="state-card-light">
                        Cet utilisateur n'a aucune annonce disponible pour le
                        moment.
                    </div>

                    <div v-else class="annonces-grid">
                        <article
                            v-for="ann in annonces"
                            :key="ann.id"
                            class="annonce-card"
                        >
                            <div
                                class="annonce-card__image-wrapper"
                                @click="goToAnnonce(ann.id)"
                            >
                                <img
                                    :src="ann.imageUrl || imageParDefaut"
                                    alt="Image annonce"
                                    class="annonce-card__image"
                                />
                                <div class="annonce-card__badges">
                                    <span
                                        :class="
                                            ann.type === 'Vente'
                                                ? 'badge badge--orange'
                                                : 'badge badge--green'
                                        "
                                    >
                                        {{ (ann.type || "N/A").toUpperCase() }}
                                    </span>
                                </div>
                            </div>

                            <div
                                class="annonce-card__content"
                                @click="goToAnnonce(ann.id)"
                            >
                                <div class="annonce-card__header">
                                    <h3 class="annonce-card__title">
                                        {{ ann.titre }}
                                    </h3>
                                    <p class="annonce-card__price">
                                        {{ formatPrice(ann.prix, ann.type) }}
                                    </p>
                                </div>
                                <p class="annonce-card__desc">
                                    {{ ann.description }}
                                </p>
                            </div>

                            <div class="annonce-card__footer">
                                <button
                                    class="btn-view"
                                    @click="goToAnnonce(ann.id)"
                                >
                                    Voir l'annonce
                                </button>
                            </div>
                        </article>
                    </div>
                </div>

                <div v-if="activeTab === 'evaluations'">
                    <div
                        v-if="
                            isLoggedIn &&
                            currentUserId !== user?.id &&
                            !hasAlreadyCommented
                        "
                        class="review-form-card"
                    >
                        <h3>Laisser un avis sur {{ user?.prenom }}</h3>

                        <div class="rating-selector">
                            <span
                                v-for="star in 5"
                                :key="star"
                                class="star-click"
                                :class="{ 'star-active': star <= newAvis.note }"
                                @click="newAvis.note = star"
                                >★</span
                            >
                        </div>

                        <textarea
                            v-model="newAvis.commentaire"
                            placeholder="Comment s'est passée votre transaction avec ce membre ?"
                            rows="3"
                            class="review-input"
                        ></textarea>

                        <div class="review-submit-row">
                            <button
                                class="btn-submit-review"
                                @click="submitAvis"
                                :disabled="isSubmittingAvis"
                            >
                                {{
                                    isSubmittingAvis
                                        ? "Envoi..."
                                        : "Publier mon avis"
                                }}
                            </button>
                        </div>
                    </div>

                    <div
                        v-else-if="
                            isLoggedIn &&
                            currentUserId !== user?.id &&
                            hasAlreadyCommented
                        "
                        class="already-commented-msg"
                    >
                        ✓ Vous avez déjà laissé une évaluation à ce membre.
                    </div>

                    <div v-if="avisList.length === 0" class="state-card-light">
                        Aucun avis pour le moment. Soyez le premier à évaluer ce
                        membre !
                    </div>

                    <div v-else class="reviews-list">
                        <div
                            v-for="avis in avisList"
                            :key="avis.id"
                            class="review-item"
                        >
                            <div class="review-header">
                                <div
                                    class="review-author"
                                    @click="goToUser(avis.id_auteur)"
                                    title="Voir le profil"
                                >
                                    <div class="review-avatar">
                                        {{
                                            avis.prenom_auteur
                                                ? avis.prenom_auteur
                                                      .charAt(0)
                                                      .toUpperCase()
                                                : "U"
                                        }}
                                    </div>
                                    <div>
                                        <strong class="author-name">{{
                                            avis.prenom_auteur || "Utilisateur"
                                        }}</strong>
                                        <span class="review-date">
                                            •
                                            {{
                                                formatDate(avis.date_creation)
                                            }}</span
                                        >
                                    </div>
                                </div>

                                <div class="stars-wrapper">
                                    <span
                                        v-for="n in 5"
                                        :key="n"
                                        class="star-container"
                                    >
                                        <span class="star-bg">★</span>
                                        <span
                                            class="star-fill"
                                            :style="{
                                                width: getStarFillPercentage(
                                                    n,
                                                    avis.note,
                                                ),
                                            }"
                                            >★</span
                                        >
                                    </span>
                                </div>
                            </div>
                            <p class="review-text">{{ avis.commentaire }}</p>
                        </div>
                    </div>
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
const user = ref(null);
const stats = ref({});
const annonces = ref([]);
const avisList = ref([]);
const activeTab = ref("annonces");

const isFollowing = ref(false);
const followersCount = ref(0);
const followingCount = ref(0);

const newAvis = ref({ note: 5, commentaire: "" });
const isSubmittingAvis = ref(false);

const defaultAvatar = "https://cdn-icons-png.flaticon.com/512/149/149071.png";

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const currentUserId = computed(() => {
    const storedId =
        sessionStorage.getItem("id") || sessionStorage.getItem("userId");
    return Number(storedId) || 0;
});
const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const averageNote = computed(() => {
    if (avisList.value.length === 0) return 0;
    const sum = avisList.value.reduce((total, avis) => total + avis.note, 0);
    return sum / avisList.value.length;
});

const displayNote = computed(() => {
    return averageNote.value > 0 ? averageNote.value.toFixed(1) : "N/A";
});

const getStarFillPercentage = (n, rating) => {
    if (rating >= n) return "100%";
    if (rating > n - 1) return `${(rating - (n - 1)) * 100}%`;
    return "0%";
};

// Vérifier si l'utilisateur a déjà commenté
const hasAlreadyCommented = computed(() => {
    if (!avisList.value || !currentUserId.value) return false;
    return avisList.value.some(
        (avis) => avis.id_auteur === currentUserId.value,
    );
});

// Formatage
const formatDate = (value) => {
    if (!value) return "N/A";
    const d = new Date(value);
    if (isNaN(d)) return "N/A";
    return new Intl.DateTimeFormat("fr-FR", {
        day: "2-digit",
        month: "long",
        year: "numeric",
    }).format(d);
};

const formatPrice = (value, type) => {
    if (value === null || value === undefined || value === "") return "N/A";
    if ((type || "").toLowerCase() === "don" || Number(value) === 0)
        return "Gratuit";
    return new Intl.NumberFormat("fr-FR", {
        style: "currency",
        currency: "EUR",
    }).format(Number(value));
};

const goToAnnonce = (id) => router.push(`/annonce/${id}`);

const goToUser = (id) => {
    if (id) {
        router.push(`/user/${id}`);
        activeTab.value = "annonces";
    }
};

// Actions API
const toggleFollow = async () => {
    if (!isLoggedIn.value)
        return alert("Veuillez vous connecter pour suivre ce membre.");
    isFollowing.value = !isFollowing.value;
    followersCount.value += isFollowing.value ? 1 : -1;

    try {
        const res = await fetch(
            `${API_URL}/users/${route.params.id}/follow/${currentUserId.value}`,
            { method: "POST" },
        );
        if (!res.ok) throw new Error("Erreur serveur");
    } catch (e) {
        isFollowing.value = !isFollowing.value;
        followersCount.value += isFollowing.value ? 1 : -1;
        console.error(e);
    }
};

const submitAvis = async () => {
    if (!newAvis.value.commentaire.trim())
        return alert("Veuillez écrire un commentaire.");

    isSubmittingAvis.value = true;
    try {
        const token = sessionStorage.getItem("userToken");
        const res = await fetch(`${API_URL}/users/${route.params.id}/avis`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
                id_auteur: currentUserId.value,
                note: newAvis.value.note,
                commentaire: newAvis.value.commentaire,
            }),
        });

        if (res.ok) {
            newAvis.value.commentaire = "";
            newAvis.value.note = 5;
            await refreshData();
        } else {
            const err = await res.text();
            alert("Erreur serveur : " + err);
        }
    } catch (e) {
        console.error(e);
    } finally {
        isSubmittingAvis.value = false;
    }
};

const refreshData = async () => {
    const userId = route.params.id;
    try {
        const [resStats, resAnn, resAvis, resFollow] = await Promise.all([
            fetch(`${API_URL}/users/${userId}/stats`).then((r) =>
                r.ok ? r.json() : {},
            ),
            fetch(`${API_URL}/users/${userId}/annonces`).then((r) =>
                r.ok ? r.json() : [],
            ),
            fetch(`${API_URL}/users/${userId}/avis`).then((r) =>
                r.ok ? r.json() : [],
            ),
            fetch(
                `${API_URL}/users/${userId}/follow/${currentUserId.value}`,
            ).then((r) => (r.ok ? r.json() : {})),
        ]);

        stats.value = resStats;
        annonces.value = (resAnn || []).filter(
            (a) => a.statut === "Disponible" && a.est_valide === "Valide",
        );
        avisList.value = resAvis || [];
        followersCount.value = resFollow.followers || 0;
        followingCount.value = resFollow.following || 0;
        isFollowing.value = resFollow.is_following || false;
    } catch (e) {
        console.error("Erreur refresh:", e);
    }
};

const loadProfile = async () => {
    loading.value = true;
    try {
        const resUser = await fetch(`${API_URL}/users/${route.params.id}`);
        if (resUser.ok) {
            user.value = await resUser.json();
            await refreshData();
        } else {
            user.value = null;
        }
    } catch (e) {
        console.error("Erreur initiale:", e);
    } finally {
        loading.value = false;
    }
};

onMounted(loadProfile);
watch(() => route.params.id, loadProfile);
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

.sidebar__category2 {
    font-size: 0.8rem;
    color: #8fa396;
    letter-spacing: 1px;
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

.profile-container {
    display: flex;
    flex-direction: column;
    gap: 36px;
    max-width: 1400px;
    width: 100%;
    margin: 0 auto;
}

.profile-header-card {
    background: #fff;
    border-radius: 16px;
    overflow: hidden;
    border: 1px solid #e5ede7;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
}

.profile-banner {
    height: 180px;
    background: linear-gradient(135deg, #2d7a4f, #9bcbae);
}

.profile-info-wrapper {
    padding: 0 40px 20px 40px;
    position: relative;
}

.profile-top-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-top: -65px;
    margin-bottom: 16px;
}

.profile-thumbnail {
    width: 130px;
    height: 130px;
    background: #ffffff;
    border: 4px solid #fff;
    border-radius: 50%;
    object-fit: cover;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.action-buttons-wrapper {
    margin-top: 80px;
}

.profile-identity {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
    margin-bottom: 16px;
}

.user-name {
    font-size: 2rem;
    font-weight: 800;
    color: #1a1a1a;
    margin: 0;
    line-height: 1.1;
}
.badge {
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: 700;
}
.badge--green {
    background: #e9f5ed;
    color: #1e5636;
}

.btn-follow-joint {
    padding: 10px 24px;
    background-color: #2d7a4f;
    color: #ffffff;
    border: 1px solid #2d7a4f;
    border-radius: 10px;
    font-size: 0.95rem;
    font-family: inherit;
    font-weight: 700;
    cursor: pointer;
    transition: all 0.2s;
}
.btn-follow-joint:hover {
    background-color: #23653e;
    border-color: #23653e;
}
.btn-follow-joint.is-following {
    background-color: #ffffff;
    color: #2d7a4f;
    border: 1px solid #2d7a4f;
}
.btn-follow-joint.is-following:hover {
    background-color: #f0f4f1;
}

.profile-meta-details {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 16px;
}
.meta-item {
    font-size: 0.95rem;
    color: #6d7b72;
    display: flex;
    align-items: center;
    gap: 8px;
}
.rating-row {
    display: flex;
    align-items: center;
}
.reviews-count {
    font-weight: 600;
    color: #1a1a1a;
    margin-left: 8px;
}

.stars-wrapper {
    display: flex;
    gap: 2px;
}
.star-container {
    position: relative;
    display: inline-block;
    font-size: 1.2rem;
}
.star-bg {
    color: #e5ede7;
}
.star-fill {
    position: absolute;
    top: 0;
    left: 0;
    overflow: hidden;
    color: #f39c12;
    white-space: nowrap;
}

.profile-tabs {
    display: flex;
    justify-content: flex-start;
    gap: 40px;
    border-top: 1px solid #f0f4f1;
    padding: 0 40px;
}

.tab-btn {
    background: none;
    border: none;
    padding: 20px 0;
    font-size: 1.05rem;
    font-family: inherit;
    font-weight: 700;
    color: #8fa396;
    cursor: pointer;
    position: relative;
    transition: color 0.2s;
}

.tab-btn:hover {
    color: #2d7a4f;
}
.tab-btn.active {
    color: #2d7a4f;
}
.tab-btn.active::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 3px;
    background: #2d7a4f;
}

.tab-content-area {
    margin-top: -10px;
}

.state-card-light {
    padding: 40px;
    text-align: center;
    color: #8fa396;
    font-style: italic;
    background: #fff;
    border-radius: 12px;
    border: 1px solid #e5ede7;
}

.annonces-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
}
.annonce-card {
    background: #ffffff;
    border-radius: 12px;
    border: 1px solid #e5ede7;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    transition: 0.2s ease;
}
.annonce-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 16px rgba(44, 126, 79, 0.06);
}
.annonce-card__image-wrapper {
    position: relative;
    width: 100%;
    aspect-ratio: 4/3;
    background-color: #f0f4f1;
    cursor: pointer;
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
}
.badge--orange {
    background: #fff4e6;
    color: #cc6600;
}
.annonce-card__content {
    padding: 16px;
    flex: 1;
    cursor: pointer;
}
.annonce-card__header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 8px;
    margin-bottom: 6px;
}
.annonce-card__title {
    font-size: 0.95rem;
    font-weight: 700;
    margin: 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
.annonce-card__price {
    font-size: 1rem;
    font-weight: 800;
    color: #2c7e4f;
    margin: 0;
}
.annonce-card__desc {
    font-size: 0.8rem;
    color: #6d7b72;
    margin: 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
.annonce-card__footer {
    padding: 14px 16px;
    border-top: 1px solid #f0f4f1;
}
.btn-view {
    box-sizing: border-box;
    width: 100%;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 8px;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    transition: 0.2s;
    background-color: #f0f4f1;
    color: #2c7e4f;
    border: none;
}
.btn-view:hover {
    background-color: #e1ede5;
}

.review-form-card {
    background: #fff;
    padding: 24px;
    border-radius: 12px;
    border: 1px solid #e5ede7;
    margin-bottom: 24px;
}
.review-form-card h3 {
    margin: 0 0 16px 0;
    font-size: 1.1rem;
    color: #1a1a1a;
}
.rating-selector {
    font-size: 2rem;
    margin-bottom: 16px;
    user-select: none;
}
.star-click {
    color: #e5ede7;
    cursor: pointer;
    transition: color 0.2s;
    margin-right: 4px;
}
.star-click:hover,
.star-active {
    color: #f39c12;
}
.review-input {
    width: 100%;
    box-sizing: border-box;
    padding: 12px;
    border: 1px solid #e5ede7;
    border-radius: 8px;
    font-family: inherit;
    font-size: 0.9rem;
    resize: vertical;
    background: #fbfdfb;
}
.review-input:focus {
    outline: none;
    border-color: #9bcbae;
}
.review-submit-row {
    display: flex;
    justify-content: flex-end;
    margin-top: 12px;
}
.btn-submit-review {
    padding: 10px 24px;
    background-color: #2d7a4f;
    color: #ffffff;
    border: none;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: 0.2s;
}
.btn-submit-review:hover:not(:disabled) {
    background-color: #23653e;
}
.btn-submit-review:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.already-commented-msg {
    padding: 16px;
    background-color: #e9f5ed;
    color: #1e5636;
    border: 1px solid #cfe0d4;
    border-radius: 8px;
    margin-bottom: 24px;
    font-weight: 600;
}

.reviews-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.review-item {
    background: #fff;
    padding: 20px;
    border-radius: 12px;
    border: 1px solid #e5ede7;
}
.review-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
}
.review-author {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    transition: opacity 0.2s;
}
.review-author:hover {
    opacity: 0.7;
}
.author-name {
    color: #1a1a1a;
    text-decoration: underline transparent;
    transition: 0.2s;
}
.review-author:hover .author-name {
    text-decoration-color: #1a1a1a;
}

.review-avatar {
    width: 36px;
    height: 36px;
    background: #f0f4f1;
    color: #2c7e4f;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    font-size: 0.9rem;
}
.review-date {
    color: #8fa396;
    font-size: 0.8rem;
}
.review-text {
    margin: 0;
    color: #555;
    font-size: 0.95rem;
    line-height: 1.5;
}

@media (max-width: 1024px) {
    .annonces-grid {
        grid-template-columns: repeat(3, 1fr);
    }
}
@media (max-width: 768px) {
    .profile-top-row {
        flex-direction: column;
        align-items: center;
    }
    .profile-thumbnail {
        margin-top: -65px;
        margin-bottom: 16px;
    }
    .action-buttons-wrapper {
        margin-top: 10px;
    }
    .profile-tabs {
        padding: 0 20px;
        gap: 20px;
        justify-content: space-around;
    }
    .annonces-grid {
        grid-template-columns: repeat(2, 1fr);
    }
}
@media (max-width: 500px) {
    .annonces-grid {
        grid-template-columns: 1fr;
    }
}
</style>
