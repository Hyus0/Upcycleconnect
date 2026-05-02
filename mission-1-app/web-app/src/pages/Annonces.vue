<template>
  <main class="public-dashboard">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <header class="content-header annonces-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > CATALOGUE</p>
        <h1 class="hero-title1">Annonces disponibles</h1>
        <p class="classic-text">
          Consulte les objets proposes sur UpcycleConnect, filtre par type et trouve rapidement ce qui est reutilisable.
        </p>
      </div>
      <RouterLink class="btn-main-action" to="/profil/annonces">+ Deposer une annonce</RouterLink>
    </header>

    <div class="stats-grid annonces-stats">
      <div class="card card--score">
        <p class="tag-score">CATALOGUE PUBLIC</p>
        <div class="score-value">{{ filteredAnnonces.length }} <span>annonces</span></div>
        <p class="score-level">{{ sourceLabel }}</p>
        <div class="score-footer">
          <div class="mini-stat">
            <strong>{{ donationsCount }}</strong><br />Dons
          </div>
          <div class="mini-stat">
            <strong>{{ salesCount }}</strong><br />Ventes
          </div>
          <div class="mini-stat">
            <strong>{{ availableCount }}</strong><br />Disponibles
          </div>
        </div>
      </div>
      <div class="card card--white">
        <div class="card-num">{{ donationsCount }}</div>
        <p class="text-dm">Objets gratuits</p>
        <span class="badge badge--green">DON</span>
      </div>
      <div class="card card--white">
        <div class="card-num2">{{ salesCount }}</div>
        <p class="text-dm">Objets en vente</p>
        <span class="badge badge--orange">VENTE</span>
      </div>
    </div>

    <section class="section-container">
      <div class="section-header">
        <div>
          <h2>Catalogue des annonces</h2>
          <p class="classic-text">{{ loading ? "Chargement..." : `${filteredAnnonces.length} resultat${filteredAnnonces.length > 1 ? "s" : ""}` }}</p>
        </div>
        <div class="header-actions">
          <input
            v-model="filters.search"
            type="search"
            placeholder="Rechercher..."
            class="search-input"
          />
          <select v-model="filters.type" class="btn-secondary">
            <option value="">Tous types</option>
            <option value="don">Don</option>
            <option value="vente">Vente</option>
          </select>
          <select v-model="filters.status" class="btn-secondary">
            <option value="">Tous statuts</option>
            <option value="en ligne">En ligne</option>
            <option value="reserve">Reserve</option>
            <option value="vendu">Vendu</option>
          </select>
        </div>
      </div>

      <div v-if="loading" class="state-card">Chargement des annonces...</div>
      <div v-else-if="filteredAnnonces.length === 0" class="state-card">
        Aucune annonce ne correspond a ces filtres.
      </div>

      <table v-else class="data-table annonces-table">
        <thead>
          <tr>
            <th>Objet</th>
            <th>Localisation</th>
            <th>Type</th>
            <th>Statut</th>
            <th>Prix</th>
            <th>Date</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="annonce in filteredAnnonces" :key="annonce.id">
            <td>
              <strong>{{ displayValue(annonce.titre) }}</strong>
              <span class="table-subtext">{{ displayValue(annonce.description) }}</span>
            </td>
            <td>
              {{ displayValue(annonce.ville) }}
              <span class="table-subtext">{{ displayValue(annonce.code_postal) }}</span>
            </td>
            <td>
              <span :class="annonce.type === 'vente' ? 'tag-vente' : 'tag-don'">
                {{ displayValue(annonce.type).toUpperCase() }}
              </span>
            </td>
            <td>
              <span :class="statusClass(annonce.statut)">
                {{ displayValue(annonce.statut || 'en ligne').toUpperCase() }}
              </span>
            </td>
            <td>{{ formatPrice(annonce.prix, annonce.type) }}</td>
            <td>{{ formatDate(annonce.date_creation) }}</td>
            <td class="actions-cell">
              <button class="btn-view" type="button">Voir</button>
              <button
                :class="isInCartItem(annonce.id) ? 'btn-remove' : 'btn-modify'"
                type="button"
                @click="toggleCart(annonce)"
              >
                {{ isInCartItem(annonce.id) ? "Retirer" : "Ajouter au panier" }}
              </button>
              <button class="btn-view" type="button" @click="contactSeller(annonce)">
                Contacter
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </section>
  </main>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import { fetchAnnonces } from "../services/annoncesApi";
import { addToCart, isInCart, onCartChange, removeFromCart } from "../services/cartService";
import { startConversation } from "../services/messageService";

const router = useRouter();
const loading = ref(true);
const source = ref("api");
const annonces = ref([]);
const isLoggedIn = computed(() => {
  return !!localStorage.getItem("userToken");
});

const userName = computed(() => {
  const prenom = localStorage.getItem("userPrenom") || "";
  const nom = localStorage.getItem("userNom") || "";
  
  return (prenom || nom) ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const filters = reactive({
  search: "",
  type: "",
  status: ""
});

const filteredAnnonces = computed(() => {
  const search = filters.search.trim().toLowerCase();
  return annonces.value.filter((annonce) => {
    const matchesSearch =
      !search ||
      [
        annonce.titre,
        annonce.description,
        annonce.ville,
        annonce.adresse,
        annonce.etat_objet
      ]
        .join(" ")
        .toLowerCase()
        .includes(search);

    const matchesType = !filters.type || annonce.type === filters.type;
    const matchesStatus =
      !filters.status || (annonce.statut || "").toLowerCase() === filters.status.toLowerCase();

    return matchesSearch && matchesType && matchesStatus;
  });
});

const donationsCount = computed(() => filteredAnnonces.value.filter((item) => item.type === "don").length);
const salesCount = computed(() => filteredAnnonces.value.filter((item) => item.type === "vente").length);
const availableCount = computed(() =>
  filteredAnnonces.value.filter((item) => ["", "en ligne", "validee"].includes((item.statut || "en ligne").toLowerCase())).length
);
const sourceLabel = computed(() =>
  source.value === "api" ? "Donnees issues de l'API annonces" : "Aucune donnee disponible pour le moment"
);

let stopCartSync = null;

function displayValue(value) {
  return value === null || value === undefined || value === "" ? "NULL" : value;
}

function statusClass(value) {
  const status = (value || "en ligne").toLowerCase();
  return status === "reserve" || status === "en attente" ? "status-pending" : "status-valid";
}

function formatDate(value) {
  if (!value) return "NULL";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return "NULL";
  }
  return new Intl.DateTimeFormat("fr-FR", {
    day: "2-digit",
    month: "short",
    year: "numeric"
  }).format(date);
}

function formatPrice(value, type) {
  if (value === null || value === undefined || value === "") return "NULL";
  if (type === "don" || Number(value) === 0) {
    return "Gratuit";
  }
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR"
  }).format(Number(value));
}

function isInCartItem(annonceId) {
  return isInCart(annonceId);
}

function toggleCart(annonce) {
  if (isInCartItem(annonce.id)) {
    removeFromCart(annonce.id);
    return;
  }
  addToCart(annonce);
}

function contactSeller(annonce) {
  if (!localStorage.getItem("userToken")) {
    router.push("/connexion");
    return;
  }

  const conversation = startConversation({
    kind: "vendeur",
    targetId: annonce.id_vendeur,
    name: annonce.vendeur || `Vendeur annonce #${annonce.id}`,
    subject: annonce.titre,
    contextId: annonce.id,
    contextLabel: `Annonce - ${annonce.titre}`
  });
  router.push({ path: "/messages", query: { conversation: conversation.id } });
}

onMounted(async () => {
  loading.value = true;
  try {
    annonces.value = await fetchAnnonces();
    source.value = "api";
  } catch {
    annonces.value = [];
    source.value = "empty";
  } finally {
    loading.value = false;
  }

  stopCartSync = onCartChange(() => {
    annonces.value = [...annonces.value];
  });
});

onBeforeUnmount(() => {
  stopCartSync?.();
});
</script>

<style scoped>
.public-dashboard {
  min-height: 100vh;
  padding: 20px;
  background: var(--bg-light, #f7f9f7);
}

.annonces-header .btn-main-action {
  display: inline-flex;
  align-items: center;
  text-decoration: none;
}

.annonces-stats {
  grid-template-columns: 1.4fr 0.8fr 0.8fr;
}

.score-value span {
  font-size: 1.2rem;
}

.table-subtext {
  display: block;
  margin-top: 4px;
  max-width: 360px;
  overflow: hidden;
  color: var(--text-grey);
  font-size: 0.78rem;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.annonces-table .btn-modify {
  text-decoration: none;
}

.state-card {
  border: 1px dashed #cfe0d4;
  border-radius: 14px;
  padding: 26px;
  color: var(--text-grey);
  background: #fbfdfb;
}

@media (max-width: 920px) {
  .annonces-stats {
    grid-template-columns: 1fr;
  }
}
</style>
