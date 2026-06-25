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
                <p class="sidebar__category2">ACCUEIL > RÉCLAMER UN OBJET</p>
                <h1 class="hero-title1">RÉCUPÉRER UN OBJET</h1>
                <p class="classic-text">
                    Sélectionnez le site de dépôt et saisissez le token de retrait.
                </p>
            </div>

            <div class="header-actions">
                <button class="btn-secondary" @click="$router.back()">🠔 Retour</button>
            </div>
        </header>

        <div class="section-container">
            <div class="split-layout">
                <div class="form-card">
                    <h2 class="card-title">1. Sélection du lieu</h2>
                    
                    <div class="form-group">
                        <label>Site de récupération</label>
                        <select v-model="selectedSiteId" @change="fetchConteneurs">
                            <option value="" disabled>Choisir un site...</option>
                            <option v-for="site in sites" :key="site.id" :value="site.id">
                                {{ site.nom }} - {{ site.ville }} + {{site.id}}
                            </option>
                        </select>
                    </div>

                    <div class="form-group" v-if="conteneurs.length > 0">
                        <label>Conteneur concerné</label>
                        <select v-model="selectedConteneurId">
                            <option value="" disabled>Choisir un conteneur...</option>
                            <option v-for="c in conteneurs" :key="c.id" :value="c.id">
                                Conteneur #{{ c.id }} - {{ c.statut }}
                            </option>
                        </select>
                    </div>
                </div>

                <div class="form-card">
                    <h2 class="card-title">2. Validation</h2>
                    
                    <div class="form-group">
                        <label>Token de retrait</label>
                        <input type="text" v-model="token" placeholder="Saisissez le token unique" />
                    </div>

                    <button class="btn-save" @click="handleClaim" :disabled="!selectedConteneurId || !token">
                        Valider la récupération
                    </button>
                </div>
            </div>
        </div>
    </main>
    <SiteFooter />
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import SiteNavbar from '../components/SiteNavbar.vue';
import SiteFooter from "../components/SiteFooter.vue";

const router = useRouter();
const API_URL = "/go";

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});
const userRole = computed(() => sessionStorage.getItem("userRole") || "Salarie");
const userScore = ref(sessionStorage.getItem("userScore") || 0); 

const sites = ref([]);
const conteneurs = ref([]);
const selectedSiteId = ref("");
const selectedConteneurId = ref("");
const token = ref("");

onMounted(async () => {
    try {
        const res = await fetch(`${API_URL}/sites`);
        if (res.ok) sites.value = await res.json();
    } catch (err) { console.error("Erreur chargement sites:", err); }
});

const fetchConteneurs = async () => {
    if (!selectedSiteId.value) return;
    try {
        const res = await fetch(`${API_URL}/sites/${selectedSiteId.value}/conteneurs`);
        if (res.ok) conteneurs.value = await res.json();
    } catch (err) { console.error("Erreur conteneurs:", err); }
};

const handleClaim = async () => {
    try {
        const response = await fetch(`${API_URL}/annonces/${selectedSiteId.value}/retirer`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem("userToken")
            },
            body: JSON.stringify({ token: token.value })
        });

        if (response.ok) {
            alert("Récupération validée avec succès !");
            token.value = "";
            selectedConteneurId.value = "";
        } else {
            const msg = await response.text();
            alert("Erreur : " + msg);
        }
    } catch (err) {
        console.error("Erreur de requête:", err);
        alert("Impossible de joindre le serveur.");
    }
};
</script>

<style scoped>
.page-main-content {
    min-height: 100vh;
    padding: 20px;
    background: #f7f9f7;
    max-width: 1600px; 
    margin: 0 auto;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #eee;
}

.header-left {
    display: flex;
    flex-direction: column;
}

.sidebar__category2 {
    font-size: 0.65rem;
    color: #8fa396;
    letter-spacing: 1px;
    margin: 0 0 0.5rem 0;
    text-transform: uppercase;
}

.hero-title1 {
    font-size: 2rem;
    font-weight: 800;
    margin: 1.5rem 0 0.5rem; 
    color: #1a1a1a;
}

.classic-text {
    color: #666;
    margin: 0;
}

.header-actions {
    display: flex;
    align-items: center;
}

.btn-secondary {
    padding: 8px 16px;
    border-radius: 10px;
    border: 1px solid #ddd;
    background: white;
    cursor: pointer;
    font-weight: 500;
}

.section-container {
    margin-top: 1rem;
}

.split-layout { 
    display: grid; 
    grid-template-columns: 1fr 1fr; 
    gap: 2rem; 
}

.form-card {
    background: #ffffff;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
    box-shadow: 0 4px 12px rgba(0,0,0,0.03);
}

.card-title { 
    font-size: 1.2rem; 
    font-weight: 700; 
    margin-top: 0;
    margin-bottom: 1.5rem; 
    color: #1a1a1a;
}

.form-group { 
    display: flex; 
    flex-direction: column; 
    gap: 0.5rem; 
    margin-bottom: 1.5rem; 
}

.form-group label { 
    font-size: 0.95rem; 
    font-weight: 600; 
    color: #333; 
}

input, select {
    padding: 0.9rem;
    border: 1px solid #ddd;
    border-radius: 10px;
    width: 100%;
    box-sizing: border-box;
    font-family: inherit;
    font-size: 0.95rem;
    transition: border-color 0.2s;
    background-color: #fcfdfc;
}

input:focus, select:focus {
    outline: none;
    border-color: #2d7a4f;
}

.btn-save {
    background-color: #2d7a4f;
    color: white;
    padding: 1rem;
    border: none;
    border-radius: 12px;
    font-weight: 700;
    cursor: pointer;
    width: 100%;
    transition: background 0.2s;
}

.btn-save:hover:not(:disabled) {
    background-color: #246343;
}

.btn-save:disabled { 
    background-color: #ccc; 
    cursor: not-allowed;
}

@media (max-width: 768px) {
    .split-layout {
        grid-template-columns: 1fr;
    }
}
</style>