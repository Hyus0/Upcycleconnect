<template>
  <div class="page-container">
    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > MES PROJETS</p>
        <h1 class="hero-title1">VOS PROJETS</h1>
        <p class="classic-text">
          Gérez vos créations d'upcycling, suivez votre impact CO2 et
          votre visibilité auprès de la communauté.
        </p>

        <div v-if="!loading" class="vitrine-counter mt-3">
          <span class="status-valid" :class="{ 'limit-reached': projetsEnLigne >= maxProjets }">
            Vitrine publique : {{ projetsEnLigne }} / {{ maxProjets }} projets en ligne
          </span>
          <span v-if="projetsEnLigne >= maxProjets && !isPremium" class="upgrade-text">
            <router-link to="/abonnement" style='text-decoration:none;'>Passez Pro pour débloquer 20 emplacements</router-link>
          </span>
        </div>
      </div>

      <router-link 
        v-if="projetsEnLigne < maxProjets"
        to="/profil/createProjet" 
        class="btn-main-action" 
        style="text-decoration: none;"
      >
        + Créer un projet
      </router-link>
      <button 
        v-else 
        class="btn-main-action disabled-btn" 
        disabled
        title="Limite de vitrine atteinte"
      >
        Vitrine pleine
      </button>
    </header>

    <div class="section-container">
      <div v-if="loading" class="loading-state">
        Chargement de vos projets...
      </div>

      <table v-else-if="projets.length > 0" class="data-table">
        <thead>
          <tr>
            <th>PROJET</th>
            <th>IMPACT CO2</th>
            <th>ENGAGEMENT</th>
            <th>STATUT</th>
            <th>DATE</th>
            <th>ACTIONS</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="projet in projets" :key="projet.id">
            <td>
              <strong>{{ projet.titre }}</strong><br>
              <small class="text-truncate">{{ projet.description_courte || "Aucune description" }}</small>
            </td>
            <td>
              <span class="material-tag">
                {{ projet.co2_evite_kg || 0 }} kg
              </span>
            </td>
            <td>
              <span class="stats-text">{{ projet.nb_vues || 0 }} vues, {{ projet.nb_likes || 0 }} likes</span><br>
            </td>
            <td>
              <span :class="projet.visible_public ? 'status-valid' : 'status-neutral'">
                {{ projet.visible_public ? 'PUBLIC' : 'PRIVÉ' }}
              </span>
            </td>
            <td>{{ formatDate(projet.date_creation) }}</td>
            <td class="actions-cell">
              <button class="btn-view" type="button" @click="goToProjet(projet.id)">Voir</button>
              <button class="btn-modify" type="button" @click="goToModify(projet.id)">Modifier</button>
              <button class="btn-remove" type="button" @click="removeProjet(projet.id)">Retirer</button>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="!loading && projets.length === 0" class="empty-state">
        <p>Vous n'avez pas encore publié de projets d'upcycling.</p>
        <router-link
          to="/profil/createProjet"
          class="btn-secondary"
          style="text-decoration: none"
        >
          Créer mon premier projet
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const API_URL = "http://localhost:8081";
const loading = ref(true);
const projets = ref([]);

const isPremium = ref(false);

const currentUserId = computed(() => {
  const storedId = sessionStorage.getItem("id") || sessionStorage.getItem("userId");
  return Number(storedId) || 0;
});

const projetsEnLigne = computed(() => {
  return projets.value.filter(p => p.visible_public).length;
});

const maxProjets = computed(() => {
  return isPremium.value ? 3 : 2;
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

const goToProjet = (id) => {
  router.push({ name: "projet-detail", params: { id: id } });
};

const goToModify = (id) => {
  router.push({ name: "modify-projet", params: { id: id } });
};

const removeProjet = async (id) => {
  if (!confirm("Voulez-vous vraiment supprimer ce projet ? Cette action est irréversible.")) return;

  const token = sessionStorage.getItem("userToken");

  try {
    const res = await fetch(`${API_URL}/projets/${id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
    });

    if (res.ok) {
      projets.value = projets.value.filter((p) => p.id !== id);
    } else {
      const msg = await res.text();
      alert("Erreur lors de la suppression : " + msg);
    }
  } catch (e) {
    console.error("Erreur réseau :", e);
    alert("Impossible de joindre le serveur.");
  }
};

onMounted(async () => {
  if (currentUserId.value === 0) return;
  const token = sessionStorage.getItem("userToken");

  loading.value = true;
  try {
    const [resProjets, resSub] = await Promise.all([
        fetch(`${API_URL}/users/${currentUserId.value}/projets`, { headers: { Authorization: `Bearer ${token}` } }),
        fetch(`${API_URL}/users/${currentUserId.value}/abonnement`, { headers: { Authorization: `Bearer ${token}` } })
    ]);

    if (resProjets.ok) {
      projets.value = (await resProjets.json()) || [];
    }
    
    if (resSub.ok) {
      const sub = await resSub.json();
      isPremium.value = sub.is_premium || false;
    }
  } catch (error) {
    console.error("Erreur projets/abonnement :", error);
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
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

.btn-main-action {
    display: inline-flex;
    align-items: center;
    background: #2d7a4f;
    color: white;
    padding: 10px 20px;
    border-radius: 10px;
    text-decoration: none;
    font-weight: bold;
    border: none;
    cursor: pointer;
}

.btn-main-action:hover:not(.disabled-btn) {
  background-color: #246343;
}

.disabled-btn {
  opacity: 0.5;
  cursor: not-allowed !important;
  background-color: #6d7b72;
}

.section-container {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e5ede7;
  padding: 20px;
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}
.data-table th, .data-table td {
  padding: 12px;
  border-bottom: 1px solid #eee;
}

.material-tag {
  background: #e9f5ed;
  color: #1e5636;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 700;
}

.stats-text {
  font-size: 0.85rem;
  color: #555;
}

.text-truncate {
  display: inline-block;
  max-width: 250px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #6d7b72;
}

.loading-state {
  text-align: center;
  padding: 2rem;
  color: #8fa396;
  font-style: italic;
}

.status-neutral {
  background: #f5f5f5;
  color: #666;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: bold;
}
.status-valid {
  background: #e9f5ed;
  color: #1e5636;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: bold;
}

.vitrine-counter {
  display: flex;
  align-items: center;
  gap: 12px;
}
.mt-3 {
  margin-top: 1rem;
}
.limit-reached {
  background: #fff4e6 !important;
  color: #cc6600 !important;
}
.upgrade-text a {
  font-size: 0.85rem;
  font-weight: bold;
  color: #2d7a4f;
  text-decoration: underline;
}

.actions-cell {
  display: flex;
  gap: 8px;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: #888;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.btn-secondary {
  display: inline-block;
  margin-top: 12px;
  padding: 8px 16px;
  border-radius: 10px;
  border: 1px solid #ddd;
  background: white;
  color: #1a1a1a;
  cursor: pointer;
  font-weight: 500;
  transition: 0.2s;
}
.btn-secondary:hover {
  background: #f0f4f1;
}
</style>