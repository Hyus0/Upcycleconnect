<template>
  <main class="public-dashboard">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > ABONNEMENT</p>
        <h1 class="hero-title1">DM Plus : Messagerie illimitée</h1>
        <p class="classic-text">
          Sans abonnement, tu peux contacter 5 vendeurs via les annonces. Avec DM Plus,
          tu peux discuter avec tous les vendeurs, artisans et membres du site en illimité.
        </p>
      </div>
    </header>

    <section class="subscription-layout">
      <div class="subscription-details">
        <article class="state-card">
          <span class="step-num">01</span>
          <h2 class="step-title">Gratuit</h2>
          <p class="classic-text">
            Tu peux ouvrir des discussions avec 5 vendeurs différents via leurs annonces. Idéal pour des achats ponctuels.
          </p>
        </article>

        <article class="state-card">
          <span class="step-num">02</span>
          <h2 class="step-title">Abonné</h2>
          <p class="classic-text">
            Tu peux contacter autant de vendeurs que nécessaire et envoyer un DM direct à tout membre de la plateforme.
          </p>
        </article>

        <article class="state-card">
          <span class="step-num">03</span>
          <h2 class="step-title">Paiement</h2>
          <p class="classic-text">
            L'abonnement s'ajoute directement à ton panier. Le processus de checkout sécurisé active instantanément ton accès.
          </p>
        </article>
      </div>

      <aside class="pricing-card">
        <p class="badge-plan">Abonnement mensuel</p>
        <div class="price-container">
          <strong>2,99 €</strong>
          <span> / mois</span>
        </div>
        
        <ul class="pricing-features">
          <li>✓ Messages illimités via annonces</li>
          <li>✓ DM direct vers n'importe quel utilisateur</li>
          <li>✓ Historique centralisé et sécurisé</li>
          <li>✓ Activation immédiate après paiement</li>
        </ul>
        
        <button 
          class="btn-main-action pricing-btn" 
          type="button" 
          :disabled="loading || isSubscriber" 
          @click="subscribe"
        >
          {{ buttonLabel }}
        </button>
        <RouterLink class="btn-secondary pricing-link" to="/messages">
          Voir ma messagerie
        </RouterLink>
      </aside>
    </section>
  </main>
  
  <SiteFooter />
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

const router = useRouter();
const loading = ref(false);
const status = ref({ is_subscriber: false });
const API_URL = "/go";

const currentUserId = () => Number(sessionStorage.getItem("userId")) || 0;
const getHeaders = () => ({
    "Content-Type": "application/json",
    Authorization: sessionStorage.getItem("userToken") || "",
});

async function fetchSubscriptionStatus(userId) {
    const res = await fetch(`${API_URL}/users/${userId}/subscription`, { headers: getHeaders() });
    if (!res.ok) throw new Error("Erreur");
    return await res.json();
}

async function addSubscriptionToCart(userId) {
    const res = await fetch(`${API_URL}/users/${userId}/panier`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify({
            type_item: "Abonnement",
            reference_id: 1,
            prix_unitaire: 2.99
        })
    });
    if (!res.ok) throw new Error("Erreur lors de l'ajout au panier");
    return await res.json();
}

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const isSubscriber = computed(() => !!status.value.is_subscriber);
const userName = computed(() => {
  const prenom = sessionStorage.getItem("userPrenom") || "";
  const nom = sessionStorage.getItem("userNom") || "";
  return `${prenom} ${nom}`.trim() || "Utilisateur";
});

const buttonLabel = computed(() => {
  if (isSubscriber.value) return "DM Plus déjà actif";
  return loading.value ? "Ajout en cours..." : "Ajouter au panier - 2,99 €";
});

async function loadStatus() {
  const userId = currentUserId();
  if (!userId) return;
  try {
      status.value = await fetchSubscriptionStatus(userId);
  } catch(e) {
      console.error(e);
  }
}

async function subscribe() {
  if (!isLoggedIn.value) {
    router.push("/connexion");
    return;
  }
  loading.value = true;
  try {
    await addSubscriptionToCart(currentUserId());
    router.push("/panier");
  } catch (error) {
    alert(error.message || "Impossible d'ajouter l'abonnement au panier.");
  } finally {
    loading.value = false;
  }
}

onMounted(loadStatus);
</script>

<style scoped>

.public-dashboard {
  min-height: 100vh;
  padding: 20px;
  background: var(--bg-light, #f7f9f7);
  display: flex;
  flex-direction: column;
}

.content-header {
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-top: 20px;
}

.header-left {
  flex: 1;
  max-width: 800px;
  padding-right: 20px;
}

.sidebar__category2 {
  color: #2f8f5b;
  font-size: 0.8rem;
  font-weight: 800;
  letter-spacing: 0.1em;
  margin: 0 0 8px 0;
  text-transform: uppercase;
}

.hero-title1 {
  font-size: 2.2rem;
  font-weight: 900;
  color: #1a1a1a;
  margin: 5px 0 0 0;
  font-family: "Syne", sans-serif;
  line-height: 1.1;
}

.classic-text {
  font-size: 0.95rem;
  color: #6d7b72;
  margin-top: 8px;
  line-height: 1.5;
}


.subscription-layout {
  display: grid;
  grid-template-columns: 1fr 420px;
  gap: 30px;
  align-items: start;
  margin-top: 24px;
}

.subscription-details {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.state-card {
  border: 1px dashed #cfe0d4;
  border-radius: 14px;
  padding: 26px;
  background: #fbfdfb;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.state-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 24px rgba(44, 126, 79, 0.08);
  border-color: #9bcbae;
}

.step-num {
  color: #2c7e4f;
  font-weight: 900;
  font-size: 1.2rem;
  letter-spacing: 0.1em;
}

.step-title {
  font-family: "Syne", sans-serif;
  font-size: 1.5rem;
  font-weight: 800;
  color: #1a1a1a;
  margin: 10px 0;
}


.pricing-card {
  padding: 32px;
  border-radius: 16px;
  color: #fff;
  background: linear-gradient(150deg, #102018, #1f5035 58%, #2c7e4f);
  box-shadow: 0 28px 80px rgba(17, 44, 29, 0.24);
  display: flex;
  flex-direction: column;
}

.badge-plan {
  display: inline-block;
  align-self: flex-start;
  margin: 0 0 16px 0;
  padding: 8px 16px;
  border-radius: 8px;
  color: #102018;
  background: #a9f1c6;
  font-size: 0.85rem;
  font-weight: 900;
  font-family: "Syne", sans-serif;
  text-transform: uppercase;
}

.price-container {
  margin: 10px 0;
}

.price-container strong {
  font-family: "Syne", sans-serif;
  font-size: 3.5rem;
  line-height: 1;
}

.price-container span {
  color: rgba(255, 255, 255, 0.7);
  font-weight: 800;
  font-size: 1.2rem;
}

.pricing-features {
  margin: 20px 0 30px 0;
  padding: 0;
  list-style: none;
}

.pricing-features li {
  padding: 12px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.12);
  font-size: 0.95rem;
  color: #eef8f0;
}


.btn-main-action {
  box-sizing: border-box;
  width: 100%;
  height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 0.95rem;
  font-family: "Syne", sans-serif;
  font-weight: 700;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s;
  border: none;
}

.btn-main-action:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.pricing-btn {
  background-color: #a9f1c6;
  color: #102018;
}

.pricing-btn:hover:not(:disabled) {
  background-color: #8ce6b0;
}

.btn-secondary {
  box-sizing: border-box;
  width: 100%;
  height: 44px;
  margin-top: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 0.95rem;
  font-family: "Syne", sans-serif;
  font-weight: 700;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.pricing-link {
  background: transparent;
  color: #fff;
  border-color: rgba(255, 255, 255, 0.3);
}

.pricing-link:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.6);
}


@media (max-width: 920px) {
  .subscription-layout {
    grid-template-columns: 1fr;
  }
}
</style>