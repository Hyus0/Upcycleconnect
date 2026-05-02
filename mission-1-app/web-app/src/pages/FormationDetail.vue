<template>
    <main class="page-container">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            user-role="Particulier"
            :user-score="userScore"
        />

        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL &gt; FORMATIONS &gt; {{ formation?.titre || "NULL" }}</p>
                <h1 class="hero-title1">{{ formation?.titre || "Chargement..." }}</h1>
            </div>
            <div class="header-actions">
                <button class="btn-secondary" @click="$router.back()">Retour</button>
            </div>
        </header>

        <div v-if="loading" class="loading-state">Recuperation des donnees...</div>
        <div v-else-if="formation?.id" class="split-layout">
            <div class="info-card">
                <div class="card-header-flex">
                    <h2 class="card-title">Presentation de la session</h2>
                    <span class="type-badge">{{ formation.type || "NULL" }}</span>
                </div>

                <div class="description-box">{{ formation.description || "Description NULL" }}</div>

                <div class="specs-grid">
                    <div class="spec-item"><label>Date de debut</label><p>{{ formatDate(formation.date_debut) }}</p></div>
                    <div class="spec-item"><label>Date de fin</label><p>{{ formatDate(formation.date_fin) }}</p></div>
                    <div class="spec-item"><label>Capacite</label><p>{{ formation.capacite_max || 0 }} participants</p></div>
                    <div class="spec-item"><label>Adresse</label><p>{{ formation.adresse || "NULL" }}</p></div>
                </div>
            </div>

            <aside class="right-column">
                <div class="info-card">
                    <h3>Organisateur</h3>
                    <p>{{ formation.prenom_formateur || "NULL" }} {{ formation.nom_formateur || "" }}</p>
                    <button class="btn-secondary contact-button" type="button" @click="contactOrganizer">
                        Ouvrir une discussion
                    </button>
                </div>

                <div class="price-card">
                    <label>PRIX</label>
                    <div class="price-value">{{ formation.prix_unitaire > 0 ? `${formation.prix_unitaire} EUR` : "GRATUIT" }}</div>
                    <p class="price-hint">par personne</p>
                </div>

                <div class="card registration-card">
                    <button class="btn-secondary" type="button" @click="toggleCart">
                        {{ isInCart ? "Retirer du panier" : "Ajouter au panier" }}
                    </button>

                    <button
                        class="btn-main-action"
                        :class="{ 'btn-registered': isRegistered }"
                        :disabled="isRegistered || formation.statut !== 'Ouvert' || isRegistering"
                        @click="handleInscription"
                    >
                        <span v-if="isRegistered">Deja inscrit</span>
                        <span v-else>{{ formation.statut === "Ouvert" ? "Reserver ma place" : "Complet" }}</span>
                    </button>

                    <button v-if="isRegistered" class="btn-secondary" :disabled="isLeaving" @click="handleQuit">
                        Se desister de la session
                    </button>
                </div>
            </aside>
        </div>
    </main>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import { fetchFormation, joinFormation, quitFormation } from "../services/publicApi";
import { addItemToCart, isItemInCart, removeItemFromCart } from "../services/cartService";
import { startConversation } from "../services/messageService";

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const formation = ref(null);
const userScore = ref(0);
const isRegistered = ref(false);
const isRegistering = ref(false);
const isLeaving = ref(false);
const isInCart = computed(() => formation.value?.id ? isItemInCart("formation", formation.value.id) : false);

const isLoggedIn = computed(() => !!localStorage.getItem("userToken"));

const userName = computed(() => {
    const prenom = localStorage.getItem("userPrenom") || "";
    const nom = localStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const formatDate = (value) => {
    if (!value || String(value).startsWith("0001")) return "Non definie";
    return new Date(value).toLocaleString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit"
    });
};

const loadDetail = async () => {
    loading.value = true;
    try {
        const userId = Number(localStorage.getItem("userId") || 0);
        const payload = await fetchFormation(route.params.id, userId);
        formation.value = payload;
        isRegistered.value = Boolean(payload?.is_registered);
    } catch (error) {
        console.error("Erreur detail formation :", error);
        formation.value = null;
    } finally {
        loading.value = false;
    }
};

const handleInscription = async () => {
    const userId = Number(localStorage.getItem("userId"));
    if (!userId || !localStorage.getItem("userToken")) {
        alert("Vous devez etre connecte pour vous inscrire.");
        router.push("/connexion");
        return;
    }

    isRegistering.value = true;
    try {
        await joinFormation(formation.value.id, userId);
        isRegistered.value = true;
        await loadDetail();
    } catch (error) {
        console.error("Erreur inscription formation :", error);
        alert(error.message || "Inscription impossible.");
    } finally {
        isRegistering.value = false;
    }
};

const handleQuit = async () => {
    if (!confirm("Voulez-vous vraiment vous desinscrire de cette formation ?")) return;
    const userId = Number(localStorage.getItem("userId"));
    if (!userId) return;

    isLeaving.value = true;
    try {
        await quitFormation(formation.value.id, userId);
        isRegistered.value = false;
        await loadDetail();
    } catch (error) {
        console.error("Erreur desinscription formation :", error);
        alert(error.message || "Desinscription impossible.");
    } finally {
        isLeaving.value = false;
    }
};

const toggleCart = () => {
    if (!formation.value?.id) return;
    if (isInCart.value) {
        removeItemFromCart("formation", formation.value.id);
        return;
    }
    addItemToCart(formation.value, "formation");
};

const contactOrganizer = () => {
    if (!localStorage.getItem("userToken")) {
        router.push("/connexion");
        return;
    }

    const conversation = startConversation({
        kind: "formateur",
        targetId: formation.value?.id_formateur,
        name: `${formation.value?.prenom_formateur || ""} ${formation.value?.nom_formateur || ""}`.trim() || "Formateur",
        subject: formation.value?.titre,
        contextId: formation.value?.id,
        contextLabel: `Formation - ${formation.value?.titre}`
    });
    router.push({ path: "/messages", query: { conversation: conversation.id } });
};

onMounted(loadDetail);
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
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #f0f0f0;
}

.split-layout {
    display: grid;
    grid-template-columns: 1.6fr 1fr;
    gap: 1.5rem;
}

.info-card,
.price-card,
.registration-card {
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

.specs-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
}

.price-value {
    font-size: 2rem;
    font-weight: 800;
    color: #2d6a4f;
}

.contact-button,
.registration-card .btn-secondary,
.registration-card .btn-main-action {
    width: 100%;
    justify-content: center;
}

.registration-card {
    display: grid;
    gap: 10px;
}
</style>
