<template>
  <main class="page-main-content">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" />

    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > MON COMPTE > PANIER</p>
        <h1 class="hero-title1">MON PANIER</h1>
        <p class="classic-text">
          Vérifiez vos articles avant de procéder au paiement sécurisé.
        </p>
      </div>
      <div class="header-actions">
        <button class="btn-secondary-back" @click="$router.back()">
          🠔 Continuer mes achats
        </button>
      </div>
    </header>

    <div v-if="loading" class="loading-state">
      Chargement de votre panier...
    </div>

    <div v-else-if="checkoutResult" class="invoice-success-card">
      <span class="invoice-kicker">Paiement valide</span>
      <h2>Votre commande est confirmee.</h2>
      <p>
        Facture <strong>{{ checkoutResult.numero_facture }}</strong> generee pour la commande
        #{{ checkoutResult.commande_id }}.
      </p>
      <div class="invoice-actions">
        <button class="btn-main-action" @click="downloadInvoice">Telecharger la facture</button>
        <button class="btn-secondary-back" :disabled="sendingInvoice" @click="sendInvoiceByMail">
          {{ sendingInvoice ? 'Envoi...' : 'Envoyer par mail' }}
        </button>
        <button class="btn-secondary-back" @click="$router.push('/profil/factures')">
          Voir mes factures
        </button>
      </div>
      <p class="secure-payment-text">
        Sans SMTP configure, l'envoi mail cree une notification locale.
      </p>
    </div>

    <div v-else-if="cartItems.length === 0" class="empty-cart-card">
      <div class="empty-icon">🛍️</div>
      <h2>Votre panier est vide</h2>
      <p>Découvrez nos formations et nos annonces premium pour lui donner vie !</p>
      <button class="btn-main-action mt-4" @click="$router.push('/catalogue')">Explorer la plateforme</button>
    </div>

    <div v-else class="cart-layout-wrapper">
      <div class="cart-layout">
        
        <div class="cart-items-column">
          <div class="cart-items-list">
            <div v-for="item in cartItems" :key="item.id" class="cart-item-card">
              
              <div class="item-icon" :class="getIconClass(item.type_item)">
                {{ getIcon(item.type_item) }}
              </div>
              
              <div class="item-details">
                <span class="item-type">{{ item.type_item }}</span>
                <h3 class="item-title">
                  {{ getItemName(item) }}
                </h3>
                <p class="item-ref">
                  <a href="#" @click.prevent="goToItem(item)" class="item-link">
                    Voir {{ getLabelSuffix(item.type_item) }}
                  </a> 
                  • Ajouté le {{ formatDate(item.date_ajout) }}
                </p>
              </div>
              
              <div class="item-price-actions">
                <span class="item-price">{{ formatPrice(item.prix_unitaire) }}</span>
                <button class="btn-remove" @click="removeItem(item.id)" title="Supprimer l'article">
                  🗑️ Supprimer
                </button>
              </div>
              
            </div>
          </div>
        </div>

        <div class="cart-summary-column">
          <div class="summary-card">
            <h2>Résumé de la commande</h2>
            
            <div class="summary-line">
              <span>Sous-total ({{ cartItems.length }} articles)</span>
              <span>{{ formatPrice(cartTotal) }}</span>
            </div>
            
            <div class="summary-line">
              <span>Frais de gestion</span>
              <span>Gratuit</span>
            </div>
            
            <div class="summary-divider"></div>
            
            <div class="summary-line total-line">
              <span>Total TTC</span>
              <span>{{ formatPrice(cartTotal) }}</span>
            </div>
            
            <button 
              class="btn-checkout" 
              :disabled="isCheckingOut" 
              @click="handleCheckout"
            >
              {{ isCheckingOut ? 'Traitement...' : 'Payer la commande' }}
            </button>
            
            <p class="secure-payment-text">🔒 Paiement 100% sécurisé via Stripe</p>
          </div>
        </div>

      </div>
    </div>
  </main>
  <SiteFooter />
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

const router = useRouter();
const API_URL = "http://localhost:8081";

const loading = ref(true);
const cartItems = ref([]);
const isCheckingOut = ref(false);
const checkoutResult = ref(null);
const sendingInvoice = ref(false);

const getHeaders = () => ({
    "Content-Type": "application/json",
    Authorization: sessionStorage.getItem("userToken") || "",
});

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const currentUserId = computed(() => {
  const storedId = sessionStorage.getItem("id") || sessionStorage.getItem("userId");
  return Number(storedId) || 0;
});
const userName = computed(() => sessionStorage.getItem("userPrenom") || "Utilisateur");

const cartTotal = computed(() => {
  return cartItems.value.reduce((total, item) => total + Number(item.prix_unitaire || 0), 0);
});

const formatPrice = (price) => {
  return new Intl.NumberFormat("fr-FR", { style: "currency", currency: "EUR" }).format(Number(price) || 0);
};

const formatDate = (dateStr) => {
  if (!dateStr) return "";
  const date = new Date(dateStr.replace(/Z$/, ""));
  if(Number.isNaN(date.getTime())) return "";
  return new Intl.DateTimeFormat("fr-FR", { day: '2-digit', month: 'short', year: 'numeric'}).format(date);
};

const getItemName = (item) => {
  if (item.titre) return item.titre; 
  if (item.type_item === 'Formation') return `Billet pour la session de formation`;
  if (item.type_item === 'Annonce') return `Mise en avant d'annonce Premium`;
  if (item.type_item === 'Abonnement') return `Abonnement DM Plus (1 mois)`;
  return `Article #${item.reference_id}`;
};

const getIcon = (type) => {
  if (type === 'Formation') return '🎓';
  if (type === 'Annonce') return '📢';
  if (type === 'Abonnement') return '⭐';
  return '📦';
};

const getIconClass = (type) => {
  if (type === 'Formation') return 'icon-blue';
  if (type === 'Annonce') return 'icon-orange';
  if (type === 'Abonnement') return 'icon-gold';
  return 'icon-green';
};

const getLabelSuffix = (type) => {
  if (type === 'Formation') return 'la formation';
  if (type === 'Annonce') return "l'annonce";
  if (type === 'Abonnement') return "l'abonnement";
  return "l'article";
};

const goToItem = (item) => {
  if (item.type_item === 'Formation') {
    router.push(`/formations/${item.reference_id}`);
  } else if (item.type_item === 'Annonce') {
    router.push(`/annonce/${item.reference_id}`);
  }
};

const fetchPanier = async () => {
  if (!currentUserId.value) {
    router.push("/connexion");
    return;
  }

  loading.value = true;
  try {
    const res = await fetch(`${API_URL}/users/${currentUserId.value}/panier`, {
      headers: getHeaders()
    });
    
    if (res.ok) {
      cartItems.value = await res.json() || [];
    } else {
      cartItems.value = [];
    }
  } catch (error) {
    console.error("Erreur de chargement du panier:", error);
  } finally {
    loading.value = false;
  }
};

const removeItem = async (itemId) => {
  if (!confirm("Retirer cet article du panier ?")) return;

  try {
    const res = await fetch(`${API_URL}/users/${currentUserId.value}/panier/${itemId}`, {
      method: "DELETE",
      headers: getHeaders()
    });

    if (res.ok) {
      cartItems.value = cartItems.value.filter(i => i.id !== itemId);
    } else {
      alert("Erreur lors de la suppression.");
    }
  } catch (error) {
    console.error("Erreur DELETE:", error);
  }
};

const handleCheckout = () => {
  if (cartItems.value.length === 0) return;
  
  sessionStorage.setItem("pending_payment_amount", cartTotal.value);
  sessionStorage.setItem("payment_source", "panier");

  router.push({ 
    path: '/paiement', 
    query: { 
      source: 'panier', 
      amount: cartTotal.value 
    } 
  });
};

const downloadInvoice = () => {
  if (!checkoutResult.value?.facture_id) return;
  window.open(`${API_URL}/users/${currentUserId.value}/factures/${checkoutResult.value.facture_id}/download`, "_blank");
};

const sendInvoiceByMail = async () => {
  if (!checkoutResult.value?.facture_id) return;

  sendingInvoice.value = true;
  try {
    const res = await fetch(`${API_URL}/users/${currentUserId.value}/factures/${checkoutResult.value.facture_id}/send`, {
      method: "POST",
      headers: getHeaders()
    });
    const payload = res.ok ? await res.json() : { message: await res.text() };
    alert(payload.message || "Demande d'envoi traitee.");
  } catch (error) {
    console.error("Erreur envoi facture:", error);
    alert("Impossible de demander l'envoi de la facture.");
  } finally {
    sendingInvoice.value = false;
  }
};

onMounted(fetchPanier);
</script>

<style scoped>
.page-main-content {
  min-height: 100vh;
  padding: 20px;
  background: var(--bg-light, #f7f9f7);
  max-width: 1600px;
  margin: 0 auto;
  font-family: "Syne", sans-serif;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #f0f0f0;
}

.sidebar__category2 { 
  font-size: 0.75rem;
  color: #999;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 5px; 
}

.hero-title1 { 
  font-size: 2.2rem; 
  font-weight: 900; 
  color: #1a1a1a; 
  margin: 0; 
}

.btn-secondary-back {
  padding: 10px 20px;
  border-radius: 20px;
  border: 1px solid #ddd;
  background: white;
  cursor: pointer;
  font-weight: 700;
  transition: 0.2s;
  margin-bottom: 10px;
}
.btn-secondary-back:hover { background: #f0f4f1; }

.mt-4 { margin-top: 20px; }

.loading-state {
  text-align: center;
  padding: 3rem;
  color: #999;
  font-size: 0.95rem;
}

.empty-cart-card {
  background: #fff;
  border: 1px dashed #cfe0d4;
  border-radius: 16px;
  padding: 80px 20px;
  text-align: center;
  max-width: 600px;
  margin: 40px auto;
}
.empty-icon { font-size: 4rem; margin-bottom: 20px; }
.empty-cart-card h2 { font-size: 1.8rem; color: #1a1a1a; margin-bottom: 10px; }
.empty-cart-card p { color: #6d7b72; font-size: 1.1rem; }

.invoice-success-card {
  background: #fff;
  border: 1px solid #cfe8d8;
  border-radius: 22px;
  padding: 42px;
  max-width: 780px;
  margin: 40px auto;
  box-shadow: 0 18px 48px rgba(19, 87, 52, 0.08);
}
.invoice-kicker {
  display: inline-flex;
  padding: 8px 12px;
  border-radius: 999px;
  background: #e9f7ee;
  color: #2d7a4f;
  font-size: 0.78rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 1px;
}
.invoice-success-card h2 {
  margin: 18px 0 10px;
  font-size: 2rem;
  color: #16221c;
}
.invoice-success-card p {
  color: #5d6d64;
  line-height: 1.6;
}
.invoice-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin: 24px 0 8px;
}

.cart-layout-wrapper {
  max-width: 1200px;
  margin: 0 auto;
}

.cart-layout {
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: 40px;
  align-items: start;
}

.cart-items-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.cart-item-card {
  background: #fff;
  border: 1px solid #e5ede7;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  gap: 20px;
  align-items: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
  transition: 0.2s;
}
.cart-item-card:hover { border-color: #9bcbae; }

.item-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.8rem;
}
.icon-blue { background: #e8eaf6; color: #3f51b5; }
.icon-orange { background: #fff4e6; color: #cc6600; }
.icon-gold { background: #fcf4db; color: #f39c12; }
.icon-green { background: #e9f5ed; color: #2d7a4f; }

.item-details { flex: 1; }
.item-type { font-size: 0.75rem; font-weight: 800; color: #8fa396; text-transform: uppercase; letter-spacing: 1px; }
.item-title { font-size: 1.2rem; font-weight: 800; color: #1a1a1a; margin: 4px 0 8px 0; }
.item-ref { font-size: 0.85rem; color: #8fa396; margin: 0; }

.item-link {
  color: #2d7a4f;
  text-decoration: none;
  font-weight: 700;
  transition: color 0.2s;
}
.item-link:hover {
  color: #1b4332;
  text-decoration: underline;
}

.item-price-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 12px;
}

.item-price { font-size: 1.4rem; font-weight: 900; color: #2d7a4f; }
.btn-remove {
  background: none; border: none; color: #e74c3c; font-size: 0.9rem; font-weight: 700;
  cursor: pointer; padding: 4px 8px; border-radius: 6px; transition: 0.2s;
}
.btn-remove:hover { background: #fceaea; }

.summary-card {
  background: #fff;
  border: 1px solid #e5ede7;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.04);
  position: sticky;
  top: 100px;
}

.summary-card h2 { font-size: 1.4rem; color: #1a1a1a; margin-top: 0; margin-bottom: 24px; border-bottom: 1px solid #f0f4f1; padding-bottom: 15px;}

.summary-line {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
  color: #6d7b72;
  font-size: 1rem;
}

.summary-divider {
  height: 1px;
  background: #f0f4f1;
  margin: 20px 0;
}

.total-line {
  font-size: 1.3rem;
  color: #1a1a1a;
  font-weight: 900;
  margin-bottom: 30px;
}
.total-line span:last-child { color: #2d7a4f; }

.btn-checkout {
  width: 100%;
  padding: 16px;
  background: #1a1a1a;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 1.1rem;
  font-weight: 800;
  font-family: inherit;
  cursor: pointer;
  transition: 0.2s;
}
.btn-checkout:hover:not(:disabled) { background: #333; transform: translateY(-2px); }
.btn-checkout:disabled { background: #ccc; cursor: not-allowed; }

.secure-payment-text {
  text-align: center;
  font-size: 0.85rem;
  color: #8fa396;
  margin-top: 16px;
  font-weight: 600;
}

.btn-main-action {
  background: #2d7a4f;
  color: white;
  padding: 12px 24px;
  border-radius: 10px;
  font-weight: 700;
  border: none;
  cursor: pointer;
}

@media (max-width: 1024px) {
  .cart-layout { grid-template-columns: 1fr; }
  .summary-card { position: static; }
}
@media (max-width: 600px) {
  .cart-item-card { flex-direction: column; align-items: flex-start; }
  .item-price-actions { width: 100%; flex-direction: row; justify-content: space-between; align-items: center; border-top: 1px solid #f0f4f1; padding-top: 16px; }
}
</style>