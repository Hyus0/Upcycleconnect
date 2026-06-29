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
                        ACCUEIL > FORMATIONS > DÉTAIL
                    </p>
                    <h1 class="hero-title1">
                        {{ formation?.titre || formation?.Titre || "Chargement..." }}
                    </h1>
                    <div v-if="formation" class="project-sub-header">
                        <p class="classic-text">
                            Du {{ formatDateLong(formation.date_debut || formation.Date_debut) }}
                            <span class="dot-separator">•</span>
                            au {{ formatDateLong(formation.date_fin || formation.Date_fin) }}
                            <span class="dot-separator">•</span>
                            {{ formation.nb_inscrit || formation.Nb_inscrit || 0 }} /
                            {{ formation.capacite_max || formation.Capacite_max || 0 }} participants (Total)
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

            <div v-else-if="formation?.id || formation?.ID" class="split-layout">
                <div class="left-column">
                    <div class="form-card">
                        <div class="card-header-flex">
                            <h2 class="card-title">
                                Présentation de la formation
                            </h2>
                            <span
                                :class="[
                                    'type-badge',
                                    'type-' + (formation.type || formation.Type || '').toLowerCase(),
                                ]"
                            >
                                {{ (formation.type || formation.Type || '').toUpperCase() }}
                            </span>
                        </div>

                        <div class="description-section">
                            <label class="info-label"
                                >Description & Programme</label
                            >
                            <div class="description-box">
                                {{ formation.description || formation.Description }}
                            </div>
                        </div>

                        <div class="specs-section">
                            <label class="info-label"
                                >Informations pratiques</label
                            >
                            <div class="specs-grid">
                                <div class="spec-item">
                                    <span class="spec-label"
                                        >Date de début (globale)</span
                                    >
                                    <p class="spec-value highlight-val">
                                        {{
                                            formatDateLong(formation.date_debut || formation.Date_debut)
                                        }}
                                    </p>
                                </div>
                                <div class="spec-item">
                                    <span class="spec-label">Date de fin (globale)</span>
                                    <p class="spec-value">
                                        {{ formatDateLong(formation.date_fin || formation.Date_fin) }}
                                    </p>
                                </div>
                                <div class="spec-item">
                                    <span class="spec-label">Capacité max par session</span>
                                    <p class="spec-value">
                                        {{
                                            formation.capacite_max || formation.Capacite_max
                                        }}
                                        participants
                                    </p>
                                </div>
                                <div class="spec-item">
                                    <span class="spec-label">Format</span>
                                    <p class="spec-value">
                                        {{ formation.type || formation.Type }}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <div class="description-section">
                            <label class="info-label"
                                >Lieu exact de rendez-vous</label
                            >
                            <div class="description-box address-box">
                                <strong>📍 {{ formation.adresse || formation.Adresse }}</strong
                                ><br />
                                {{ formation.code_postal || formation.CodePostal }}
                                {{ formation.ville || formation.Ville }}
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
                        <div class="participants-list" v-else>
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
                                    <div
                                        style="font-size: 0.75rem; color: #666"
                                    >
                                        {{ p.mail }}
                                    </div>
                                    <div 
                                        v-if="canViewParticipants && p.email"
                                        style="font-size: 0.75rem; color: #2d7a4f; font-weight: 600; margin-top: 2px;"
                                    >
                                        {{ p.email }}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="right-column">
                    <div class="form-card side-card">
                        <h2 class="card-title-side">
                            Organisateur
                        </h2>
                        <div
                            class="trainer-preview cursor-pointer"
                            @click="viewProfile(formation.id_formateur || formation.ID_formateur)"
                        >
                            <div class="mini-avatar">
                                <img
                                    v-if="formation.image_formateur || formation.Image_formateur"
                                    :src="formation.image_formateur || formation.Image_formateur"
                                    class="mini-avatar-img"
                                />
                                <span v-else>
                                    {{ (formation.prenom_formateur || formation.Prenom_formateur)?.charAt(0)
                                    }}{{ (formation.nom_formateur || formation.Nom_formateur)?.charAt(0) }}
                                </span>
                            </div>
                            <div>
                                <p class="trainer-name hover-underline">
                                    {{ formation.prenom_formateur || formation.Prenom_formateur }}
                                    {{ formation.nom_formateur || formation.Nom_formateur }}
                                </p>
                                <button class="link-btn">
                                    Voir le profil expert
                                </button>
                            </div>
                        </div>
                    </div>

                    <div class="form-card side-card price-card">
                        <h2 class="card-title-side">Prix par session</h2>
                        <div class="price-value">
                            {{
                                (formation.prix_unitaire || formation.Prix_unitaire) > 0
                                    ? (formation.prix_unitaire || formation.Prix_unitaire) + " €"
                                    : "GRATUIT"
                            }}
                        </div>
                        <p class="price-hint">par personne</p>

                        <div v-if="availableSessions.length > 0" class="session-selector mt-4 text-left">
                            <label class="info-label" style="display: block; margin-bottom: 8px;">Choisir une session :</label>
                            <select v-model="selectedSessionId" class="session-select w-full">
                                <option v-for="session in availableSessions" :key="session.id" :value="session.id" :disabled="session.statut !== 'Ouvert'">
                                    {{ session.nom }} {{ session.statut !== 'Ouvert' ? `(${session.statut})` : '' }}
                                </option>
                            </select>
                        </div>

                        <div class="main-actions mt-4">
                            <button 
                                v-if="!isRegistered"
                                @click="handleMainAction" 
                                class="btn-save w-full"
                                :disabled="isRegistering || isAddingToCart || !selectedSessionId || !isSessionOpen"
                            >
                                <template v-if="!isSessionOpen">❌ Session non disponible</template>
                                <template v-else-if="(formation.prix_unitaire || formation.Prix_unitaire) > 0">Ajouter au panier</template>
                                <template v-else>S'inscrire</template>
                            </button>

                            <button
                                v-else
                                @click="handleQuit(selectedSessionId)"
                                class="btn-quit w-full"
                                :disabled="isLeaving || !selectedSessionId"
                            >
                                Se désister de cette session
                            </button>

                            <button 
                                v-if="!isRegistered"
                                @click="handleTestInscription" 
                                class="btn-test w-full mt-2"
                                :disabled="isRegistering || !selectedSessionId || !isSessionOpen"
                            >
                                Inscription TEST (Directe)
                            </button>
                        </div>
                    </div>


                    <div class="form-card side-card">
                        <h2 class="card-title-side">Détail des sessions</h2>
                        
                        <div v-if="availableSessions.length > 0" class="sessions-list">
                            <div v-for="session in availableSessions" :key="session.id" class="session-item">
                                <div class="session-info">
                                    <strong>📅 {{ session.nom }}</strong>
                                    <div class="session-dates">
                                        Début: {{ formatDateLong(session.date_debut) }}<br/>
                                        Fin: {{ formatDateLong(session.date_fin) }}
                                    </div>
                                    <span class="status-badge" :class="session.statut === 'Ouvert' ? 'badge--green' : ''">
                                        {{ session.statut }}
                                    </span>
                                </div>
                            </div>
                        </div>

                        <div v-else class="description-box" style="text-align: center; color: #666; margin-top: 1rem;">
                            Aucune session n'est disponible pour cette formation.
                        </div>
                    </div>

                </div>
            </div>
        </main>
    </div>
    <SiteFooter />
</template>

<script setup>
import { ref, onMounted, computed, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const isRegistering = ref(false);
const isAddingToCart = ref(false);
const formation = ref(null);
const userScore = ref(0);

const isRegistered = ref(false);
const isLeaving = ref(false);

const participants = ref([]);
const participantsLoading = ref(false);

const selectedSessionId = ref(null);

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

// SOLUTION AU PROBLÈME DE CASSE JSON DE GO
const availableSessions = computed(() => {
    if (!formation.value) return [];
    
    const rawSessions = formation.value.sessions || formation.value.Sessions || [];
    
    return rawSessions.map(s => ({
        id: s.id || s.ID,
        nom: s.nom || s.Nom,
        date_debut: s.date_debut || s.DateDebut,
        date_fin: s.date_fin || s.DateFin,
        statut: s.statut || s.Statut
    }));
});

// Auto-sélection de la première session ouverte
watch(availableSessions, (newSessions) => {
    if (newSessions.length > 0 && !selectedSessionId.value) {
        const openSession = newSessions.find(s => s.statut === 'Ouvert');
        selectedSessionId.value = openSession ? openSession.id : newSessions[0].id;
    }
}, { immediate: true });

// Vérifie si la session sélectionnée est ouverte
const isSessionOpen = computed(() => {
    if (!selectedSessionId.value) return false;
    const session = availableSessions.value.find(s => s.id === selectedSessionId.value);
    return session && session.statut === 'Ouvert';
});

const canViewParticipants = computed(() => {
    if (!formation.value) return false;
    const currentUserId = parseInt(sessionStorage.getItem("userId") || "0");
    const currentUserRole = sessionStorage.getItem("userRole") || "";
    return (
        currentUserId === (formation.value.id_formateur || formation.value.ID_formateur) ||
        currentUserRole === "Admin"
    );
});

const formatDateLong = (d) => {
    if (!d || d.startsWith("0001") || d.startsWith("1970")) return "Non définie";
    const dateObj = new Date(d);
    return dateObj.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit"
    }).replace(":", "h");
};

const fetchDetail = async () => {
    const id = route.params.id;
    const userId = sessionStorage.getItem("userId") || 0;
    try {
        const res = await fetch(
            `/go/formations/${id}?user_id=${userId}`,
        );
        if (res.ok) {
            const data = await res.json();
            formation.value = data;
            isRegistered.value = data.is_registered || data.IsRegistered || false;

            if (canViewParticipants.value) {
                fetchParticipants(id);
            }
        }
    } catch (error) {
        console.error("Erreur fetch :", error);
    } finally {
        loading.value = false;
    }
};

const fetchParticipants = async (formationId) => {
    participantsLoading.value = true;
    try {
        const res = await fetch(
            `/go/api/formations/${formationId}/participants`,
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

const viewProfile = (id) => {
    if (id) {
        router.push(`/user/${id}`);
    } else {
        console.error("ID manquant pour la redirection");
    }
};

// Actions globales basées sur la session sélectionnée
const handleMainAction = () => {
    if (!selectedSessionId.value) return;
    const prix = formation.value.prix_unitaire || formation.value.Prix_unitaire || 0;
    if (prix > 0) {
        handleAddToCart(selectedSessionId.value);
    } else {
        handleInscription(selectedSessionId.value);
    }
};

const handleTestInscription = () => {
    if (!selectedSessionId.value) return;
    handleInscription(selectedSessionId.value);
};

const handleAddToCart = async (sessionId) => {
    const token = sessionStorage.getItem("userToken");
    const userId = sessionStorage.getItem("userId");
    if (!token || !userId) return router.push("/connexion");
    
    isAddingToCart.value = true;
    try {
        const res = await fetch(
            `/go/users/${userId}/panier`,
            {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    type_item: "Formation",
                    reference_id: parseInt(formation.value.id || formation.value.ID),
                    prix_unitaire: parseFloat(formation.value.prix_unitaire || formation.value.Prix_unitaire),
                }),
            },
        );
        if (res.ok) {
            alert("Session de formation ajoutée à votre panier !");
        } else {
            const errorMsg = await res.text();
            alert("Erreur lors de l'ajout au panier : " + errorMsg);
        }
    } catch (error) {
        console.error("Erreur panier :", error);
        alert("Impossible de joindre le serveur.");
    } finally {
        isAddingToCart.value = false;
    }
};

// S'inscrire à la session
const handleInscription = async (sessionId) => {
    const token = sessionStorage.getItem("userToken");
    const userId = sessionStorage.getItem("userId");
    if (!token || !userId) return router.push("/connexion");
    
    isRegistering.value = true;
    try {
        const formationId = formation.value.id || formation.value.ID;
        const res = await fetch(
            `/go/api/formations/${formationId}/join`,
            {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ 
                    id_utilisateur: parseInt(userId),
                    id_session: parseInt(sessionId) // Utilise la session spécifique
                }),
            },
        );
        if (res.status === 201 || res.status === 200) {
            alert("Inscription à la session réussie !");
            isRegistered.value = true;
            fetchDetail();
        } else if (res.status === 409) {
            alert("Désolé, cette session est déjà complète.");
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

// Se désister
const handleQuit = async (sessionId) => {
    if (!confirm("Voulez-vous vraiment vous désinscrire de cette session ?"))
        return;
    const token = sessionStorage.getItem("userToken");
    const userId = sessionStorage.getItem("userId");
    isLeaving.value = true;
    try {
        const formationId = formation.value.id || formation.value.ID;
        const res = await fetch(
            `/go/api/formations/${formationId}/quit`,
            {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ 
                    id_utilisateur: parseInt(userId),
                    id_session: parseInt(sessionId)
                }),
            },
        );
        if (res.ok) {
            alert("Vous avez bien été désinscrit.");
            isRegistered.value = false;
            fetchDetail();
        } else {
            const msg = await res.text();
            alert("Erreur lors de la désinscription : " + msg);
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
.type-cours {
    background: #e8eaf6;
    color: #3f51b5;
}
.type-webinaire {
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

.price-card {
    text-align: center;
}

.price-value {
    font-size: 2.8rem;
    font-weight: 900;
    color: #2d7a4f;
    margin: 0.2rem 0;
}

.price-hint {
    font-size: 0.85rem;
    color: #999;
    margin: 0;
}

.session-selector {
    width: 100%;
}

.session-select {
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-family: inherit;
    font-size: 0.95rem;
    background-color: white;
    outline: none;
    color: #333;
}
.session-select:focus {
    border-color: #2d7a4f;
}

.text-left {
    text-align: left;
}

.sessions-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-top: 1rem;
}

.session-item {
    background: #fbfdfb;
    border: 1px solid #e5ede7;
    border-radius: 10px;
    padding: 16px;
}

.session-info strong {
    display: block;
    color: #2d7a4f;
    margin-bottom: 6px;
    font-size: 1rem;
}

.session-dates {
    font-size: 0.85rem;
    color: #666;
    line-height: 1.5;
    margin-bottom: 8px;
}

.w-full {
    width: 100%;
}

.mt-2 {
    margin-top: 0.5rem;
}
.mt-4 {
    margin-top: 1rem;
}

.btn-save {
    background-color: #2d7a4f;
    color: white;
    padding: 12px;
    border: none;
    border-radius: 8px;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.2s;
}

.btn-save:hover:not(:disabled) {
    background-color: #246343;
}

.btn-save:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-test {
    padding: 12px;
    background-color: #f3f4f6;
    color: #374151;
    border: 1px solid #d1d5db;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    transition: background 0.2s;
}

.btn-test:hover:not(:disabled) {
    background-color: #e5e7eb;
}

.btn-test:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-quit {
    background: transparent;
    border: 1px solid #ef4444;
    color: #ef4444;
    padding: 12px;
    border-radius: 8px;
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
</style>