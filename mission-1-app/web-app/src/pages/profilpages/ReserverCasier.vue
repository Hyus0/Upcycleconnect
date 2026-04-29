<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL &gt; DEPOTS &gt; RESERVER</p>
            <h1 class="hero-title1">CHOISIR UN POINT DE COLLECTE</h1>
            <p class="classic-text">La reservation du casier passe par l'API annonce et s'enregistre en base.</p>
            <div v-if="errors.length" class="error-box"><ul><li v-for="(err, index) in errors" :key="index">{{ err }}</li></ul></div>
            <div v-if="successMsg" class="success-box">{{ successMsg }}</div>
        </div>
        <button class="btn-secondary" @click="$router.back()">Annuler</button>
    </header>

    <div v-if="loading" class="loading-state">Recherche des sites...</div>

    <div v-else class="section-container">
        <div class="split-layout">
            <div class="left-column">
                <div class="search-container">
                    <input v-model="searchQuery" type="text" placeholder="Rechercher une ville ou un site..." class="search-input" />
                </div>

                <div class="sites-list-vertical">
                    <div
                        v-for="site in filteredSites"
                        :key="site.id"
                        class="site-item"
                        :class="{ 'selected-item': selectedSiteId === site.id }"
                        @click="selectedSiteId = site.id"
                    >
                        <div class="item-body">
                            <span class="mini-type">{{ site.type || "Site" }}</span>
                            <h3 class="mini-name">{{ site.nom }}</h3>
                            <p class="mini-addr">{{ site.adresse }}, {{ site.ville }} ({{ site.code_postal }})</p>
                        </div>
                    </div>
                    <p v-if="filteredSites.length === 0" class="empty-text">Aucun site ne correspond a votre recherche.</p>
                </div>
            </div>

            <div class="right-column">
                <div class="status-card">
                    <h3 class="card-title">Resume de la reservation</h3>

                    <div v-if="annonce" class="annonce-preview-box">
                        <label class="field-label">OBJET</label>
                        <p class="preview-titre">{{ annonce.titre }}</p>
                        <div class="preview-tags">
                            <span class="tag-sm">{{ annonce.type }}</span>
                            <span class="tag-sm">{{ annonce.ville }}</span>
                        </div>
                    </div>

                    <div class="recap-section">
                        <label class="field-label">POINT DE COLLECTE</label>
                        <div v-if="selectedSiteId" class="selected-detail-box">
                            <p class="selected-site-name">{{ selectedSiteName }}</p>
                            <p class="selected-site-addr">
                                {{ selectedSite?.adresse }}<br />
                                {{ selectedSite?.code_postal }} {{ selectedSite?.ville }}
                            </p>
                        </div>
                        <div v-else class="warning-alert">Veuillez selectionner un site.</div>
                    </div>

                    <button class="btn-main-action full-width" :disabled="!selectedSiteId || submitting" @click="confirmReservation">
                        {{ submitting ? "Traitement..." : "Confirmer la reservation" }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { fetchAnnonce, fetchSites, reserveCasier } from "../../services/publicApi";

const route = useRoute();
const router = useRouter();
const sites = ref([]);
const annonce = ref(null);
const loading = ref(true);
const submitting = ref(false);
const selectedSiteId = ref(null);
const searchQuery = ref("");
const errors = ref([]);
const successMsg = ref("");

const filteredSites = computed(() => {
    const query = searchQuery.value.trim().toLowerCase();
    if (!query) return sites.value;
    return sites.value.filter((site) =>
        [site.nom, site.ville, site.adresse].some((value) => (value || "").toLowerCase().includes(query))
    );
});

const selectedSite = computed(() => sites.value.find((site) => site.id === selectedSiteId.value) || null);
const selectedSiteName = computed(() => selectedSite.value?.nom || "");

const loadPage = async () => {
    loading.value = true;
    try {
        const [sitesPayload, annoncePayload] = await Promise.all([fetchSites(), fetchAnnonce(route.params.id)]);
        sites.value = Array.isArray(sitesPayload) ? sitesPayload : [];
        annonce.value = annoncePayload || null;
    } catch (error) {
        console.error("Erreur reservation casier :", error);
    } finally {
        loading.value = false;
    }
};

const confirmReservation = async () => {
    if (!selectedSiteId.value) return;
    submitting.value = true;
    errors.value = [];
    try {
        await reserveCasier(route.params.id, selectedSiteId.value);
        successMsg.value = "Reservation reussie ! Redirection...";
        setTimeout(() => router.push("/profil/depots"), 1200);
    } catch (error) {
        errors.value = [error.message || "Serveur injoignable"];
    } finally {
        submitting.value = false;
    }
};

onMounted(loadPage);
</script>

<style scoped>
.split-layout { display: grid; grid-template-columns: 1.6fr 1fr; gap: 2.5rem; align-items: start; }
.search-input { width: 100%; padding: 12px 16px; border-radius: 12px; border: 1px solid #eee; font-size: 0.95rem; }
.sites-list-vertical { display: flex; flex-direction: column; gap: 1rem; margin-top: 1rem; }
.site-item { background: white; border: 1px solid #eee; padding: 1.2rem; border-radius: 16px; cursor: pointer; }
.selected-item { border: 2px solid #2d6a4f; background: #f0fdf4; }
.error-box { background-color: #fef2f2; border-left: 4px solid #ef4444; color: #b91c1c; padding: 1rem; margin: 1rem 0; border-radius: 8px; }
.success-box { background-color: #f0fdf4; border-left: 4px solid #22c55e; color: #15803d; padding: 1rem; margin: 1rem 0; border-radius: 8px; }
.status-card { background: white; border: 1px solid #eee; padding: 2rem; border-radius: 20px; display: flex; flex-direction: column; gap: 2rem; }
.annonce-preview-box { background: #fcfcfc; border: 1px solid #f0f0f0; padding: 1.2rem; border-radius: 12px; }
.field-label { display: block; font-size: 0.7rem; font-weight: 800; color: #bbb; letter-spacing: 1px; margin-bottom: 0.8rem; }
</style>
