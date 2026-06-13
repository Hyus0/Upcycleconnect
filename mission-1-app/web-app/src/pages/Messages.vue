<template>
    <main class="messages-page public-dashboard">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            variant="public"
        />

        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > MESSAGERIE</p>
                <h1 class="hero-title1">Discussions entre membres</h1>
                <p class="classic-text">
                    Contacte les vendeurs, artisans et utilisateurs autour d'une
                    annonce ou en direct.
                </p>
            </div>
            <RouterLink class="hero-subscription" to="/abonnement">
                <span>DM Plus</span>
                <strong>2,99 € / mois</strong>
                <small>{{ subscriptionLabel }}</small>
            </RouterLink>
        </header>

        <section
            v-if="limitPopup"
            class="limit-popup"
            role="dialog"
            aria-live="polite"
        >
            <button
                type="button"
                class="limit-popup__close"
                @click="limitPopup = null"
            >
                ×
            </button>
            <p class="sidebar__category2">LIMITE ATTEINTE</p>
            <h2
                style="
                    font-family: &quot;Syne&quot;, sans-serif;
                    margin: 10px 0;
                "
            >
                {{ limitPopup.message }}
            </h2>
            <p class="classic-text">
                Les comptes gratuits peuvent contacter 5 vendeurs via annonce.
                L'abonnement DM Plus debloque les messages illimites.
            </p>
            <div class="limit-popup__actions">
                <RouterLink class="btn-main-action" to="/abonnement"
                    >Prendre DM Plus</RouterLink
                >
                <button
                    class="btn-secondary"
                    type="button"
                    @click="limitPopup = null"
                >
                    Plus tard
                </button>
            </div>
        </section>

        <section class="messenger-shell">
            <aside class="conversation-list">
                <div class="panel-heading">
                    <div>
                        <p
                            class="sidebar__category2"
                            style="color: #2f8f5b; margin-bottom: 5px"
                        >
                            BOITE DE RECEPTION
                        </p>
                        <h2 style="font-size: 1.2rem; margin: 0">
                            {{ conversations.length }} conversations
                        </h2>
                    </div>
                    <button
                        class="btn-secondary btn-small"
                        type="button"
                        @click="loadConversations"
                    >
                        ↻
                    </button>
                </div>

                <div v-if="loadingConversations" class="state-card">
                    Chargement...
                </div>
                <div v-else-if="conversations.length === 0" class="state-card">
                    Aucune discussion pour le moment.
                </div>

                <div v-else class="conversation-scroll">
                    <button
                        v-for="conversation in conversations"
                        :key="conversation.id"
                        class="conversation-item"
                        :class="{
                            active: conversation.id === activeConversationId,
                        }"
                        type="button"
                        @click="selectConversation(conversation.id)"
                    >
                        <span class="conversation-avatar">{{
                            initials(conversation.other_user_name)
                        }}</span>
                        <span class="conversation-main">
                            <strong>{{ conversation.other_user_name }}</strong>
                            <small>{{ conversation.annonce_title }}</small>
                            <em>{{
                                conversation.last_message ||
                                "Nouvelle discussion"
                            }}</em>
                        </span>
                        <span
                            v-if="conversation.unread_count"
                            class="unread-badge"
                            >{{ conversation.unread_count }}</span
                        >
                    </button>
                </div>
            </aside>

            <section class="chat-panel">
                <div v-if="!activeConversation" class="empty-chat">
                    <p class="sidebar__category2">SELECTION</p>
                    <h2
                        style="
                            font-family: &quot;Syne&quot;, sans-serif;
                            font-size: 1.5rem;
                            margin: 10px 0;
                        "
                    >
                        Choisis une conversation
                    </h2>
                    <p class="classic-text">
                        Le fil de discussion apparaitra ici avec l'annonce liee
                        et l'historique des messages.
                    </p>
                </div>

                <template v-else>
                    <header class="chat-header">
                        <div
                            style="
                                display: flex;
                                gap: 16px;
                                align-items: center;
                            "
                        >
                            <div
                                class="conversation-avatar conversation-avatar--large"
                            >
                                {{
                                    initials(activeConversation.other_user_name)
                                }}
                            </div>
                            <div>
                                <p
                                    class="sidebar__category2"
                                    style="color: #64766c; margin: 0"
                                >
                                    {{
                                        activeConversation.other_user_role ||
                                        "MEMBRE"
                                    }}
                                </p>
                                <h2
                                    style="
                                        margin: 4px 0;
                                        font-family:
                                            &quot;Syne&quot;, sans-serif;
                                        font-size: 1.3rem;
                                    "
                                >
                                    {{ activeConversation.other_user_name }}
                                    <span
                                        v-if="
                                            activeConversation.other_user_premium
                                        "
                                        class="premium-check"
                                        title="Membre DM Plus"
                                        >✓</span
                                    >
                                </h2>
                                <p
                                    class="classic-text"
                                    style="margin: 0; font-weight: 600"
                                >
                                    {{ activeConversation.annonce_title }}
                                </p>
                            </div>
                        </div>

                        <div
                            v-if="activeConversation.annonce_id"
                            class="header-actions-right"
                        >
                            <span class="annonce-price">{{
                                formatPrice(activeConversation.annonce_price)
                            }}</span>
                            <RouterLink
                                class="btn-secondary btn-small"
                                :to="`/annonce/${activeConversation.annonce_id}`"
                            >
                                Voir l'annonce ❯
                            </RouterLink>
                        </div>
                    </header>

                    <div 
                        v-if="pendingPaymentSale" 
                        class="payment-banner"
                    >
                        <div class="payment-banner-content">
                            <strong>{{ pendingPaymentSale.buyer_id === currentUser ? "Payer l'offre" : "En attente de paiement" }}</strong>
                            <p>{{ pendingPaymentSale.buyer_id === currentUser ? "Le vendeur a accepté votre proposition. Finalisez la commande." : "Vous avez accepté l'offre. L'acheteur doit finaliser le paiement." }}</p>
                        </div>
                        <button 
                            v-if="pendingPaymentSale.buyer_id === currentUser"
                            class="btn-main-action" 
                            style="height: auto; padding: 12px 20px;"
                            type="button"
                            @click="payerOffre(pendingPaymentSale)"
                        >
                            Ajouter au panier et payer {{ formatPrice(pendingPaymentSale.amount) }}
                        </button>
                    </div>

                    <div class="message-thread" ref="messagesContainer">
                        <div
                            v-if="loadingMessages"
                            class="state-card"
                            style="text-align: center"
                        >
                            Chargement des messages...
                        </div>

                        <template
                            v-for="item in chatFeed"
                            :key="item.type + item.id"
                        >
                            <article
                                v-if="item.type === 'text'"
                                class="message-bubble"
                                :class="{
                                    mine: item.sender_id === currentUser,
                                }"
                            >
                                <p>{{ item.content }}</p>
                                <span>{{ formatDate(item.created_at) }}</span>
                            </article>

                            <article
                                v-else-if="item.type === 'offer'"
                                class="offer-bubble"
                                :class="{ mine: item.buyer_id === currentUser }"
                            >
                                <h4 class="offer-title">
                                    🤝 Offre de prix :
                                    {{ formatPrice(item.amount) }}
                                </h4>

                                <span
                                    class="offer-badge"
                                    :class="getOfferBadgeClass(item.status)"
                                >
                                    {{ item.status }}
                                </span>

                                <div
                                    v-if="
                                        item.status === 'En attente' &&
                                        item.seller_id === currentUser
                                    "
                                    class="offer-actions-row"
                                >
                                    <button
                                        type="button"
                                        class="btn-offer-accept"
                                        @click="
                                            handleOfferResponse(
                                                item.id,
                                                'accept',
                                            )
                                        "
                                    >
                                        Accepter
                                    </button>
                                    <button
                                        type="button"
                                        class="btn-offer-refuse"
                                        @click="
                                            handleOfferResponse(
                                                item.id,
                                                'refuse',
                                            )
                                        "
                                    >
                                        Refuser
                                    </button>
                                </div>

                                <span class="offer-date">{{
                                    formatDate(item.created_at)
                                }}</span>
                            </article>

                            <article
                                v-else-if="item.type === 'sale'"
                                class="offer-bubble"
                                :class="{ mine: item.buyer_id === currentUser }"
                            >
                                <h4 class="offer-title">
                                    ✅ Vente conclue :
                                    {{ formatPrice(item.amount) }}
                                </h4>
                                <span class="offer-badge acceptee">{{
                                    item.status
                                }}</span>

                                <button
                                    v-if="
                                        item.buyer_id === currentUser &&
                                        item.status === 'Offre acceptee'
                                    "
                                    type="button"
                                    class="btn-main-action"
                                    style="
                                        margin-top: 10px;
                                        width: 100%;
                                        border-radius: 12px;
                                    "
                                    @click="handleConfirmReception(item.id)"
                                >
                                    J'ai reçu la commande
                                </button>

                                <form
                                    v-if="
                                        item.buyer_id === currentUser &&
                                        item.status === 'Recue'
                                    "
                                    class="review-vertical"
                                    @submit.prevent="handleReviewSale(item.id)"
                                >
                                    <select v-model.number="reviewDraft.note">
                                        <option
                                            v-for="n in 5"
                                            :key="n"
                                            :value="n"
                                        >
                                            {{ n }} / 5 Étoiles
                                        </option>
                                    </select>
                                    <input
                                        v-model="reviewDraft.commentaire"
                                        placeholder="Votre avis sur la transaction"
                                    />
                                    <button
                                        class="btn-main-action"
                                        style="border-radius: 12px"
                                        type="submit"
                                    >
                                        Noter le vendeur
                                    </button>
                                </form>

                                <span class="offer-date">{{
                                    formatDate(item.updated_at)
                                }}</span>
                            </article>
                        </template>
                    </div>

                    <div class="composer-wrapper">
                        <div v-if="showNegotiation" class="negotiation-flyout">
                            <button
                                class="btn-close-flyout"
                                type="button"
                                @click="showNegotiation = false"
                            >
                                ×
                            </button>
                            <p class="flyout-eyebrow">NEGOCIATION</p>
                            <h3 class="flyout-title">
                                {{ activeConversation.annonce_title }}
                            </h3>
                            <p class="flyout-price">
                                Prix annonce :
                                <strong>{{
                                    formatPrice(
                                        activeConversation.annonce_price,
                                    )
                                }}</strong>
                            </p>

                            <form
                                class="flyout-form"
                                @submit.prevent="submitOffer"
                            >
                                <input
                                    v-model="offerAmount"
                                    type="number"
                                    min="1"
                                    step="0.01"
                                    class="flyout-input"
                                    placeholder="Votre offre (€)"
                                />
                                <button
                                    class="flyout-btn"
                                    type="submit"
                                    :disabled="!offerAmount || creatingOffer"
                                >
                                    {{ creatingOffer ? "..." : "Proposer" }}
                                </button>
                            </form>
                        </div>

                        <form class="composer" @submit.prevent="handleSend">
                            <button
                                v-if="
                                    activeConversation.annonce_id &&
                                    activeConversation.annonce_seller_id !==
                                        currentUser
                                "
                                type="button"
                                class="btn-composer btn-composer-neg"
                                @click="showNegotiation = !showNegotiation"
                            >
                                Négocier
                            </button>

                            <textarea
                                v-model="draft"
                                placeholder="Ecris ton message..."
                            />

                            <button
                                class="btn-composer btn-composer-send"
                                type="submit"
                                :disabled="sending || !draft.trim()"
                            >
                                Envoyer
                            </button>
                        </form>
                    </div>
                </template>
            </section>
        </section>
    </main>
</template>

<script setup>
import { computed, onMounted, ref, watch, nextTick } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";

const route = useRoute();
const router = useRouter();
const API_URL = "http://localhost:8081";

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
async function fetchConversations(userId) {
    const res = await fetch(`${API_URL}/users/${userId}/messages`, { headers: getHeaders() });
    if (!res.ok) throw new Error("Erreur");
    return await res.json();
}
async function startConversation(payload) {
    const userId = currentUserId();
    const res = await fetch(`${API_URL}/users/${userId}/messages/start`, {
        method: "POST", headers: getHeaders(), body: JSON.stringify(payload),
    });
    if (!res.ok) throw await res.json();
    return await res.json();
}
async function fetchMessages(conversationId, userId) {
    const res = await fetch(`${API_URL}/users/${userId}/messages/${conversationId}`, { headers: getHeaders() });
    if (!res.ok) throw new Error("Erreur");
    return await res.json();
}
async function sendMessage(conversationId, texteMessage, userId) {
    const res = await fetch(`${API_URL}/users/${userId}/messages/${conversationId}`, {
        method: "POST", headers: getHeaders(), body: JSON.stringify({ content: texteMessage }),
    });
    if (!res.ok) throw new Error("Erreur");
    return await res.json();
}
async function fetchConversationState(conversationId, userId) {
    const res = await fetch(`${API_URL}/users/${userId}/messages/${conversationId}/state`, { headers: getHeaders() });
    if (!res.ok) throw new Error("Erreur");
    return await res.json();
}
async function createOffer(conversationId, amount, userId) {
    const res = await fetch(`${API_URL}/users/${userId}/messages/${conversationId}/offers`, {
        method: "POST", headers: getHeaders(), body: JSON.stringify({ amount: Number(amount) }),
    });
    if (!res.ok) throw new Error("Erreur");
    return await res.json();
}
async function respondOffer(offerId, action, userId) {
    const res = await fetch(`${API_URL}/users/${userId}/messages/offers/${offerId}`, {
        method: "PATCH", headers: getHeaders(), body: JSON.stringify({ action }),
    });
    if (!res.ok) throw new Error("Erreur");
    return await res.json();
}
async function confirmSaleReception(saleId, userId) {
    const res = await fetch(`${API_URL}/users/${userId}/messages/sales/${saleId}/reception`, {
        method: "POST", headers: getHeaders(),
    });
    if (!res.ok) throw new Error("Erreur");
    return await res.json();
}
async function reviewSale(saleId, reviewDraft, userId) {
    const res = await fetch(`${API_URL}/users/${userId}/messages/sales/${saleId}/review`, {
        method: "POST", headers: getHeaders(), body: JSON.stringify(reviewDraft),
    });
    if (!res.ok) throw new Error("Erreur");
    return await res.json();
}

const conversations = ref([]);
const messages = ref([]);
const activeConversationId = ref(Number(route.query.conversation) || 0);
const loadingConversations = ref(false);
const loadingMessages = ref(false);
const sending = ref(false);
const draft = ref("");
const offerAmount = ref("");
const creatingOffer = ref(false);
const showNegotiation = ref(false);
const threadState = ref({ offers: [], sales: [] });
const reviewDraft = ref({ note: 5, commentaire: "" });
const subscription = ref({ is_subscriber: false, used: 0, limit: 5 });
const limitPopup = ref(null);
const messagesContainer = ref(null);

const currentUser = computed(() => currentUserId());
const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return `${prenom} ${nom}`.trim() || "Utilisateur";
});

const activeConversation = computed(() =>
    conversations.value.find((c) => Number(c.id) === Number(activeConversationId.value)),
);

const subscriptionLabel = computed(() =>
    subscription.value.is_subscriber
        ? "Messagerie illimitée"
        : `${subscription.value.used || 0}/${subscription.value.limit || 5} gratuits`,
);

const chatFeed = computed(() => {
    const feed = [];
    
    messages.value.forEach((m) => {
        const content = m.content || "";
        if (
            content.startsWith("Offre proposee") ||
            content.startsWith("Offre acceptee") ||
            content.startsWith("Offre refusee") ||
            content.startsWith("Reception confirmee") ||
            content.startsWith("Vente notee")
        ) {
            return;
        }
        feed.push({
            ...m,
            type: "text",
            timestamp: new Date(m.created_at.replace(/Z$/, "")).getTime(),
        });
    });

    threadState.value.offers.forEach((o) => {
        if (o.status !== "Acceptee") {
            feed.push({
                ...o,
                type: "offer",
                timestamp: new Date(o.created_at.replace(/Z$/, "")).getTime() + 10,
            });
        }
    });

    threadState.value.sales.forEach((s) => {
        feed.push({
            ...s,
            type: "sale",
            timestamp: new Date((s.reviewed_at || s.received_at || s.updated_at || s.created_at || new Date()).replace(/Z$/, "")).getTime() + 20,
        });
    });

    return feed.sort((a, b) => a.timestamp - b.timestamp);
});

function initials(name) {
    return (name || "UC").split(" ").filter(Boolean).slice(0, 2).map((part) => part[0]?.toUpperCase()).join("");
}

function formatDate(value) {
    if (!value) return "";
    const cleanDate = value.replace(/Z$/, "");
    const date = new Date(cleanDate);
    if (Number.isNaN(date.getTime())) return "";

    const day = date.getDate();
    const month = date.toLocaleString("fr-FR", { month: "long" });
    const time = date.toLocaleString("fr-FR", { hour: "2-digit", minute: "2-digit" }).replace(":", "h");
    return `${day} ${month}, ${time}`;
}

function formatPrice(value) {
    return new Intl.NumberFormat("fr-FR", { style: "currency", currency: "EUR" }).format(Number(value) || 0);
}

function getOfferBadgeClass(status) {
    const s = (status || "").toLowerCase();
    if (s.includes("attente")) return "en-attente";
    if (s.includes("accept") || s.includes("recue") || s.includes("evalue")) return "acceptee";
    if (s.includes("refus") || s.includes("annul")) return "refusee";
    return "";
}

function scrollToBottom() {
    nextTick(() => {
        if (messagesContainer.value) {
            messagesContainer.value.scrollTo({
                top: messagesContainer.value.scrollHeight,
                behavior: 'smooth'
            });
        }
    });
}

const pendingPaymentSale = computed(() => {
    if (!threadState.value || !threadState.value.sales) return null;
    return threadState.value.sales.find(
        (sale) => sale.status === 'Offre acceptee'
    );
});

async function loadSubscription() {
    if (!currentUser.value) return;
    try {
        subscription.value = await fetchSubscriptionStatus(currentUser.value);
    } catch (e) {}
}

async function loadConversations() {
    if (!currentUser.value) return router.push("/connexion");
    
    if(conversations.value.length === 0) loadingConversations.value = true;
    try {
        conversations.value = await fetchConversations(currentUser.value);
    } catch (e) {
    } finally {
        loadingConversations.value = false;
    }
}

async function selectConversation(id) {
    if (activeConversationId.value === id) {
       await fetchOnlyMessages(id);
       return;
    }
    
    activeConversationId.value = id;
    showNegotiation.value = false;
    
    router.replace({ path: "/messages", query: { conversation: id } }).catch(() => {});
    await loadMessages();
}

async function fetchOnlyMessages(id) {
     try {
        messages.value = await fetchMessages(id, currentUser.value);
        threadState.value = await fetchConversationState(id, currentUser.value);
        conversations.value = await fetchConversations(currentUser.value);
        scrollToBottom();
    } catch (e) { console.error(e); }
}

async function loadMessages() {
    if (!activeConversationId.value) return;
    loadingMessages.value = true;
    try {
        messages.value = await fetchMessages(activeConversationId.value, currentUser.value);
        threadState.value = await fetchConversationState(activeConversationId.value, currentUser.value);
        conversations.value = await fetchConversations(currentUser.value);
        scrollToBottom();
    } catch (e) {
    } finally {
        loadingMessages.value = false;
    }
}

async function payerOffre(sale) {
    const annonceId = activeConversation.value.annonce_id;
    if (!annonceId) return;

    try {
        const res = await fetch(`${API_URL}/users/${currentUser.value}/panier`, {
            method: "POST",
            headers: getHeaders(),
            body: JSON.stringify({
                type_item: "Annonce",
                reference_id: annonceId,
                prix_unitaire: sale.amount
            })
        });

        if (!res.ok) throw new Error("Impossible d'ajouter au panier.");

        router.push("/profil");
        
    } catch (error) {
        alert(error.message);
    }
}

async function submitOffer() {
    if (!activeConversationId.value || !offerAmount.value) return;
    creatingOffer.value = true;
    try {
        await createOffer(activeConversationId.value, offerAmount.value, currentUser.value);
        offerAmount.value = "";
        showNegotiation.value = false;
        await fetchOnlyMessages(activeConversationId.value);
    } catch (error) {
        alert(error.message || "Impossible de proposer cette offre.");
    } finally {
        creatingOffer.value = false;
    }
}

async function handleOfferResponse(offerId, action) {
    try {
        await respondOffer(offerId, action, currentUser.value);
        await fetchOnlyMessages(activeConversationId.value);
    } catch (error) {
        alert(error.message);
    }
}

async function handleConfirmReception(saleId) {
    try {
        await confirmSaleReception(saleId, currentUser.value);
        await fetchOnlyMessages(activeConversationId.value);
    } catch (error) {
        alert(error.message);
    }
}

async function handleReviewSale(saleId) {
    try {
        await reviewSale(saleId, reviewDraft.value, currentUser.value);
        reviewDraft.value = { note: 5, commentaire: "" };
        await fetchOnlyMessages(activeConversationId.value);
    } catch (error) {
        alert(error.message);
    }
}

async function startFromQuery() {
    const target = Number(route.query.user);
    const annonce = Number(route.query.annonce);
    if (!target) return;
    try {
        const result = await startConversation({
            target_user_id: target,
            annonce_id: annonce || null,
        });
        activeConversationId.value = result.conversation_id;
        await selectConversation(result.conversation_id);
    } catch (error) {
        if (error.status === 402) {
            limitPopup.value = error.payload || { message: error.message };
        } else {
            alert(error.message || "Impossible.");
        }
    }
}

async function handleSend() {
    if (!draft.value.trim() || !activeConversationId.value) return;
    sending.value = true;
    try {
        const message = await sendMessage(
            activeConversationId.value,
            draft.value.trim(),
            currentUser.value,
        );
        messages.value.push(message);
        draft.value = "";
        scrollToBottom();
        conversations.value = await fetchConversations(currentUser.value);
    } catch (e) {
    } finally {
        sending.value = false;
    }
}

watch(
    () => route.query.conversation,
    async (value) => {
        const id = Number(value);
        if (id && id !== activeConversationId.value) {
            activeConversationId.value = id;
            showNegotiation.value = false;
            await loadMessages();
        }
    },
);

onMounted(async () => {
    if (!isLoggedIn.value) return router.push("/connexion");
    await loadSubscription();
    await loadConversations();

    if (route.query.user && route.query.annonce) {
        await startFromQuery();
    } else if (activeConversationId.value) {
        await loadMessages();
    } else if (conversations.value.length > 0) {
        await selectConversation(conversations.value[0].id);
    }
});
</script>

<style scoped>
.public-dashboard {
    min-height: 100vh; 
    display: flex;
    flex-direction: column;
    padding: 0 20px 20px 20px;
    box-sizing: border-box;
    background: var(--bg-light, #f7f9f7);
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

.hero-title1 {
    font-size: 2.2rem;
    font-weight: 900;
    color: #1a1a1a;
    margin: 5px 0 0 0;
    font-family: "Syne", sans-serif;
    white-space: nowrap;
}

.classic-text {
    font-size: 0.95rem;
    color: #6d7b72;
    margin-top: 8px;
    line-height: 1.5;
}

.hero-subscription {
    min-width: 220px;
    padding: 18px 24px;
    border-radius: 20px;
    color: #eef8f0;
    background: linear-gradient(145deg, #102018, #2f8f5b);
    text-decoration: none;
    box-shadow: 0 10px 30px rgba(23, 55, 40, 0.15);
    text-align: center;
}
.hero-subscription span,
.hero-subscription small {
    display: block;
    color: rgba(255, 255, 255, 0.8);
}
.hero-subscription strong {
    display: block;
    margin: 6px 0;
    font-size: 1.5rem;
    font-weight: 900;
}

.messenger-shell {
    height: 75vh;         
      min-height: 650px;
    display: grid;
    grid-template-columns: 360px 1fr;
    border: 1px solid rgba(29, 56, 42, 0.12);
    border-radius: 24px;
    background: rgba(255, 255, 255, 0.76);
    box-shadow: 0 14px 40px rgba(26, 52, 38, 0.08);
}

.conversation-list {
    background: #102018;
    color: #fff;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    border-right: 1px solid rgba(255, 255, 255, 0.1);
    border-top-left-radius: 24px;
    border-bottom-left-radius: 24px;
}

.panel-heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 24px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
    flex-shrink: 0;
}
.panel-heading h2 {
    margin: 0;
    font-family: "Syne", sans-serif;
}

.conversation-scroll {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
}

.btn-secondary {
    border: 1px solid rgba(47, 143, 91, 0.28);
    border-radius: 12px;
    padding: 10px 16px;
    color: #1d3528;
    background: #eef6f0;
    font-weight: 800;
    text-decoration: none;
    cursor: pointer;
    transition: 0.2s;
}
.btn-secondary:hover {
    background: #e1ede5;
}
.btn-small {
    padding: 6px 12px;
    font-size: 0.85rem;
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
    border-radius: 18px;
    color: #dce9e1;
    background: rgba(255, 255, 255, 0.04);
    text-align: left;
    cursor: pointer;
    transition: 0.2s;
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
    width: 56px;
    height: 56px;
    border-radius: 18px;
    font-size: 1.2rem;
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
    flex-direction: column;
    background: #fbfdfb;
    height: 100%;
    min-height: 0;
    overflow: hidden;
    border-top-right-radius: 24px;
    border-bottom-right-radius: 24px;
}

.chat-header {
    flex-shrink: 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 26px;
    border-bottom: 1px solid #dfe9e2;
    background: #fff;
}
.header-actions-right {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 8px;
}
.annonce-price {
    font-size: 1.3rem;
    font-weight: 900;
    color: #2f8f5b;
}
.premium-check {
    display: inline-grid;
    width: 20px;
    height: 20px;
    place-items: center;
    margin-left: 8px;
    border-radius: 999px;
    color: #102018;
    background: #8ef0a8;
    font-size: 0.75rem;
    font-weight: 900;
    vertical-align: middle;
}

.message-thread {
    flex: 1 1 0; 
    min-height: 0;
    padding: 28px;
    overflow-y: auto; 
    background:
        linear-gradient(rgba(16, 32, 24, 0.025) 1px, transparent 1px),
        linear-gradient(90deg, rgba(16, 32, 24, 0.025) 1px, transparent 1px);
    background-size: 28px 28px;
    display: flex;
    flex-direction: column;
}

.message-bubble {
    max-width: 68%;
    margin-bottom: 14px;
    padding: 14px 18px;
    border-radius: 20px 20px 20px 6px;
    background: #ffffff;
    border: 1px solid #dfe9e2;
    box-shadow: 0 8px 20px rgba(32, 58, 43, 0.05);
    align-self: flex-start;
}
.message-bubble.mine {
    align-self: flex-end;
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
    font-size: 0.75rem;
}

   .offer-bubble {
     min-width: 280px;
     padding: 24px; 
     background: #f4f9f6;
     border: 1px solid #e1ede5;
     border-radius: 24px;
     color: #122018;
     align-self: flex-start;
     margin-bottom: 16px;
     display: flex;
     flex-direction: column;
     gap: 16px; 
   }
   
   .offer-bubble.mine {
     align-self: flex-end;
     background: #ffffff;
   }
   
   .offer-title {
     font-size: 1.25rem;
     font-weight: 800;
     margin: 0;
     display: flex;
     align-items: center;
     gap: 8px;
   }
   
   .offer-badge {
     align-self: flex-start;
     display: inline-block;
     padding: 8px 16px;
     border-radius: 12px;
     font-size: 0.9rem;
     font-weight: 800;
   }
   
   .offer-badge.en-attente { background: #fdf3e1; color: #c2883a; }
   .offer-badge.acceptee { background: #e3f2e8; color: #2a7d4d; }
   .offer-badge.refusee { background: #fceaea; color: #c94b4b; }
   
   .offer-actions-row {
     display: flex;
     gap: 12px;
   }
   
   .btn-offer-accept {
     background: #308d58;
     color: #ffffff;
     border: none;
     border-radius: 12px;
     padding: 12px 20px;
     font-size: 0.95rem;
     font-weight: 800;
     cursor: pointer;
     transition: 0.2s;
   }
   .btn-offer-accept:hover { background: #23653e; }
   
   .btn-offer-refuse {
     background: #fdfdfd;
     color: #102018;
     border: 1px solid #d7e5dc;
     border-radius: 12px;
     padding: 12px 20px;
     font-size: 0.95rem;
     font-weight: 800;
     cursor: pointer;
     transition: 0.2s;
   }
   .btn-offer-refuse:hover { background: #f0f6f2; }
   
   .offer-date {
     display: block;
     font-size: 0.85rem;
     color: #819588;
     margin: 0;
   }

.review-vertical {
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin-top: 14px;
}
.review-vertical select,
.review-vertical input {
    width: 100%;
    border: 1px solid #d4e3d8;
    border-radius: 12px;
    padding: 12px;
    font: inherit;
}

.btn-main-action {
    border: 0;
    border-radius: 12px;
    padding: 10px 18px;
    color: #fff;
    background: #2f8f5b;
    font-weight: 900;
    text-decoration: none;
    cursor: pointer;
    transition: 0.2s;
}
.btn-main-action:disabled {
    opacity: 0.55;
    cursor: not-allowed;
}
.btn-main-action:hover:not(:disabled) {
    background: #23653e;
}

.composer-wrapper {
    position: relative;
    flex-shrink: 0;
    background: #ffffff;
    border-top: 1px solid #dfe9e2;
    padding: 16px 24px;
}

.composer {
    display: flex;
    flex-direction: row;
    gap: 12px;
    align-items: center;
}

.composer textarea {
    flex: 1;
    height: 54px;
    min-height: 54px;
    resize: none;
    border: 1px solid #d4e3d8;
    border-radius: 16px;
    padding: 16px 18px;
    font: inherit;
    outline: none;
    background: #fafafa;
    margin: 0;
}
.composer textarea:focus {
    border-color: #2f8f5b;
    background: #fff;
}

.btn-composer {
    flex: 0 0 auto;
    height: 54px;
    padding: 0 24px;
    border-radius: 16px;
    font-weight: 900;
    font-size: 0.95rem;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    border: 1px solid transparent;
    transition: 0.2s;
    white-space: nowrap;
}

.btn-composer-neg {
    background: #fff4d8;
    color: #9a6400;
    border-color: #ffe8b3;
}
.btn-composer-neg:hover {
    background: #ffe8b3;
}
.btn-composer-send {
    background: #2f8f5b;
    color: #fff;
}
.btn-composer-send:disabled {
    opacity: 0.55;
    cursor: not-allowed;
}

.negotiation-flyout {
    position: absolute;
    bottom: calc(100% + 15px);
    left: 24px;
    width: 400px;
    padding: 24px;
    background: #ffffff;
    border: 1px solid #e1ede5;
    border-radius: 24px;
    box-shadow: 0 12px 36px rgba(29, 56, 42, 0.12);
    z-index: 10;
    animation: slideUp 0.2s ease-out;
}

.btn-close-flyout {
    position: absolute;
    top: 16px;
    right: 16px;
    background: transparent;
    border: none;
    font-size: 1.4rem;
    color: #8fa396;
    cursor: pointer;
}

.flyout-eyebrow {
    font-size: 0.75rem;
    color: #3b8b5d;
    font-weight: 800;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    margin: 0 0 8px 0;
}

.flyout-title {
    font-size: 1.4rem;
    color: #122018;
    font-weight: 800;
    font-family: "Syne", sans-serif;
    margin: 0 0 4px 0;
}

.flyout-price {
    font-size: 1rem;
    color: #63746a;
    margin: 0 0 20px 0;
}
.flyout-price strong {
    color: #63746a;
    font-weight: 800;
}

.flyout-form {
    display: flex;
    gap: 12px;
}

.flyout-input {
    flex: 1;
    height: 44px;
    padding: 0 16px;
    border: 1px solid #d7e5dc;
    border-radius: 12px;
    background: #fbfdfb;
    font-size: 0.95rem;
    outline: none;
}

.flyout-btn {
    height: 44px;
    padding: 0 20px;
    background: #2f8f5b;
    color: #ffffff;
    border: none;
    border-radius: 12px;
    font-weight: 800;
    cursor: pointer;
}

.state-card,
.empty-chat {
    padding: 32px;
    border-radius: 20px;
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

@keyframes slideUp {
    from {
        opacity: 0;
        transform: translateY(10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

@media (max-width: 920px) {
    .public-dashboard {
        height: auto;
        min-height: 100vh;
        overflow: auto;
    }
    .messenger-shell {
        grid-template-columns: 1fr;
        height: 800px;
    }
    .conversation-list {
        max-height: 360px;
        overflow-y: auto;
        border-right: none;
    }
    .hero-title1 {
        white-space: normal;
    }
}

.payment-banner {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #fff4e6;
    border-bottom: 1px solid #ffe8cc;
    padding: 16px 26px;
    flex-shrink: 0; 
}

.payment-banner-content strong {
    display: block;
    color: #cc6600;
    font-size: 1.1rem;
    font-family: "Syne", sans-serif;
    margin-bottom: 4px;
}

.payment-banner-content p {
    margin: 0;
    color: #994d00;
    font-size: 0.9rem;
}
</style>
