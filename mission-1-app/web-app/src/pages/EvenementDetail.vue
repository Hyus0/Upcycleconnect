<template>
  <main class="page-container">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" user-role="Particulier" :user-score="userScore" />

    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > EVENEMENTS > {{ evenement?.titre || "NULL" }}</p>
        <h1 class="hero-title1">{{ evenement?.titre || "Chargement..." }}</h1>
        <p class="classic-text">{{ evenement?.type || "Evenement" }}</p>
      </div>
      <button class="btn-secondary" type="button" @click="$router.back()">Retour</button>
    </header>

    <div v-if="loading" class="loading-state">Chargement de l'evenement...</div>

    <div v-else-if="evenement?.id" class="event-detail-layout">
      <section class="info-card">
        <h2>Informations</h2>
        <p class="event-description">{{ evenement.description || "Description NULL" }}</p>

        <div class="specs-grid">
          <div class="spec-item"><label>Date</label><p>{{ formatDate(evenement.date_evenement) }}</p></div>
          <div class="spec-item"><label>Adresse</label><p>{{ evenement.adresse || "NULL" }}</p></div>
          <div class="spec-item"><label>Ville</label><p>{{ evenement.ville || "NULL" }}</p></div>
          <div class="spec-item"><label>Code postal</label><p>{{ evenement.code_postal || "NULL" }}</p></div>
        </div>
      </section>

      <aside class="info-card event-actions">
        <h3>Organisateur</h3>
        <p>Equipe UpcycleConnect</p>
        <button class="btn-secondary" type="button" @click="contactOrganizer">Ouvrir une discussion</button>

        <button
          v-if="!isRegistered"
          class="btn-main-action"
          type="button"
          :disabled="busy"
          @click="register"
        >
          {{ busy ? "Inscription..." : "S'inscrire" }}
        </button>
        <button v-else class="btn-secondary" type="button" :disabled="busy" @click="unregister">
          {{ busy ? "Desinscription..." : "Se desinscrire" }}
        </button>
      </aside>
    </div>

    <div v-else class="state-card">Evenement introuvable.</div>
  </main>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import {
  fetchEvenement,
  fetchEvenementInscriptionStatus,
  joinEvenement,
  quitEvenement
} from "../services/publicApi";
import { startConversation } from "../services/messageService";

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const busy = ref(false);
const evenement = ref(null);
const isRegistered = ref(false);
const userScore = ref(0);

const isLoggedIn = computed(() => !!localStorage.getItem("userToken"));
const userName = computed(() => {
  const prenom = localStorage.getItem("userPrenom") || "";
  const nom = localStorage.getItem("userNom") || "";
  return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

function formatDate(value) {
  if (!value) return "NULL";
  return new Date(value).toLocaleString("fr-FR", {
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });
}

async function loadDetail() {
  loading.value = true;
  try {
    evenement.value = await fetchEvenement(route.params.id);
    const userId = Number(localStorage.getItem("userId"));
    if (userId) {
      const status = await fetchEvenementInscriptionStatus(route.params.id, userId);
      isRegistered.value = Boolean(status?.inscrit);
    }
  } catch (error) {
    console.error("Erreur detail evenement :", error);
    evenement.value = null;
  } finally {
    loading.value = false;
  }
}

async function register() {
  const userId = Number(localStorage.getItem("userId"));
  if (!userId || !localStorage.getItem("userToken")) {
    router.push("/connexion");
    return;
  }
  busy.value = true;
  try {
    await joinEvenement(evenement.value.id, userId);
    isRegistered.value = true;
  } catch (error) {
    alert(error.message || "Inscription impossible.");
  } finally {
    busy.value = false;
  }
}

async function unregister() {
  const userId = Number(localStorage.getItem("userId"));
  if (!userId) return;
  busy.value = true;
  try {
    await quitEvenement(evenement.value.id, userId);
    isRegistered.value = false;
  } catch (error) {
    alert(error.message || "Desinscription impossible.");
  } finally {
    busy.value = false;
  }
}

function contactOrganizer() {
  if (!localStorage.getItem("userToken")) {
    router.push("/connexion");
    return;
  }
  const conversation = startConversation({
    kind: "organisateur",
    targetId: "upcycleconnect",
    name: "Equipe UpcycleConnect",
    subject: evenement.value?.titre,
    contextId: evenement.value?.id,
    contextLabel: `Evenement - ${evenement.value?.titre}`
  });
  router.push({ path: "/messages", query: { conversation: conversation.id } });
}

onMounted(loadDetail);
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  padding: 20px;
  background: #f7f9f7;
  max-width: 1600px;
  margin: 0 auto;
}

.event-detail-layout {
  display: grid;
  grid-template-columns: 1.6fr minmax(300px, 420px);
  gap: 18px;
}

.info-card {
  border: 1px solid var(--border);
  border-radius: 14px;
  background: #fff;
  padding: 22px;
}

.event-description {
  color: var(--text-secondary);
  line-height: 1.7;
}

.specs-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  margin-top: 18px;
}

.spec-item {
  border: 1px solid var(--border);
  border-radius: 12px;
  background: #f8fbf8;
  padding: 14px;
}

.spec-item label {
  display: block;
  color: #8a958e;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
}

.spec-item p {
  margin: 6px 0 0;
}

.event-actions {
  display: grid;
  align-content: start;
  gap: 12px;
}

.event-actions .btn-secondary,
.event-actions .btn-main-action {
  width: 100%;
  justify-content: center;
}

@media (max-width: 900px) {
  .event-detail-layout,
  .specs-grid {
    grid-template-columns: 1fr;
  }
}
</style>
