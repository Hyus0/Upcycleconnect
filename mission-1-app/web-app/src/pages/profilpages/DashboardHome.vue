<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL > TABLEAU DE BORD</p>
            <h1 class="hero-title1">Bonjour {{ prenom }} 👋</h1>
            <p class="classic-text">
                Voici un résumé de votre activité sur UpcycleConnect
            </p>
        </div>
        <router-link to="/profil/createAnnonce" class="btn-main-action">+ Déposer une annonce</router-link>
    </header>

    <div class="stats-grid">
        <div class="card card--score">
            <p class="tag-score">♻ UPCYCLING SCORE</p>

            <div class="score-value">
                {{ stats.total_points }} <span>pts</span>
            </div>

            <p class="score-level">Niveau : {{ stats.niveau }}</p>

            <div class="score-footer">
                <div class="mini-stat">
                    <strong>{{ stats.co2_total_evite_kg }} kg</strong><br />
                    CO2 évité
                </div>

                <div class="mini-stat">
                    <strong>{{ stats.nb_objets_recycles }}</strong
                    ><br />
                    Objets recyclés
                </div>

                <div class="mini-stat">
                    <strong>€ {{ stats.ressources_economisees }}</strong
                    ><br />
                    Économisé
                </div>
            </div>
        </div>
        <div class="card card--white">
            <div class="card-num">{{ annoncesActivesCount }}</div>
            <p class="text-dm">Annonces actives</p>
            <span class="badge badge--green" v-if="annoncesActivesCount > 0">+1 ce mois</span>
        </div>

        <div class="card card--white">
            <div class="card-num2">{{ depotsEnAttenteCount }}</div>
            <p class="text-dm">Dépôts en attente</p>
            <span class="badge badge--orange">
                {{ depotsEnAttenteCount > 0 ? 'EN COURS' : 'AUCUN' }}
            </span>
        </div>
    </div>

    <div class="section-container">
        <div class="section-header">
            <h2>Mes dernières annonces</h2>
            <div class="header-actions">
                <input
                    type="text"
                    placeholder="Rechercher..."
                    class="search-input"
                />
                <button class="btn-secondary">Tous statuts</button>
                <router-link to="/profil/createAnnonce" class="btn-main-action">+ Nouvelle annonce</router-link>
            </div>
        </div>
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
                <tr v-for="annonce in annonces.slice(0, 4)" :key="annonce.id">
                    <td>
                        <strong>{{ annonce.titre }}</strong><br>
                        <small>{{ annonce.type_materiau }}</small>
                    </td>

                    <td>
                        <span :class="annonce.type === 'Don' ? 'tag-don' : 'tag-vente'">
                            {{ annonce.type === 'Don' ? '🎁 DON' : '💰 VENTE ' + annonce.prix + '€' }}
                        </span>
                    </td>

                    <td>
                        <span :class="annonce.est_valide === 'Valide' ? 'status-valid' : 'status-pending'">
                            {{ annonce.est_valide === 'Valide' ? '✓ APPROUVÉ' : '⌛ EN ANALYSE' }}
                        </span>
                    </td>

                    <td>
                        <span class="status-neutral">
                            📍 {{ annonce.statut }}
                        </span>
                    </td>

                    <td>{{ formatDate(annonce.date_creation) }}</td>

                    <td class="actions-cell">
                        <button 
                            v-if="annonce.est_valide === 'Valide' && annonce.statut === 'Disponible'"
                            class="btn-plan"
                            @click="goToPlanning(annonce.id)"
                        >
                            Planifier dépôt
                        </button>
                        
                        <button
                            class="btn-view"
                            @click="goToAnnonce(annonce.id)"
                        >
                            Voir
                        </button>
                        
                        <button 
                            v-if="annonce.est_valide === 'En attente'" 
                            class="btn-modify" 
                            @click="goToModify(annonce.id)"
                        >
                            Modifier
                        </button>

                        <button class="btn-remove" @click="removeAnnonce(annonce.id)">Retirer</button>
                    </td>
                </tr>
            </tbody>
        </table>
        <p v-if="annonces.length === 0" class="empty-msg">Vous n'avez pas encore déposé d'annonces.</p>
    </div>

    <div class="section-container">
        <div class="section-header">
            <h2>📅 Mon planning — semaine du 23 fév.</h2>
            <button class="btn-secondary">Vue mensuelle</button>
        </div>
        <div class="planning-row">
            <div class="day-box active">
                <span class="day-num">23</span>
                <div class="event event--green">Atelier bois - 14h</div>
            </div>
            <div class="day-box"><span class="day-num">24</span></div>
            <div class="day-box active">
                <span class="day-num">25</span>
                <div class="event event--orange">Dépôt conteneur 11h</div>
            </div>
            <div class="day-box"><span class="day-num">26</span></div>
        </div>
    </div>
    <div class="end-grid">
        <div class="section-container-tips">
            <p class="tag-vente">💡 Conseil du jour</p>

            <div class="section-header">
                <h2>Transformer un vieux jean en sac</h2>
            </div>
            <p>
                Apprenez à confectionner un sac tote en 30 minutes avec un jean
                usé. Matériel nécessaire : aiguille, fil, ciseaux.
            </p>
            <button class="btn-view">Lire la suite →</button>
        </div>
        <div class="section-container-tips">
            <p class="tag-don">🔔 Notification</p>

            <div class="section-header">
                <h2>Votre dépôt a été récupéré !</h2>
            </div>
            <p>
                La chaise vintage que vous avez déposée le 12 fév. a été
                récupérée par un artisan. +50 points Upcycling Score !
            </p>
            <button class="btn-view">Voir le projet →</button>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";

import { useRouter } from 'vue-router';

const router = useRouter();

const prenom = ref(localStorage.getItem("userPrenom") || "Invité");
const annonces = ref([]);

const stats = ref({
    total_points: 0,
    niveau: "Chargement...",
    co2_total_evite_kg: 0,
    nb_objets_recycles: 0,
    ressources_economisees: 0,
});

const formatDate = (dateString) => {
    if (!dateString) return "...";
    const date = new Date(dateString);
    return date.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
    });
};

const goToAnnonce = (id) => {
    router.push({ name: 'see-annonce', params: { id: id } });
};

const goToModify = (id) => {
    router.push({ name: 'modification-annonce', params: { id: id } });
};

const removeAnnonce = (id) => {
    console.log("Demande de suppression pour :", id);
};

const annoncesActivesCount = computed(() => {
    return annonces.value.filter(a => a.est_valide === 'Valide' && a.statut === 'Disponible').length;
});

const depotsEnAttenteCount = computed(() => {
    return annonces.value.filter(a => a.statut === 'Reserve').length;
});

onMounted(async () => {
    const id = localStorage.getItem("userId");
    const token = localStorage.getItem("userToken");

    if (!id || !token) return;

    try {
            const response = await fetch(
                `http://localhost:8081/users/${id}/stats`,
                {
                    headers: { Authorization: token },
                },
            );
            if (response.ok) {
                const data = await response.json();
                stats.value = data;
                if (data.total_points !== undefined) {
                    localStorage.setItem("userScore", data.total_points.toString());
                    window.dispatchEvent(new Event("auth-change"));
                }
            }
        } catch (error) {
            console.error("Erreur stats :", error);
        }

    try {
        const resAnnonces = await fetch(
            `http://localhost:8081/users/${id}/annonces`,
            {
                headers: { Authorization: token },
            },
        );
        if (resAnnonces.ok) {
            annonces.value = await resAnnonces.json();
        }
    } catch (error) {
        console.error("Erreur annonces :", error);
    }
});
</script>
