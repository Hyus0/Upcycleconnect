<template>
    <main class="public-dashboard">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            variant="public"
        />

        <header class="content-header annonces-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > CATALOGUE</p>
                <h1 class="hero-title1">Annonces disponibles</h1>
                <p class="classic-text">
                    Consulte les objets proposés sur UpcycleConnect, filtre par
                    type et trouve rapidement ce qui est réutilisable.
                </p>
            </div>
            <RouterLink
                v-if="isParticulier"
                class="btn-main-action"
                to="/profil/annonces"
                >+ Déposer une annonce</RouterLink
            >
        </header>

        <div class="stats-grid annonces-stats">
            <div class="card card--score">
                <p class="tag-score">CATALOGUE PUBLIC</p>
                <div class="score-value">
                    {{ filteredAnnonces.length }} <span>annonces</span>
                </div>
                <p class="score-level">{{ sourceLabel }}</p>
                <div class="score-footer">
                    <div class="mini-stat">
                        <strong>{{ donationsCount }}</strong
                        ><br />Dons
                    </div>
                    <div class="mini-stat">
                        <strong>{{ salesCount }}</strong
                        ><br />Ventes
                    </div>
                    <div class="mini-stat">
                        <strong>{{ filteredAnnonces.length }}</strong
                        ><br />Disponibles
                    </div>
                </div>
            </div>
            <div class="card card--white">
                <div class="card-num">{{ donationsCount }}</div>
                <p class="text-dm">Objets gratuits</p>
                <span class="badge badge--green">DON</span>
            </div>
            <div class="card card--white">
                <div class="card-num2">{{ salesCount }}</div>
                <p class="text-dm">Objets en vente</p>
                <span class="badge badge--orange">VENTE</span>
            </div>
        </div>

        <section class="section-container">
            <div class="section-header" style="flex-wrap: wrap; gap: 1rem">
                <div>
                    <h2>Catalogue des annonces</h2>
                    <p class="classic-text">
                        {{
                            loading
                                ? "Chargement..."
                                : `${filteredAnnonces.length} résultat${filteredAnnonces.length > 1 ? "s" : ""}`
                        }}
                    </p>
                </div>

                <div class="header-actions">
                    <input
                        v-model="filters.search"
                        type="search"
                        placeholder="Rechercher (titre, ville)..."
                        class="search-input"
                    />
                    <button
                        class="btn-secondary"
                        @click="showFilters = !showFilters"
                        :class="{ 'active-filter-btn': showFilters }"
                    >
                        <i class="ti ti-filter"></i> Filtres avancés
                    </button>
                </div>
            </div>

            <div v-if="showFilters" class="filters-panel">
                <div class="filter-group">
                    <label>Type de transaction</label>
                    <select v-model="filters.type" class="filter-input">
                        <option value="">Tous les types</option>
                        <option value="Don">Don</option>
                        <option value="Vente">Vente</option>
                    </select>
                </div>

                <div class="filter-group">
                    <label>Catégorie</label>
                    <select v-model="filters.categorie" class="filter-input">
                        <option value="">Toutes les catégories</option>
                        <option
                            v-for="cat in categories"
                            :key="cat.id"
                            :value="cat.id"
                        >
                            {{ cat.nom }}
                        </option>
                    </select>
                </div>

                <div class="filter-group">
                    <label>État de l'objet</label>
                    <select v-model="filters.etat" class="filter-input">
                        <option value="">Tous les états</option>
                        <option value="Neuf">Neuf</option>
                        <option value="Bon etat">Bon état</option>
                        <option value="Usage">À restaurer / Usagé</option>
                    </select>
                </div>

                <div class="filter-group">
                    <label>Poids maximum (kg)</label>
                    <input
                        v-model.number="filters.poidsMax"
                        type="number"
                        min="0"
                        placeholder="Ex: 50"
                        class="filter-input"
                    />
                </div>

                <div class="filter-group reset-group">
                    <button class="btn-text-danger" @click="resetFilters">
                        Réinitialiser
                    </button>
                </div>
            </div>

            <div v-if="loading" class="state-card">
                Chargement des annonces...
            </div>
            <div v-else-if="filteredAnnonces.length === 0" class="state-card">
                Aucune annonce ne correspond à ces filtres. Essayez de les
                modifier.
            </div>

            <div v-else class="annonces-grid">
                <article v-for="annonce in paginatedAnnonces" :key="annonce.id" class="annonce-card">
                    <div class="annonce-card__image-wrapper">
                        <img
                            :src="
                                annonce.image && annonce.image.trim() !== ''
                                    ? annonce.image
                                    : imageParDefaut
                            "
                            alt="Image de l'annonce"
                            class="annonce-card__image"
                        />
                        <div class="annonce-card__badges">
                            <span
                                :class="
                                    annonce.type === 'Vente'
                                        ? 'badge badge--orange'
                                        : 'badge badge--green'
                                "
                            >
                                {{ displayValue(annonce.type).toUpperCase() }}
                            </span>
                        </div>
                    </div>

                    <div class="annonce-card__content">
                        <div class="annonce-card__header">
                            <h3 class="annonce-card__title">
                                {{ displayValue(annonce.titre) }}
                            </h3>
                            <p class="annonce-card__price">
                                {{ formatPrice(annonce.prix, annonce.type) }}
                            </p>
                        </div>

                        <p class="annonce-card__desc">
                            {{ displayValue(annonce.description) }}
                        </p>

                        <div class="annonce-card__meta">
                            <span
                                >📍 {{ displayValue(annonce.ville) }} ({{
                                    displayValue(annonce.code_postal)
                                }})</span
                            >
                            <span
                                >⚖️
                                {{
                                    annonce.poids_estime_kg
                                        ? annonce.poids_estime_kg + " kg"
                                        : "Poids non précisé"
                                }}</span
                            >
                            <span
                                >📅
                                {{ formatDate(annonce.date_creation) }}</span
                            >
                        </div>
                    </div>

                    <div class="annonce-card__footer">
                        <button
                            class="btn-view btn-small"
                            type="button"
                            @click="goToAnnonce(annonce.id)"
                        >
                            Voir
                        </button>

                        <template
                            v-if="
                                parseInt(annonce.id_vendeur) === currentUserId
                            "
                        >
                            <button
                                class="btn-main-action btn-small"
                                type="button"
                                @click="modifierAnnonce(annonce.id)"
                            >
                                Modifier
                            </button>
                        </template>
                        <template v-else>
                            <button
                                v-if="annonce.statut === 'Disponible'"
                                class="btn-main-action btn-small"
                                @click="contacterVendeur(annonce)"
                            >
                                Contacter
                            </button>
                        </template>
                    </div>
                </article>
            </div>

            <div v-if="totalPages > 1" class="pagination">
                <button 
                  class="page-btn" 
                  :disabled="currentPage === 1" 
                  @click="changePage(currentPage - 1)"
                >
                  Précédent
                </button>
            
                <div class="page-numbers">
                  <button 
                    v-for="page in totalPages" 
                    :key="page"
                    class="page-btn"
                    :class="{ 'active': page === currentPage }"
                    @click="changePage(page)"
                  >
                    {{ page }}
                  </button>
                </div>
            
                <button 
                  class="page-btn" 
                  :disabled="currentPage === totalPages" 
                  @click="changePage(currentPage + 1)"
                >
                  Suivant
                </button>
              </div>
        </section>
    </main>
    <SiteFooter />
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from "vue";
import { RouterLink, useRouter, useRoute } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";
import imageParDefaut from "../components/upcycling-concept.jpg";

const router = useRouter();
const route = useRoute();
const loading = ref(true);
const source = ref("api");
const annonces = ref([]);
const categories = ref([]);
const showFilters = ref(false);
const API_URL = "/go";

const currentPage = ref(1);
const itemsPerPage = 15; 

const isParticulier = computed(() => {
    const role = sessionStorage.getItem("userRole");
    return role === "Particulier";
});

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

const filters = reactive({
    search: "",
    type: "",
    categorie: "",
    etat: "",
    poidsMax: null,
});

const resetFilters = () => {
    filters.search = "";
    filters.type = "";
    filters.categorie = "";
    filters.etat = "";
    filters.poidsMax = null;
    currentPage.value = 1; 
};

const filteredAnnonces = computed(() => {
    const search = filters.search.trim().toLowerCase();

    return annonces.value.filter((annonce) => {
        const matchesSearch =
            !search ||
            [annonce.titre, annonce.description, annonce.ville]
                .join(" ")
                .toLowerCase()
                .includes(search);
        const matchesType =
            !filters.type ||
            (annonce.type || "").toLowerCase() === filters.type.toLowerCase();
        const matchesCat =
            !filters.categorie ||
            String(annonce.id_categorie) === String(filters.categorie);
        const matchesEtat =
            !filters.etat || (annonce.etat_objet || "") === filters.etat;
        const matchesPoids =
            !filters.poidsMax ||
            (annonce.poids_estime_kg !== null &&
                parseFloat(annonce.poids_estime_kg) <= filters.poidsMax);

        return (
            matchesSearch &&
            matchesType &&
            matchesCat &&
            matchesEtat &&
            matchesPoids
        );
    });
});

watch(filteredAnnonces, () => {
    currentPage.value = 1;
});

const totalPages = computed(() => {
    return Math.ceil(filteredAnnonces.value.length / itemsPerPage);
});

const paginatedAnnonces = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage;
    const end = start + itemsPerPage;
    return filteredAnnonces.value.slice(start, end);
});

const changePage = (page) => {
    if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page;
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }
};

const donationsCount = computed(
    () =>
        filteredAnnonces.value.filter(
            (item) => (item.type || "").toLowerCase() === "don",
        ).length,
);
const salesCount = computed(
    () =>
        filteredAnnonces.value.filter(
            (item) => (item.type || "").toLowerCase() === "vente",
        ).length,
);

const sourceLabel = computed(() =>
    source.value === "api"
        ? "Données issues de l'API annonces"
        : "Aucune donnée disponible",
);

function displayValue(value) {
    return value === null || value === undefined || value === ""
        ? "N/A"
        : value;
}

const modifierAnnonce = (id) => {
    router.push(`/profil/modifyAnnonce/${id}`);
};

function formatDate(value) {
    if (!value) return "N/A";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "N/A";

    return new Intl.DateTimeFormat("fr-FR", {
        day: "2-digit",
        month: "short",
        year: "numeric",
    }).format(date);
}

function goToAnnonce(id) {
    router.push(`/annonce/${id}`);
}

function formatPrice(value, type) {
    if (value === null || value === undefined || value === "") return "N/A";
    if ((type || "").toLowerCase() === "don" || Number(value) === 0)
        return "Gratuit";

    return new Intl.NumberFormat("fr-FR", {
        style: "currency",
        currency: "EUR",
    }).format(Number(value));
}

const contacterVendeur = (annonce) => {
    if (!isLoggedIn.value) {
        alert("Veuillez vous connecter pour contacter le vendeur.");
        return router.push({
            path: "/connexion",
            query: { redirect: route.fullPath },
        });
    }

    if (currentUserId.value === 0) {
        alert("Erreur d'identification. Veuillez vous reconnecter.");
        return;
    }

    if (annonce && annonce.id_vendeur) {
        router.push({
            path: "/messages",
            query: {
                user: annonce.id_vendeur,
                annonce: annonce.id,
            },
        });
    }
};

onMounted(async () => {
    loading.value = true;

    try {
        const resAnnonces = await fetch(`${API_URL}/annonces`);
        if (resAnnonces.ok) {
            annonces.value = (await resAnnonces.json()) || [];
            source.value = "api";
        }

        const resCategories = await fetch(`${API_URL}/categories`);
        if (resCategories.ok) {
            categories.value = (await resCategories.json()) || [];
        }
    } catch (error) {
        console.error("Erreur de récupération des données :", error);
        annonces.value = [];
        source.value = "empty";
    } finally {
        loading.value = false;
    }
});
</script>

<style scoped>
.public-dashboard {
    min-height: 100vh;
    padding: 20px;
    background: var(--bg-light, #f7f9f7);
}

.annonces-header .btn-main-action {
    display: inline-flex;
    align-items: center;
    text-decoration: none;
}

.annonces-stats {
    grid-template-columns: 1.4fr 0.8fr 0.8fr;
}

.state-card {
    border: 1px dashed #cfe0d4;
    border-radius: 14px;
    padding: 26px;
    color: var(--text-grey);
    background: #fbfdfb;
    text-align: center;
}

.header-actions {
    display: flex;
    gap: 10px;
    align-items: center;
    flex-wrap: wrap;
}

.search-input {
    padding: 8px 16px;
    border: 1px solid #ddd;
    border-radius: 10px;
    font-family: inherit;
    font-size: 0.95rem;
    width: 250px;
}

.search-input:focus {
    outline: none;
    border-color: #2d7a4f;
}

.btn-secondary.active-filter-btn {
    background-color: #eaf4ed;
    color: #2d7a4f;
    border-color: #2d7a4f;
}

.filters-panel {
    display: flex;
    flex-wrap: wrap;
    gap: 1.5rem;
    background: #fafdfb;
    border: 1px solid #e5ede7;
    border-radius: 12px;
    padding: 1.5rem;
    margin-bottom: 2rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
}

.filter-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
    flex: 1;
    min-width: 180px;
}

.filter-group label {
    font-size: 0.8rem;
    font-weight: 700;
    color: #63746a;
    text-transform: uppercase;
}

.filter-input {
    padding: 10px 14px;
    border: 1px solid #ddd;
    border-radius: 8px;
    background: white;
    font-size: 0.95rem;
    font-family: inherit;
    color: #333;
}

.filter-input:focus {
    outline: none;
    border-color: #2d7a4f;
}

.reset-group {
    justify-content: flex-end;
    flex: 0 0 auto;
}

.btn-text-danger {
    background: transparent;
    border: none;
    color: #d32f2f;
    font-weight: bold;
    cursor: pointer;
    padding: 10px;
    border-radius: 8px;
    transition: 0.2s;
}

.btn-text-danger:hover {
    background: #fceaea;
}

.annonces-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 20px;
}

.annonce-card {
    background: #ffffff;
    border-radius: 16px;
    border: 1px solid #e5ede7;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    transition:
        transform 0.2s ease,
        box-shadow 0.2s ease;
    text-align: left;
}

.annonce-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 24px rgba(44, 126, 79, 0.08);
    border-color: #9bcbae;
}

.annonce-card__image-wrapper {
    position: relative;
    width: 100%;
    aspect-ratio: 4/3;
    background-color: #f0f4f1;
    overflow: hidden;
}

.annonce-card__image {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center;
}

.annonce-card__badges {
    position: absolute;
    top: 12px;
    left: 12px;
    display: flex;
    gap: 8px;
}

.annonce-card__content {
    padding: 16px;
    display: flex;
    flex-direction: column;
    flex: 1;
}

.annonce-card__header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 10px;
    margin-bottom: 8px;
}

.annonce-card__title {
    font-family: "Syne", sans-serif;
    font-size: 1.1rem;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.annonce-card__price {
    font-size: 1.1rem;
    font-weight: 800;
    color: #2c7e4f;
    margin: 0;
    white-space: nowrap;
}

.annonce-card__desc {
    font-size: 0.85rem;
    color: #6d7b72;
    margin: 0 0 16px 0;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.annonce-card__meta {
    margin-top: auto;
    display: flex;
    flex-direction: column;
    gap: 6px;
    font-size: 0.75rem;
    color: #8fa396;
}

.annonce-card__footer {
    padding: 14px 16px;
    border-top: 1px solid #f0f4f1;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
}

.annonce-card__footer > button,
.annonce-card__footer > a {
    box-sizing: border-box !important;
    width: 100% !important;
    height: 32px !important;
    min-height: 32px !important;
    max-height: 32px !important;
    margin: 0 !important;
    padding: 0 8px !important;

    display: inline-flex !important;
    align-items: center !important;
    justify-content: center !important;

    border-radius: 8px !important;
    font-size: 0.85rem !important;
    font-family: "Syne", sans-serif !important;
    font-weight: 600 !important;
    cursor: pointer !important;
    text-decoration: none !important;
    transition: all 0.2s !important;
    line-height: 1 !important;
    white-space: nowrap !important;
    overflow: hidden !important;
    text-overflow: ellipsis !important;

    border: 1px solid transparent !important;
}

.btn-view {
    background-color: #f0f4f1 !important;
    color: #2c7e4f !important;
    border-color: #f0f4f1 !important;
}

.btn-view:hover {
    background-color: #e1ede5 !important;
    border-color: #e1ede5 !important;
}

.annonce-card__footer .btn-main-action {
    background-color: #2c7e4f !important;
    color: #ffffff !important;
    border-color: #2c7e4f !important;
}

.annonce-card__footer .btn-main-action:hover {
    background-color: #23653e !important;
    border-color: #23653e !important;
}

.pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 15px;
    margin-top: 2.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid #eee;
}

.page-numbers {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
    justify-content: center;
}

.page-btn {
    background: white;
    border: 1px solid #ddd;
    color: #333;
    padding: 8px 14px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    transition: all 0.2s;
    font-family: inherit;
    font-size: 0.9rem;
}

.page-btn:hover:not(:disabled) {
    background: #f0f4f1;
    border-color: #2d7a4f;
    color: #2d7a4f;
}

.page-btn.active {
    background: #2d7a4f;
    color: white;
    border-color: #2d7a4f;
}

.page-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    background: #f5f5f5;
}

@media (max-width: 1500px) {
    .annonces-grid {
        grid-template-columns: repeat(4, 1fr);
    }
}

@media (max-width: 1200px) {
    .annonces-grid {
        grid-template-columns: repeat(3, 1fr);
    }
}

@media (max-width: 920px) {
    .annonces-stats {
        grid-template-columns: 1fr;
    }
    .annonces-grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .pagination {
        flex-direction: column;
    }
}

@media (max-width: 600px) {
    .annonces-grid {
        grid-template-columns: 1fr;
    }
}
</style>