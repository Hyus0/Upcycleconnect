<template>
    <main class="page-main-content">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            :user-role="userRole"
            :user-score="userScore"
        />

        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > ESPACE CONSEILS</p>
                <h1 class="hero-title1">GUIDES & ASTUCES POUR {{ userRole.toUpperCase() }}</h1>
                <p class="classic-text">
                    Découvrez nos meilleurs conseils pour donner une seconde vie à vos objets.
                </p>
            </div>

            <div class="header-actions">
                <div class="search-section">
                    <div class="search-box">
                        <input 
                            v-model="searchQuery" 
                            type="text" 
                            placeholder="Rechercher une astuce..."
                            class="input-search"
                        />
                    </div>
                </div>
            </div>
        </header>

        <div class="section-container">
            <div v-if="loading" class="loading-state">
                Chargement des conseils...
            </div>

            <div v-else-if="filteredTips.length === 0" class="empty-msg">
                Aucun conseil ne correspond à votre recherche pour votre profil.
            </div>

            <div v-else class="annonces-grid">
                <div
                    v-for="tip in filteredTips"
                    :key="tip.id"
                    class="annonce-card formation-card"
                    @click="goToDetail(tip.id)"
                    style="cursor: pointer;"
                >
                    <div class="card-header">
                        <span :class="['tag-type', getRoleClass(tip.role_cible)]">
                            {{ tip.role_cible }}
                        </span>
                        <span class="tag-price">
                            Astuce
                        </span>
                    </div>

                    <div class="card-body">
                        <h3 class="item-title">{{ tip.titre }}</h3>
                        <p class="host-name">📍 Ajouté le {{ formatDate(tip.date_creation) }}</p>

                        <div class="info-list" style="margin-bottom: 0;">
                            <div class="info-item">
                                <span class="article-text">{{ truncateText(tip.description, 120) }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="card-footer">
                        <button class="btn-main-action-full btn-outline">
                            Lire la suite →
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </main>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";

const API_URL = "http://localhost:8081";
const router = useRouter();

const tips = ref([]);
const loading = ref(true);

const searchQuery = ref("");

const userScore = ref(parseInt(localStorage.getItem("userScore")) || 0);
const userRole = ref(localStorage.getItem("userRole") || "Particulier");

const isLoggedIn = computed(() => {
    return !!localStorage.getItem("userToken");
});

const userName = computed(() => {
    const prenom = localStorage.getItem("userPrenom") || "";
    const nom = localStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const filteredTips = computed(() => {
    return tips.value.filter(tip => {
        // Filtrage strict : on ne montre QUE les conseils dédiés au rôle de l'utilisateur
        const matchesRole = tip.role_cible === userRole.value;
        const term = searchQuery.value.toLowerCase();
        const matchesSearch = !term || tip.titre?.toLowerCase().includes(term) || tip.description?.toLowerCase().includes(term);

        return matchesRole && matchesSearch;
    });
});

const getRoleClass = (role) => {
    if (role === 'Particulier') return 'type-cours';
    if (role === 'Prestataire') return 'type-webinaire'; 
    if (role === 'Salarie') return 'type-atelier';
    return '';
};

const formatDate = (dateStr) => {
    if (!dateStr) return "N/C";
    const date = new Date(dateStr);
    return date.toLocaleDateString('fr-FR', {
        day: 'numeric',
        month: 'short',
        year: 'numeric'
    });
};

const truncateText = (text, length) => {
    if (!text) return "";
    if (text.length <= length) return text;
    return text.substring(0, length) + "...";
};

const goToDetail = (id) => {
    router.push({ name: "conseil-detail", params: { id } });
};

const fetchTips = async () => {
    loading.value = true;
    try {
        const res = await fetch(`${API_URL}/tips`);
        if (res.ok) {
            tips.value = await res.json() || [];
        }
    } catch (error) {
        console.error("Erreur lors du chargement des conseils: ", error);
    } finally {
        loading.value = false;
    }
};

onMounted(() => {
    fetchTips();
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

.section-container {
    margin-top: 1rem;
}

.loading-state, .empty-msg {
    text-align: center;
    padding: 3rem;
    color: #666;
    font-size: 1.1rem;
    background: white;
    border-radius: 12px;
    border: 1px dashed #ddd;
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
    font-weight: 800;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.7rem;
    text-transform: uppercase;
}
.type-cours { background: #eaf4ed; color: #2d7a4f; }
.type-webinaire { background: #e0f2fe; color: #0369a1; }
.type-atelier { background: #f3e5f5; color: #7b1fa2; }

.tag-price {
    color: #2d6a4f;
    font-weight: 800;
    font-size: 0.9rem;
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
    line-height: 1.4;
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
    font-size: 0.9rem;
    color: #555;
    line-height: 1.5;
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
    transition: background 0.2s;
}
.btn-outline {
    background: transparent;
    color: #2d6a4f;
    border: 2px solid #2d6a4f;
}
.formation-card:hover .btn-outline {
    background: #2d6a4f;
    color: white;
}

.input-search {
    padding: 10px 20px;
    border-radius: 12px;
    border: 1px solid #ddd;
    width: 280px;
    outline: none;
    font-size: 0.9rem;
}
.input-search:focus {
    border-color: #2d6a4f;
}
</style>