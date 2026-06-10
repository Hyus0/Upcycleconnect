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
                        ACCUEIL > ÉVÉNEMENTS > DÉTAIL
                    </p>
                    <h1 class="hero-title1">
                        {{ evenement?.titre || "Chargement..." }}
                    </h1>
                    <div v-if="evenement" class="project-sub-header">
                        <p class="classic-text">
                            Publié le
                            {{ formatDateLong(evenement.date_creation) }}
                            <span class="dot-separator">•</span>
                            {{ formatDateLong(evenement.date_evenement) }}
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

            <div v-else-if="evenement?.id" class="split-layout">
                <div class="left-column">
                    <div class="form-card">
                        <div class="card-header-flex">
                            <h2 class="card-title">
                                Présentation de l'événement
                            </h2>
                            <span
                                :class="[
                                    'type-badge',
                                    'type-' + evenement.type?.toLowerCase(),
                                ]"
                            >
                                {{ evenement.type?.toUpperCase() }}
                            </span>
                        </div>

                        <div class="description-section">
                            <label class="info-label"
                                >Description & Programme</label
                            >
                            <div class="description-box">
                                {{ evenement.description }}
                            </div>
                        </div>

                        <div class="specs-section">
                            <label class="info-label"
                                >Informations pratiques</label
                            >
                            <div class="specs-grid">
                                <div class="spec-item">
                                    <span class="spec-label"
                                        >Date de l'événement</span
                                    >
                                    <p class="spec-value highlight-val">
                                        {{
                                            formatDateLong(
                                                evenement.date_evenement,
                                            )
                                        }}
                                    </p>
                                </div>
                                <div class="spec-item">
                                    <span class="spec-label">Format</span>
                                    <p class="spec-value">
                                        {{ evenement.type }}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <div class="description-section">
                            <label class="info-label"
                                >Lieu exact de rendez-vous</label
                            >
                            <div class="description-box address-box">
                                <strong>📍 {{ evenement.adresse }}</strong
                                ><br />
                                {{ evenement.code_postal }}
                                {{ evenement.ville }}
                            </div>
                        </div>
                    </div>
                    <div v-if="canViewParticipants" class="form-card">
                        <div class="card-header-flex">
                            <h2 class="card-title">Liste des participants</h2>
                            <span
                                class="type-badge"
                                style="background: #e8eaf6; color: #3f51b5"
                                >Privé</span
                            >
                        </div>

                        <div
                            v-if="participantsLoading"
                            class="loading-state"
                            style="padding: 1rem"
                        >
                            Chargement de la liste...
                        </div>
                        <div
                            v-else-if="participants.length === 0"
                            class="description-box"
                            style="text-align: center; color: #666"
                        >
                            Aucun participant inscrit pour le moment.
                        </div>
                        <div v-else class="participants-list">
                            <div
                                v-for="p in participants"
                                :key="p.id"
                                class="participant-item cursor-pointer"
                                @click="viewProfile(p.id)"
                            >
                                <div class="mini-avatar">
                                    <img
                                        v-if="p.image_profil"
                                        :src="p.image_profil"
                                        class="mini-avatar-img"
                                    />
                                    <span v-else
                                        >{{ p.prenom?.charAt(0)
                                        }}{{ p.nom?.charAt(0) }}</span
                                    >
                                </div>
                                <div>
                                    <strong class="hover-underline"
                                        >{{ p.prenom }} {{ p.nom }}</strong
                                    >
                                    <div
                                        style="font-size: 0.75rem; color: #666"
                                    >
                                        {{ p.role }}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="right-column">
                    <div class="form-card side-card">
                        <h2 class="card-title-side">
                            Organisateur de l'événement
                        </h2>
                        <div
                            class="trainer-preview cursor-pointer"
                            @click="viewProfile(evenement.id_createur)"
                        >
                            <div class="mini-avatar">
                                <img
                                    v-if="createur?.image_profil"
                                    :src="createur.image_profil"
                                    class="mini-avatar-img"
                                />
                                <span v-else>
                                    {{ createur?.prenom?.charAt(0)
                                    }}{{ createur?.nom?.charAt(0) }}
                                </span>
                            </div>
                            <div>
                                <p class="trainer-name hover-underline">
                                    {{ createur?.prenom }} {{ createur?.nom }}
                                </p>
                                <button class="link-btn">
                                    Voir le profil expert
                                </button>
                            </div>
                        </div>
                    </div>

                    <div class="form-card side-card">
                        <h2 class="card-title-side">Participation</h2>
                        <div class="data-row">
                            <span class="data-label">Accès :</span>
                            <span class="text-success">Gratuit</span>
                        </div>
                        <div class="data-row">
                            <span class="data-label">Publié le :</span>
                            <span class="status-badge">{{
                                formatDateLong(evenement.date_creation)
                            }}</span>
                        </div>
                    </div>

                    <div class="form-actions-card">
                        <button
                            @click="handleInscription"
                            class="btn-save"
                            :class="{ 'btn-liked-active': isRegistered }"
                            :disabled="isRegistering"
                        >
                            <span v-if="isRegistered"
                                >Inscrit à l'événement</span
                            >
                            <span v-else>Participer à l'événement</span>
                        </button>

                        <button
                            v-if="isRegistered"
                            @click="handleQuit"
                            class="btn-quit"
                            :disabled="isLeaving"
                        >
                            Annuler ma participation
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

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const evenement = ref(null);
const createur = ref(null);
const isRegistered = ref(false);
const isRegistering = ref(false);
const isLeaving = ref(false);
const userScore = ref(0);
const participants = ref([]);
const participantsLoading = ref(false);

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const formatDateLong = (d) => {
    if (!d || d.startsWith("0001")) return "Date inconnue";
    return new Date(d).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
        year: "numeric",
    });
};

const fetchCreateur = async (id) => {
    try {
        const res = await fetch(`http://localhost:8081/users/${id}`);
        if (res.ok) createur.value = await res.json();
    } catch (error) {
        console.error("Erreur créateur: ", error);
    }
};

const canViewParticipants = computed(() => {
    if (!evenement.value) return false;
    const currentUserId = parseInt(sessionStorage.getItem("userId") || "0");
    const currentUserRole = sessionStorage.getItem("userRole") || "";

    return (
        currentUserId === evenement.value.id_createur ||
        currentUserRole === "Admin"
    );
});

const fetchParticipants = async (evenementId) => {
    participantsLoading.value = true;
    try {
        const res = await fetch(
            `http://localhost:8081/api/evenements/${evenementId}/participants`,
        );
        if (res.ok) {
            participants.value = await res.json();
        }
    } catch (error) {
        console.error("Erreur chargement participants:", error);
    } finally {
        participantsLoading.value = false;
    }
};

const fetchDetail = async () => {
    const id = route.params.id;
    const userId = sessionStorage.getItem("userId") || 0;
    try {
        const res = await fetch(
            `http://localhost:8081/evenements/${id}?user_id=${userId}`,
        );
        if (res.ok) {
            const data = await res.json();
            evenement.value = data;
            isRegistered.value = data.is_registered;
            if (data.id_createur) {
                await fetchCreateur(data.id_createur);
            }
            if (canViewParticipants.value) {
                fetchParticipants(id);
            }
        }
    } catch (error) {
        console.error("Erreur événement detail :", error);
    } finally {
        loading.value = false;
    }
};

const viewProfile = (id) => {
    if (id) {
        router.push(`/user/${id}`);
    } else {
        console.error("ID manquant");
    }
};

const handleInscription = async () => {
    if (!isLoggedIn.value) {
        router.push("/connexion");
        return;
    }
    const token = sessionStorage.getItem("userToken");
    const userId = sessionStorage.getItem("userId");
    isRegistering.value = true;
    try {
        const res = await fetch(
            `http://localhost:8081/api/evenements/${evenement.value.id}/join`,
            {
                method: "POST",
                headers: {
                    Authorization: token,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ id_utilisateur: parseInt(userId) }),
            },
        );
        if (res.status === 201 || res.ok) {
            alert("Votre participation a bien été enregistrée !");
            isRegistered.value = true;
            fetchDetail();
        } else {
            const errorMsg = await res.text();
            alert("Erreur : " + errorMsg);
        }
    } catch (error) {
        console.error("Erreur lors de l'inscription :", error);
        alert("Impossible de joindre le serveur.");
    } finally {
        isRegistering.value = false;
    }
};

const handleQuit = async () => {
    if (
        !confirm(
            "Voulez-vous vraiment annuler votre participation à cet événement ?",
        )
    )
        return;
    const token = sessionStorage.getItem("userToken");
    const userId = sessionStorage.getItem("userId");
    isLeaving.value = true;
    try {
        const res = await fetch(
            `http://localhost:8081/api/evenements/${evenement.value.id}/quit`,
            {
                method: "POST",
                headers: {
                    Authorization: token,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ id_utilisateur: parseInt(userId) }),
            },
        );
        if (res.ok) {
            alert("Votre participation a été annulée.");
            isRegistered.value = false;
            fetchDetail();
        } else {
            alert("Erreur lors de la désinscription.");
        }
    } catch (e) {
        alert("Erreur de connexion.");
    } finally {
        isLeaving.value = false;
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
    padding: 6px 12px;
    border-radius: 15px;
    font-size: 0.75rem;
    font-weight: 800;
}
.type-atelier {
    background: #e9f5ed;
    color: #2d7a4f;
}
.type-collecte {
    background: #fff3e0;
    color: #e65100;
}
.type-conference {
    background: #e8eaf6;
    color: #3f51b5;
}
.type-echange {
    background: #f3e5f5;
    color: #7b1fa2;
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

.address-box {
    line-height: 1.8;
}

.specs-section {
    display: flex;
    flex-direction: column;
    gap: 0.8rem;
}

.specs-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.2rem;
    background: #fafafa;
    padding: 1.2rem;
    border-radius: 10px;
    border: 1px solid #eee;
}

.spec-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.spec-label {
    font-size: 0.75rem;
    color: #aaa;
    text-transform: uppercase;
    font-weight: 700;
    letter-spacing: 0.5px;
}

.spec-value {
    font-weight: 700;
    color: #333;
    margin: 0;
    font-size: 0.95rem;
}

.highlight-val {
    color: #2d7a4f !important;
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
    flex-shrink: 0;
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

/* BOUTONS */
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

.btn-save:hover:not(:disabled) {
    background-color: #246343;
}

.btn-save:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-liked-active {
    background-color: #2d7a4f;
    opacity: 0.6;
    cursor: default;
}

.btn-quit {
    width: 100%;
    background: none;
    border: 1px solid #ef4444;
    color: #ef4444;
    padding: 10px;
    border-radius: 12px;
    margin-top: 12px;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.2s;
}

.btn-quit:hover:not(:disabled) {
    background: #fff5f5;
}

.loading-state {
    text-align: center;
    padding: 3rem;
    color: #8fa396;
    font-style: italic;
}

.participants-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 15px;
}

.participant-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px;
    background: #fafafa;
    border: 1px solid #eee;
    border-radius: 8px;
    transition:
        transform 0.2s,
        box-shadow 0.2s;
}

.cursor-pointer {
    cursor: pointer;
    transition: all 0.2s;
}

.cursor-pointer:hover {
    transform: translateY(-2px);
}

.participant-item:hover {
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
    border-color: #dcdcdc;
}
</style>
