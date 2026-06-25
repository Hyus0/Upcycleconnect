<template>
    <div class="layout-wrapper bg-light">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            variant="public"
        />

        <main class="page-container">
            <header class="content-header">
                <div class="header-left">
                    <p class="sidebar__category2">ACCUEIL > PAIEMENT</p>
                    <h1 class="hero-title1">Paiement sécurisé</h1>
                    <p class="classic-text">
                        Finalisez votre transaction en toute sécurité.
                    </p>
                </div>
                <button
                    class="btn-secondary"
                    type="button"
                    @click="$router.back()"
                >
                    🠔 Retour
                </button>
            </header>

            <div v-if="paymentMode" class="checkout-grid">
                <div class="checkout-form card">
                    <h2 class="section-title">Moyen de paiement</h2>

                    <div class="stripe-card-wrapper">
                        <div class="cc-header">
                            <i class="ti ti-credit-card" style="font-size: 1.5rem"></i>
                            <span>Carte Bancaire</span>
                            <img src="https://js.stripe.com/v3/fingerprinted/img/visa-365725566f9578a9589553aa9296d178.svg" alt="Visa" class="card-logo" />
                            <img src="https://js.stripe.com/v3/fingerprinted/img/mastercard-4d8844094130711885b5e41b28c9848f.svg" alt="Mastercard" class="card-logo" />
                        </div>

                        <div v-if="stripeError" class="stripe-error">
                            <i class="ti ti-alert-circle"></i> {{ stripeError }}
                        </div>

                        <div v-if="!stripeReady" class="stripe-loading">
                            <span>Chargement du formulaire sécurisé...</span>
                        </div>

                        <div class="stripe-element-container">
                            <label class="stripe-label">Informations de carte</label>
                            <div id="stripe-card-element" class="stripe-element"></div>
                        </div>
                    </div>

                    <button
                        class="btn-main-action btn-pay w-full mt-6"
                        :disabled="isProcessing || !stripeReady"
                        @click="validerPaiement"
                    >
                        <i class="ti ti-lock" style="margin-right: 8px"></i>
                        {{
                            isProcessing
                                ? "Traitement en cours..."
                                : `Payer ${totalAmount.toFixed(2)} €`
                        }}
                    </button>
                    <p class="text-center text-xs mt-3 text-gray-500">
                        <i class="ti ti-shield-check" style="color: #2d7a4f; margin-right:4px"></i>
                        Paiement chiffré SSL par Stripe. Vos données bancaires ne transitent pas par nos serveurs.
                    </p>
                </div>

                <div class="checkout-summary card">
                    <h2 class="section-title">Récapitulatif</h2>

                    <div class="summary-item mt-4">
                        <div class="item-info">
                            <h3 class="item-title">{{ title }}</h3>
                            <p class="item-desc">{{ displayLabel }}</p>
                        </div>
                        <div class="item-price">
                            {{ baseAmount.toFixed(2) }} €
                        </div>
                    </div>

                    <div
                        v-if="paymentMode === 'annonce'"
                        class="summary-item mt-3"
                    >
                        <div class="item-info">
                            <p
                                class="item-desc"
                                style="color: #2d7a4f; font-weight: 600"
                            >
                                Commission Upcycle (5%)
                            </p>
                        </div>
                        <div class="item-price" style="font-size: 1.1rem">
                            + {{ commissionAmount.toFixed(2) }} €
                        </div>
                    </div>

                    <div class="summary-divider"></div>

                    <div class="summary-total">
                        <span>Total TTC</span>
                        <span class="total-price">{{ totalAmount.toFixed(2) }} €</span>
                    </div>

                    <div
                        v-if="paymentMode === 'subscription'"
                        class="perks-list mt-6"
                    >
                        <h4 class="text-sm font-bold text-gray-700 mb-2">
                            Inclus dans votre offre :
                        </h4>
                        <ul>
                            <li>✔️ Tableau de bord avancé</li>
                            <li>✔️ Analyse d'impact</li>
                            <li>✔️ Alertes priorisées</li>
                        </ul>
                    </div>

                    <div class="stripe-badge mt-6">
                        <span>Propulsé par</span>
                        <svg viewBox="0 0 60 25" xmlns="http://www.w3.org/2000/svg" class="stripe-logo">
                            <path d="M59.64 14.28h-8.06c.19 1.93 1.6 2.55 3.2 2.55 1.64 0 2.96-.37 4.05-.95v3.32a8.33 8.33 0 0 1-4.56 1.1c-4.01 0-6.83-2.5-6.83-7.48 0-4.19 2.39-7.52 6.3-7.52 3.92 0 5.96 3.28 5.96 7.5 0 .4-.04 1.26-.06 1.48zm-5.92-5.62c-1.03 0-2.17.73-2.17 2.58h4.25c0-1.85-1.07-2.58-2.08-2.58zM40.95 20.3c-1.44 0-2.32-.6-2.9-1.04l-.02 4.63-4.12.87V5.57h3.76l.08 1.02a4.7 4.7 0 0 1 3.23-1.29c2.9 0 5.62 2.6 5.62 7.4 0 5.23-2.7 7.6-5.65 7.6zM40 8.95c-.95 0-1.54.34-1.97.81l.02 6.12c.4.44.98.78 1.95.78 1.52 0 2.54-1.65 2.54-3.87 0-2.15-1.04-3.84-2.54-3.84zM28.24 5.57h4.13v14.44h-4.13V5.57zm0-4.7L32.37 0v3.36l-4.13.88V.88zm-4.32 9.35v9.79H19.8V5.57h3.7l.12 1.22c1-1.77 2.96-1.5 3.5-1.32v3.8c-.53-.17-2.35-.47-3.2.98zm-8.55 4.72c0 2.43 2.6 1.68 3.12 1.46v3.36c-.55.3-1.54.54-2.89.54a4.15 4.15 0 0 1-4.27-4.24l.01-13.17 4.02-.86v3.54h3.14V9.1h-3.13v5.85zm-4.91.7c0 2.97-2.31 4.66-5.73 4.66a11.2 11.2 0 0 1-4.46-.93v-3.93c1.38.75 3.1 1.31 4.46 1.31.92 0 1.53-.24 1.53-1C6.26 13.77 0 14.51 0 9.95 0 7.04 2.28 5.3 5.62 5.3c1.36 0 2.72.2 4.09.75v3.88a9.23 9.23 0 0 0-4.1-1.06c-.86 0-1.44.25-1.44.9 0 1.85 6.29.97 6.29 5.88z" fill="#6772e5"/>
                        </svg>
                    </div>
                </div>
            </div>

            <div v-else class="empty-state card">
                <i class="ti ti-shopping-cart-x empty-icon"></i>
                <h2>Aucun article à payer</h2>
                <p class="classic-text mt-2">
                    Votre panier est vide ou le lien est expiré.
                </p>
                <button
                    class="btn-main-action mt-6"
                    @click="$router.push('/profil')"
                >
                    Retourner au profil
                </button>
            </div>
        </main>

        <SiteFooter />
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { loadStripe } from "@stripe/stripe-js";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

// ─── CONFIG ────────────────────────────────────────────────────────────────────
// Remplacez par votre clé publique Stripe (pk_test_... ou pk_live_...)
const STRIPE_PUBLIC_KEY = import.meta.env.VITE_STRIPE_PUBLIC_KEY || "pk_test_REMPLACEZ_PAR_VOTRE_CLE_STRIPE";
const API_URL = "/go";

// ─── ROUTE & ROUTER ────────────────────────────────────────────────────────────
const route = useRoute();
const router = useRouter();

// ─── AUTH ──────────────────────────────────────────────────────────────────────
const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => sessionStorage.getItem("userPrenom") || "Utilisateur");
const currentUserId = computed(
    () => Number(sessionStorage.getItem("id") || sessionStorage.getItem("userId")) || 0
);

// ─── PAYMENT MODE ──────────────────────────────────────────────────────────────
const paymentMode = computed(() => {
    if (route.query.plan_id) return "subscription";
    if (route.query.annonce_id) return "annonce";
    if (route.query.source === "panier") return "panier";
    return null;
});

const displayLabel = computed(() => {
    switch (paymentMode.value) {
        case "subscription": return `Abonnement ${route.query.nom_plan || "Pro"}`;
        case "panier":       return "Commande globale";
        case "annonce":      return "Paiement annonce";
        default:             return "Transaction";
    }
});

const title = computed(() => route.query.title || route.query.nom_plan || "Article");
const rawAmount = computed(() => route.query.amount || route.query.prix || "0");
const baseAmount = computed(() => parseFloat(rawAmount.value) || 0);
const commissionAmount = computed(() => paymentMode.value === "annonce" ? baseAmount.value * 0.05 : 0);
const totalAmount = computed(() => baseAmount.value + commissionAmount.value);

// ─── STRIPE ────────────────────────────────────────────────────────────────────
const stripeReady = ref(false);
const stripeError = ref("");
const isProcessing = ref(false);

let stripe = null;
let cardElement = null;

onMounted(async () => {
    if (!paymentMode.value) return;

    try {
        stripe = await loadStripe(STRIPE_PUBLIC_KEY);
        if (!stripe) throw new Error("Impossible de charger Stripe.");

        const elements = stripe.elements({
            locale: "fr",
            appearance: {
                theme: "stripe",
                variables: {
                    colorPrimary: "#2d7a4f",
                    colorText: "#1a1a1a",
                    colorDanger: "#df1b41",
                    fontFamily: "inherit",
                    borderRadius: "8px",
                },
            },
        });

        cardElement = elements.create("card", {
            style: {
                base: {
                    fontSize: "15px",
                    color: "#1a1a1a",
                    "::placeholder": { color: "#aab7c4" },
                },
                invalid: { color: "#df1b41" },
            },
            hidePostalCode: true,
        });

        cardElement.mount("#stripe-card-element");

        cardElement.on("ready", () => { stripeReady.value = true; });
        cardElement.on("change", (event) => {
            stripeError.value = event.error ? event.error.message : "";
        });
    } catch (err) {
        stripeError.value = "Erreur lors du chargement du module de paiement.";
        console.error(err);
    }
});

onUnmounted(() => {
    if (cardElement) cardElement.destroy();
});

// ─── PAYMENT FLOW ──────────────────────────────────────────────────────────────
const validerPaiement = async () => {
    if (!currentUserId.value) {
        alert("Veuillez vous reconnecter.");
        router.push("/connexion");
        return;
    }

    if (!stripe || !cardElement) {
        stripeError.value = "Le module de paiement n'est pas prêt. Veuillez patienter.";
        return;
    }

    isProcessing.value = true;
    stripeError.value = "";

    try {
        // 1. Demander un PaymentIntent au backend
        const intentRes = await fetch(`${API_URL}/paiement/intent`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${sessionStorage.getItem("userToken") || ""}`,
            },
            body: JSON.stringify({
                montant: Math.round(totalAmount.value * 100), // en centimes
                mode: paymentMode.value,
                plan_id: parseInt(route.query.plan_id, 10) || 0,
                annonce_id: parseInt(route.query.annonce_id, 10) || 0,
                user_id: currentUserId.value,
            }),
        });

        if (!intentRes.ok) {
            const msg = await intentRes.text();
            throw new Error(msg || `Erreur serveur (${intentRes.status})`);
        }

        const { client_secret } = await intentRes.json();

        // 2. Confirmer le paiement côté Stripe
        const { error: stripeConfirmError, paymentIntent } = await stripe.confirmCardPayment(
            client_secret,
            { payment_method: { card: cardElement } }
        );

        if (stripeConfirmError) {
            throw new Error(stripeConfirmError.message);
        }

        if (paymentIntent.status !== "succeeded") {
            throw new Error(`Statut inattendu : ${paymentIntent.status}`);
        }

        // 3. Notifier le backend pour finaliser la commande
        await finaliserCommande(paymentIntent.id);

        router.push("/profil?paiement=success");

    } catch (err) {
        stripeError.value = err.message || "Une erreur est survenue lors du paiement.";
        console.error("Erreur paiement:", err);
    } finally {
        isProcessing.value = false;
    }
};

const finaliserCommande = async (stripePaymentId) => {
    const planId    = parseInt(route.query.plan_id, 10) || 0;
    const annonceId = parseInt(route.query.annonce_id, 10) || 0;

    let endpoint = "";
    let body = {};

    if (paymentMode.value === "annonce") {
        endpoint = `/annonces/${annonceId}/acheter`;
        body = { id_acheteur: currentUserId.value, montant_paye: baseAmount.value, stripe_payment_id: stripePaymentId };
    } else if (paymentMode.value === "subscription") {
        endpoint = `/users/${currentUserId.value}/abonnement/souscrire`;
        body = { plan_id: planId, stripe_payment_id: stripePaymentId };
    } else if (paymentMode.value === "panier") {
        endpoint = `/users/${currentUserId.value}/checkout`;
        body = { stripe_payment_id: stripePaymentId };
    } else {
        throw new Error("Mode de paiement inconnu.");
    }

    const res = await fetch(`${API_URL}${endpoint}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${sessionStorage.getItem("userToken") || ""}`,
        },
        body: JSON.stringify(body),
    });

    if (!res.ok) {
        const msg = await res.text();
        throw new Error(msg || `Erreur lors de la finalisation (${res.status})`);
    }
};
</script>

<style scoped>
.bg-light {
    background-color: #f7f9f7;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
}
.page-container {
    flex: 1;
    padding: 20px;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
}
.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
    border-bottom: 1px solid #e5ede7;
    padding-bottom: 1.5rem;
}
.sidebar__category2 {
    margin: 0;
    color: #a0ada7;
    font-size: 0.65rem;
    font-weight: bold;
    letter-spacing: 1px;
    text-transform: uppercase;
}
.hero-title1 {
    font-size: 2.2rem;
    font-weight: 900;
    color: #1a1a1a;
    margin: 0.5rem 0 0 0;
}
.classic-text {
    font-size: 0.95rem;
    color: #666;
    margin: 0;
}
.btn-secondary {
    padding: 8px 16px;
    border-radius: 8px;
    border: 1px solid #ddd;
    background: white;
    cursor: pointer;
    font-weight: bold;
}
.btn-secondary:hover {
    background: #f5f5f5;
}
.btn-main-action {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: #2d7a4f;
    color: white;
    padding: 12px 20px;
    border-radius: 10px;
    font-weight: bold;
    border: none;
    cursor: pointer;
    transition: 0.2s;
}
.btn-main-action:hover:not(:disabled) {
    background: #1b4d31;
}
.btn-main-action:disabled {
    opacity: 0.7;
    cursor: not-allowed;
}
.w-full { width: 100%; }
.mt-6   { margin-top: 1.5rem; }
.mt-4   { margin-top: 1rem; }
.mt-3   { margin-top: 0.75rem; }
.text-center { text-align: center; }
.text-xs { font-size: 0.75rem; }
.text-gray-500 { color: #6b7280; }

.card {
    background: white;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.checkout-grid {
    display: grid;
    grid-template-columns: 1.5fr 1fr;
    gap: 2rem;
    align-items: start;
}
@media (max-width: 800px) {
    .checkout-grid { grid-template-columns: 1fr; }
}

.section-title {
    font-size: 1.2rem;
    font-weight: 800;
    color: #1a1a1a;
    margin-bottom: 1rem;
    border-bottom: 2px solid #f0f0f0;
    padding-bottom: 10px;
}

/* ── Stripe wrapper ── */
.stripe-card-wrapper {
    background: #fafdfb;
    border: 1px solid #e5ede7;
    padding: 1.5rem;
    border-radius: 12px;
}
.cc-header {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #2d7a4f;
    font-weight: bold;
    margin-bottom: 1.25rem;
}
.card-logo {
    height: 20px;
    margin-left: 4px;
}
.stripe-label {
    display: block;
    font-size: 0.85rem;
    font-weight: 600;
    color: #555;
    margin-bottom: 8px;
}
.stripe-element {
    padding: 12px;
    border: 1px solid #ddd;
    border-radius: 8px;
    background: white;
    transition: border-color 0.2s, box-shadow 0.2s;
}
.stripe-element:focus-within {
    border-color: #2d7a4f;
    box-shadow: 0 0 0 3px rgba(45, 122, 79, 0.1);
}
.stripe-loading {
    text-align: center;
    color: #999;
    font-size: 0.9rem;
    padding: 12px;
}
.hidden { display: none; }

.stripe-error {
    background: #fff0f0;
    border: 1px solid #fca5a5;
    color: #dc2626;
    border-radius: 8px;
    padding: 10px 14px;
    font-size: 0.875rem;
    margin-bottom: 1rem;
    display: flex;
    align-items: center;
    gap: 6px;
}

.btn-pay {
    font-size: 1.1rem;
    padding: 14px;
}

/* ── Summary ── */
.summary-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.item-title {
    font-weight: 800;
    font-size: 1.1rem;
    color: #1a1a1a;
    margin: 0;
}
.item-desc {
    font-size: 0.85rem;
    color: #666;
    margin: 4px 0 0 0;
}
.item-price {
    font-weight: bold;
    font-size: 1.2rem;
    color: #2d7a4f;
}
.summary-divider {
    height: 1px;
    background: #eee;
    margin: 1.5rem 0;
}
.summary-total {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 1.2rem;
    font-weight: 800;
}
.total-price {
    font-size: 1.8rem;
    color: #1a1a1a;
}
.perks-list ul {
    list-style: none;
    padding: 0;
    margin: 0;
}
.perks-list li {
    margin-bottom: 8px;
    font-size: 0.9rem;
    color: #555;
}

.stripe-badge {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 0.75rem;
    color: #999;
    justify-content: center;
}
.stripe-logo {
    width: 48px;
    height: auto;
}

/* ── Empty state ── */
.empty-state {
    text-align: center;
    padding: 4rem 2rem;
}
.empty-icon {
    font-size: 4rem;
    color: #ccc;
    margin-bottom: 1rem;
}
.empty-state h2 {
    font-size: 1.5rem;
    color: #333;
    margin: 0;
}
</style>
