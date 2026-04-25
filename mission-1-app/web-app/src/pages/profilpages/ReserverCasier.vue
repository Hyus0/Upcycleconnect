<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL > DEPOTS > RÉSERVER</p>
            <h1 class="hero-title1">CHOISIR UN POINT DE COLLECTE</h1>
            <p class="classic-text">Sélectionnez le site où vous souhaitez déposer votre objet.</p>
        </div>
        <button class="btn-secondary" @click="$router.back()">🠔 Annuler</button>
    </header>

    <div v-if="loading" class="loading-state">Recherche des sites disponibles...</div>

    <div v-else class="section-container">
        <div class="split-layout">
            <div class="left-column">
                <div class="sites-grid">
                    <div 
                        v-for="site in sites" 
                        :key="site.id" 
                        class="site-card"
                        :class="{ 'selected-site': selectedSiteId === site.id }"
                        @click="selectedSiteId = site.id"
                    >
                        <div class="site-info">
                            <span class="site-type">{{ site.type }}</span>
                            <h3 class="site-name">{{ site.nom }}</h3>
                            <p class="site-addr">📍 {{ site.adresse }}, {{ site.ville }}</p>
                            <p class="site-tel">📞 {{ site.telephone }}</p>
                        </div>
                        <div class="selection-indicator">
                            <div class="radio-circle"></div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="right-column">
                <div class="status-card sticky-card">
                    <h3>Récapitulatif</h3>
                    <div class="recap-item">
                        <label>Objet à déposer :</label>
                        <p><strong>{{ annonceTitle }}</strong></p>
                    </div>

                    <div v-if="selectedSiteId" class="recap-item">
                        <label>Site sélectionné :</label>
                        <p class="text-success">{{ selectedSiteName }}</p>
                    </div>

                    <div v-else class="warning-msg">
                        Veuillez sélectionner un site pour continuer.
                    </div>

                    <button 
                        class="btn-main-action" 
                        :disabled="!selectedSiteId || submitting"
                        @click="confirmReservation"
                    >
                        {{ submitting ? 'Réservation en cours...' : 'Confirmer la réservation' }}
                    </button>
                    
                    <p class="disclaimer">
                        En confirmant, un casier vous sera réservé pour 48h.
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();

const sites = ref([]);
const loading = ref(true);
const submitting = ref(false);
const selectedSiteId = ref(null);
const annonceTitle = ref("Chargement...");

// On récupère le nom du site sélectionné pour l'affichage
const selectedSiteName = computed(() => {
    const site = sites.value.find(s => s.id === selectedSiteId.value);
    return site ? site.nom : "";
});

const fetchData = async () => {
    const token = localStorage.getItem("userToken");
    try {
        // 1. On récupère les sites
        const resSites = await fetch("http://localhost:8081/sites", {
            headers: { "Authorization": token }
        });
        if (resSites.ok) sites.value = await resSites.json();

        // 2. On récupère les infos de l'annonce pour le récap
        const resAnnonce = await fetch(`http://localhost:8081/annonces/${route.params.id}`, {
            headers: { "Authorization": token }
        });
        if (resAnnonce.ok) {
            const data = await resAnnonce.json();
            annonceTitle.value = data.titre;
        }
    } catch (e) {
        console.error("Erreur chargement:", e);
    } finally {
        loading.value = false;
    }
};

const confirmReservation = async () => {
    submitting.value = true;
    const token = localStorage.getItem("userToken");
    
    try {
        const res = await fetch(`http://localhost:8081/annonces/${route.params.id}/reserver`, {
            method: "POST",
            headers: { 
                "Authorization": token,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ site_id: selectedSiteId.value })
        });

        if (res.ok) {
            // Succès ! On redirige vers la page des dépôts
            router.push('/profil/mes-depots');
        } else {
            alert("Erreur : Aucun casier disponible sur ce site.");
        }
    } catch (e) {
        alert("Erreur serveur.");
    } finally {
        submitting.value = false;
    }
};

onMounted(fetchData);
</script>

<style scoped>
/* REPRISE DE LA STRUCTURE SPLIT-LAYOUT */
.split-layout { display: grid; grid-template-columns: 1.5fr 1fr; gap: 2rem; }

.sites-grid { display: flex; flex-direction: column; gap: 1rem; }

.site-card {
    background: white;
    border: 1px solid #eee;
    padding: 1.5rem;
    border-radius: 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    transition: all 0.2s;
}

.site-card:hover { border-color: #2d6a4f; background: #f9fdfb; }

.selected-site {
    border: 2px solid #2d6a4f;
    background: #f0fdf4;
}

.site-type {
    font-size: 0.65rem;
    font-weight: 800;
    text-transform: uppercase;
    color: #2d6a4f;
    background: #e8f5e9;
    padding: 2px 8px;
    border-radius: 4px;
}

.site-name { margin: 8px 0 4px 0; font-size: 1.1rem; }
.site-addr { font-size: 0.9rem; color: #666; margin: 0; }
.site-tel { font-size: 0.8rem; color: #999; margin-top: 4px; }

/* INDICATEUR DE SÉLECTION (Petit cercle radio) */
.radio-circle {
    width: 20px; height: 20px; border: 2px solid #ddd; border-radius: 50%;
    position: relative;
}
.selected-site .radio-circle { border-color: #2d6a4f; }
.selected-site .radio-circle::after {
    content: ''; position: absolute; top: 3px; left: 3px;
    width: 10px; height: 10px; background: #2d6a4f; border-radius: 50%;
}

.sticky-card { position: sticky; top: 20px; }

.recap-item { margin-bottom: 1.5rem; }
.recap-item label { font-size: 0.75rem; color: #999; text-transform: uppercase; font-weight: bold; }
.recap-item p { font-size: 1.1rem; margin-top: 5px; }

.warning-msg { color: #f57c00; font-size: 0.85rem; margin-bottom: 1rem; font-style: italic; }

.disclaimer { font-size: 0.7rem; color: #aaa; text-align: center; margin-top: 1rem; }

.btn-main-action:disabled { background: #ccc; cursor: not-allowed; }
</style>