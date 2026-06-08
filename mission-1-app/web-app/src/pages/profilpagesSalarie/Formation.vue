<template>
    <div class="page-container">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ESPACE SALARIÉ > MES FORMATIONS
                </p>
                <h1 class="hero-title1">VOS FORMATIONS</h1>
                <p class="classic-text">
                    Gérez vos parcours d'apprentissage, vos modules et vos cours
                    dédiés à l'upcycling.
                </p>
            </div>
            <router-link
                to="/profil/createFormation"
                class="btn-main-action"
                style="text-decoration: none"
            >
                + Créer une formation
            </router-link>
        </header>

        <div class="section-container">
            <div v-if="loading" class="loading-state">
                Chargement de vos formations...
            </div>

            <table v-else-if="formations.length > 0" class="data-table">
                <thead>
                    <tr>
                        <th>FORMATION</th>
                        <th>TYPE</th>
                        <th>LIEU</th>
                        <th>STATUT</th>
                        <th>VALIDATION</th>
                        <th>DÉBUT</th>
                        <th>ACTIONS</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="form in formations" :key="form.id">
                        <td>
                            <strong>{{ form.titre }}</strong><br />
                            <small class="text-truncate">{{
                                form.description || "Aucune description"
                            }}</small>
                        </td>
                        <td>
                            <span class="material-tag">
                                {{ form.type || "Non défini" }}
                            </span>
                        </td>
                        <td>
                            <span class="stats-text">
                                {{ form.ville || "En ligne" }}
                            </span>
                        </td>
                        <td>
                            <span :class="form.statut === 'Ouvert' ? 'status-valid' : 'status-neutral'">
                                {{ form.statut || "En attente" }}
                            </span>
                        </td>
                        <td>
                            <span v-if="form.Est_valide === 'Valide'" class="status-valid">
                                VALIDÉ
                            </span>
                            <span v-else-if="form.Est_valide === 'Refuse'" class="status-refused">
                                REFUSÉ
                            </span>
                            <span v-else class="status-pending">
                                EN ATTENTE
                            </span>
                        </td>
                        <td>{{ formatDate(form.date_debut) }}</td>
                        <td class="actions-cell">
                            <button
                                class="btn-view"
                                type="button"
                                @click="goToFormation(form.id)"
                            >
                                Voir
                            </button>
                            <button
                                class="btn-modify"
                                type="button"
                                @click="goToModifyFormation(form.id)"
                            >
                                Modifier
                            </button>
                            <button
                                class="btn-remove"
                                type="button"
                                @click="removeFormation(form.id)"
                            >
                                Retirer
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>

            <div v-if="!loading && formations.length === 0" class="empty-state">
                <p>Vous n'avez pas encore créé de formations.</p>
                <router-link
                    to="/profil/formations/nouveau"
                    class="btn-secondary"
                    style="text-decoration: none"
                >
                    Créer ma première formation
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
const formations = ref([]);

const currentUserId = computed(() => {
    const storedId =
        sessionStorage.getItem("id") || sessionStorage.getItem("userId");
    return Number(storedId) || 0;
});

const formatDate = (dateString) => {
    if (!dateString) return "...";
    const date = new Date(dateString);
    return date.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit"
    }).replace(":", "h");
};

const goToFormation = (id) => {
    router.push({ name: "formation-detail", params: { id: id } });
};

const goToModifyFormation = (id) => {
    router.push({ name: "modify-formation", params: { id: id } });
};

const removeFormation = async (id) => {
    if (
        !confirm(
            "Voulez-vous vraiment supprimer cette formation ? Cette action est irréversible.",
        )
    )
        return;

    const token = sessionStorage.getItem("userToken");

    try {
        const res = await fetch(`${API_URL}/formation/${id}`, {
            method: "DELETE",
            headers: { Authorization: token },
        });

        if (res.ok) {
            formations.value = formations.value.filter((f) => f.id !== id);
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
        const res = await fetch(`${API_URL}/formations`, {
            method: "GET",
            headers: { Authorization: token },
        });

        if (res.ok) {
            const allFormations = (await res.json()) || [];
            
            formations.value = allFormations.filter(
                (f) => Number(f.id_formateur) === currentUserId.value,
            );
        }
    } catch (error) {
        console.error("Erreur formations :", error);
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

.status-refused {
  background: #fef2f2;
  color: #dc2626;
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