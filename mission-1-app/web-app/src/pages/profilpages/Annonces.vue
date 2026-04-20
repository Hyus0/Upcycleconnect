<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL > ANNONCES</p>
            <h1 class="hero-title1">VOS ANNONCES</h1>
            <p class="classic-text">
                Voici un résumé de vos annonces sur UpcycleConnect
            </p>
        </div>
        <button class="btn-main-action">+ Déposer une annonce</button>
    </header>

    <div class="section-container">
        <table class="data-table">
            <thead>
                <tr>
                    <th>OBJET</th>
                    <th>CATÉGORIE</th>
                    <th>TYPE</th>
                    <th>STATUT</th>
                    <th>DATE</th>
                    <th>ACTIONS</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="annonce in annonces" :key="annonce.id">
                    <td>{{ annonce.titre }}</td>
                    <td>{{ annonce.etat_objet }}</td>
                    <td>
                        <span :class="annonce.type === 'Don' ? 'tag-vente' : 'tag-don'">
                            {{ annonce.type === 'Don' ? 'DON' : 'VENTE ' + annonce.prix + '€' }}
                        </span>
                    </td>
                    <td>
                        <span :class="annonce.est_valide === 'Validé' ? 'status-valid' : 'status-pending'">
                            {{ annonce.est_valide === 'Validé' ? '✓ VALIDÉE' : '⌛ EN ATTENTE' }}
                        </span>
                    </td>
                    <td>{{ formatDate(annonce.date_creation) }}</td>
                    <td class="actions-cell">
                        <button class="btn-view" @click="goToAnnonce(annonce.id)">Voir</button>
                        <button class="btn-remove" @click="removeAnnonce(annonce.id)">Retirer</button>
                        <button 
                          v-if="annonce.est_valide === 'Non validé'" 
                          class="btn-modify" 
                          @click="goToAnnonce(annonce.id)"
                        >
                          Modifier
                        </button>                    </td>
                </tr>
            </tbody>
        </table>

        <p v-if="annonces.length === 0" class="empty-msg">Vous n'avez pas encore déposé d'annonces.</p>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from 'vue-router';

const prenom = ref(localStorage.getItem("userPrenom") || 'Invité');
const annonces = ref([]); 

const formatDate = (dateString) => {
    if (!dateString) return "...";
    const date = new Date(dateString);
    return date.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
    });
};

const router = useRouter();

const goToAnnonce = (id) => {
    router.push({ name: 'modification-annonce', params: { id: id } });
};

const removeAnnonce = (id) => {
    console.log("Demande de suppression pour :", id);
};

onMounted(async () => {
    const id = localStorage.getItem("userId");
    const token = localStorage.getItem("userToken");

    if (!id || !token) return;

    try {
        const resAnnonces = await fetch(`http://localhost:8081/users/${id}/annonces`, {
            headers: { "Authorization": token }
        });
        if (resAnnonces.ok) {
            annonces.value = await resAnnonces.json();
        }
    } catch (error) {
        console.error("Erreur annonces :", error);
    }
});
</script>