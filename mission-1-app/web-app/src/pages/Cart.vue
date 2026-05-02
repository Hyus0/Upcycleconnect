<template>
  <main class="public-dashboard">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > PANIER</p>
        <h1 class="hero-title1">Mon panier</h1>
        <p class="classic-text">
          Retrouve les annonces et formations que tu souhaites acheter ou reserver avant validation.
        </p>
      </div>
      <button v-if="items.length" class="btn-main-action cart-clear-button" type="button" @click="handleClearCart">
        Vider le panier
      </button>
    </header>

    <div class="stats-grid annonces-stats">
      <div class="card card--score">
        <p class="tag-score">PANIER ACTIF</p>
        <div class="score-value">{{ items.length }} <span>article{{ items.length > 1 ? "s" : "" }}</span></div>
        <p class="score-level">Selection prete pour validation</p>
        <div class="score-footer">
          <div class="mini-stat">
            <strong>{{ paidItemsCount }}</strong><br />Payants
          </div>
          <div class="mini-stat">
            <strong>{{ freeItemsCount }}</strong><br />Dons
          </div>
          <div class="mini-stat">
            <strong>{{ totalPriceLabel }}</strong><br />Total
          </div>
        </div>
      </div>
      <div class="card card--white">
        <div class="card-num">{{ paidItemsCount }}</div>
        <p class="text-dm">Articles payants</p>
        <span class="badge badge--orange">PANIER</span>
      </div>
      <div class="card card--white">
        <div class="card-num2">{{ freeItemsCount }}</div>
        <p class="text-dm">Articles gratuits</p>
        <span class="badge badge--green">DON</span>
      </div>
    </div>

    <section class="section-container">
      <div class="section-header">
        <div>
          <h2>Articles selectionnes</h2>
          <p class="classic-text">{{ items.length }} element{{ items.length > 1 ? "s" : "" }} dans le panier</p>
        </div>
        <RouterLink class="btn-secondary cart-continue-link" to="/annonces">
          Continuer mes recherches
        </RouterLink>
      </div>

      <div v-if="items.length === 0" class="state-card">
        Le panier est vide pour le moment. Ajoute des annonces ou des formations pour commencer.
      </div>

      <template v-else>
        <table class="data-table annonces-table">
          <thead>
            <tr>
              <th>Article</th>
              <th>Localisation</th>
              <th>Nature</th>
              <th>Statut</th>
              <th>Prix</th>
              <th>Date</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in items" :key="`${item.itemType || 'annonce'}-${item.id}`">
              <td>
                <strong>{{ item.titre }}</strong>
                <span class="table-subtext">{{ item.description }}</span>
              </td>
              <td>
                {{ item.ville }}
                <span class="table-subtext">{{ item.code_postal }}</span>
              </td>
              <td>
                <span :class="item.itemType === 'formation' || item.type === 'vente' ? 'tag-vente' : 'tag-don'">
                  {{ itemLabel(item) }}
                </span>
              </td>
              <td>
                <span :class="statusClass(item.statut)">
                  {{ (item.statut || "en ligne").toUpperCase() }}
                </span>
              </td>
              <td>{{ formatPrice(item.prix, item.type) }}</td>
              <td>{{ formatDate(item.date_creation) }}</td>
              <td class="actions-cell">
                <button class="btn-remove" type="button" @click="handleRemove(item)">
                  Retirer
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <div class="cart-summary">
          <div class="cart-summary__content">
            <div>
              <strong>Total panier</strong>
              <p>{{ paidItemsCount }} achat{{ paidItemsCount > 1 ? "s" : "" }} et {{ freeItemsCount }} don{{ freeItemsCount > 1 ? "s" : "" }}</p>
            </div>
            <div class="cart-summary__price">{{ totalPriceLabel }}</div>
          </div>
          <div class="cart-summary__actions">
            <button class="site-navbar__button site-navbar__button--ghost" type="button" @click="handleClearCart">
              Vider
            </button>
            <button class="site-navbar__button site-navbar__button--primary" type="button" @click="handleCheckout">
              Confirmer l'achat
            </button>
          </div>
        </div>
      </template>
    </section>

    <section v-if="purchases.length" class="section-container">
      <div class="section-header">
        <div>
          <h2>Achats recents</h2>
          <p class="classic-text">Historique local des paniers confirmes.</p>
        </div>
      </div>
      <div class="purchase-list">
        <article v-for="purchase in purchases" :key="purchase.id" class="purchase-row">
          <strong>{{ purchase.id }}</strong>
          <span>{{ formatDate(purchase.created_at) }}</span>
          <span>{{ purchase.items.length }} element{{ purchase.items.length > 1 ? "s" : "" }}</span>
          <span>{{ formatPrice(purchase.total, 'vente') }}</span>
        </article>
      </div>
    </section>
  </main>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import { clearCart, getCartItems, getPurchases, onCartChange, removeItemFromCart, savePurchase } from "../services/cartService";

const items = ref([]);
const purchases = ref([]);
const isLoggedIn = computed(() => !!localStorage.getItem("userToken"));
const userName = computed(() => {
  const prenom = localStorage.getItem("userPrenom") || "";
  const nom = localStorage.getItem("userNom") || "";
  return (prenom || nom) ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const paidItemsCount = computed(() => items.value.filter((item) => Number(item.prix) > 0 && String(item.type).toLowerCase() !== "don").length);
const freeItemsCount = computed(() => items.value.length - paidItemsCount.value);
const totalPrice = computed(() =>
  items.value.reduce((sum, item) => sum + (String(item.type).toLowerCase() === "don" ? 0 : Number(item.prix) || 0), 0)
);
const totalPriceLabel = computed(() =>
  new Intl.NumberFormat("fr-FR", { style: "currency", currency: "EUR" }).format(totalPrice.value)
);

let stopCartSync = null;

function syncCart() {
  items.value = getCartItems();
  purchases.value = getPurchases();
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
  if (String(type).toLowerCase() === "don" || Number(value) === 0) {
    return "Gratuit";
  }
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR"
  }).format(Number(value) || 0);
}

function statusClass(value) {
  const status = (value || "en ligne").toLowerCase();
  return status === "reserve" || status === "en attente" ? "status-pending" : "status-valid";
}

function itemLabel(item) {
  if (item.itemType === "formation") return "FORMATION";
  return String(item.type || "annonce").toUpperCase();
}

function handleRemove(item) {
  removeItemFromCart(item.itemType || "annonce", item.id);
}

function handleClearCart() {
  clearCart();
}

function handleCheckout() {
  if (!localStorage.getItem("userToken")) {
    window.alert("Connectez-vous pour confirmer un panier.");
    return;
  }
  const purchase = savePurchase(items.value);
  clearCart();
  window.alert(`Panier confirme. Reference ${purchase.id}.`);
}

onMounted(() => {
  syncCart();
  stopCartSync = onCartChange(syncCart);
});

onBeforeUnmount(() => {
  stopCartSync?.();
});
</script>

<style scoped>
.cart-clear-button {
  padding-inline: 20px;
}

.cart-continue-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
}

.cart-summary {
  margin-top: 22px;
  display: flex;
  justify-content: space-between;
  gap: 18px;
  padding: 20px;
  border: 1px solid var(--border);
  border-radius: 18px;
  background: linear-gradient(180deg, #fcfffd, #f3f8f4);
}

.cart-summary__content {
  display: grid;
  gap: 6px;
}

.cart-summary__content p {
  margin: 0;
  color: var(--text-secondary);
}

.cart-summary__price {
  font-family: "Syne", sans-serif;
  font-size: 2rem;
  font-weight: 800;
  color: var(--brand-green-deep);
}

.cart-summary__actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.purchase-list {
  display: grid;
  gap: 10px;
}

.purchase-row {
  display: grid;
  grid-template-columns: 1.2fr 1fr 1fr 1fr;
  gap: 12px;
  padding: 14px;
  border: 1px solid var(--border);
  border-radius: 12px;
  background: #f8fbf8;
}

@media (max-width: 920px) {
  .cart-summary {
    flex-direction: column;
  }

  .cart-summary__actions {
    justify-content: flex-start;
    flex-wrap: wrap;
  }
}
</style>
