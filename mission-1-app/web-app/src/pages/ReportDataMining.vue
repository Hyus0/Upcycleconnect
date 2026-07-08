<template>
  <main class="analytics-dashboard">
      <SiteNavbar
        :is-authenticated="isLoggedIn"
        :user-name="userName"
        user-role="Particulier"
        :user-score="userScore"
      />
    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">DASHBOARD > DATA MINING & IA</p>
        <h1 class="hero-title1">Analytique & Prédictions</h1>
        <p class="classic-text">
          Module de suivi des activités et recommandations d'intelligence artificielle.
        </p>
      </div>
      <button class="btn-secondary" @click="fetchData" :disabled="loading">
        {{ loading ? "Actualisation..." : "Actualiser" }}
      </button>
    </header>

    <div v-if="loading" class="state-card">
      Récupération des données depuis l'API Go...
    </div>

    <div v-else class="dashboard-grid">
      
      <section class="card">
        <h2 class="section-title">📊 Répartition des Acteurs</h2>
        <p class="subtitle">Analyse démographique des profils inscrits.</p>
        
        <div class="chart-container mt-6">
          <div v-for="acteur in data.acteurs" :key="acteur.role" class="bar-row">
            <div class="bar-label">
              <strong>{{ acteur.role }}</strong>
              <span>{{ acteur.count }} membres</span>
            </div>
            <div class="bar-track">
              <div class="bar-fill" :style="{ width: getPercentage(acteur.count, totalActeurs) + '%' }"></div>
            </div>
          </div>
        </div>
      </section>

      <section class="card">
        <h2 class="section-title">Succès des Prestations</h2>
        <p class="subtitle">Répartition des annonces par type d'échange.</p>
        
        <div class="chart-container mt-6">
          <div v-for="prest in data.prestations" :key="prest.type_item" class="bar-row">
            <div class="bar-label">
              <strong>Annonce : {{ prest.type_item }}</strong>
              <span>{{ prest.count }} transactions</span>
            </div>
            <div class="bar-track">
              <div class="bar-fill fill-orange" :style="{ width: getPercentage(prest.count, totalPrestations) + '%' }"></div>
            </div>
          </div>
        </div>
      </section>

      <section class="card ml-card">
        <h2 class="section-title">Prédictions Machine Learning</h2>
        <p class="subtitle">Modèle de classification entraîné en Python sur les 500+ utilisateurs du dataset.</p>
        
        <div class="table-wrapper mt-6">
          <table class="data-table">
            <thead>
              <tr>
                <th>Utilisateur</th>
                <th>Rôle actuel</th>
                <th>Prestation Prédite (Besoin futur)</th>
                <th>Confiance de l'IA</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="pred in data.predictions" :key="pred.id_utilisateur">
                <td><strong>{{ pred.prenom }} {{ pred.nom }}</strong></td>
                <td><span class="badge">{{ pred.role }}</span></td>
                <td><span class="badge badge--green">{{ pred.prestation_recommandee }}</span></td>
                <td>
                  <div class="confiance-wrapper">
                    <div class="confiance-track">
                      <div class="confiance-fill" :style="{ width: pred.confiance + '%' }"></div>
                    </div>
                    <span class="confiance-text">{{ pred.confiance }}%</span>
                  </div>
                </td>
              </tr>
              <tr v-if="data.predictions.length === 0">
                <td colspan="4" class="text-center py-4">Aucune prédiction générée. Lancez le script Python !</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

    </div>
  </main>
  <SiteFooter />
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';

import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));

const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});


const API_URL = "/go";
const loading = ref(true);
const data = ref({
  acteurs: [],
  prestations: [],
  predictions: []
});

const totalActeurs = computed(() => {
  return data.value.acteurs.reduce((acc, curr) => acc + curr.count, 0);
});

const totalPrestations = computed(() => {
  return data.value.prestations.reduce((acc, curr) => acc + curr.count, 0);
});

const getPercentage = (value, total) => {
  if (total === 0) return 0;
  return Math.round((value / total) * 100);
};

const fetchData = async () => {
  loading.value = true;
  try {
    const res = await fetch(`${API_URL}/api/analytics`, {
      headers: { Authorization: `Bearer ${sessionStorage.getItem("userToken") || ""}` }
    });
    if (res.ok) {
      data.value = await res.json();
    }
  } catch (err) {
    console.error("Erreur Analytics:", err);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.analytics-dashboard {
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
}

.sidebar__category2 { 
  margin: 0;
  color: #a0ada7;
  font-size: 0.75rem;
  font-weight: bold;
  letter-spacing: 1px;
}

.hero-title1 { 
  font-size: 2.2rem; 
  font-weight: 900; 
  color: #122018; 
  margin: 0.5rem 0 0 0; 
}

.classic-text {
  color: #6d7b72;
  margin-top: 10px;
}

.btn-secondary {
  padding: 10px 20px;
  border-radius: 10px;
  border: 1px solid #ddd;
  background: white;
  cursor: pointer;
  font-weight: bold;
}

.state-card {
  padding: 40px;
  text-align: center;
  border: 1px dashed #cfe0d4;
  border-radius: 16px;
  color: #63746a;
  background: #fbfdfb;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.card {
  background: white;
  border-radius: 16px;
  border: 1px solid #e5ede7;
  padding: 30px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
}

.ml-card {
  grid-column: 1 / -1; 
}

.section-title {
  margin: 0;
  font-size: 1.4rem;
  font-weight: 800;
  color: #1a1a1a;
}

.subtitle {
  color: #8fa396;
  font-size: 0.9rem;
  margin-top: 5px;
}

.chart-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.bar-row {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.bar-label {
  display: flex;
  justify-content: space-between;
  font-size: 0.95rem;
  color: #333;
}

.bar-track {
  width: 100%;
  height: 12px;
  background: #f0f4f1;
  border-radius: 10px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: #2d7a4f;
  border-radius: 10px;
  transition: width 1s ease-out;
}

.fill-orange {
  background: #f59e0b;
}

.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th, .data-table td {
  padding: 14px;
  text-align: left;
  border-bottom: 1px solid #f0f4f1;
}

.data-table th {
  color: #8fa396;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.8rem;
  letter-spacing: 1px;
}

.badge {
  background: #f0f4f1;
  color: #555;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: bold;
}

.badge--green {
  background: #d1e7dd;
  color: #0f5132;
}

.confiance-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
}

.confiance-track {
  width: 100px;
  height: 8px;
  background: #eee;
  border-radius: 4px;
}

.confiance-fill {
  height: 100%;
  background: #3b82f6;
  border-radius: 4px;
}

.confiance-text {
  font-size: 0.85rem;
  font-weight: 600;
  color: #3b82f6;
}
</style>