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

                    <div class="credit-card-mock">
                        <div class="cc-header">
                            <i
                                class="ti ti-credit-card"
                                style="font-size: 1.5rem"
                            ></i>
                            <span>Carte Bancaire</span>
                        </div>

                        <div class="form-group mt-4">
                            <label>Nom sur la carte</label>
                            <input
                                type="text"
                                placeholder="Ex: Jean Dupont"
                                class="form-input"
                            />
                        </div>

                        <div class="form-group mt-3">
                            <label>Numéro de carte</label>
                            <input
                                type="text"
                                placeholder="0000 0000 0000 0000"
                                class="form-input font-mono"
                                maxlength="19"
                            />
                        </div>

                        <div class="grid-2 mt-3">
                            <div class="form-group">
                                <label>Date d'expiration</label>
                                <input
                                    type="text"
                                    placeholder="MM/AA"
                                    class="form-input"
                                />
                            </div>
                            <div class="form-group">
                                <label>CVC</label>
                                <input
                                    type="text"
                                    placeholder="123"
                                    class="form-input"
                                    maxlength="3"
                                />
                            </div>
                        </div>
                    </div>

                    <button
                        class="btn-main-action btn-pay w-full mt-6"
                        :disabled="isProcessing"
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
                        Paiement chiffré et sécurisé.
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
                                Comission Upcycle (5%)
                            </p>
                        </div>
                        <div class="item-price" style="font-size: 1.1rem">
                            + {{ commissionAmount.toFixed(2) }} €
                        </div>
                    </div>

                    <div class="summary-divider"></div>

                    <div class="summary-total">
                        <span>Total TTC</span>
                        <span class="total-price"
                            >{{ totalAmount.toFixed(2) }} €</span
                        >
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
import { ref, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";

const route = useRoute();
const router = useRouter();
const API_URL = "/go";

const paymentMode = computed(() => {
    if (route.query.plan_id) return "subscription";
    if (route.query.annonce_id) return "annonce";
    if (route.query.source === "panier") return "panier";
    return null;
});

const displayLabel = computed(() => {
    switch (paymentMode.value) {
        case "subscription":
            return `Abonnement ${route.query.nom_plan || "Pro"}`;
        case "panier":
            return "Commande globale";
        case "annonce":
            return "Paiement annonce";
        default:
            return "Transaction";
    }
});

const title = computed(
    () => route.query.title || route.query.nom_plan || "Article",
);

const rawAmount = computed(() => route.query.amount || route.query.prix || "0");
const baseAmount = computed(() => parseFloat(rawAmount.value) || 0);

const commissionAmount = computed(() => {
    if (paymentMode.value === "annonce") {
        return baseAmount.value * 0.05;
    }
    return 0;
});

const totalAmount = computed(() => baseAmount.value + commissionAmount.value);

const isProcessing = ref(false);
const currentUserId = computed(
    () =>
        Number(
            sessionStorage.getItem("id") || sessionStorage.getItem("userId"),
        ) || 0,
);
const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(
    () => sessionStorage.getItem("userPrenom") || "Utilisateur",
);

const validerPaiement = async () => {
    if (!currentUserId.value) {
        alert("Veuillez vous reconnecter.");
        router.push("/connexion");
        return;
    }

    isProcessing.value = true;

    try {
        let endpoint = "";
        let requestBody = {};

        const planId = parseInt(route.query.plan_id, 10) || 0;
        const annonceId = parseInt(route.query.annonce_id, 10) || 0;

        if (paymentMode.value === "annonce") {
            if (!annonceId) throw new Error("ID de l'annonce manquant.");
            endpoint = `/annonces/${annonceId}/acheter`;

            requestBody = {
                id_acheteur: currentUserId.value,
                montant_paye: baseAmount.value,
            };
        } else if (paymentMode.value === "subscription") {
            if (!planId) throw new Error("ID de l'abonnement manquant.");
            endpoint = `/users/${currentUserId.value}/abonnement/souscrire`;
            requestBody = { plan_id: planId };
        } else if (paymentMode.value === "panier") {
            endpoint = `/users/${currentUserId.value}/checkout`;
            requestBody = {};
        } else {
            throw new Error("Mode de paiement inconnu.");
        }

        const res = await fetch(`${API_URL}${endpoint}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${sessionStorage.getItem("userToken") || ""}`,
            },
            body: JSON.stringify(requestBody),
        });

        if (res.ok) {
            alert("Paiement réussi !");
            router.push("/profil");
        } else {
            const errMsg = await res.text();
            throw new Error(errMsg || `Erreur serveur (${res.status})`);
        }
    } catch (error) {
        console.error("Erreur paiement:", error);
        alert("Impossible de finaliser le paiement: " + error.message);
    } finally {
        isProcessing.value = false;
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
.w-full {
    width: 100%;
}
.mt-6 {
    margin-top: 1.5rem;
}
.mt-4 {
    margin-top: 1rem;
}
.mt-3 {
    margin-top: 0.75rem;
}
.text-center {
    text-align: center;
}

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
    .checkout-grid {
        grid-template-columns: 1fr;
    }
}

.section-title {
    font-size: 1.2rem;
    font-weight: 800;
    color: #1a1a1a;
    margin-bottom: 1rem;
    border-bottom: 2px solid #f0f0f0;
    padding-bottom: 10px;
}

.credit-card-mock {
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
    margin-bottom: 1rem;
}
.form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.form-group label {
    font-size: 0.85rem;
    font-weight: 600;
    color: #555;
}
.form-input {
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 0.95rem;
}
.form-input:focus {
    outline: none;
    border-color: #2d7a4f;
    box-shadow: 0 0 0 3px rgba(45, 122, 79, 0.1);
}
.font-mono {
    font-family: monospace;
    letter-spacing: 1px;
}
.grid-2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
}
.btn-pay {
    font-size: 1.1rem;
    padding: 14px;
}

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
