<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL > DEPOTS > RÉSERVER</p>
            <h1 class="hero-title1">CHOISIR UN POINT DE COLLECTE</h1>
            <p class="classic-text">Sélectionnez le site de dépôt idéal pour votre objet.</p>
            <div v-if="errors.length > 0" class="error-box">
                <ul style="margin: 0; padding-left: 20px;">
                    <li v-for="(err, index) in errors" :key="index">{{ err }}</li>
                </ul>
            </div>
            <div v-if="successMsg" class="success-box">
                {{ successMsg }}
            </div>
        </div>
        <button class="btn-secondary" @click="$router.back()">🠔 Annuler</button>
    </header>

    <div v-if="loading" class="loading-state">Recherche des sites...</div>

    <div v-else class="section-container">
        <div class="split-layout">
            <div class="left-column">
                <div class="search-container">
                    <input 
                        type="text" 
                        v-model="searchQuery" 
                        placeholder="Rechercher une ville ou un nom de site..."
                        class="search-input"
                    >
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
                            <div class="item-header">
                                <span class="mini-type">{{ site.type }}</span>
                                <h3 class="mini-name">{{ site.nom }}</h3>
                            </div>
                            <p class="mini-addr">📍 {{ site.adresse }}, {{ site.ville }} ({{ site.code_postal }})</p>
                        </div>
                        <div class="selection-status">
                            <div class="custom-radio"></div>
                        </div>
                    </div>
                    <p v-if="filteredSites.length === 0" class="empty-text">Aucun site ne correspond à votre recherche.</p>
                </div>
            </div>

            <div class="right-column">
                <div class="status-card sticky-card">
                    <h3 class="card-title">Résumé de la réservation</h3>

                    <div v-if="annonce" class="annonce-preview-box">
                        <label class="field-label">OBJET À DÉPOSER</label>
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
                                {{ sites.find(s => s.id === selectedSiteId)?.adresse }}<br>
                                {{ sites.find(s => s.id === selectedSiteId)?.code_postal }} {{ sites.find(s => s.id === selectedSiteId)?.ville }}
                            </p>
                        </div>
                        <div v-else class="warning-alert">
                            Veuillez sélectionner un site dans la liste de gauche pour continuer.
                        </div>
                    </div>

                    <div class="action-footer">
                        <button
                            class="btn-main-action full-width"
                            :disabled="!selectedSiteId || submitting"
                            @click="confirmReservation"
                        >
                            {{ submitting ? "Traitement..." : "Confirmer la réservation" }}
                        </button>
                        <p class="info-note">Un code PIN vous sera attribué pour l'ouverture du casier.</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

const sites = ref([]);
const loading = ref(true);
const submitting = ref(false);
const selectedSiteId = ref(null);
const annonce = ref(null);
const searchQuery = ref("");

const errors = ref([]);
const successMsg = ref("");

const filteredSites = computed(() => {
    if (!searchQuery.value) return sites.value;
    const q = searchQuery.value.toLowerCase();
    return sites.value.filter(s => 
        s.nom.toLowerCase().includes(q) || 
        s.ville.toLowerCase().includes(q) ||
        s.adresse.toLowerCase().includes(q)
    );
});

const selectedSiteName = computed(() => {
    const site = sites.value.find((s) => s.id === selectedSiteId.value);
    return site ? site.nom : "";
});

const fetchData = async () => {
    const token = localStorage.getItem("userToken");
    try {
        const resSites = await fetch("http://localhost:8081/sites", {
            headers: { Authorization: token },
        });
        if (resSites.ok) sites.value = await resSites.json();

        const resAnnonce = await fetch(`http://localhost:8081/annonces/${route.params.id}`, {
            headers: { Authorization: token },
        });
        if (resAnnonce.ok) annonce.value = await resAnnonce.json();
    } catch (e) {
        console.error(e);
    } finally {
        loading.value = false;
    }
};

const confirmReservation = async () => {
    submitting.value = true;
    errors.value = []; 
    const token = localStorage.getItem("userToken");

    try {
        const res = await fetch(`http://localhost:8081/annonces/${route.params.id}/reserver`, {
            method: "POST",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ site_id: selectedSiteId.value }),
        });

        if (res.ok) {
            successMsg.value = "Réservation réussie ! Redirection...";
            setTimeout(() => {
                router.push("/profil/depots");
            }, 2000);   
        } else {
            const data = await res.json();
            errors.value.push(data.error || "Erreur inconnue");

            setTimeout(() => {
                errors.value = [];
            }, 5000);
        }
    } catch (e) {
        errors.value.push("Serveur injoignable");
        setTimeout(() => { errors.value = [] }, 5000);
    } finally {
        submitting.value = false;
    }
};

onMounted(fetchData);
</script>

<style scoped>
.split-layout {
    display: grid;
    grid-template-columns: 1.6fr 1fr;
    gap: 2.5rem;
    align-items: start;
}

.search-container {
    margin-bottom: 1.5rem;
}

.search-input {
    width: 100%;
    padding: 12px 16px;
    border-radius: 12px;
    border: 1px solid #eee;
    font-size: 0.95rem;
    outline: none;
    transition: border-color 0.2s;
}

.search-input:focus {
    border-color: #2d6a4f;
}

.sites-list-vertical {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.site-item {
    background: white;
    border: 1px solid #eee;
    padding: 1.2rem;
    border-radius: 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    transition: all 0.2s;
}

.site-item:hover {
    border-color: #2d6a4f;
    transform: translateX(4px);
}

.selected-item {
    border: 2px solid #2d6a4f;
    background: #f0fdf4;
}

.mini-type {
    font-size: 0.6rem;
    font-weight: 900;
    text-transform: uppercase;
    background: #e8f5e9;
    color: #2d6a4f;
    padding: 2px 8px;
    border-radius: 4px;
}

.mini-name {
    margin: 6px 0 2px 0;
    font-size: 1.1rem;
    font-weight: 700;
}

.mini-addr {
    font-size: 0.85rem;
    color: #666;
    margin: 0;
}

.custom-radio {
    width: 20px;
    height: 20px;
    border: 2px solid #ddd;
    border-radius: 50%;
    position: relative;
}

.selected-item .custom-radio {
    border-color: #2d6a4f;
}

.selected-item .custom-radio::after {
    content: "";
    position: absolute;
    top: 3px;
    left: 3px;
    width: 10px;
    height: 10px;
    background: #2d6a4f;
    border-radius: 50%;
}

.error-box {
    background-color: #fef2f2;
    border-left: 4px solid #ef4444;
    color: #b91c1c;
    padding: 1rem;
    margin: 1rem 0;
    border-radius: 8px;
    font-size: 0.9rem;
}

.success-box {
    background-color: #f0fdf4;
    border-left: 4px solid #22c55e;
    color: #15803d;
    padding: 1rem;
    margin: 1rem 0;
    border-radius: 8px;
    font-size: 0.9rem;
}

.status-card {
    background: white;
    border: 1px solid #eee;
    padding: 2rem;
    border-radius: 20px;
    display: flex;
    flex-direction: column;
    gap: 2rem;
}

.card-title {
    margin: 0;
    font-size: 1.2rem;
    color: #333;
}

.field-label {
    display: block;
    font-size: 0.7rem;
    font-weight: 800;
    color: #bbb;
    letter-spacing: 1px;
    margin-bottom: 0.8rem;
}

.annonce-preview-box {
    background: #fcfcfc;
    border: 1px solid #f0f0f0;
    padding: 1.2rem;
    border-radius: 12px;
}

.preview-titre {
    font-size: 1.1rem;
    font-weight: 800;
    margin: 0 0 8px 0;
}

.tag-sm {
    font-size: 0.65rem;
    font-weight: 700;
    background: #eee;
    padding: 2px 8px;
    border-radius: 4px;
    margin-right: 6px;
    text-transform: uppercase;
}

.selected-site-name {
    font-size: 1.1rem;
    font-weight: 800;
    color: #2d6a4f;
    margin: 0 0 4px 0;
}

.selected-site-addr {
    font-size: 0.9rem;
    color: #666;
    line-height: 1.4;
    margin: 0;
}

.warning-alert {
    color: #e67e22;
    font-size: 0.9rem;
    font-style: italic;
    background: #fff8f1;
    padding: 1rem;
    border-radius: 8px;
    border-left: 4px solid #e67e22;
}

.action-footer {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.full-width {
    width: 100%;
    padding: 14px;
    font-size: 1rem;
}

.info-note {
    font-size: 0.75rem;
    color: #aaa;
    text-align: center;
    margin: 0;
}

.sticky-card {
    position: sticky;
    top: 2rem;
}

.empty-text {
    text-align: center;
    color: #ccc;
    padding: 2rem;
}
</style>