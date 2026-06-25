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
                <p class="sidebar__category2">ACCUEIL > DÉPÔTS</p>
                <h1 class="hero-title1">DÉPOSER UN OBJET</h1>
                <p class="classic-text">
                    Simulez votre présence devant une borne et saisissez votre token de dépôt.
                </p>
            </div>

            <div class="header-actions">
                <button class="btn-secondary" @click="$router.back()">🠔 Retour</button>
            </div>
        </header>

        <div class="section-container">
            <div class="split-layout">
                
                <div class="form-card">
                    <h2 class="card-title">1. Présence physique (Simulation)</h2>
                    
                    <div class="form-group">
                        <label>Je me trouve devant la borne de :</label>
                        <select v-model="selectedSiteId">
                            <option value="" disabled>Choisir un site...</option>
                            <option v-for="site in sites" :key="site.id" :value="site.id">
                                {{ site.nom }} - {{ site.ville }}
                            </option>
                        </select>
                    </div>
                </div>

                <div class="form-card">
                    <h2 class="card-title">2. Validation du dépôt</h2>
                    
                    <div class="form-group">
                        <label>Token de dépôt unique</label>
                        <input 
                            type="text" 
                            v-model="token" 
                            placeholder="Saisissez le token (ex: ABCD123...)" 
                        />
                    </div>

                    <button class="btn-save" @click="handleDepot" :disabled="!selectedSiteId || !token">
                        Ouvrir le casier & Valider
                    </button>

                    <div v-if="successMessage" class="success-message">
                        {{ successMessage }}
                    </div>
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
const selectedSiteId = ref("");
const token = ref("");
const successMessage = ref("");

onMounted(async () => {
    try {
        const res = await fetch(`${API_URL}/sites`);
        if (res.ok) sites.value = await res.json();
    } catch (err) { console.error("Erreur chargement sites:", err); }
});

const handleDepot = async () => {
    successMessage.value = "";

    try {
        const response = await fetch(`${API_URL}/depot`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem("userToken")
            },
            body: JSON.stringify({ 
                token: token.value,
                site_id: Number(selectedSiteId.value)
            })
        });
        if (response.ok) {
            successMessage.value = "Dépôt validé ! Le casier se déverrouille. Redirection en cours...";
            token.value = "";
            selectedSiteId.value = "";
            
            setTimeout(() => {
                router.push('/profil/depots');
            }, 2500);
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

.success-message {
    margin-top: 1.5rem;
    padding: 1rem;
    background-color: #e8f5e9;
    color: #2e7d32;
    border-radius: 8px;
    text-align: center;
    font-weight: 700;
    border: 1px solid #c8e6c9;
}

@media (max-width: 768px) {
    .split-layout {
        grid-template-columns: 1fr;
    }
}
</style>