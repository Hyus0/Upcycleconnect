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
                <p class="sidebar__category2">ACCUEIL > FORMATIONS</p>
                <h1 class="hero-title1">APPRENDRE L'UPCYCLING</h1>
                <p class="classic-text">
                    Développez vos compétences avec nos experts et formateurs passionnés.
                </p>
            </div>

            <div class="header-actions">
                <div class="search-section">
                    <div class="search-box">
                        <input 
                            v-model="searchQuery" 
                            type="text" 
                            placeholder="Rechercher par titre de formation..."
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
                    Tous les formats
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
                    :class="{ active: selectedType === 'Cours' }"
                    @click="selectedType = 'Cours'"
                >
                    Cours
                </button>
                <button 
                    class="tag-filter" 
                    :class="{ active: selectedType === 'Webinaire' }"
                    @click="selectedType = 'Webinaire'"
                >
                    Webinaires
                </button>
            </div>

            <div v-if="loading" class="loading-state">
                Chargement des formations...
            </div>

            <div v-else-if="filteredFormations.length === 0" class="empty-msg">
                Aucune formation de ce type disponible pour le moment.
            </div>

            <div v-else class="annonces-grid">
                <div
                    v-for="formation in filteredFormations"
                    :key="formation.id"
                    class="annonce-card formation-card"
                >
                    <div class="card-header">
                        <span :class="['tag-type', 'type-' + formation.type.toLowerCase()]">
                            {{ formation.type }}
                        </span>
                        <span class="tag-price">
                            {{ formation.prix_unitaire > 0 ? formation.prix_unitaire + "€" : "GRATUIT" }}
                        </span>
                    </div>

                    <div class="card-body">
                        <h3 class="item-title">{{ formation.titre }}</h3>
                        <p class="host-name">📍 {{ formation.ville }} ({{ formation.code_postal }})</p>

                        <div class="info-list">
                            <div class="info-item">
                                <span class="icon">📅</span>
                                <span>Début : {{ formatDate(formation.date_debut) }}</span>
                            </div>
                            <div class="info-item">
                                <span class="icon">🏁</span>
                                <span>Fin : {{ formatDate(formation.date_fin) }}</span>
                            </div>
                            <div class="info-item">
                                <span class="icon">🏠</span>
                                <span class="address-text">{{ formation.adresse }}</span>
                            </div>
                        </div>

                        <div class="status-zone">
                            <div class="slots-info">
                                <span class="info-label">Capacité :</span>
                                <span class="val">{{ formation.capacite_max }} pers.</span>
                            </div>
                            <span :class="['status-badge', formation.statut === 'Ouvert' ? 'status-open' : 'status-full']">
                                {{ formation.statut.toUpperCase() }}
                            </span>
                        </div>
                    </div>

                    <div class="card-footer">
                        <button
                            class="btn-main-action-full"
                            :disabled="formation.statut !== 'Ouvert'"
                            @click="goToFormation(formation.id)"
                        >
                            {{ formation.statut === "Ouvert" ? "Réserver ma place" : "Complet" }}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </main>
    <SiteFooter />
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";

import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

const formations = ref([]);
const loading = ref(true);
const selectedType = ref("All"); 
const searchQuery = ref(""); 
const userScore = ref(0); 
const router = useRouter();

const isLoggedIn = computed(() => {
    return !!sessionStorage.getItem("userToken");
});

const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const filteredFormations = computed(() => {
    return formations.value.filter(f => {
        const isValidated = (f.Est_valide || f.est_valide) === "Valide";
        const matchesType = selectedType.value === "All" || f.type === selectedType.value;
        const term = searchQuery.value.toLowerCase();
        const matchesSearch = !term || f.titre?.toLowerCase().includes(term);

        return isValidated && matchesType && matchesSearch; 
    });
});

const formatDate = (dateStr) => {
    if (!dateStr) return "N/C";
    const date = new Date(dateStr);
    return date.toLocaleDateString('fr-FR', {
        day: 'numeric',
        month: 'short',
        hour: '2-digit',
        minute: '2-digit'
    });
};

const fetchFormations = async () => {
    loading.value = true;
    try {
        const res = await fetch(`http://localhost:8081/formations`, {
            method: "GET"
        });
        if (res.ok) {
            formations.value = await res.json();
        }
    } catch (error) {
        console.error("Erreur lors du chargement des formations: ", error);
    } finally {
        loading.value = false;
    }
}

const goToFormation = (id) => {
    router.push({ name: 'formation-detail', params: { id: id } });
};

onMounted(fetchFormations);
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
.formation-card {
    border: 1px solid #f0f0f0;
    border-radius: 16px;
    background: white;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    transition: all 0.3s ease;
}
.formation-card:hover {
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
    font-size: 1rem;
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
    align-items: center;
    gap: 10px;
    font-size: 0.85rem;
    color: #444;
}
.status-zone {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 15px;
    border-top: 1px solid #f9f9f9;
    margin-bottom: 1.2rem;
}
.info-label {
    font-size: 0.75rem;
    color: #999;
}
.val {
    font-size: 0.8rem;
    font-weight: 700;
    color: #333;
}
.status-badge {
    padding: 4px 10px;
    border-radius: 5px;
    font-size: 0.7rem;
    font-weight: 900;
}
.status-open {
    background: #e8f5e9;
    color: #2e7d32;
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
}
.input-search {
    padding: 10px 20px;
    border-radius: 25px;
    border: 1px solid #ddd;
    width: 280px;
    outline: none;
    font-size: 0.9rem;
}
</style>
