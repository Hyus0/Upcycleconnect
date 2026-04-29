<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL &gt; MES ANNONCES</p>
            <h1 class="hero-title1">VOS ANNONCES</h1>
            <p class="classic-text">Cette liste vient directement du back et des validations admin.</p>
        </div>
        <router-link to="/profil/createAnnonce" class="btn-main-action">+ Deposer une annonce</router-link>
    </header>

    <div class="section-container">
        <table class="data-table">
            <thead>
                <tr>
                    <th>OBJET</th>
                    <th>TYPE</th>
                    <th>VERIFICATION</th>
                    <th>STATUT</th>
                    <th>DATE</th>
                    <th>ACTIONS</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="annonce in annonces" :key="annonce.id">
                    <td><strong>{{ annonce.titre }}</strong><br /><small>{{ annonce.type_materiau || "NULL" }}</small></td>
                    <td>
                        <span :class="annonce.type === 'Don' ? 'tag-don' : 'tag-vente'">
                            {{ annonce.type === "Don" ? "DON" : `VENTE ${annonce.prix || 0} EUR` }}
                        </span>
                    </td>
                    <td>
                        <span :class="annonce.est_valide === 'Valide' ? 'status-valid' : 'status-pending'">
                            {{ annonce.est_valide === "Valide" ? "APPROUVEE" : "EN ANALYSE" }}
                        </span>
                    </td>
                    <td><span class="status-neutral">{{ annonce.statut || "NULL" }}</span></td>
                    <td>{{ formatDate(annonce.date_creation) }}</td>
                    <td class="actions-cell">
                        <button v-if="annonce.est_valide === 'Valide' && annonce.statut === 'Disponible'" class="btn-plan" @click="goToPlanning(annonce.id)">
                            Planifier depot
                        </button>
                        <button class="btn-view" @click="goToAnnonce(annonce.id)">Voir</button>
                        <button v-if="annonce.est_valide === 'En attente'" class="btn-modify" @click="goToModify(annonce.id)">Modifier</button>
                        <button class="btn-remove" @click="removeAnnonce(annonce.id)">Retirer</button>
                    </td>
                </tr>
            </tbody>
        </table>

        <div v-if="annonces.length === 0" class="empty-state">
            <p>Vous n'avez pas encore d'annonces en cours.</p>
            <router-link to="/profil/createAnnonce" class="btn-secondary">Creer ma premiere annonce</router-link>
        </div>
    </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { deleteAnnonce, fetchUserAnnonces } from "../../services/publicApi";

const router = useRouter();
const annonces = ref([]);

const formatDate = (dateString) => {
    if (!dateString) return "NULL";
    return new Date(dateString).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric"
    });
};

const goToAnnonce = (id) => router.push({ name: "see-annonce", params: { id } });
const goToModify = (id) => router.push({ name: "modification-annonce", params: { id } });
const goToPlanning = (id) => router.push({ name: "mes-depots", query: { selectedAnnonce: id } });

const loadAnnonces = async () => {
    const id = localStorage.getItem("userId");
    if (!id) return;
    try {
        const payload = await fetchUserAnnonces(id);
        annonces.value = Array.isArray(payload) ? payload : [];
    } catch (error) {
        console.error("Erreur annonces :", error);
        annonces.value = [];
    }
};

const removeAnnonce = async (id) => {
    if (!confirm("Voulez-vous vraiment retirer cette annonce ?")) return;
    try {
        await deleteAnnonce(id);
        annonces.value = annonces.value.filter((item) => item.id !== id);
    } catch (error) {
        console.error("Erreur suppression annonce :", error);
        alert(error.message || "Suppression impossible.");
    }
};

onMounted(loadAnnonces);
</script>

<style scoped>
.btn-plan {
    background-color: #2d7a4f;
    color: white;
    border: none;
    padding: 6px 12px;
    border-radius: 6px;
    font-weight: bold;
    cursor: pointer;
    font-size: 0.85rem;
}

.status-neutral {
    background: #f0f0f0;
    color: #666;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: bold;
}

.empty-state {
    text-align: center;
    padding: 3rem;
    color: #888;
}
</style>
