<template>
  <main class="subscription-page">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <section class="subscription-hero">
      <div class="subscription-copy">
        <p class="eyebrow">DM PLUS</p>
        <h1>Messagerie illimitee pour tes echanges UpcycleConnect.</h1>
        <p>
          Sans abonnement, tu peux contacter 5 vendeurs via les annonces. Avec DM Plus,
          tu peux discuter avec tous les vendeurs, artisans et membres du site.
        </p>
      </div>

      <aside class="pricing-card">
        <p class="badge-plan">Abonnement mensuel</p>
        <strong>2,99 €</strong>
        <span>/ mois</span>
        <ul>
          <li>Messages illimites via annonces</li>
          <li>DM direct vers n'importe quel utilisateur</li>
          <li>Historique centralise dans la page messagerie</li>
          <li>Activation apres paiement du panier</li>
        </ul>
        <button class="btn-main-action" type="button" :disabled="loading || isSubscriber" @click="subscribe">
          {{ buttonLabel }}
        </button>
        <RouterLink class="btn-secondary" to="/messages">Voir ma messagerie</RouterLink>
      </aside>
    </section>

    <section class="subscription-details">
      <article>
        <span>01</span>
        <h2>Gratuit</h2>
        <p>Tu peux ouvrir des discussions avec 5 vendeurs differents via leurs annonces.</p>
      </article>
      <article>
        <span>02</span>
        <h2>Abonne</h2>
        <p>Tu peux contacter autant de vendeurs que necessaire et envoyer un DM direct a tout membre.</p>
      </article>
      <article>
        <span>03</span>
        <h2>Paiement</h2>
        <p>L'abonnement est ajoute au panier. Le checkout local simule le paiement Stripe en attendant les cles API.</p>
      </article>
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import { addSubscriptionToCart, currentUserId, fetchSubscriptionStatus } from "../services/messagesApi";

const router = useRouter();
const loading = ref(false);
const status = ref({ is_subscriber: false });

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const isSubscriber = computed(() => !!status.value.is_subscriber);
const userName = computed(() => {
  const prenom = sessionStorage.getItem("userPrenom") || "";
  const nom = sessionStorage.getItem("userNom") || "";
  return `${prenom} ${nom}`.trim() || "Utilisateur";
});
const buttonLabel = computed(() => {
  if (isSubscriber.value) return "DM Plus deja actif";
  return loading.value ? "Ajout en cours..." : "Ajouter au panier - 2,99 €";
});

async function loadStatus() {
  if (!currentUserId()) return;
  status.value = await fetchSubscriptionStatus(currentUserId());
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
.subscription-page {
  min-height: 100vh;
  padding: 20px;
  background:
    radial-gradient(circle at top right, rgba(47, 143, 91, 0.16), transparent 34rem),
    #f5f8f4;
}

.subscription-hero,
.subscription-details {
  max-width: 1440px;
  margin: 0 auto;
}

.subscription-hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 420px;
  gap: 38px;
  align-items: center;
  padding: 72px 0 36px;
}

.eyebrow {
  margin: 0 0 12px;
  color: #2f8f5b;
  font-size: 0.75rem;
  font-weight: 900;
  letter-spacing: 0.24em;
}

.subscription-copy h1 {
  max-width: 880px;
  margin: 0;
  color: #102018;
  font-family: "Syne", sans-serif;
  font-size: clamp(3rem, 7vw, 6.8rem);
  line-height: 0.9;
}

.subscription-copy p {
  max-width: 720px;
  color: #60736a;
  font-size: 1.12rem;
  line-height: 1.7;
}

.pricing-card {
  padding: 32px;
  border-radius: 34px;
  color: #fff;
  background: linear-gradient(150deg, #102018, #1f5035 58%, #2f8f5b);
  box-shadow: 0 28px 80px rgba(17, 44, 29, 0.24);
}

.badge-plan {
  display: inline-block;
  margin: 0 0 22px;
  padding: 8px 12px;
  border-radius: 999px;
  color: #a9f1c6;
  background: rgba(255, 255, 255, 0.1);
  font-weight: 900;
}

.pricing-card strong {
  display: inline-block;
  font-family: "Syne", sans-serif;
  font-size: 4.8rem;
  line-height: 1;
}

.pricing-card span {
  color: rgba(255, 255, 255, 0.7);
  font-weight: 800;
}

.pricing-card ul {
  margin: 24px 0;
  padding: 0;
  list-style: none;
}

.pricing-card li {
  padding: 12px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.12);
}

.btn-main-action,
.btn-secondary {
  width: 100%;
  display: inline-flex;
  justify-content: center;
  box-sizing: border-box;
  border-radius: 18px;
  padding: 16px 20px;
  font-weight: 900;
  text-decoration: none;
}

.btn-main-action {
  border: 0;
  color: #102018;
  background: #a9f1c6;
  cursor: pointer;
}

.btn-main-action:disabled {
  opacity: 0.62;
  cursor: not-allowed;
}

.btn-secondary {
  margin-top: 12px;
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.22);
}

.subscription-details {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 18px;
  padding-bottom: 60px;
}

.subscription-details article {
  padding: 28px;
  border: 1px solid #dfe9e2;
  border-radius: 28px;
  background: #fff;
}

.subscription-details span {
  color: #2f8f5b;
  font-weight: 900;
  letter-spacing: 0.2em;
}

.subscription-details h2 {
  margin: 14px 0 8px;
  color: #102018;
  font-family: "Syne", sans-serif;
}

.subscription-details p {
  color: #60736a;
}

@media (max-width: 920px) {
  .subscription-hero,
  .subscription-details {
    grid-template-columns: 1fr;
  }
}
</style>
