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
                <p class="sidebar__category2">ACCUEIL &gt; FORMATIONS</p>
                <h1 class="hero-title1">APPRENDRE L'UPCYCLING</h1>
                <p class="classic-text">
                    Les sessions affichees ici viennent de l'API et donc du back admin.
                </p>
            </div>

            <div class="header-actions">
                <div class="search-box">
                    <input
                        v-model="searchQuery"
                        type="text"
                        placeholder="Rechercher une formation..."
                        class="search-input"
                    />
                </div>
            </div>
        </header>

        <div class="section-container">
            <div class="filter-tags">
                <button class="tag-filter" :class="{ active: selectedType === 'All' }" @click="selectedType = 'All'">
                    Tous les formats
                </button>
                <button class="tag-filter" :class="{ active: selectedType === 'Atelier' }" @click="selectedType = 'Atelier'">
                    Ateliers
                </button>
                <button class="tag-filter" :class="{ active: selectedType === 'Cours' }" @click="selectedType = 'Cours'">
                    Cours
                </button>
                <button class="tag-filter" :class="{ active: selectedType === 'Webinaire' }" @click="selectedType = 'Webinaire'">
                    Webinaires
                </button>
            </div>

            <div v-if="loading" class="loading-state">Chargement des formations...</div>
            <div v-else-if="filteredFormations.length === 0" class="empty-msg">
                Aucune formation disponible pour ce filtre.
            </div>

            <div v-else class="annonces-grid">
                <article v-for="formation in filteredFormations" :key="formation.id" class="annonce-card formation-card">
                    <div class="card-header">
                        <span class="tag-type">{{ formation.type || "Formation" }}</span>
                        <span class="tag-price">{{ formation.prix_unitaire > 0 ? `${formation.prix_unitaire} EUR` : "GRATUIT" }}</span>
                    </div>

                    <div class="card-body">
                        <h3 class="item-title">{{ formation.titre }}</h3>
                        <p class="host-name">{{ formation.ville || "Ville NULL" }} ({{ formation.code_postal || "NULL" }})</p>

                        <div class="info-list">
                            <div class="info-item">Debut : {{ formatDate(formation.date_debut) }}</div>
                            <div class="info-item">Fin : {{ formatDate(formation.date_fin) }}</div>
                            <div class="info-item">{{ formation.adresse || "Adresse NULL" }}</div>
                        </div>

                        <div class="status-zone">
                            <div class="slots-info">Capacite : {{ formation.capacite_max || 0 }} pers.</div>
                            <span class="status-badge" :class="formation.statut === 'Ouvert' ? 'status-open' : 'status-full'">
                                {{ formation.statut || "NULL" }}
                            </span>
                        </div>
                    </div>

                    <div class="card-footer">
                        <button
                            class="btn-secondary btn-full-width"
                            type="button"
                            @click="toggleFormationCart(formation)"
                        >
                            {{ isFormationInCart(formation.id) ? "Retirer du panier" : "Ajouter au panier" }}
                        </button>
                        <button class="btn-main-action-full" :disabled="formation.statut !== 'Ouvert'" @click="goToFormation(formation.id)">
                            Voir les details
                        </button>
                    </div>
                </article>
            </div>
        </div>
    </main>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import { fetchFormations } from "../services/publicApi";
import { addItemToCart, isItemInCart, onCartChange, removeItemFromCart } from "../services/cartService";

const formations = ref([]);
const loading = ref(true);
const selectedType = ref("All");
const searchQuery = ref("");
const userScore = ref(0);
const router = useRouter();
let stopCartSync = null;

const isLoggedIn = computed(() => !!localStorage.getItem("userToken"));

const userName = computed(() => {
    const prenom = localStorage.getItem("userPrenom") || "";
    const nom = localStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const filteredFormations = computed(() =>
    formations.value.filter((item) => {
        const isValidated = (item.Est_valide || item.est_valide || "").toLowerCase() === "valide";
        const matchesType = selectedType.value === "All" || item.type === selectedType.value;
        const term = searchQuery.value.toLowerCase();
        const matchesSearch = !term || (item.titre || "").toLowerCase().includes(term);
        return isValidated && matchesType && matchesSearch;
    })
);

const formatDate = (dateStr) => {
    if (!dateStr) return "NULL";
    return new Date(dateStr).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short",
        hour: "2-digit",
        minute: "2-digit"
    });
};

const loadFormations = async () => {
    loading.value = true;
    try {
        const payload = await fetchFormations();
        formations.value = Array.isArray(payload) ? payload : [];
    } catch (error) {
        console.error("Erreur chargement formations :", error);
        formations.value = [];
    } finally {
        loading.value = false;
    }
};

const goToFormation = (id) => {
    router.push({ name: "formation-detail", params: { id } });
};

const isFormationInCart = (id) => isItemInCart("formation", id);

const toggleFormationCart = (formation) => {
    if (isFormationInCart(formation.id)) {
        removeItemFromCart("formation", formation.id);
        return;
    }
    addItemToCart(formation, "formation");
};

onMounted(() => {
    loadFormations();
    stopCartSync = onCartChange(() => {
        formations.value = [...formations.value];
    });
});

onBeforeUnmount(() => {
    stopCartSync?.();
});
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
}

.card-header,
.card-body,
.card-footer {
    padding: 1.2rem;
}

.card-header {
    display: flex;
    justify-content: space-between;
}

.card-footer {
    display: grid;
    gap: 10px;
}

.tag-type {
    background: #f3e5f5;
    color: #7b1fa2;
    font-weight: 800;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.7rem;
}

.tag-price {
    color: #2d6a4f;
    font-weight: 800;
}

.item-title {
    font-size: 1.15rem;
    font-weight: 700;
    margin-bottom: 0.8rem;
    color: #222;
}

.info-item,
.host-name {
    color: #666;
    margin-bottom: 0.45rem;
}

.status-zone {
    margin-top: 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.status-open {
    color: #2d6a4f;
}

.status-full {
    color: #c2410c;
}
</style>
