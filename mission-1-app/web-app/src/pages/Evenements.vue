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
                <p class="sidebar__category2">ACCUEIL &gt; EVENEMENTS</p>
                <h1 class="hero-title1">DECOUVRIR LES EVENEMENTS</h1>
                <p class="classic-text">
                    Participez aux ateliers, collectes, conferences et echanges publies depuis le back office.
                </p>
            </div>

            <div class="header-actions">
                <div class="search-box">
                    <input
                        v-model="searchQuery"
                        type="text"
                        placeholder="Rechercher un evenement..."
                        class="search-input"
                    />
                </div>
            </div>
        </header>

        <div class="section-container">
            <div class="filter-tags">
                <button class="tag-filter" :class="{ active: selectedType === 'All' }" @click="selectedType = 'All'">
                    Tous les types
                </button>
                <button class="tag-filter" :class="{ active: selectedType === 'Atelier' }" @click="selectedType = 'Atelier'">
                    Ateliers
                </button>
                <button class="tag-filter" :class="{ active: selectedType === 'Collecte' }" @click="selectedType = 'Collecte'">
                    Collectes
                </button>
                <button class="tag-filter" :class="{ active: selectedType === 'Conference' }" @click="selectedType = 'Conference'">
                    Conferences
                </button>
                <button class="tag-filter" :class="{ active: selectedType === 'Echange' }" @click="selectedType = 'Echange'">
                    Echanges
                </button>
            </div>

            <div v-if="loading" class="loading-state">Chargement des evenements...</div>
            <div v-else-if="filteredEvenements.length === 0" class="empty-msg">
                Aucun evenement disponible pour le moment.
            </div>

            <div v-else class="annonces-grid">
                <article v-for="evenement in filteredEvenements" :key="evenement.id" class="annonce-card formation-card">
                    <div class="card-header">
                        <span class="tag-type">{{ evenement.type || "Evenement" }}</span>
                        <span class="tag-price">{{ formatDate(evenement.date_evenement) }}</span>
                    </div>

                    <div class="card-body">
                        <h3 class="item-title">{{ evenement.titre }}</h3>
                        <p class="host-name">{{ evenement.ville || "Ville inconnue" }} ({{ evenement.code_postal || "NULL" }})</p>

                        <div class="info-list">
                            <div class="info-item">
                                <span>{{ formatDateLong(evenement.date_evenement) }}</span>
                            </div>
                            <div class="info-item">
                                <span>{{ evenement.adresse || "Adresse non renseignee" }}</span>
                            </div>
                            <div class="info-item">
                                <span>{{ truncate(evenement.description) || "Description NULL" }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="card-footer">
                        <button class="btn-secondary btn-full-width" type="button" @click="goToEvenement(evenement.id)">
                            Voir les details
                        </button>
                        <button
                            v-if="!inscrit.has(evenement.id)"
                            class="btn-main-action-full"
                            :disabled="inscriptionsEnCours.has(evenement.id)"
                            @click="inscrire(evenement)"
                        >
                            {{ inscriptionsEnCours.has(evenement.id) ? "Inscription..." : "S'inscrire a l'evenement" }}
                        </button>
                        <button
                            v-else
                            class="btn-secondary btn-full-width"
                            :disabled="inscriptionsEnCours.has(evenement.id)"
                            @click="desinscrire(evenement)"
                        >
                            {{ inscriptionsEnCours.has(evenement.id) ? "Desinscription..." : "Inscrit - Se desinscrire" }}
                        </button>
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
import {
    fetchEvenementInscriptionStatus,
    fetchEvenements,
    joinEvenement,
    quitEvenement
} from "../services/publicApi";

const evenements = ref([]);
const loading = ref(true);
const selectedType = ref("All");
const searchQuery = ref("");
const userScore = ref(0);
const inscriptionsEnCours = ref(new Set());
const inscrit = ref(new Set());
const router = useRouter();

const isLoggedIn = computed(() => !!localStorage.getItem("userToken"));

const userName = computed(() => {
    const prenom = localStorage.getItem("userPrenom") || "";
    const nom = localStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const filteredEvenements = computed(() =>
    evenements.value.filter((eventItem) => {
        const matchesType = selectedType.value === "All" || eventItem.type === selectedType.value;
        const term = searchQuery.value.toLowerCase();
        const matchesSearch = !term || (eventItem.titre || "").toLowerCase().includes(term);
        return matchesType && matchesSearch;
    })
);

const formatDate = (dateStr) => {
    if (!dateStr) return "NULL";
    return new Date(dateStr).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short"
    });
};

const formatDateLong = (dateStr) => {
    if (!dateStr) return "Date NULL";
    return new Date(dateStr).toLocaleDateString("fr-FR", {
        weekday: "long",
        day: "numeric",
        month: "long",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit"
    });
};

const truncate = (text, max = 100) => {
    if (!text) return "";
    return text.length > max ? `${text.slice(0, max)}...` : text;
};

const syncInscriptions = async () => {
    const userId = Number(localStorage.getItem("userId"));
    if (!userId) {
        inscrit.value = new Set();
        return;
    }

    const nextRegistered = new Set();
    await Promise.all(
        evenements.value.map(async (eventItem) => {
            try {
                const payload = await fetchEvenementInscriptionStatus(eventItem.id, userId);
                if (payload?.inscrit) {
                    nextRegistered.add(eventItem.id);
                }
            } catch {
                // Ignore status read failures item by item.
            }
        })
    );
    inscrit.value = nextRegistered;
};

const loadEvenements = async () => {
    loading.value = true;
    try {
        const payload = await fetchEvenements();
        evenements.value = Array.isArray(payload) ? payload : [];
        await syncInscriptions();
    } catch (error) {
        console.error("Erreur chargement evenements :", error);
        evenements.value = [];
    } finally {
        loading.value = false;
    }
};

const inscrire = async (evenement) => {
    const userId = Number(localStorage.getItem("userId"));
    if (!userId || !localStorage.getItem("userToken")) {
        alert("Vous devez etre connecte pour vous inscrire.");
        return;
    }
    if (inscriptionsEnCours.value.has(evenement.id)) return;

    inscriptionsEnCours.value.add(evenement.id);
    inscriptionsEnCours.value = new Set(inscriptionsEnCours.value);

    try {
        await joinEvenement(evenement.id, userId);
        inscrit.value = new Set([...inscrit.value, evenement.id]);
        alert(`Inscription confirmee pour "${evenement.titre}".`);
    } catch (error) {
        console.error("Erreur inscription evenement :", error);
        alert(error.message || "Inscription impossible.");
    } finally {
        inscriptionsEnCours.value.delete(evenement.id);
        inscriptionsEnCours.value = new Set(inscriptionsEnCours.value);
    }
};

const desinscrire = async (evenement) => {
    const userId = Number(localStorage.getItem("userId"));
    if (!userId || !localStorage.getItem("userToken")) return;
    if (inscriptionsEnCours.value.has(evenement.id)) return;

    inscriptionsEnCours.value.add(evenement.id);
    inscriptionsEnCours.value = new Set(inscriptionsEnCours.value);

    try {
        await quitEvenement(evenement.id, userId);
        inscrit.value.delete(evenement.id);
        inscrit.value = new Set(inscrit.value);
        alert(`Desinscription confirmee pour "${evenement.titre}".`);
    } catch (error) {
        console.error("Erreur desinscription evenement :", error);
        alert(error.message || "Desinscription impossible.");
    } finally {
        inscriptionsEnCours.value.delete(evenement.id);
        inscriptionsEnCours.value = new Set(inscriptionsEnCours.value);
    }
};

const goToEvenement = (id) => {
    router.push({ name: "evenement-detail", params: { id } });
};

onMounted(loadEvenements);
</script>

<style scoped>
.page-main-content {
    min-height: 100vh;
    padding: 20px;
    background: var(--bg-light, #f7f9f7);
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

.filter-tags {
    display: flex;
    gap: 10px;
    margin-bottom: 2rem;
    overflow-x: auto;
    padding-bottom: 5px;
}

.tag-filter {
    padding: 8px 20px;
    border-radius: 20px;
    border: 1px solid #eee;
    background: white;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 600;
    color: #666;
    white-space: nowrap;
}

.tag-filter.active {
    background: #2d6a4f;
    color: white;
    border-color: #2d6a4f;
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
}

.card-header,
.card-body,
.card-footer {
    padding: 1.2rem;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.card-footer {
    display: grid;
    gap: 10px;
}

.tag-type {
    background: #eaf4ed;
    color: #2d6a4f;
    font-weight: 800;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.7rem;
}

.tag-price {
    color: #2d6a4f;
    font-weight: 800;
    font-size: 0.9rem;
}

.item-title {
    font-size: 1.15rem;
    font-weight: 700;
    margin-bottom: 0.8rem;
    color: #222;
}

.host-name,
.info-item {
    color: #666;
    margin-bottom: 0.45rem;
}

.btn-main-action-full,
.btn-full-width {
    width: 100%;
}
</style>
