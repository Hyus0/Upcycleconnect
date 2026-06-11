<template>
  <main class="messages-page">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <section class="messages-hero">
      <div>
        <p class="eyebrow">MESSAGERIE</p>
        <h1>Discussions entre membres</h1>
        <p>
          Contacte les vendeurs, artisans et utilisateurs autour d'une annonce ou en direct avec DM Plus.
        </p>
      </div>
      <RouterLink class="hero-subscription" to="/abonnement">
        <span>DM Plus</span>
        <strong>2,99 € / mois</strong>
        <small>{{ subscriptionLabel }}</small>
      </RouterLink>
    </section>

    <section v-if="limitPopup" class="limit-popup" role="dialog" aria-live="polite">
      <button type="button" class="limit-popup__close" @click="limitPopup = null">×</button>
      <p class="eyebrow">LIMITE ATTEINTE</p>
      <h2>{{ limitPopup.message }}</h2>
      <p>
        Les comptes gratuits peuvent contacter 5 vendeurs via annonce. L'abonnement DM Plus debloque les messages illimites.
      </p>
      <div class="limit-popup__actions">
        <RouterLink class="btn-main-action" to="/abonnement">Prendre DM Plus</RouterLink>
        <button class="btn-secondary" type="button" @click="limitPopup = null">Plus tard</button>
      </div>
    </section>

    <section class="messenger-shell">
      <aside class="conversation-list">
        <div class="panel-heading">
          <div>
            <p class="eyebrow">BOITE DE RECEPTION</p>
            <h2>{{ conversations.length }} conversations</h2>
          </div>
          <button class="refresh-btn" type="button" @click="loadConversations">Actualiser</button>
        </div>

        <div v-if="loadingConversations" class="state-card">Chargement...</div>
        <div v-else-if="conversations.length === 0" class="state-card">
          Aucune discussion pour le moment. Clique sur “Contacter” depuis une annonce.
        </div>

        <template v-else>
          <button
            v-for="conversation in conversations"
            :key="conversation.id"
            class="conversation-item"
            :class="{ active: conversation.id === activeConversationId }"
            type="button"
            @click="selectConversation(conversation.id)"
          >
            <span class="conversation-avatar">{{ initials(conversation.other_user_name) }}</span>
            <span class="conversation-main">
              <strong>{{ conversation.other_user_name }}</strong>
              <small>{{ conversation.annonce_title }}</small>
              <em>{{ conversation.last_message || "Nouvelle discussion" }}</em>
            </span>
            <span v-if="conversation.unread_count" class="unread-badge">{{ conversation.unread_count }}</span>
          </button>
        </template>
      </aside>

      <section class="chat-panel">
        <div v-if="!activeConversation" class="empty-chat">
          <p class="eyebrow">SELECTION</p>
          <h2>Choisis une conversation</h2>
          <p>Le fil de discussion apparaitra ici avec l'annonce liee et l'historique des messages.</p>
        </div>

        <template v-else>
          <header class="chat-header">
            <div class="conversation-avatar conversation-avatar--large">
              {{ initials(activeConversation.other_user_name) }}
            </div>
            <div>
              <p class="eyebrow">{{ activeConversation.other_user_role || "MEMBRE" }}</p>
              <h2>
                {{ activeConversation.other_user_name }}
                <span v-if="activeConversation.other_user_premium" class="premium-check" title="Membre DM Plus">✓</span>
              </h2>
              <p>{{ activeConversation.annonce_title }}</p>
            </div>
            <RouterLink class="btn-secondary" :to="`/user/${activeConversation.other_user_id}`">
              Profil
            </RouterLink>
          </header>

          <div class="message-thread">
            <section v-if="activeConversation.annonce_id" class="deal-panel">
              <div>
                <p class="eyebrow">NEGOCIATION</p>
                <h3>{{ activeConversation.annonce_title }}</h3>
                <p>Prix annonce : {{ formatPrice(activeConversation.annonce_price) }}</p>
              </div>
              <form class="offer-form" @submit.prevent="handleCreateOffer">
                <input v-model="offerAmount" type="number" min="1" step="0.01" placeholder="Votre offre (€)" />
                <button class="btn-main-action" type="submit" :disabled="!offerAmount || creatingOffer">
                  {{ creatingOffer ? "Offre..." : "Proposer" }}
                </button>
              </form>
            </section>

            <section v-if="threadState.offers.length" class="offer-stack">
              <article v-for="offer in threadState.offers" :key="offer.id" class="offer-card">
                <div>
                  <strong>{{ formatPrice(offer.amount) }}</strong>
                  <span :class="['offer-status', statusClass(offer.status)]">{{ offer.status }}</span>
                </div>
                <p>Offre de {{ offer.buyer_id === currentUser ? "vous" : "l'acheteur" }}</p>
                <div v-if="offer.status === 'En attente' && offer.seller_id === currentUser" class="offer-actions">
                  <button type="button" class="btn-main-action" @click="handleOfferResponse(offer.id, 'accept')">Accepter</button>
                  <button type="button" class="btn-secondary" @click="handleOfferResponse(offer.id, 'refuse')">Refuser</button>
                </div>
              </article>
            </section>

            <section v-if="threadState.sales.length" class="offer-stack">
              <article v-for="sale in threadState.sales" :key="sale.id" class="sale-card">
                <div>
                  <strong>Vente {{ formatPrice(sale.amount) }}</strong>
                  <span :class="['offer-status', statusClass(sale.status)]">{{ sale.status }}</span>
                </div>
                <p>Validez la reception puis notez le vendeur pour finaliser la vente.</p>
                <button
                  v-if="sale.buyer_id === currentUser && sale.status === 'Offre acceptee'"
                  type="button"
                  class="btn-main-action"
                  @click="handleConfirmReception(sale.id)"
                >
                  Confirmer la reception
                </button>
                <form
                  v-if="sale.buyer_id === currentUser && sale.status === 'Recue'"
                  class="review-inline"
                  @submit.prevent="handleReviewSale(sale.id)"
                >
                  <select v-model.number="reviewDraft.note">
                    <option v-for="n in 5" :key="n" :value="n">{{ n }} / 5</option>
                  </select>
                  <input v-model="reviewDraft.commentaire" placeholder="Votre avis sur la transaction" />
                  <button class="btn-main-action" type="submit">Noter</button>
                </form>
              </article>
            </section>

            <div v-if="loadingMessages" class="state-card">Chargement des messages...</div>
            <article
              v-for="message in messages"
              v-else
              :key="message.id"
              class="message-bubble"
              :class="{ mine: message.sender_id === currentUser }"
            >
              <p>{{ message.content }}</p>
              <span>{{ formatDate(message.created_at) }}</span>
            </article>
          </div>

          <form class="composer" @submit.prevent="handleSend">
            <textarea
              v-model="draft"
              rows="2"
              placeholder="Ecris ton message..."
            />
            <button class="btn-main-action" type="submit" :disabled="sending || !draft.trim()">
              {{ sending ? "Envoi..." : "Envoyer" }}
            </button>
          </form>
        </template>
      </section>
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import {
  currentUserId,
  confirmSaleReception,
  createOffer,
  fetchConversations,
  fetchConversationState,
  fetchMessages,
  fetchSubscriptionStatus,
  respondOffer,
  reviewSale,
  sendMessage,
  startConversation
} from "../services/messagesApi";

const route = useRoute();
const router = useRouter();

const conversations = ref([]);
const messages = ref([]);
const activeConversationId = ref(Number(route.query.conversation) || 0);
const loadingConversations = ref(false);
const loadingMessages = ref(false);
const sending = ref(false);
const draft = ref("");
const offerAmount = ref("");
const creatingOffer = ref(false);
const threadState = ref({ offers: [], sales: [] });
const reviewDraft = ref({ note: 5, commentaire: "" });
const subscription = ref({ is_subscriber: false, used: 0, limit: 5 });
const limitPopup = ref(null);

const currentUser = computed(() => currentUserId());
const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => {
  const prenom = sessionStorage.getItem("userPrenom") || "";
  const nom = sessionStorage.getItem("userNom") || "";
  return `${prenom} ${nom}`.trim() || "Utilisateur";
});
const activeConversation = computed(() =>
  conversations.value.find((conversation) => conversation.id === activeConversationId.value)
);
const subscriptionLabel = computed(() =>
  subscription.value.is_subscriber
    ? "Messagerie illimitee active"
    : `${subscription.value.used || 0}/${subscription.value.limit || 5} vendeurs gratuits`
);

function initials(name) {
  return (name || "UC")
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");
}

function formatDate(value) {
  if (!value) return "";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";
  return new Intl.DateTimeFormat("fr-FR", {
    day: "2-digit",
    month: "short",
    hour: "2-digit",
    minute: "2-digit"
  }).format(date);
}

function formatPrice(value) {
  return new Intl.NumberFormat("fr-FR", { style: "currency", currency: "EUR" }).format(Number(value) || 0);
}

function statusClass(status) {
  const value = (status || "").toLowerCase();
  if (value.includes("accept") || value.includes("recue") || value.includes("evalue")) return "is-ok";
  if (value.includes("refus") || value.includes("annul")) return "is-ko";
  return "is-pending";
}

async function loadSubscription() {
  if (!currentUser.value) return;
  subscription.value = await fetchSubscriptionStatus(currentUser.value);
}

async function loadConversations() {
  if (!currentUser.value) {
    router.push("/connexion");
    return;
  }
  loadingConversations.value = true;
  try {
    conversations.value = await fetchConversations(currentUser.value);
    if (!activeConversationId.value && conversations.value.length) {
      activeConversationId.value = conversations.value[0].id;
    }
  } finally {
    loadingConversations.value = false;
  }
}

async function selectConversation(id) {
  activeConversationId.value = id;
  router.replace({ path: "/messages", query: { conversation: id } });
  await loadMessages();
}

async function loadMessages() {
  if (!activeConversationId.value) return;
  loadingMessages.value = true;
  try {
    messages.value = await fetchMessages(activeConversationId.value, currentUser.value);
    threadState.value = await fetchConversationState(activeConversationId.value, currentUser.value);
    await loadConversations();
  } finally {
    loadingMessages.value = false;
  }
}

async function handleCreateOffer() {
  if (!activeConversationId.value || !offerAmount.value) return;
  creatingOffer.value = true;
  try {
    await createOffer(activeConversationId.value, offerAmount.value, currentUser.value);
    offerAmount.value = "";
    await loadMessages();
  } catch (error) {
    alert(error.message || "Impossible de proposer cette offre.");
  } finally {
    creatingOffer.value = false;
  }
}

async function handleOfferResponse(offerId, action) {
  try {
    await respondOffer(offerId, action, currentUser.value);
    await loadMessages();
  } catch (error) {
    alert(error.message || "Impossible de traiter cette offre.");
  }
}

async function handleConfirmReception(saleId) {
  try {
    await confirmSaleReception(saleId, currentUser.value);
    await loadMessages();
  } catch (error) {
    alert(error.message || "Impossible de confirmer la reception.");
  }
}

async function handleReviewSale(saleId) {
  try {
    await reviewSale(saleId, reviewDraft.value, currentUser.value);
    reviewDraft.value = { note: 5, commentaire: "" };
    await loadMessages();
  } catch (error) {
    alert(error.message || "Impossible d'enregistrer l'avis.");
  }
}

async function startFromQuery() {
  const target = Number(route.query.user);
  const annonce = Number(route.query.annonce);
  if (!target) return;
  try {
    const result = await startConversation({
      targetUserId: target,
      annonceId: annonce || null
    });
    activeConversationId.value = result.conversation_id;
    await loadConversations();
    await selectConversation(result.conversation_id);
  } catch (error) {
    if (error.status === 402) {
      limitPopup.value = error.payload || { message: error.message };
    } else {
      alert(error.message || "Impossible d'ouvrir la conversation.");
    }
  }
}

async function handleSend() {
  if (!draft.value.trim() || !activeConversationId.value) return;
  sending.value = true;
  try {
    const message = await sendMessage(activeConversationId.value, draft.value.trim(), currentUser.value);
    messages.value.push(message);
    draft.value = "";
    await loadConversations();
  } finally {
    sending.value = false;
  }
}

watch(() => route.query.conversation, async (value) => {
  const id = Number(value);
  if (id && id !== activeConversationId.value) {
    activeConversationId.value = id;
    await loadMessages();
  }
});

onMounted(async () => {
  if (!isLoggedIn.value) {
    router.push("/connexion");
    return;
  }
  await loadSubscription();
  await loadConversations();
  await startFromQuery();
  await loadMessages();
});
</script>

<style scoped>
.messages-page {
  min-height: 100vh;
  padding: 20px;
  background:
    radial-gradient(circle at top left, rgba(47, 143, 91, 0.14), transparent 34rem),
    #f5f8f4;
}

.messages-hero,
.messenger-shell,
.limit-popup {
  max-width: 1480px;
  margin: 0 auto;
}

.messages-hero {
  display: flex;
  align-items: stretch;
  justify-content: space-between;
  gap: 24px;
  padding: 44px 4px 28px;
}

.eyebrow {
  margin: 0 0 8px;
  color: #2f8f5b;
  font-size: 0.75rem;
  font-weight: 800;
  letter-spacing: 0.22em;
}

.messages-hero h1 {
  margin: 0;
  color: #15231d;
  font-family: "Syne", sans-serif;
  font-size: clamp(2.5rem, 6vw, 5.5rem);
  line-height: 0.92;
}

.messages-hero p {
  color: #5f7068;
  max-width: 680px;
  font-size: 1.05rem;
}

.hero-subscription {
  min-width: 260px;
  padding: 24px;
  border-radius: 28px;
  color: #eef8f0;
  background: linear-gradient(145deg, #102018, #2f8f5b);
  text-decoration: none;
  box-shadow: 0 24px 70px rgba(23, 55, 40, 0.22);
}

.hero-subscription span,
.hero-subscription small {
  display: block;
  color: rgba(255, 255, 255, 0.72);
}

.hero-subscription strong {
  display: block;
  margin: 12px 0;
  font-size: 2rem;
}

.messenger-shell {
  display: grid;
  grid-template-columns: 380px minmax(0, 1fr);
  min-height: 650px;
  border: 1px solid rgba(29, 56, 42, 0.12);
  border-radius: 34px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.76);
  box-shadow: 0 28px 80px rgba(26, 52, 38, 0.1);
}

.conversation-list {
  padding: 24px;
  background: #102018;
  color: #fff;
}

.panel-heading {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 22px;
}

.panel-heading h2,
.chat-header h2,
.empty-chat h2,
.limit-popup h2 {
  margin: 0;
  font-family: "Syne", sans-serif;
}

.refresh-btn,
.btn-secondary {
  border: 1px solid rgba(47, 143, 91, 0.28);
  border-radius: 16px;
  padding: 12px 16px;
  color: #1d3528;
  background: #eef6f0;
  font-weight: 800;
  text-decoration: none;
  cursor: pointer;
}

.conversation-item {
  width: 100%;
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 14px;
  margin-bottom: 12px;
  padding: 14px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 22px;
  color: #dce9e1;
  background: rgba(255, 255, 255, 0.04);
  text-align: left;
  cursor: pointer;
}

.conversation-item.active,
.conversation-item:hover {
  background: rgba(47, 143, 91, 0.26);
  border-color: rgba(89, 210, 130, 0.34);
}

.conversation-avatar {
  width: 46px;
  height: 46px;
  display: inline-grid;
  place-items: center;
  border-radius: 16px;
  color: #fff;
  background: #2f8f5b;
  font-weight: 900;
}

.conversation-avatar--large {
  width: 64px;
  height: 64px;
  border-radius: 22px;
}

.conversation-main strong,
.conversation-main small,
.conversation-main em {
  display: block;
}

.conversation-main small {
  margin: 2px 0;
  color: #a7b9ae;
}

.conversation-main em {
  color: #d7e5dc;
  font-size: 0.88rem;
  font-style: normal;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.unread-badge {
  min-width: 24px;
  height: 24px;
  display: inline-grid;
  place-items: center;
  border-radius: 999px;
  background: #f1a321;
  color: #102018;
  font-size: 0.8rem;
  font-weight: 900;
}

.chat-panel {
  display: flex;
  min-width: 0;
  flex-direction: column;
  background: #fbfdfb;
}

.chat-header {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 26px;
  border-bottom: 1px solid #dfe9e2;
}

.chat-header p {
  margin: 4px 0 0;
  color: #64766c;
}

.premium-check {
  display: inline-grid;
  width: 24px;
  height: 24px;
  place-items: center;
  margin-left: 8px;
  border-radius: 999px;
  color: #102018;
  background: #8ef0a8;
  font-size: 0.85rem;
  font-weight: 900;
  vertical-align: middle;
}

.chat-header .btn-secondary {
  margin-left: auto;
}

.message-thread {
  flex: 1;
  padding: 28px;
  overflow-y: auto;
  background:
    linear-gradient(rgba(16, 32, 24, 0.025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(16, 32, 24, 0.025) 1px, transparent 1px);
  background-size: 28px 28px;
}

.deal-panel,
.offer-card,
.sale-card {
  margin-bottom: 16px;
  padding: 18px;
  border: 1px solid #d7e6dc;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 14px 36px rgba(32, 58, 43, 0.07);
}

.deal-panel {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 16px;
  align-items: end;
}

.deal-panel h3,
.offer-card strong,
.sale-card strong {
  margin: 0;
  color: #122018;
  font-family: "Syne", sans-serif;
}

.deal-panel p,
.offer-card p,
.sale-card p {
  margin: 6px 0 0;
  color: #63746a;
}

.offer-form,
.review-inline {
  display: flex;
  gap: 10px;
  align-items: center;
}

.offer-form input,
.review-inline input,
.review-inline select {
  min-height: 46px;
  border: 1px solid #d4e3d8;
  border-radius: 14px;
  padding: 0 14px;
  font: inherit;
  background: #f8fbf9;
}

.offer-stack {
  display: grid;
  gap: 12px;
}

.offer-card,
.sale-card {
  display: grid;
  gap: 12px;
}

.offer-card > div:first-child,
.sale-card > div:first-child,
.offer-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.offer-status {
  padding: 6px 10px;
  border-radius: 999px;
  font-size: 0.76rem;
  font-weight: 900;
}

.offer-status.is-ok {
  background: #dff5e6;
  color: #1f7a46;
}

.offer-status.is-ko {
  background: #fde6e2;
  color: #b3261e;
}

.offer-status.is-pending {
  background: #fff4d8;
  color: #9a6400;
}

.message-bubble {
  max-width: 68%;
  margin-bottom: 14px;
  padding: 14px 16px;
  border-radius: 20px 20px 20px 6px;
  background: #ffffff;
  border: 1px solid #dfe9e2;
  box-shadow: 0 12px 30px rgba(32, 58, 43, 0.07);
}

.message-bubble.mine {
  margin-left: auto;
  border-radius: 20px 20px 6px 20px;
  color: #fff;
  background: #2f8f5b;
  border-color: #2f8f5b;
}

.message-bubble p {
  margin: 0;
  line-height: 1.5;
}

.message-bubble span {
  display: block;
  margin-top: 8px;
  color: inherit;
  opacity: 0.7;
  font-size: 0.78rem;
}

.composer {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 14px;
  padding: 20px;
  border-top: 1px solid #dfe9e2;
  background: #ffffff;
}

.composer textarea {
  resize: none;
  min-height: 54px;
  border: 1px solid #d4e3d8;
  border-radius: 18px;
  padding: 14px;
  font: inherit;
  outline: none;
}

.btn-main-action {
  border: 0;
  border-radius: 18px;
  padding: 14px 22px;
  color: #fff;
  background: #2f8f5b;
  font-weight: 900;
  text-decoration: none;
  cursor: pointer;
}

.btn-main-action:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.state-card,
.empty-chat {
  padding: 28px;
  border-radius: 24px;
  color: #63746a;
  background: rgba(255, 255, 255, 0.08);
}

.empty-chat {
  margin: auto;
  max-width: 420px;
  text-align: center;
  background: #f0f6f2;
}

.limit-popup {
  position: fixed;
  z-index: 40;
  right: 28px;
  top: 120px;
  max-width: 440px;
  padding: 28px;
  border-radius: 26px;
  color: #102018;
  background: #fff;
  box-shadow: 0 24px 80px rgba(23, 55, 40, 0.24);
}

.limit-popup__close {
  position: absolute;
  top: 16px;
  right: 18px;
  border: 0;
  background: transparent;
  font-size: 1.5rem;
  cursor: pointer;
}

.limit-popup__actions {
  display: flex;
  gap: 12px;
  margin-top: 18px;
}

@media (max-width: 920px) {
  .messages-hero,
  .messenger-shell {
    grid-template-columns: 1fr;
  }

  .messages-hero {
    flex-direction: column;
  }

  .conversation-list {
    max-height: 360px;
    overflow-y: auto;
  }

  .composer {
    grid-template-columns: 1fr;
  }
}
</style>
