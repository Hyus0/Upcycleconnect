<template>
    <main class="page-main-content">
        <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" />

        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > PAIEMENT</p>
                <h1 class="hero-title1">Paiement sécurisé</h1>
                <p class="classic-text">
                    Finalisez votre commande en toute sécurité.
                </p>
            </div>
            <div class="header-actions">
                <button class="btn-secondary-back" @click="$router.back()">
                    🠔 Retour
                </button>
            </div>
        </header>

        <div v-if="loading" class="loading-state">
            Préparation de votre paiement...
        </div>

        <div v-else-if="checkoutResult" class="invoice-success-card">
            <span class="invoice-kicker">Paiement validé</span>
            <h2>Votre commande est confirmée.</h2>
            <p>
                Facture
                <strong>{{ checkoutResult.numero_facture }}</strong> générée
                pour la commande #{{ checkoutResult.commande_id }}.
            </p>
            <div class="invoice-actions">
                <button class="btn-main-action" @click="downloadInvoice">
                    Télécharger la facture
                </button>
                <button
                    class="btn-secondary-back"
                    :disabled="sendingInvoice"
                    @click="sendInvoiceByMail"
                >
                    {{ sendingInvoice ? "Envoi..." : "Envoyer par mail" }}
                </button>
                <button
                    class="btn-secondary-back"
                    @click="$router.push('/profil/factures')"
                >
                    Voir mes factures
                </button>
            </div>
        </div>

        <div v-else-if="cartItems.length === 0" class="empty-cart-card">
            <h2>Aucun article à payer</h2>
            <p>Votre panier est vide.</p>
            <button
                class="btn-main-action mt-4"
                @click="$router.push('/catalogue')"
            >
                Retourner au catalogue
            </button>
        </div>

        <div v-else class="checkout-layout-wrapper">
            <div class="checkout-layout">
                <div class="checkout-form-column">
                    <div class="payment-card">
                        <div class="payment-header">
                            <h2>Carte Bancaire</h2>
                            <div class="card-icons">💳</div>
                        </div>

                        <form
                            @submit.prevent="submitPayment"
                            class="checkout-form"
                        >
                            <div class="form-group">
                                <label>Nom sur la carte</label>
                                <input
                                    type="text"
                                    v-model="paymentData.name"
                                    required
                                    placeholder="Jean Dupont"
                                />
                            </div>

                            <div class="form-group">
                                <label>Numéro de carte</label>
                                <input
                                    type="text"
                                    v-model="paymentData.card"
                                    required
                                    placeholder="0000 0000 0000 0000"
                                    maxlength="19"
                                />
                            </div>

                            <div class="form-row">
                                <div class="form-group">
                                    <label>Date d'expiration</label>
                                    <input
                                        type="text"
                                        v-model="paymentData.exp"
                                        required
                                        placeholder="MM/AA"
                                        maxlength="5"
                                    />
                                </div>
                                <div class="form-group">
                                    <label>Code de sécurité (CVC)</label>
                                    <input
                                        type="text"
                                        v-model="paymentData.cvc"
                                        required
                                        placeholder="123"
                                        maxlength="3"
                                    />
                                </div>
                            </div>

                            <button
                                type="submit"
                                class="btn-checkout"
                                :disabled="isCheckingOut"
                            >
                                <span v-if="isCheckingOut"
                                    >Traitement sécurisé en cours...</span
                                >
                                <span v-else
                                    >Payer {{ formatPrice(cartTotal) }}</span
                                >
                            </button>
                            <p class="secure-payment-text">
                                🔒 Connexion chiffrée de bout en bout
                            </p>
                        </form>
                    </div>
                </div>

                <div class="checkout-summary-column">
                    <div class="summary-card">
                        <h2>Résumé de la commande</h2>

                        <div class="summary-items-list">
                            <div
                                v-for="item in cartItems"
                                :key="item.id"
                                class="summary-item"
                            >
                                <div class="summary-item-info">
                                    <span class="summary-item-type">{{
                                        item.type_item
                                    }}</span>
                                    <span class="summary-item-title">{{
                                        getItemName(item)
                                    }}</span>
                                </div>
                                <div class="summary-item-price">
                                    {{ formatPrice(item.prix_unitaire) }}
                                </div>
                            </div>
                        </div>

                        <div class="summary-divider"></div>

                        <div class="summary-line">
                            <span>Sous-total</span>
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
                    </div>
                </div>
            </div>
        </div>
    </main>
    <SiteFooter />
</template>

<script setup>
import { ref, onMounted, computed, reactive } from "vue";
import { useRouter, useRoute } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

const router = useRouter();
const route = useRoute();
const API_URL = "http://localhost:8081";

const loading = ref(true);
const cartItems = ref([]);
const isCheckingOut = ref(false);
const checkoutResult = ref(null);
const sendingInvoice = ref(false);

const paymentData = reactive({
    name: "Test User",
    card: "4000 0000 0000 0000",
    exp: "12/25",
    cvc: "123",
});

const getHeaders = () => ({
    "Content-Type": "application/json",
    Authorization: sessionStorage.getItem("userToken") || "",
});

const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const currentUserId = computed(
    () =>
        Number(
            sessionStorage.getItem("id") || sessionStorage.getItem("userId"),
        ) || 0,
);
const userName = computed(
    () => sessionStorage.getItem("userPrenom") || "Utilisateur",
);

const isDirectBuy = computed(() => !!route.query.annonce_id);

const cartTotal = computed(() =>
    cartItems.value.reduce(
        (total, item) => total + Number(item.prix_unitaire || 0),
        0,
    ),
);

const formatPrice = (price) =>
    new Intl.NumberFormat("fr-FR", {
        style: "currency",
        currency: "EUR",
    }).format(Number(price) || 0);

const getItemName = (item) => {
    if (item.titre) return item.titre;
    if (item.type_item === "Formation")
        return `Billet pour la session de formation`;
    if (item.type_item === "Abonnement") return `Abonnement DM Plus`;
    if (item.type_item === "Annonce") return `Achat d'annonce`;
    return `Article #${item.reference_id}`;
};

const fetchCheckoutData = async () => {
    if (!currentUserId.value) return router.push("/connexion");

    loading.value = true;

    if (isDirectBuy.value) {
        cartItems.value = [
            {
                id: "direct",
                type_item: "Annonce",
                reference_id: route.query.annonce_id,
                titre: route.query.title || "Achat négocié",
                prix_unitaire: Number(route.query.amount || 0),
            },
        ];
        loading.value = false;
    } else {
        try {
            const res = await fetch(
                `${API_URL}/users/${currentUserId.value}/panier`,
                { headers: getHeaders() },
            );
            if (res.ok) {
                cartItems.value = (await res.json()) || [];
            }
        } catch (error) {
            console.error("Erreur panier:", error);
        } finally {
            loading.value = false;
        }
    }
};

const submitPayment = async () => {
    isCheckingOut.value = true;

    await new Promise((resolve) => setTimeout(resolve, 1500));

    try {
        let res;

        if (isDirectBuy.value) {
            res = await fetch(
                `${API_URL}/annonces/${route.query.annonce_id}/acheter`,
                {
                    method: "POST",
                    headers: getHeaders(),
                    body: JSON.stringify({
                        montant: Number(route.query.amount),
                        user_id: currentUserId.value,
                    }),
                },
            );
        } else {
            res = await fetch(
                `${API_URL}/users/${currentUserId.value}/checkout`,
                {
                    method: "POST",
                    headers: getHeaders(),
                },
            );
        }

        if (res.ok) {
            checkoutResult.value = (await res.json()) || { message: "Succès" };
            cartItems.value = [];
        } else {
            alert("Erreur lors de la commande : " + (await res.text()));
        }
    } catch (error) {
        alert("Impossible de contacter le serveur.");
    } finally {
        isCheckingOut.value = false;
    }
};

const downloadInvoice = () => {
    if (!checkoutResult.value?.facture_id) return;
    window.open(
        `${API_URL}/users/${currentUserId.value}/factures/${checkoutResult.value.facture_id}/download`,
        "_blank",
    );
};

const sendInvoiceByMail = async () => {
    if (!checkoutResult.value?.facture_id) return;
    sendingInvoice.value = true;
    try {
        const res = await fetch(
            `${API_URL}/users/${currentUserId.value}/factures/${checkoutResult.value.facture_id}/send`,
            {
                method: "POST",
                headers: getHeaders(),
            },
        );
        const payload = res.ok
            ? await res.json()
            : { message: await res.text() };
        alert(payload.message);
    } finally {
        sendingInvoice.value = false;
    }
};

onMounted(fetchCheckoutData);
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
    border-bottom: 1px solid #dce9e1;
}

.sidebar__category2 {
    font-size: 0.75rem;
    color: #2d7a4f;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 1px;
    margin-bottom: 5px;
}
.hero-title1 {
    font-size: 2.2rem;
    font-weight: 900;
    color: #1a1a1a;
    margin: 0;
}
.classic-text {
    color: #6d7b72;
    margin-top: 8px;
}

.btn-secondary-back {
    padding: 10px 20px;
    border-radius: 12px;
    border: 1px solid #dce9e1;
    background: white;
    color: #2d7a4f;
    cursor: pointer;
    font-weight: 700;
    transition: 0.2s;
}
.btn-secondary-back:hover {
    background: #f0f4f1;
}

.checkout-layout-wrapper {
    max-width: 1100px;
    margin: 0 auto;
}
.checkout-layout {
    display: grid;
    grid-template-columns: 1fr 400px;
    gap: 40px;
    align-items: start;
}

.payment-card {
    background: #fff;
    border: 1px solid #e5ede7;
    border-radius: 20px;
    padding: 36px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.03);
}

.payment-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 28px;
}
.payment-header h2 {
    margin: 0;
    font-size: 1.6rem;
    color: #1a1a1a;
}
.card-icons {
    font-size: 2rem;
}

.checkout-form {
    display: flex;
    flex-direction: column;
    gap: 20px;
}
.form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.form-group label {
    font-size: 0.9rem;
    font-weight: 700;
    color: #1a1a1a;
}
.form-group input {
    padding: 16px;
    border: 1px solid #dce9e1;
    border-radius: 12px;
    font-size: 1.05rem;
    font-family: inherit;
    background: #fbfdfb;
    transition: 0.2s;
}
.form-group input:focus {
    border-color: #2d7a4f;
    background: #fff;
    outline: none;
    box-shadow: 0 0 0 3px rgba(45, 122, 79, 0.1);
}
.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
}

.btn-checkout {
    margin-top: 10px;
    padding: 18px;
    background: #2d7a4f;
    color: white;
    border: none;
    border-radius: 14px;
    font-size: 1.2rem;
    font-weight: 800;
    font-family: inherit;
    cursor: pointer;
    transition: 0.2s;
}
.btn-checkout:hover:not(:disabled) {
    background: #23653e;
    transform: translateY(-2px);
}
.btn-checkout:disabled {
    background: #9bcbae;
    cursor: not-allowed;
}

.secure-payment-text {
    text-align: center;
    font-size: 0.85rem;
    color: #8fa396;
    margin-top: 10px;
    font-weight: 600;
}

.summary-card {
    background: #fcfdfc;
    border: 1px solid #e5ede7;
    border-radius: 20px;
    padding: 30px;
    position: sticky;
    top: 40px;
}
.summary-card h2 {
    font-size: 1.3rem;
    margin: 0 0 24px;
    border-bottom: 1px solid #dce9e1;
    padding-bottom: 16px;
}

.summary-items-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin-bottom: 24px;
}
.summary-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
}
.summary-item-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
    flex: 1;
}
.summary-item-type {
    font-size: 0.7rem;
    font-weight: 800;
    color: #8fa396;
    text-transform: uppercase;
}
.summary-item-title {
    font-size: 0.95rem;
    font-weight: 700;
    color: #1a1a1a;
    line-height: 1.3;
}
.summary-item-price {
    font-size: 1.1rem;
    font-weight: 800;
    color: #2d7a4f;
}

.summary-line {
    display: flex;
    justify-content: space-between;
    margin-bottom: 12px;
    color: #6d7b72;
    font-size: 0.95rem;
}
.summary-divider {
    height: 1px;
    background: #dce9e1;
    margin: 20px 0;
}
.total-line {
    font-size: 1.4rem;
    color: #1a1a1a;
    font-weight: 900;
    margin-bottom: 0;
}
.total-line span:last-child {
    color: #2d7a4f;
}

.loading-state {
    text-align: center;
    padding: 4rem;
    color: #8fa396;
    font-weight: 600;
}
.empty-cart-card {
    background: #fff;
    border: 1px dashed #cfe0d4;
    border-radius: 16px;
    padding: 60px 20px;
    text-align: center;
    max-width: 500px;
    margin: 40px auto;
}
.btn-main-action {
    background: #2d7a4f;
    color: white;
    padding: 12px 24px;
    border-radius: 10px;
    font-weight: 700;
    border: none;
    cursor: pointer;
    transition: 0.2s;
}
.btn-main-action:hover {
    background: #23653e;
}

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
}
.invoice-success-card h2 {
    margin: 18px 0 10px;
    font-size: 2rem;
    color: #16221c;
}
.invoice-success-card p {
    color: #5d6d64;
}
.invoice-actions {
    display: flex;
    gap: 12px;
    margin: 24px 0 0;
}

@media (max-width: 1024px) {
    .checkout-layout {
        grid-template-columns: 1fr;
    }
    .summary-card {
        position: static;
    }
}
@media (max-width: 600px) {
    .form-row {
        grid-template-columns: 1fr;
    }
    .invoice-actions {
        flex-direction: column;
    }
}
</style>
