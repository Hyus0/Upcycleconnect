<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL &gt; MES ANNONCES &gt; {{ annonce?.titre || "NULL" }}</p>
            <h1 class="hero-title1">{{ annonce?.titre || "Chargement..." }}</h1>
        </div>
        <div class="header-actions">
            <button class="btn-secondary" @click="$router.back()">Retour</button>
            <button v-if="annonce?.est_valide === 'En attente'" class="btn-modify" @click="goToEdit">Modifier l'annonce</button>
        </div>
    </header>

    <div v-if="loading" class="loading-state">Recuperation des donnees...</div>

    <div v-else-if="annonce?.id" class="split-layout">
        <div class="info-card">
            <div class="card-header-flex">
                <h2 class="card-title">Details de l'objet</h2>
                <span :class="annonce.type === 'Don' ? 'tag-don' : 'tag-vente'" class="type-badge">
                    {{ annonce.type === "Don" ? "DON" : `VENTE ${annonce.prix || 0} EUR` }}
                </span>
            </div>
            <div class="description-box">{{ annonce.description || "Description NULL" }}</div>
            <div class="specs-grid">
                <div class="spec-item"><label>Categorie</label><p>{{ categoryName || "NULL" }}</p></div>
                <div class="spec-item"><label>Materiau</label><p>{{ annonce.type_materiau || "NULL" }}</p></div>
                <div class="spec-item"><label>Poids</label><p>{{ annonce.poids_estime_kg || 0 }} kg</p></div>
                <div class="spec-item"><label>Etat</label><p>{{ annonce.etat_objet || "NULL" }}</p></div>
            </div>
        </div>

        <aside class="right-column">
            <div class="status-card">
                <h3>Position et validation</h3>
                <div class="data-row"><span>Verification :</span><span :class="annonce.est_valide === 'Valide' ? 'text-success' : 'text-pending'">{{ annonce.est_valide || "NULL" }}</span></div>
                <div class="data-row"><span>Etat actuel :</span><span class="status-badge">{{ annonce.statut || "NULL" }}</span></div>
                <div v-if="annonce.id_casier" class="data-row"><span>Casier reserve :</span><span class="value-casier">N{{ annonce.id_casier }}</span></div>
            </div>

            <div class="dates-card">
                <h3>Historique logistique</h3>
                <div class="timeline-item"><label>Mise en ligne</label><p>{{ formatDate(annonce.date_creation) }}</p></div>
                <div class="timeline-item"><label>Depot</label><p>{{ annonce.date_depot_effective ? formatDate(annonce.date_depot_effective) : "Pas encore depose" }}</p></div>
                <div class="timeline-item"><label>Recuperation</label><p>{{ annonce.date_recuperation_effective ? formatDate(annonce.date_recuperation_effective) : "En attente de retrait" }}</p></div>
            </div>
        </aside>
    </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { fetchAnnonce, fetchCategory } from "../../services/publicApi";

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const annonce = ref(null);
const categoryName = ref("");

const formatDate = (value) => {
    if (!value || String(value).startsWith("0001")) return "NULL";
    return new Date(value).toLocaleString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit"
    });
};

const goToEdit = () => {
    router.push({ name: "modification-annonce", params: { id: annonce.value.id } });
};

const loadAnnonce = async () => {
    loading.value = true;
    try {
        const payload = await fetchAnnonce(route.params.id);
        annonce.value = payload;
        if (payload?.id_categorie) {
            const categoryPayload = await fetchCategory(payload.id_categorie);
            categoryName.value = categoryPayload?.nom || "";
        }
    } catch (error) {
        console.error("Erreur detail annonce :", error);
        annonce.value = null;
    } finally {
        loading.value = false;
    }
};

onMounted(loadAnnonce);
</script>

<style scoped>
.info-card, .status-card, .dates-card { background: white; padding: 1.5rem; border-radius: 16px; border: 1px solid #eee; margin-bottom: 1.5rem; }
.split-layout { display: grid; grid-template-columns: 1.5fr 1fr; gap: 2rem; }
.description-box { background: #f9f9f9; padding: 1.2rem; border-radius: 12px; color: #444; line-height: 1.6; font-size: 1rem; border: 1px solid #f0f0f0; }
.specs-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; background: #fcfcfc; padding: 1.5rem; border-radius: 12px; border: 1px solid #f5f5f5; margin-top: 1.5rem; }
</style>
