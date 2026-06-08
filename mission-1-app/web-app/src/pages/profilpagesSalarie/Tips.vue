<template>
  <div class="page-container">
    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ESPACE SALARIÉ > MES TIPS</p>
        <h1 class="hero-title1">VOS TIPS & ASTUCES</h1>
        <p class="classic-text">
          Gérez vos tutoriels, vos conseils d'upcycling et choisissez à quel public ils s'adressent.
        </p>
      </div>
      <router-link to="/profil/createTips" class="btn-main-action" style="text-decoration: none;">
        + Créer un tip
      </router-link>
    </header>

    <div class="section-container">
      <div v-if="loading" class="loading-state">
        Chargement de vos tips...
      </div>

      <table v-else-if="tips.length > 0" class="data-table">
        <thead>
          <tr>
            <th>TIP</th>
            <th>CIBLE</th>
            <th>VIDÉO</th>
            <th>STATUT</th>
            <th>DATE</th>
            <th>ACTIONS</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="tip in tips" :key="tip.id">
            <td>
              <strong>{{ tip.titre }}</strong><br>
              <small class="text-truncate">{{ tip.description || "Aucune description" }}</small>
            </td>
            <td>
              <span class="material-tag">
                {{ tip.role_cible || "Tous" }}
              </span>
            </td>
            <td>
              <span class="stats-text">{{ tip.video_url ? '🎥 Incluse' : 'Aucune' }}</span><br>
            </td>
            <td>
              <span :class="tip.actif ? 'status-valid' : 'status-neutral'">
                {{ tip.actif ? 'ACTIF' : 'INACTIF' }}
              </span>
            </td>
            <td>{{ formatDate(tip.date_creation) }}</td>
            <td class="actions-cell">
              <button class="btn-view" type="button" @click="goToTip(tip.id)">Voir</button>
              <button class="btn-modify" type="button" @click="goToModifyTip(tip.id)">Modifier</button>
              <button class="btn-remove" type="button" @click="removeTip(tip.id)">Retirer</button>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="!loading && tips.length === 0" class="empty-state">
        <p>Vous n'avez pas encore publié de tips d'upcycling.</p>
        <router-link
          to="/profil/tips/nouveau"
          class="btn-secondary"
          style="text-decoration: none"
        >
          Créer mon premier tip
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
const tips = ref([]);

const currentUserId = computed(() => {
  const storedId = sessionStorage.getItem("id") || sessionStorage.getItem("userId");
  return Number(storedId) || 0;
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

const goToTip = (id) => {
  router.push({ name: "conseil-detail", params: { id: id } });
};

const goToModifyTip = (id) => {
  router.push({ name: "modify-tips", params: { id: id } });
};

const removeTip = async (id) => {
  if (!confirm("Voulez-vous vraiment supprimer ce tip ? Cette action est irréversible.")) return;

  const token = sessionStorage.getItem("userToken");

  try {
    const res = await fetch(`${API_URL}/tips/${id}`, {
      method: "DELETE",
      headers: { Authorization: token },
    });

    if (res.ok) {
      tips.value = tips.value.filter((t) => t.id !== id);
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
    const res = await fetch(`${API_URL}/tips`, {
      method: "GET",
      headers: { Authorization: token },
    });
    
    if (res.ok) {
      const allTips = (await res.json()) || [];
      
      tips.value = allTips.filter((t) => t.id_createur === currentUserId.value);
    }
  } catch (error) {
    console.error("Erreur tips :", error);
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

.btn-main-action:hover {
  background-color: #246343;
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
.status-pending {
  background: #fff4e6;
  color: #cc6600;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: bold;
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