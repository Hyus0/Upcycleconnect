<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL &gt; DEPOTS</p>
            <h1 class="hero-title1">GESTION DES FLUX</h1>
            <p class="classic-text">Les statuts de depot et de reservation sont lus depuis vos annonces backend.</p>
        </div>
        <button class="btn-main-action" @click="$router.push('/profil/createAnnonce')">+ Deposer une annonce</button>
    </header>

    <div v-if="loading" class="loading-state">Chargement de vos flux...</div>

    <div v-else class="section-container">
        <div class="dash-block">
            <h2 class="block-title">ANNONCES A PLANIFIER <span class="badge">{{ aPlanifier.length }}</span></h2>
            <table class="data-table">
                <thead><tr><th>OBJET</th><th>TYPE</th><th>VILLE</th><th class="text-right">ACTION</th></tr></thead>
                <tbody>
                    <tr v-for="item in aPlanifier" :key="item.id">
                        <td><strong>{{ item.titre }}</strong></td>
                        <td><span :class="item.type === 'Don' ? 'tag-don' : 'tag-vente'">{{ item.type === "Don" ? "DON" : `${item.prix || 0} EUR` }}</span></td>
                        <td>{{ item.ville || "NULL" }}</td>
                        <td class="text-right"><button class="btn-plan" @click="planifier(item.id)">Reserver un casier</button></td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div class="dash-block">
            <h2 class="block-title">FLUX ACTIFS <span class="badge">{{ actifs.length }}</span></h2>
            <table class="data-table">
                <thead><tr><th>OBJET</th><th>SITE</th><th>CODE PIN</th><th>STATUT</th><th class="text-right">ACTION</th></tr></thead>
                <tbody>
                    <tr v-for="item in actifs" :key="item.id">
                        <td><strong>{{ item.titre }}</strong></td>
                        <td>{{ sitesCache[item.id_site] || "Site inconnu" }}</td>
                        <td>{{ item.code_pin_depot || "---" }}</td>
                        <td><span class="status-neutral">{{ item.statut }}</span></td>
                        <td class="text-right"><button class="btn-remove" @click="annulerFlux(item.id)">{{ item.statut === "Reserve" ? "Annuler" : "Retirer l'objet" }}</button></td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchSite, fetchSites, fetchUserAnnonces, retirerAnnonceDuCasier } from "../../services/publicApi";

const router = useRouter();
const annonces = ref([]);
const sitesCache = ref({});
const loading = ref(true);

const aPlanifier = computed(() => annonces.value.filter((item) => item.statut === "Disponible"));
const actifs = computed(() => annonces.value.filter((item) => ["Reserve", "Depose"].includes(item.statut)));

const loadSites = async () => {
    try {
        const payload = await fetchSites();
        if (Array.isArray(payload)) {
            payload.forEach((site) => {
                sitesCache.value[site.id] = site.nom;
            });
        }
    } catch (error) {
        console.error("Erreur sites :", error);
    }
};

const loadAnnonces = async () => {
    const id = localStorage.getItem("userId");
    if (!id) return;
    loading.value = true;
    try {
        const payload = await fetchUserAnnonces(id);
        annonces.value = Array.isArray(payload) ? payload : [];
        const missingSiteIds = [...new Set(annonces.value.map((item) => item.id_site).filter(Boolean))].filter(
            (siteId) => !sitesCache.value[siteId]
        );
        await Promise.all(
            missingSiteIds.map(async (siteId) => {
                try {
                    const site = await fetchSite(siteId);
                    sitesCache.value[siteId] = site?.nom || "Site inconnu";
                } catch {
                    sitesCache.value[siteId] = "Site inconnu";
                }
            })
        );
    } catch (error) {
        console.error("Erreur depots :", error);
        annonces.value = [];
    } finally {
        loading.value = false;
    }
};

const planifier = (id) => router.push({ name: "reserve-casier", params: { id } });

const annulerFlux = async (id) => {
    if (!confirm("Voulez-vous vraiment annuler ou retirer ce flux ?")) return;
    try {
        await retirerAnnonceDuCasier(id);
        await loadAnnonces();
    } catch (error) {
        console.error("Erreur retrait casier :", error);
        alert(error.message || "Operation impossible.");
    }
};

onMounted(async () => {
    await loadSites();
    await loadAnnonces();
});
</script>

<style scoped>
.dash-block { background: white; border: 1px solid #eee; border-radius: 16px; padding: 1.5rem; margin-bottom: 2rem; }
.block-title { font-size: 0.9rem; font-weight: 800; color: #666; margin-bottom: 1rem; }
.badge { background: #eaf4ed; color: #2d7a4f; border-radius: 999px; padding: 0.15rem 0.6rem; }
.btn-plan { background: #2d7a4f; color: white; border: none; padding: 0.55rem 0.9rem; border-radius: 8px; }
.status-neutral { background: #f0f0f0; color: #666; padding: 4px 8px; border-radius: 4px; font-size: 0.8rem; font-weight: bold; }
</style>
