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
                <p class="sidebar__category2">ACCUEIL > ÉVÈNEMENTS</p>
                <h1 class="hero-title1">DÉCOUVRIR LES ÉVÈNEMENTS</h1>
                <p class="classic-text">
                    Participez à nos ateliers, collectes, conférences et échanges autour de l'upcycling.
                </p>
            </div>

            <div class="header-actions">
                <div class="search-section">
                    <div class="search-box">
                        <input
                            v-model="searchQuery"
                            type="text"
                            placeholder="Rechercher un évènement..."
                            class="search-input"
                        />
                    </div>
                </div>
            </div>
        </header>

        <div class="section-container">
            <div class="filter-tags">
                <button
                    class="tag-filter"
                    :class="{ active: selectedType === 'All' }"
                    @click="selectedType = 'All'"
                >
                    Tous les types
                </button>
                <button
                    class="tag-filter"
                    :class="{ active: selectedType === 'Atelier' }"
                    @click="selectedType = 'Atelier'"
                >
                    Ateliers
                </button>
                <button
                    class="tag-filter"
                    :class="{ active: selectedType === 'Collecte' }"
                    @click="selectedType = 'Collecte'"
                >
                    Collectes
                </button>
                <button
                    class="tag-filter"
                    :class="{ active: selectedType === 'Conference' }"
                    @click="selectedType = 'Conference'"
                >
                    Conférences
                </button>
                <button
                    class="tag-filter"
                    :class="{ active: selectedType === 'Echange' }"
                    @click="selectedType = 'Echange'"
                >
                    Échanges
                </button>
            </div>

            <div v-if="loading" class="loading-state">
                Chargement des évènements...
            </div>

            <div v-else-if="filteredEvenements.length === 0" class="empty-msg">
                Aucun évènement disponible pour le moment.
            </div>

            <div v-else class="annonces-grid">
                <div
                    v-for="evenement in filteredEvenements"
                    :key="evenement.id"
                    class="annonce-card evenement-card"
                >
                    <div class="card-header">
                        <span :class="['tag-type', 'type-' + evenement.type.toLowerCase()]">
                            {{ evenement.type }}
                        </span>
                        <span class="tag-date">
                            {{ formatDate(evenement.date_evenement) }}
                        </span>
                    </div>

                    <div class="card-body">
                        <h3 class="item-title">{{ evenement.titre }}</h3>
                        <p class="host-name">📍 {{ evenement.ville }} ({{ evenement.code_postal }})</p>

                        <div class="info-list">
                            <div class="info-item">
                                <span class="icon">📅</span>
                                <span>{{ formatDateLong(evenement.date_evenement) }}</span>
                            </div>
                            <div class="info-item">
                                <span class="icon">🏠</span>
                                <span class="address-text">{{ evenement.adresse }}</span>
                            </div>
                            <div class="info-item">
                                <span class="icon">📝</span>
                                <span class="desc-text">{{ truncate(evenement.description) }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="card-footer">
                        <button
                            class="btn-main-action-full"
                            @click="inscrire(evenement.id)"
                        >
                            S'inscrire à l'évènement
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </main>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import SiteNavbar from "../components/SiteNavbar.vue";

const evenements = ref([]);
const loading = ref(true);
const selectedType = ref("All");
const searchQuery = ref("");
const userScore = ref(0);

const isLoggedIn = computed(() => {
    return !!localStorage.getItem("userToken");
});

const userName = computed(() => {
    const prenom = localStorage.getItem("userPrenom") || "";
    const nom = localStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const filteredEvenements = computed(() => {
    return evenements.value.filter(e => {
        const matchesType = selectedType.value === "All" || e.type === selectedType.value;
        const term = searchQuery.value.toLowerCase();
        const matchesSearch = !term || e.titre?.toLowerCase().includes(term);
        return matchesType && matchesSearch;
    });
});

const formatDate = (dateStr) => {
    if (!dateStr) return "N/C";
    const date = new Date(dateStr);
    return date.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short",
    });
};

const formatDateLong = (dateStr) => {
    if (!dateStr) return "N/C";
    const date = new Date(dateStr);
    return date.toLocaleDateString("fr-FR", {
        weekday: "long",
        day: "numeric",
        month: "long",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
};

const truncate = (text, max = 100) => {
    if (!text) return "";
    return text.length > max ? text.substring(0, max) + "..." : text;
};

const fetchEvenements = async () => {
    loading.value = true;
    try {
        const res = await fetch("http://localhost:8081/evenements", {
            method: "GET",
        });
        if (res.ok) {
            evenements.value = await res.json();
        }
    } catch (error) {
        console.error("Erreur lors du chargement des évènements :", error);
    } finally {
        loading.value = false;
    }
};

const inscrire = (id) => {
    // TODO : appel API inscription + feedback utilisateur
    alert(`Inscription à l'évènement #${id} (à implémenter)`);
};

onMounted(fetchEvenements);
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

.section-container {
    margin-top: 1rem;
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

.evenement-card {
    border: 1px solid #f0f0f0;
    border-radius: 16px;
    background: white;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    transition: all 0.3s ease;
}

.evenement-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.05);
}

.card-header {
    padding: 1.2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.tag-type {
    font-weight: 800;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.7rem;
    background: #e3f2fd;
    color: #1565c0;
}

.tag-type.type-atelier    { background: #f3e5f5; color: #7b1fa2; }
.tag-type.type-collecte   { background: #e8f5e9; color: #2e7d32; }
.tag-type.type-conference { background: #fff3e0; color: #e65100; }
.tag-type.type-echange    { background: #e3f2fd; color: #1565c0; }

.tag-date {
    color: #2d6a4f;
    font-weight: 800;
    font-size: 0.85rem;
}

.card-body {
    padding: 0 1.2rem 1.2rem;
    flex-grow: 1;
}

.item-title {
    font-size: 1.1rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
    color: #333;
}

.host-name {
    font-size: 0.8rem;
    color: #888;
    margin-bottom: 1.2rem;
}

.info-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin-bottom: 1.5rem;
}

.info-item {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    font-size: 0.85rem;
    color: #444;
}

.desc-text {
    color: #666;
    font-style: italic;
}

.card-footer {
    padding: 1.2rem;
    border-top: 1px solid #f9f9f9;
}

.btn-main-action-full {
    width: 100%;
    background: #2d6a4f;
    color: white;
    border: none;
    padding: 12px;
    border-radius: 12px;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.2s ease;
}

.btn-main-action-full:hover {
    background: #1b4332;
}

.loading-state,
.empty-msg {
    text-align: center;
    padding: 3rem;
    color: #999;
    font-size: 0.95rem;
}

.search-input {
    padding: 10px 20px;
    border-radius: 25px;
    border: 1px solid #ddd;
    width: 280px;
    outline: none;
    font-size: 0.9rem;
}
</style>